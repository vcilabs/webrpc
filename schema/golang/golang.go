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

		// TODO: Change this to map[*types.Package]string so we can rename duplicated pkgs?
		resolvedImports: map[string]struct{}{
			// Initial schema file's package name artificially set by golang.org/x/tools/go/packages.
			"command-line-arguments": struct{}{},

			// The following imports are already defined in the Go template.
			"context":       struct{}{},
			"encoding/json": struct{}{},
			"fmt":           struct{}{},
			"io/ioutil":     struct{}{},
			"net/http":      struct{}{},
			"time":          struct{}{},
			"strings":       struct{}{},
			"bytes":         struct{}{},
			"errors":        struct{}{},
			"io":            struct{}{},
			"net/url":       struct{}{},
		},
	}
}

type parser struct {
	schema          *schema.WebRPCSchema
	parsedTypes     map[types.Type]*schema.VarType
	resolvedImports map[string]struct{}
}

// Parse parses a Go source file and return WebRPC schema.
func (p *parser) Parse(path string) (*schema.WebRPCSchema, error) {
	fmt.Println("============== before")
	t := time.Now()
	defer func() {
		fmt.Println("============== after ", time.Since(t))
	}()

	cfg := &packages.Config{
		// TODO: Make the Dir dynamic, parse it from the schema file's path (+current working directory if not absolute).
		Dir:  "/Users/vojtechvitek/go/src/github.com/vcilabs/hubs/contract",
		Mode: packages.NeedName | packages.NeedImports | packages.NeedTypes | packages.NeedFiles | packages.NeedDeps | packages.NeedSyntax,
	}

	schemaPkg, err := packages.Load(cfg, path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load packages")
	}
	if len(schemaPkg) != 1 {
		return nil, errors.Errorf("failed to load initial package (len=%v)", len(schemaPkg))
	}

	err = p.parsePkgInterfaces(schemaPkg[0].Types.Scope())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Go interfaces")
	}

	return p.schema, nil
}

func (p *parser) parsePkgInterfaces(scope *types.Scope) error {
	for _, name := range scope.Names() {
		iface, ok := scope.Lookup(name).Type().Underlying().(*types.Interface)
		if !ok {
			continue
		}

		service := &schema.Service{
			Name:   schema.VarName(name),
			Schema: p.schema, // denormalize/back-reference
		}

		fmt.Printf("interface %v\n", name)

		// TODO: Do we need to loop over embedded interfaces first? Try defining embeeded interface.
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
				return errors.Errorf("interface %v method %v(): failed to get method signature", name, methodName)
			}

			methodParams := methodSignature.Params()
			inputs, err := p.getMethodArguments(methodParams, true)
			if err != nil {
				return errors.Wrapf(err, "interface %v method %v(): failed to get inputs (method arguments)", name, methodName)
			}

			// First method argument must be of type context.Context.
			if methodParams.Len() == 0 {
				return errors.Errorf("interface %v method %v(): first method argument must be context.Context: no arguments defined", name, methodName)
			}
			if err := ensureContextType(methodParams.At(0).Type()); err != nil {
				return errors.Wrapf(err, "interface %v method %v(): first method argument must be context.Context", name, methodName)
			}
			inputs = inputs[1:] // Cut it off. The gen/golang adds context.Context as first method argument automatically.

			methodResults := methodSignature.Results()
			outputs, err := p.getMethodArguments(methodResults, false)
			if err != nil {
				return errors.Wrapf(err, "interface %v method %v(): failed to get outputs (method results)", name, methodName)
			}

			// Last method return value must be of type error.
			if methodResults.Len() == 0 {
				return errors.Errorf("interface %v method %v(): last return value must be context.Context: no return values defined", name, methodName)
			}
			if err := ensureErrorType(methodResults.At(methodResults.Len() - 1).Type()); err != nil {
				return errors.Wrapf(err, "interface %v method %v(): first method argument must be context.Context", name, methodName)
			}
			outputs = outputs[:len(outputs)-1] // Cut it off. The gen/golang adds error as a last return value automatically.

			service.Methods = append(service.Methods, &schema.Method{
				Name:    schema.VarName(methodName),
				Inputs:  inputs,
				Outputs: outputs,
				Service: service, // denormalize/back-reference
			})
		}

		p.schema.Services = append(p.schema.Services, service)
	}

	return nil
}
