package golang

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strings"

	"go/types"
	"io/ioutil"
	"log"

	"github.com/webrpc/webrpc/schema"
)

var (
	schemaMessageTypeStruct  = schema.MessageType("struct")
	schemaMessageTypeAdvance = schema.MessageType("advance")
)

type Parser struct {
	parent  *Parser
	imports map[string]struct{}

	reader *schema.Reader
}

func NewParser(r *schema.Reader) *Parser {
	return &Parser{
		reader: r,
		imports: map[string]struct{}{
			// this file imports itself
			r.File: struct{}{},
		},
	}
}

//Parse parses the go file
func (p *Parser) Parse(path string) (*schema.WebRPCSchema, error) {
	s, err := p.goparse(path)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (p *Parser) goparse(path string) (*schema.WebRPCSchema, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}
	var parsedFile = string(data)
	fset := token.NewFileSet()

	// Parse the input string, []byte, or io.Reader,
	// recording position information in fset.ParseFile returns an *ast.File, a syntax tree.
	f, err := parser.ParseFile(fset, "parsedFile.go", parsedFile, 0)
	if err != nil {
		log.Fatal(err) // parse error
	}

	// A Config controls various options of the type checker.
	// The defaults work fine except for one setting:
	// we must specify how to deal with imports.
	conf := types.Config{Importer: importer.Default()}

	// Type-check the package containing only file f.
	// Check returns a *types.Package.
	pkg, err := conf.Check("cmd/parsedFile.go", fset, []*ast.File{f}, nil)
	if err != nil {
		log.Fatal(err) // type error
	}

	//TODO: update the code to add proper schema
	s := &schema.WebRPCSchema{
		GoInterface: []*schema.GoInterface{},
	}
	//imports := pkg.Imports()
	splitString := strings.Split(pkg.Scope().String(), "type cmd/parsedFile.go.")
	elementMap := make(map[string]string)
	methods := []*schema.Method{}
	sort.Sort(ByLen(splitString))
	s.SchemaType = "go"
	for _, dataMap := range splitString {
		dataMap = strings.ReplaceAll(dataMap, "cmd/parsedFile.go.", "")
		if strings.Contains(dataMap, "interface") {
			elementMap["interface"] = dataMap
			interfaceNameField := strings.Fields(elementMap["interface"])
			interfaceName := interfaceNameField[0]
			s.GoInterface = append(s.GoInterface, &schema.GoInterface{Name: schema.VarName(interfaceName)})
			for _, method := range interfaceAllMethodNames(dataMap) {
				inputs, err := buildArgumentsList(s, dataMap, method, "isInputArgs")
				if err != nil {
					return nil, err
				}
				outputs, err := buildArgumentsList(s, dataMap, method, "isOutputArgs")
				if err != nil {
					return nil, err
				}
				methods = append(methods, &schema.Method{
					Name:    schema.VarName(method),
					Inputs:  inputs,
					Outputs: outputs,
				})
			}
			interfaceDef := s.GetInterfaceByName(interfaceName)
			interfaceDef.Methods = methods
		} else if strings.Contains(dataMap, "struct") {
			elementMap["struct"] = dataMap
			StructNameField := strings.Fields(elementMap["struct"])
			structName := StructNameField[0]
			s.Messages = append(s.Messages, &schema.Message{
				Name: schema.VarName(structName),
				Type: schemaMessageTypeStruct,
			})
			name := schema.VarName(structName)
			structDef := s.GetStructByName(string(name))
			if structDef == nil {
				return nil, fmt.Errorf("unexpected error, could not find definition for: %v", name)
			}
			for _, def := range fieldsOfStruct(dataMap) {
				splitField := strings.Split(def, " ")
				fieldName, fieldType := splitField[0], splitField[1]

				var varType schema.VarType
				err := schema.ParseVarTypeExpr(s, fieldType, &varType)
				if err != nil {
					return nil, fmt.Errorf("unknown data type: %v", fieldType)
				}

				field := &schema.MessageField{
					Name: schema.VarName(fieldName),
					Type: &varType,
				}
				//TODO: Metadata addition
				// for _, meta := range def.Meta() {
				// 	key, val := meta.Left().String(), meta.Right().String()
				// 	field.Meta = append(field.Meta, schema.MessageFieldMeta{
				// 		key: val,
				// 	})
				// }
				structDef.Fields = append(structDef.Fields, field)
			}
		} else if (!strings.Contains(dataMap, "interface") || !strings.Contains(dataMap, "struct")) && !strings.Contains(dataMap, "cmd/parsedFile.go") {
			splitDataMap := strings.Split(dataMap, " ")
			keyName := splitDataMap[0]
			typeArgRegex := regexp.MustCompile(`^[\w.]+`)
			splitDataMapArgument := typeArgRegex.FindAllString(splitDataMap[1], 1)
			typeName := splitDataMapArgument[0]
			var enumType schema.VarType
			err := schema.ParseVarTypeExpr(s, typeName, &enumType)
			if err != nil {
				return nil, fmt.Errorf("unknown data type: %v", typeName)
			}
			s.Messages = append(s.Messages, &schema.Message{
				Name:     schema.VarName(keyName),
				Type:     schemaMessageTypeAdvance,
				EnumType: &enumType,
			})
		}
	}
	return s, nil
}

