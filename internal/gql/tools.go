//go:build tools
// +build tools

package main

//go:generate go run ./generate.go

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)
