package golang

import (
	"fmt"
	"go/types"

	"github.com/pkg/errors"
	"github.com/webrpc/webrpc/schema"
)

func (p *parser) parseType(typeName string, typ types.Type) (varType *schema.VarType, err error) {
	if cached, ok := p.parsedTypes[typ]; ok {
		return cached, nil
	}
	defer func() {
		if err == nil {
			// Cache the return value to avoid parsing the same type multiple times.
			// No need to lock the map, as we're parsing sequentially.
			p.parsedTypes[typ] = varType
		}
	}()

	switch v := typ.(type) {
	case *types.Named:
		pkg := v.Obj().Pkg()
		if pkg != nil {
			// If the type belongs to a specific package, save the pkg reference to schema.Imports.
			pkgPath := pkg.Path()
			if _, ok := p.resolvedImports[pkgPath]; !ok {
				p.resolvedImports[pkgPath] = struct{}{}
				p.schema.Imports = append(p.schema.Imports, &schema.Import{
					Path: pkgPath,
				})
			}
		}
		return p.parseType(v.Obj().Name(), v.Underlying())

	case *types.Basic:
		return p.parseBasic(v)

	case *types.Struct:
		return p.parseStruct(typeName, v)

	case *types.Slice:
		return p.parseSlice(typeName, v)

	case *types.Interface:
		return p.parseInterface(typeName, v)

	case *types.Map:
		return p.parseMap(typeName, v)

	case *types.Pointer:
		// TODO: Consider adding schema.T_Pointer, or add metadata to Golang
		// type to distinguish between "pointer to struct" vs. "plain struct".
		varType, err = p.parseType(typeName, v.Elem())
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

func (p *parser) parseStruct(typeName string, structTyp *types.Struct) (*schema.VarType, error) {
	// TODO: Handle a special case for time.Time => schema.T_Timestamp.

	msg := &schema.Message{
		Name: schema.VarName(typeName),
		Type: schema.MessageType("struct"),
	}

	for i := 0; i < structTyp.NumFields(); i++ {
		field := structTyp.Field(i)
		if !field.Exported() {
			continue
		}

		varType, err := p.parseType(field.Name(), field.Type())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse var %v", field.Name())
		}

		msg.Fields = append(msg.Fields, &schema.MessageField{
			Name: schema.VarName(field.Name()),
			Type: varType,
		})
	}

	p.schema.Messages = append(p.schema.Messages, msg)

	varType := &schema.VarType{
		Type: schema.T_Struct,
		Struct: &schema.VarStructType{
			Name:    typeName,
			Message: msg,
		},
	}

	return varType, nil
}

func (p *parser) parseSlice(typeName string, sliceTyp *types.Slice) (*schema.VarType, error) {
	elem, err := p.parseType(typeName, sliceTyp.Elem())
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
func (p *parser) parseInterface(typeName string, iface *types.Interface) (*schema.VarType, error) {
	// TODO: A special case for error and context.Context.

	varType := &schema.VarType{
		Type: schema.T_Any,
	}

	return varType, nil
}

// Parse argument of type interface. We only allow context.Context and error.
func (p *parser) parseMap(typeName string, m *types.Map) (*schema.VarType, error) {
	key, err := p.parseType(typeName, m.Key())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse map key type")
	}

	value, err := p.parseType(typeName, m.Elem())
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
