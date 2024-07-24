package util

import (
	"github.com/azarc-io/verathread-gateway/internal/gql/graph/common/model"
	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
)

func MapFromEntity(n *model.ShellNavigation, navigation *apptypes.Navigation, available bool) {
	n.Title = navigation.Title
	n.SubTitle = navigation.SubTitle
	n.AuthRequired = navigation.AuthRequired
	n.Healthy = available
	n.ID = navigation.ID
	n.Module = &model.ShellNavigationModule{
		Path:          navigation.Module.Path,
		ExposedModule: navigation.Module.ExposedModule,
		ModuleName:    navigation.Module.ModuleName,
		Outlet:        navigation.Module.Outlet,
		RemoteEntry:   navigation.RemoteEntry,
	}

	for _, child := range navigation.Children {
		nc := &model.ShellNavigationChild{}
		MapChildFromEntity(nc, child, available)
		n.Children = append(n.Children, nc)
	}
}

func MapChildFromEntity(c *model.ShellNavigationChild, navigation *apptypes.Navigation, available bool) {
	c.Title = navigation.Title
	c.SubTitle = navigation.SubTitle
	c.AuthRequired = navigation.AuthRequired
	c.Healthy = available

	if navigation.Module != nil {
		c.Module = &model.ShellNavigationModule{
			Path:          navigation.Module.Path,
			ExposedModule: navigation.Module.ExposedModule,
			ModuleName:    navigation.Module.ModuleName,
			Outlet:        navigation.Module.Outlet,
		}
	}
}
