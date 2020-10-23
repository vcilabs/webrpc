package golang

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"

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

//GoParse parses the go file
func (p *Parser) Parse(path string) (*schema.WebRPCSchema, error) {
	s, err := p.goparse(path)
	if err != nil {
		return nil, err
	}
	fmt.Println("I am outsde s-->", s)

	return s, nil
}

func (p *Parser) goparse(path string) (*schema.WebRPCSchema, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}
	var parsedFile = string(data)
	fmt.Println("Contents of file:", parsedFile)

	fset := token.NewFileSet()

	// Parse the input string, []byte, or io.Reader,
	// recording position information in fset.
	// ParseFile returns an *ast.File, a syntax tree.
	//f, err := parser.ParseFile(fset, filepath.Base(path), strings.TrimSuffix(filepath.Base(path), ".go"), 0)
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

	fmt.Printf("Package  %q\n", pkg.Path())
	fmt.Printf("Name:    %s\n", pkg.Name())
	fmt.Printf("Imports: %s\n", pkg.Imports())
	fmt.Printf("Scope:   %s\n", pkg.Scope())

	//TODO: update the code to add proper schema
	s := &schema.WebRPCSchema{
		//GoInterface: []*schema.GoInterface{},
		GoInterfaceScope: []string{},
	}
	s.GoInterfaceScope = append(s.GoInterfaceScope, pkg.Scope().String())

	return s, nil
}

func printBasicType(kind types.BasicKind) string {
	switch kind {
	case types.Bool:
		return "bool"
	case types.Int64:
		return "int64"
	case types.String:
		return "string"
	default:
		return fmt.Sprintf("%v", kind)
	}
}
