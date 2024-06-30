package apptypes

import (
	"github.com/asynkron/protoactor-go/actor"
)

type (
	GetAppsConfiguration struct {
		Tenant string
	}

	AppAdded struct {
		App *App
	}

	AppUpdated struct {
		App *App
	}

	AppRemoved struct {
		actor.AutoRespond
		App *App
	}

	AppUnavailable struct {
		App *App
	}
)

func (*AppAdded) GetAutoResponse(_ actor.Context) interface{} {
	return "app added"
}

func (*AppRemoved) GetAutoResponse(_ actor.Context) interface{} {
	return "app removed from cache"
}
