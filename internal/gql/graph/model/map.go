package model

import (
	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
)

func (n *ShellNavigation) MapFromEntity(navigation *apptypes.Navigation, available *bool) {
	n.Title = navigation.Title
	n.SubTitle = navigation.SubTitle
	n.AuthRequired = navigation.AuthRequired
	n.Healthy = *available
	n.ID = navigation.Id
	n.Module = &ShellNavigationModule{
		Path:          navigation.Module.Path,
		ExposedModule: navigation.Module.ExposedModule,
		ModuleName:    navigation.Module.ModuleName,
		Outlet:        navigation.Module.Outlet,
	}

	for _, child := range navigation.Children {
		nc := &ShellNavigationChild{}
		nc.MapChildFromEntity(child, available)
		n.Children = append(n.Children, nc)
	}
}

func (c *ShellNavigationChild) MapChildFromEntity(navigation *apptypes.Navigation, available *bool) {
	c.Title = navigation.Title
	c.SubTitle = navigation.SubTitle
	c.AuthRequired = navigation.AuthRequired
	c.Healthy = *available

	if navigation.Module != nil {
		c.Module = &ShellNavigationModule{
			Path:          navigation.Module.Path,
			ExposedModule: navigation.Module.ExposedModule,
			ModuleName:    navigation.Module.ModuleName,
			Outlet:        navigation.Module.Outlet,
		}
	}
}

func (c *ShellNavigationChild) MapFromNavigationEntity(navigation *apptypes.Navigation, e *ShellNavigationChild) {

}
