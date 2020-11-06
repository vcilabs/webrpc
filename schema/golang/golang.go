package golang

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go/types"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/webrpc/webrpc/schema"
)

var (
	schemaMessageTypeStruct  = schema.MessageType("struct")
	schemaMessageTypeEnum    = schema.MessageType("enum")
	schemaMessageTypeAdvance = schema.MessageType("advance")
)

//Parser struct manages the parsing of go files
type Parser struct {
	parent  *Parser
	imports map[string]struct{}

	reader *schema.Reader
}

//NewParser returns Parser
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

//goparse parse the go file as below:
//Step1: Read the file added by user(via -schema flag), parse it using go/Parser package
//       and add the content to an in-memory package name same as schema file name
//Step2: Reads the imports from file and add to Imports type of WebRPCSchema
//Step3: Split the parsed file content on basis of types and iterate to read interface, structs, advanced types
//       a) Interface/Service contanis input and output arguments. These are handled via buildArgumentsList()
//       b) Struct contains messages/datatypes and goparser also handles embedded structs
//       c) If a type is neither interface nor struct then it is comes under "advanced types" and is handled with message type "advance"
//       d) All the values are returned to WebRPCSchema and is used to populate the template
func (p *Parser) goparse(path string) (*schema.WebRPCSchema, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Error in reading the file")
	}
	var parsedFile = string(data)
	fset := token.NewFileSet()
	fileName := filepath.Base(path)

	// Parse the input string, []byte, or io.Reader,
	// recording position information in fset.ParseFile returns an *ast.File, a syntax tree.
	f, err := parser.ParseFile(fset, fileName, parsedFile, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the file: "+fileName)
	}

	// A Config controls various options of the type checker.
	// The defaults work fine except for one setting:
	// we must specify how to deal with imports.
	conf := types.Config{Importer: importer.Default()}

	// Type-check the package containing only file f.
	// Check returns a *types.Package.
	pkg, err := conf.Check("cmd/"+fileName, fset, []*ast.File{f}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid type in the file: "+fileName)
	}
	s := &schema.WebRPCSchema{}
	ext := filepath.Ext(path)

	//Add the schema type to be used in var_type.go to parse the go maps
	s.SchemaType = strings.TrimPrefix(ext, ".")
	// This reads the imports from the go file and adds to Imports schema
	additionalimports := pkg.Imports()
	for _, additionalimport := range additionalimports {
		regexForListOfImports := regexp.MustCompile(`\(.*?\)`)
		listOfImports := regexForListOfImports.FindAllString(additionalimport.String(), -1)
		listOfImports[0] = strings.Trim(listOfImports[0], "[(")
		listOfImports[0] = strings.Trim(listOfImports[0], ")]")
		s.Imports = append(s.Imports, &schema.Import{
			Path: listOfImports[0],
		})
	}
	//goTypes holds the types information for a given go file
	//It includes type interface, type struct
	goTypes := strings.Split(pkg.Scope().String(), "type cmd/"+fileName+".")
	elementMap := make(map[string]string)
	methods := []*schema.Method{}
	//Sort the types in ascending order on basis of lenght of string
	sort.Sort(ByLen(goTypes))
	for _, goType := range goTypes {
		//Replace the additional string with blank so as to get the valid desired content while parsing the types
		goType = strings.ReplaceAll(goType, "cmd/"+fileName+".", "")
		//Read the type Interface and update the name, inputs and outputs
		if strings.Contains(goType, " interface") {
			elementMap["interface"] = goType
			interfaceNameField := strings.Fields(elementMap["interface"])
			interfaceName := interfaceNameField[0]
			s.Services = append(s.Services, &schema.Service{Name: schema.VarName(interfaceName)})
			for _, method := range interfaceAllMethodNames(goType) {
				inputs, err := buildArgumentsList(s, goType, method, "isInputArgs")
				if err != nil {
					return nil, err
				}
				outputs, err := buildArgumentsList(s, goType, method, "isOutputArgs")
				if err != nil {
					return nil, err
				}
				methods = append(methods, &schema.Method{
					Name:    schema.VarName(method),
					Inputs:  inputs,
					Outputs: outputs,
				})
			}
			interfaceDef := s.GetServiceByName(interfaceName)
			interfaceDef.Methods = methods
		} else if strings.Contains(goType, "struct") {
			// Read the struct name and update the Messages
			elementMap["struct"] = goType
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
			//Read the struct fields and update the MessageFields
			for _, def := range fieldsOfStruct(goType) {
				if len(def) < 1 {
					continue
				}
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
				structDef.Fields = append(structDef.Fields, field)
			}
		} else if (!strings.Contains(goType, " interface") || !strings.Contains(goType, "struct")) && !strings.Contains(goType, "cmd/"+fileName) {
			// Handle advanced types
			splitDataMap := strings.Split(goType, " ")
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
				Type:     schemaMessageTypeEnum,
				EnumType: &enumType,
			})
		}
	}
	return s, nil
}

