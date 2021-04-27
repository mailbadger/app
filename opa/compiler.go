package opa

import (
	_ "embed"

	"github.com/open-policy-agent/opa/ast"
)

//go:embed rbac_authz.rego
var module string

func NewCompiler() (*ast.Compiler, error) {
	// Compile the module. The keys are used as identifiers in error messages.
	compiler, err := ast.CompileModules(map[string]string{
		"rbac_authz.rego": module,
	})

	return compiler, err
}
