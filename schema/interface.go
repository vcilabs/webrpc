package schema

import (
	"strings"

	"github.com/pkg/errors"
)

type GoInterface struct {
	Name    VarName
	Methods []*Method

	Schema *WebRPCSchema
}

type Method struct {
	Name VarName `json:"name"`

	StreamInput  bool `json:"streamInput,omitempty"`
	StreamOutput bool `json:"streamOutput,omitempty"`
	Proxy        bool `json:"-"` // TODO: actual implementation

	Inputs  []*MethodArgument `json:"inputs"`
	Outputs []*MethodArgument `json:"outputs"`

	Service     *Service     `json:"-"`
	GoInterface *GoInterface `json:"-"` // denormalize/back-reference
}

type MethodArgument struct {
	Name     VarName  `json:"name"`
	Type     *VarType `json:"type"`
	Optional bool     `json:"optional"`

	InputArg  bool `json:"-"` // denormalize/back-reference
	OutputArg bool `json:"-"` // denormalize/back-reference
}

func (i *GoInterface) Parse(schema *WebRPCSchema) error {
	i.Schema = schema // back-ref

	// interface name
	interfaceName := string(i.Name)
	if string(i.Name) == "" {
		return errors.Errorf("schema error: Interface name cannot be empty")
	}

	// Ensure we don't have dupe interface names (w/ normalization)
	name := strings.ToLower(string(i.Name))
	for _, goInterface := range schema.GoInterface {
		if goInterface != i && name == strings.ToLower(string(goInterface.Name)) {
			return errors.Errorf("schema error: duplicate interface name detected in interface '%s'", interfaceName)
		}
	}

	// Ensure methods is defined
	if len(i.Methods) == 0 {
		return errors.Errorf("schema error: methods cannot be empty for interface '%s'", interfaceName)
	}

	// Verify method names and ensure we don't have any duplicate method names
	methodList := map[string]string{}
	for _, method := range i.Methods {
		if string(method.Name) == "" {
			return errors.Errorf("schema error: detected empty method name in interface '%s", interfaceName)
		}

		methodName := string(method.Name)
		nMethodName := strings.ToLower(methodName)

		if _, ok := methodList[nMethodName]; ok {
			return errors.Errorf("schema error: detected duplicate method name of '%s' in interface '%s'", methodName, interfaceName)
		}
		methodList[nMethodName] = methodName
	}

	// Parse+validate methods
	for _, method := range i.Methods {
		err := method.IntParse(schema, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Method) IntParse(schema *WebRPCSchema, goInterface *GoInterface) error {
	if goInterface == nil {
		return errors.Errorf("parse error, interface arg cannot be nil")
	}
	m.GoInterface = goInterface // back-ref
	interfaceName := string(goInterface.Name)

	// Parse+validate inputs
	for _, input := range m.Inputs {
		input.InputArg = true // back-ref
		if input.Name == "" {
			return errors.Errorf("schema error: detected empty input argument name for method '%s' in interface '%s'", m.Name, interfaceName)
		}
		err := input.Type.Parse(schema)
		if err != nil {
			return err
		}
	}

	// Parse+validate outputs
	for _, output := range m.Outputs {
		output.OutputArg = true // back-ref
		if output.Name == "" {
			return errors.Errorf("schema error: detected empty output name for method '%s' in interface '%s'", m.Name, interfaceName)
		}
		err := output.Type.Parse(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
