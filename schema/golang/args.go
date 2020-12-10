package golang

import (
	"fmt"
	"go/types"

	"github.com/pkg/errors"
	"github.com/webrpc/webrpc/schema"
)

func (p *parser) getMethodArguments(params *types.Tuple, isInput bool) ([]*schema.MethodArgument, error) {
	var args []*schema.MethodArgument

	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		typ := param.Type()

		name := param.Name()
		if name == "" {
			if isInput {
				name = fmt.Sprintf("arg%v", i) // 0 is `ctx context.Context`
			} else {
				name = fmt.Sprintf("ret%v", i+1)
			}
		}

		varType, err := p.parseType("", typ) // Type name will be resolved deeper down the stack.
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse argument %v", name)
		}

		arg := &schema.MethodArgument{
			Name:      schema.VarName(name),
			Type:      varType,
			InputArg:  isInput,  // denormalize/back-reference
			OutputArg: !isInput, // denormalize/back-reference
		}

		args = append(args, arg)
	}

	return args, nil
}

func ensureContextType(typ types.Type) (err error) {
	namedType, ok := typ.(*types.Named)
	if !ok {
		return errors.Errorf("expected named type: found type %T (%+v)", typ, typ)
	}

	if _, ok := namedType.Underlying().(*types.Interface); !ok {
		return errors.Errorf("expected underlying interface: found type %T (%+v)", typ, typ)
	}

	pkgName := namedType.Obj().Pkg().Name()
	typeName := namedType.Obj().Name()

	if pkgName != "context" && typeName != "Context" {
		return errors.Errorf("expected context.Context: found %v.%v", pkgName, typeName)
	}

	return nil
}

func ensureErrorType(typ types.Type) (err error) {
	namedType, ok := typ.(*types.Named)
	if !ok {
		return errors.Errorf("expected named type: found type %T (%+v)", typ, typ)
	}

	if _, ok := namedType.Underlying().(*types.Interface); !ok {
		return errors.Errorf("expected underlying interface: found type %T (%+v)", typ, typ)
	}

	pkgName := namedType.Obj().Pkg()
	typeName := namedType.Obj().Name()

	if pkgName != nil && typeName != "error" {
		return errors.Errorf("expected error: found %v.%v", pkgName, typeName)
	}

	return nil
}
