package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	_ "github.com/99designs/gqlgen/graphql/introspection"
	"github.com/99designs/gqlgen/plugin/modelgen"
)

var (
	ConfigExitCode   = 2
	GenerateExitCode = 3
)

// Defining mutation function
func constraintFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	if c := fd.Directives.ForName("ref"); c != nil {
		field := c.Arguments.ForName("field")

		f.Tag = fmt.Sprintf("%s bson:%s yaml:%s", f.Tag, field.Value.String(), field.Value.String())
	} else {
		f.Tag += " bson:\"-\""
	}

	if c := fd.Directives.ForName("refRoot"); c != nil {
		idArg := c.Arguments.ForName("id")
		if idArg == nil || idArg.Value == nil {
			panic("refRoot id not defined for " + fd.Name)
		}
	}

	qf := fd.Directives.ForName("queryFields")
	if qf != nil {
		f.Tag += " query:\"fields\""
	}

	qc := fd.Directives.ForName("queryCondition")
	if qc != nil {
		f.Tag += " query:\"condition\""
	}

	nq := fd.Directives.ForName("nestedQueries")
	if nq != nil {
		f.Tag += " query:\"nestedQueries\""
	}

	qo := fd.Directives.ForName("queryOp")
	if qo != nil {
		f.Tag += " query:\"op\""
	}

	qv := fd.Directives.ForName("queryValue")
	if qv != nil {
		f.Tag += " query:\"value\""
	}

	qfd := fd.Directives.ForName("queryField")
	if qfd != nil {
		f.Tag += " query:\"field\""
	}

	qfed := fd.Directives.ForName("queryFieldExists")
	if qfed != nil {
		f.Tag += " query:\"fieldExists\""
	}

	qq := fd.Directives.ForName("query")
	if qq != nil {
		f.Tag += " query:\"query\""
	}

	qrules := fd.Directives.ForName("queryRules")
	if qrules != nil {
		f.Tag += " query:\"queryRules\""
	}

	qd := fd.Directives.ForName("queryType")
	if qd != nil {
		qt := qd.Arguments.ForName("type")
		f.Tag += fmt.Sprintf(" queryType:\"%s\"", qt.Value.String())
	}

	val := fd.Directives.ForName("validation")
	if val != nil {
		vc := val.Arguments.ForName("constraint")
		f.Tag += " validate: " + vc.Value.String()
	}

	return f, nil
}

func main() {
	configFile := flag.String("config", "gqlgen.yml", "--config the configuration file to load")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(ConfigExitCode)
	}

	// Attaching the mutation function onto modelgen plugin
	p := modelgen.Plugin{
		FieldHook: constraintFieldHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(GenerateExitCode)
	}
}
