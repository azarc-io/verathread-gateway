package federation

import (
	"strings"
)

const ServiceDefinitionQuery = `
	{
		"query": "query __ApolloGetServiceDefinition__ { _service { sdl } }",
		"operationName": "__ApolloGetServiceDefinition__",
		"variables": {}
	}`

type GQLErrors []struct {
	Message string `json:"message"`
}

func (g GQLErrors) Error() string {
	var builder strings.Builder
	for _, m := range g {
		_ = builder.WriteByte('\t')
		_, _ = builder.WriteString(m.Message)
	}

	return builder.String()
}
