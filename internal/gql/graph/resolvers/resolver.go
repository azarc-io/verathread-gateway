package resolvers

import (
	apptypes "github.com/azarc-io/verathread-gateway/internal/api"
	"github.com/azarc-io/verathread-gateway/internal/config"

	"github.com/dapr/go-sdk/service/common"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Opts            *config.APIGatewayOptions
	InternalService apptypes.InternalService
	Service         common.Service
}
