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
		schema: &schema.WebRPCSchema{
			SchemaType: "go",
		},
	}
}

type parser struct {
	schema *schema.WebRPCSchema
}

// Parse parses a Go source file and return WebRPC schema.
func (p *parser) Parse(path string) (*schema.WebRPCSchema, error) {
	fmt.Println("============== before")
	t := time.Now()
	defer func() {
		fmt.Println("============== after ", time.Since(t))
	}()

	cfg := &packages.Config{
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

	err = p.parseInterfaces(initialPkg[0].Types.Scope())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Go interfaces")
	}

	return p.schema, nil

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
}

func (p *parser) parseInterfaces(scope *types.Scope) error {
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
				return errors.Errorf("first input argument of each interface method must be context.Context: no arguments")
			}

			// TODO: Ensure the methodParams.At(0) is indeed of type context.Context()

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

func (p *parser) parseType(name string, typ types.Type) (*schema.VarType, error) {
	switch v := typ.(type) {
	case *types.Named:
		return p.parseType(name, v.Underlying())
	case *types.Basic:
		return p.parseBasic(v)
	case *types.Struct:
		return p.parseStruct(name, v)
	case *types.Slice:
		return p.parseSlice(v)
	case *types.Interface:
		return p.parseInterface(v)
	case *types.Map:
		return p.parseMap(v)
	case *types.Pointer:
		// TODO: Consider adding schema.T_Pointer, or add metadata to Golang
		// type to distinguish between "pointer to struct" vs. "plain struct".
		varType, err := p.parseType(name, v.Elem())
		if err != nil {
			return nil, errors.Wrap(err, "failed to dereference pointer")
		}
		return varType, nil
	default:
		return nil, errors.Errorf("unknown argument type %T", typ)
	}
}

func (p *parser) parseBasic(typ *types.Basic) (*schema.VarType, error) {
	var varType schema.VarType
	err := schema.ParseVarTypeExpr(p.schema, typ.Name(), &varType)
	if err != nil {
		return nil, fmt.Errorf("unknown data type: %v", typ.Name())
	}

	return &varType, nil
}

func (p *parser) parseStruct(name string, structTyp *types.Struct) (*schema.VarType, error) {
	msg := &schema.Message{
		Name: schema.VarName(name),
		Type: schema.MessageType("struct"),
	}

	for i := 0; i < structTyp.NumFields(); i++ {
		field := structTyp.Field(i)

		varType, err := p.parseType(field.Name(), field.Type())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse var %v", field.Name())
		}

		msg.Fields = append(msg.Fields, &schema.MessageField{
			Name: schema.VarName(field.Name()),
			Type: varType,
		})
		fmt.Printf("struct field: %+v", field)
	}

	p.schema.Messages = append(p.schema.Messages, msg)

	fmt.Printf("struct: %+v", structTyp)

	varType := &schema.VarType{
		Type: schema.T_Struct,
		Struct: &schema.VarStructType{
			Name:    name,
			Message: msg,
		},
	}

	return varType, nil
}

func (p *parser) parseSlice(sliceTyp *types.Slice) (*schema.VarType, error) {
	elem, err := p.parseType("", sliceTyp.Elem())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse slice type")
	}

	varType := &schema.VarType{
		Type: schema.T_List,
		List: &schema.VarListType{
			Elem: elem,
		},
	}

	return varType, nil
}

// Parse argument of type interface. We only allow context.Context and error.
func (p *parser) parseInterface(iface *types.Interface) (*schema.VarType, error) {
	varType := &schema.VarType{
		Type: schema.T_Any,
	}

	return varType, nil
}

// Parse argument of type interface. We only allow context.Context and error.
func (p *parser) parseMap(m *types.Map) (*schema.VarType, error) {
	key, err := p.parseType("", m.Key())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse map key type")
	}

	value, err := p.parseType("", m.Elem())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse map value type")
	}

	varType := &schema.VarType{
		Type: schema.T_Map,
		Map: &schema.VarMapType{
			Key:   key.Type,
			Value: value,
		},
	}

	return varType, nil
}

func (p *parser) getMethodArguments(params *types.Tuple) ([]*schema.MethodArgument, error) {
	var args []*schema.MethodArgument

	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		typ := param.Type()

		name := param.Name()
		if name == "" {
			name = fmt.Sprintf("ret%v", i)
		}

		varType, err := p.parseType(name, typ)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse argument %v", name)
		}

		arg := &schema.MethodArgument{
			Name: schema.VarName(name),
			Type: varType,
		}

		args = append(args, arg)
	}

	return args, nil
}
