package federation_test

import (
	"testing"

	"github.com/azarc-io/verathread-gateway/internal/gateway/federation"
	"github.com/stretchr/testify/assert"
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
