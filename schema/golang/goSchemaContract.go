package golang

import (
	"fmt"
	"go/types"

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
func (p *Parser) GoParse() (*schema.WebRPCSchema, error) {
	s, err := p.goparse()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *Parser) goparse() (*schema.WebRPCSchema, error) {
	q, err := newParser(p.reader)
	fmt.Println("I am after parser")
	if err != nil {
		return nil, err
	}
	//flag.Parse()
	//TODO: update the code to add proper schema
	s := &schema.WebRPCSchema{
		GoInterface: []*schema.GoInterface{},
	}
	// Many tools pass their command-line arguments (after any flags)
	// uninterpreted to packages.Load so that it can interpret them
	// according to the conventions of the underlying build system.

	// cfg := &packages.Config{
	// 	Mode:  packages.NeedFiles | packages.NeedTypes | packages.NeedName | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedName | packages.NeedSyntax,
	// 	Tests: false,
	// }
	// fmt.Println("cfg", cfg)
	// pkgs, err := packages.Load(cfg, flag.Args()...)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "load: %v\n", err)
	// 	os.Exit(1)
	// }
	// if packages.PrintErrors(pkgs) > 0 {
	// 	os.Exit(1)
	// }
	// // Print the names of the source files
	// // for each package listed on the command line.
	// for _, pkg := range pkgs {
	// 	//fmt.Println(pkg.ID, pkg.GoFiles)
	// 	//scope := pkg.Types.
	// 	scope := pkg.Types.Scope()
	// 	//fmt.Printf("%#v\n", scope.Names())
	// 	for _, name := range scope.Names() {
	// 		obj := scope.Lookup(name)
	// 		switch item := obj.Type().Underlying().(type) {
	// 		case *types.Interface:
	// 			//fmt.Printf("interface %s: %#v\n", name, item)
	// 		case *types.Struct:
	// 			for i := 0; i < item.NumFields(); i++ {
	// 				field := item.Field(i)
	// 				typ := field.Type().Underlying()
	// 				fmt.Printf("%s.%s: %v (%v) -- underlying %v\n", name, field.Name(), field.Type(), field.Pkg(), typ)
	// 			}
	// 		case *types.Basic:
	// 			fmt.Printf("basic type %s %v %v %#v\n", name, item.Name(), printBasicType(item.Kind()), item)
	// 		default:
	// 			fmt.Printf("what is this? %s: %T\n", name, item)
	// 		}
	// 	}
	// }
	// return s, nil

	// pushing Interfaces (1st pass)
	for _, goInterface := range q.root.Interfaces() {
		// push service
		s.GoInterface = append(s.GoInterface, &schema.GoInterface{
			Name: schema.VarName(goInterface.Name().String()),
		})
	}

	for _, goInterface := range q.root.Interfaces() {
		fmt.Println("I am inside goInterface For loop")
		methods := []*schema.Method{}

		for _, method := range goInterface.Methods() {

			fmt.Println("I am inside method For loop")
			// inputs, err := buildArgumentsList(s, method.Inputs())
			// if err != nil {
			// 	return nil, err
			// }

			// outputs, err := buildArgumentsList(s, method.Outputs())
			// if err != nil {
			// 	return nil, err
			// }

			// push method
			methods = append(methods, &schema.Method{
				Name: schema.VarName(method.Name().String()),
				//StreamInput:  method.StreamInput(),
				//StreamOutput: method.StreamOutput(),
				//Inputs:       inputs,
				//Outputs:      outputs,
			})
		}

		interfaceDef := s.GetInterfaceByName(goInterface.Name().String())
		fmt.Println("interfaceDef-->", interfaceDef)
		interfaceDef.Methods = methods
	}
	fmt.Println("s", s)
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