// interfaceAllMethodNames retuns all the method names present in an interface
func interfaceAllMethodNames(goType string) []string {
	var listOfAllInterfaceMethods []string
	re := regexp.MustCompile(`\{.*\}`)
	submatchall := re.FindAllString(goType, -1)
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

//buildArgumentsList generate the list of arguments for interface. It checks for type of arguments
//For instance: Interface has method declation as "BorrowBook(ctx context.Context, BookID int64) (string, error)"
//Here input args are (ctx context.Context, BookID int64) and output args are (string, error). So this function list them on basis of types
//If checkType is "isInputArgs", trim "ctx context.Context" as it is common for  RPC and append the name and variable type for rest of the input arguments
//If checkType is "isOutputArgs", skip argument "error" as it is common for RPC and append the type i.e string, int etc for rest of the output arguments
//NOTE: for input we take care for both Name and type. But since in output we only have return type so we read only type
//Arguments used are as follows:
//         a) WebRPC schema object needed for ParseVarTypeExpr function
//         b) goType holds the input/output arguments of method from interface
//         c) method holds the method name. It is only used to filter the arguments and not parse all the arguments of interface.(Added only to save parsing time)
//         d) checkType is a string that check we need input args or output args
func buildArgumentsList(s *schema.WebRPCSchema, goType string, method string, checkType string) ([]*schema.MethodArgument, error) {
	output := []*schema.MethodArgument{}
	interfaceRegex := regexp.MustCompile(`\{.*\}`)
	argsRegex := regexp.MustCompile(`\(.*?\)`)
	argumentMatch := interfaceRegex.FindAllString(goType, -1)
	for _, argList := range argumentMatch {
		argList = strings.Trim(argList, "[{")
		argList = strings.Trim(argList, "}]")
		result := strings.Split(argList, ";")
		for _, v := range result {
			if strings.Contains(v, method) {
				if checkType == "isInputArgs" {
					//Read name and types and append to methodArgument
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
					//Read the types and append to methodArgument
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
								// Make a Regex to say we only want letters and numbers
								reg, err := regexp.Compile("[^a-zA-Z0-9]+")
								if err != nil {
									log.Fatal(err)
								}
								responseArg := reg.ReplaceAllString(resultsNew, "")
								responseArg = strings.ToLower(responseArg)
								methodArgument := &schema.MethodArgument{
									Name: schema.VarName(responseArg),
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

//fieldsOfStruct returs the content of struct.
//For example "Author struct {ID int64, ... }" will return "ID int64", "..." as list of Fields
func fieldsOfStruct(goType string) []string {
	var listOfFields []string
	structRegex := regexp.MustCompile(`\{.*\}`)
	argumentMatch := structRegex.FindAllString(goType, -1)
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
//Case: string:= ["abcd", "p", "xyz", "ab"]
//o/p:  string:= ["p", "ab", "xyz", "abcd"]
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
