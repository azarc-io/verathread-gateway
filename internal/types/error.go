package types

import "errors"

var (
	ErrRebuildNavigationFailed = errors.New("failed to rebuild navigation")
	ErrGatewayNotReady         = errors.New("gateway is not ready")
)
