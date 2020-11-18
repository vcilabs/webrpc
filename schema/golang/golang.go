package golang

import (
	"fmt"
	"go/types"
	"time"

	"github.com/pkg/errors"
	"github.com/webrpc/webrpc/schema"
	"golang.org/x/tools/go/packages"
)

func NewParser(r *schema.Reader) *parser {
	return &parser{
		schema:      &schema.WebRPCSchema{},
		parsedTypes: map[types.Type]*schema.VarType{},
	}
}

type parser struct {
	schema      *schema.WebRPCSchema
	parsedTypes map[types.Type]*schema.VarType
}

// Parse parses a Go source file and return WebRPC schema.
func (p *parser) Parse(path string) (*schema.WebRPCSchema, error) {
	fmt.Println("============== before")
	t := time.Now()
	defer func() {
		fmt.Println("============== after ", time.Since(t))
	}()

	cfg := &packages.Config{
		// TODO: Make the Dir dynamic, parse it from the CWD + path.
		Dir:  "/Users/vojtechvitek/go/src/github.com/vcilabs/hubs/contract",
		Mode: packages.NeedName | packages.NeedImports | packages.NeedTypes | packages.NeedFiles | packages.NeedDeps | packages.NeedSyntax,
	}

	initialPkg, err := packages.Load(cfg, path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load packages")
	}

	if len(initialPkg) != 1 {
		return nil, errors.Errorf("failed to load initial package (len=%v)", len(initialPkg))
	}

	err = p.parsePkgInterfaces(initialPkg[0].Types.Scope())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Go interfaces")
	}

	// TODO: Golang "import mode", where we don't generate messages for all imported structs,
	//       but instead we use imports statements from the original package.
	//
	// // Append imported packages.
	// allPkgs := initialPkg
	// for _, pkg := range initialPkg {
	// 	for _, err := range pkg.Errors {
	// 		// TODO: Is this a syntax error? Should we return?
	// 		fmt.Printf("error: %v\n", err)
	// 	}

	// 	for _, importedPkg := range pkg.Imports {
	// 		// NOTE: These package might have additional imports. ie. importedPkg.Imports.
	// 		s.Imports = append(s.Imports, &schema.Import{
	// 			Name: importedPkg.ID,
	// 			Path: importedPkg.ID,
	// 		})
	// 		allPkgs = append(allPkgs, importedPkg)
	// 	}
	// }

	return p.schema, nil
}

func (p *parser) parsePkgInterfaces(scope *types.Scope) error {
	for _, name := range scope.Names() {
		iface, ok := scope.Lookup(name).Type().Underlying().(*types.Interface)
		if !ok {
			continue
		}

		service := &schema.Service{
			Name: schema.VarName(name),
		}

		fmt.Printf("interface %v\n", name)

		// TODO: Loop over embedded interfaces first?
		// for i := 0; i < iface.NumEmbeddeds(); i++ {
		// }

		// Loop over interface's methods.
		for i := 0; i < iface.NumMethods(); i++ {
			method := iface.Method(i)
			if !method.Exported() {
				continue
			}

			methodName := method.Id()
			fmt.Printf("- %v\n", methodName)

			methodSignature, ok := method.Type().(*types.Signature)
			if !ok {
				return errors.Errorf("failed to get signature of %v interface's method %v()", name, methodName)
			}

			methodParams := methodSignature.Params()
			if methodParams.Len() == 0 {
				return errors.Errorf("interface %v method %v(): first method argument must be context.Context: no arguments defined", name, methodName)
			}

			// First method parameter must be of type context.Context.
			if err := ensureContextType(methodParams.At(0).Type()); err != nil {
				return errors.Wrapf(err, "interface %v method %v(): first method argument must be context.Context", name, methodName)
			}

			results := methodSignature.Results()

			// TODO: Ensure the last result item is of type error.

			inputs, err := p.getMethodArguments(methodParams)
			if err != nil {
				return errors.Wrapf(err, "failed to get inputs (method arguments) of %v interface's method %v()", name, methodName)
			}
			outputs, err := p.getMethodArguments(results)
			if err != nil {
				return errors.Wrapf(err, "failed to get outputs (method results) of %v interface's method %v()", name, methodName)
			}

			service.Methods = append(service.Methods, &schema.Method{
				Name:    schema.VarName(methodName),
				Inputs:  inputs,
				Outputs: outputs,
			})
		}

		p.schema.Services = append(p.schema.Services, service)
	}

	return nil
}
