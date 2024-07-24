package pubresolvers

import (
	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Opts            *apptypes.APIGatewayOptions
	InternalService apptypes.InternalService
}
