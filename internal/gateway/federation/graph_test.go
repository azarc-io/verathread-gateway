package federation_test

import (
	"github.com/azarc-io/verathread-gateway/internal/gateway/federation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsesMultipleErrors(t *testing.T) {
	errs := federation.GQLErrors{
		struct {
			Message string `json:"message"`
		}{Message: "a"},
		struct {
			Message string `json:"message"`
		}{Message: "b"},
		struct {
			Message string `json:"message"`
		}{Message: "c"},
	}
	assert.Equal(t, "\ta\tb\tc", errs.Error())
}
