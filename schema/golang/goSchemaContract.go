package golang

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"strings"

	"go/types"
	"io/ioutil"
	"log"

	"github.com/webrpc/webrpc/schema"
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
	//fmt.Println("Contents of file:", parsedFile)

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
		GoInterfaceScope: []string{},
		GoStructScope:    []string{},
		GoDataTypeScope:  []string{},
	}

	splitString := strings.Split(pkg.Scope().String(), "type cmd/parsedFile.go.")
	elementMap := make(map[string]string)
	for _, dataMap := range splitString {
		dataMap = strings.ReplaceAll(dataMap, "cmd/parsedFile.go.", "")
		if strings.Contains(dataMap, "interface") {
			elementMap["interface"] = dataMap
			s.GoInterfaceScope = append(s.GoInterfaceScope, elementMap["interface"])
		} else if strings.Contains(dataMap, "struct") {
			elementMap["struct"] = dataMap
			s.GoStructScope = append(s.GoStructScope, elementMap["struct"])
		} else {
			elementMap["datatype"] = dataMap
			s.GoDataTypeScope = append(s.GoDataTypeScope, elementMap["datatype"])
		}
	}
	return s, nil
}