func interfaceAllMethodNames(dataMap string) []string {
	var listOfAllInterfaceMethods []string
	re := regexp.MustCompile(`\{.*?\}`)
	submatchall := re.FindAllString(dataMap, -1)

	for _, element := range submatchall {
		element = strings.Trim(element, "[{")
		element = strings.Trim(element, "}]")
		result := strings.Split(element, ";")
		for _, v := range result {
			methodName := strings.Split(v, "(")[0]
			methodName = strings.TrimSpace(methodName)
			listOfAllInterfaceMethods = append(listOfAllInterfaceMethods, methodName)
		}
	}
	return listOfAllInterfaceMethods
}

func buildArgumentsList(s *schema.WebRPCSchema, dataMap string, method string, checkType string) ([]*schema.MethodArgument, error) {
	output := []*schema.MethodArgument{}
	interfaceRegex := regexp.MustCompile(`\{.*?\}`)
	argsRegex := regexp.MustCompile(`\(.*?\)`)
	argumentMatch := interfaceRegex.FindAllString(dataMap, -1)

	for _, argList := range argumentMatch {
		argList = strings.Trim(argList, "[{")
		argList = strings.Trim(argList, "}]")
		result := strings.Split(argList, ";")
		for _, v := range result {
			if strings.Contains(v, method) {
				if checkType == "isInputArgs" {
					methodArgs := argsRegex.FindAllString(v, 1)
					for _, element1 := range methodArgs {
						element1 = strings.Trim(element1, "(ctx context.Context,")
						element1 = strings.Trim(element1, ")")
						if len(element1) > 0 {
							result1 := strings.Split(element1, ",")
							for _, resultsNew := range result1 {
								resultsNew = strings.TrimSpace(resultsNew)
								resultbreak := strings.Split(resultsNew, " ")
								var varType schema.VarType
								err := schema.ParseVarTypeExpr(s, resultbreak[1], &varType)
								if err != nil {
									return nil, fmt.Errorf("unknown data type: %v", resultbreak[1])
								}
								methodArgument := &schema.MethodArgument{
									Name: schema.VarName(resultbreak[0]),
									Type: &varType,
								}
								output = append(output, methodArgument)
							}
						}
					}
				} else if checkType == "isOutputArgs" {
					methodArgs := argsRegex.FindAllString(v, -1)
					methodArgs = methodArgs[1:]
					for _, element1 := range methodArgs {
						element1 = strings.Trim(element1, "(")
						element1 = strings.Trim(element1, ")")
						if len(element1) > 0 {
							result1 := strings.Split(element1, ",")
							for _, resultsNew := range result1 {
								if strings.Contains(resultsNew, "error") {
									continue
								}
								resultsNew = strings.TrimSpace(resultsNew)
								var varType schema.VarType
								err := schema.ParseVarTypeExpr(s, resultsNew, &varType)
								if err != nil {
									return nil, fmt.Errorf("unknown data type: %v", resultsNew)
								}
								methodArgument := &schema.MethodArgument{
									Name: schema.VarName(resultsNew),
									Type: &varType,
								}
								output = append(output, methodArgument)
							}
						}
					}
					return output, nil
				}
			}
		}
	}
	return output, nil
}

func fieldsOfStruct(dataMap string) []string {
	var listOfFields []string
	structRegex := regexp.MustCompile(`\{.*?\}`)
	argumentMatch := structRegex.FindAllString(dataMap, -1)

	for _, argList := range argumentMatch {
		argList = strings.Trim(argList, "[{")
		argList = strings.Trim(argList, "}]")
		result := strings.Split(argList, ";")
		for _, v := range result {
			v = strings.TrimSpace(v)
			listOfFields = append(listOfFields, v)
		}
	}
	return listOfFields
}

//ByLen Sort the string in ascending order by count of each string
//Case: string:= ["abcd", "a", "abc", "ab"]
//o/p:  string:= ["a", "ab", "abc", "abcd"]
type ByLen []string

func (a ByLen) Len() int {
	return len(a)
}

func (a ByLen) Less(i, j int) bool {
	return len(a[i]) < len(a[j])
}

func (a ByLen) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
