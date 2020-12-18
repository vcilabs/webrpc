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
		schema:          &schema.WebRPCSchema{},
		parsedTypes:     map[types.Type]*schema.VarType{},
		parsedTypeNames: map[string]struct{}{},

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
	schema *schema.WebRPCSchema

	// Cache for already parsed types, to improve performance & so we can traverse circular dependencies.
	parsedTypes     map[types.Type]*schema.VarType
	parsedTypeNames map[string]struct{}

	inlineMode      bool // When traversing `json:",inline"`, we don't want to store the struct type as WebRPC message.
	resolvedImports map[string]struct{}

	schemaPkgName string // Shema file's package name.
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

	p.schemaPkgName = schemaPkg[0].Name

	err = p.parsePkgInterfaces(schemaPkg[0].Types.Scope())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Go interfaces")
	}

	return p.schema, nil
}
