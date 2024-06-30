package apputil

import (
	"cmp"
	"fmt"
	"github.com/azarc-io/verathread-gateway/internal/gql/graph/model"
	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
	"github.com/azarc-io/verathread-next-common/common/app"
	"github.com/azarc-io/verathread-next-common/util"
	"github.com/rs/zerolog/log"
	"slices"
	"strings"
)

type priorityWrapper struct {
	priority int
	app      *apptypes.App
}

func MapNavigationToNavigationInput(n *apptypes.Navigation, an *app.RegisterAppNavigationInput) {
	n.Title = an.Title
	n.SubTitle = util.Ptr(an.SubTitle)
	n.AuthRequired = util.Ptr(an.AuthRequired)
	n.Hidden = util.Ptr(an.Hidden)
	n.Children = make([]*apptypes.Navigation, 0)
	n.Category = an.Category

	if an.Module != nil {
		n.Module = &apptypes.NavigationModule{
			Path:          an.Module.Path,
			ExposedModule: an.Module.ExposedModule,
			ModuleName:    an.Module.ModuleName,
			Outlet:        an.Module.Outlet,
		}
	}

	for _, child := range an.Children {
		nc := &apptypes.Navigation{}
		MapChildNavigationToNavigationInput(nc, child)
		n.Children = append(n.Children, nc)
	}
}

func MapChildNavigationToNavigationInput(n *apptypes.Navigation, an *app.RegisterChildAppNavigationInput) {
	n.Title = an.Title
	n.SubTitle = util.Ptr(an.SubTitle)
	n.AuthRequired = util.Ptr(an.AuthRequired)
	n.Children = make([]*apptypes.Navigation, 0)

	if an.Module != nil {
		n.Module = &apptypes.NavigationModule{
			Path:          an.Module.Path,
			ExposedModule: an.Module.ExposedModule,
			ModuleName:    an.Module.ModuleName,
			Outlet:        an.Module.Outlet,
		}
	}

	for _, child := range an.Children {
		nc := &apptypes.Navigation{}
		MapChildNavigationToNavigationInput(nc, child)
		n.Children = append(n.Children, nc)
	}
}

func MapAppsToNavigation(data []*apptypes.App) *model.ShellConfiguration {
	result := &model.ShellConfiguration{
		DefaultRoute: util.Ptr(""),
		Categories:   []*model.ShellNavigationCategory{},
	}

	var (
		appCategory = &model.ShellNavigationCategory{
			Title:    "Apps",
			Priority: 0,
			Category: app.CategoryApp,
			Entries:  make([]*model.ShellNavigation, 0),
		}
		settingsCategory = &model.ShellNavigationCategory{
			Title:    "Settings",
			Priority: 1,
			Category: app.CategorySetting,
			Entries:  make([]*model.ShellNavigation, 0),
		}
		dashboardCategory = &model.ShellNavigationCategory{
			Title:    "Dashboards",
			Priority: 2,
			Category: app.CategoryDashboard,
			Entries:  make([]*model.ShellNavigation, 0),
		}

		slotIndex = 0
		slots     []*model.ShellNavigationSlot
		addSlotFn func(slot *apptypes.NavigationSlot, a *apptypes.App)

		prioritizedApps []*priorityWrapper
	)

	addSlotFn = func(slot *apptypes.NavigationSlot, a *apptypes.App) {
		slots = append(slots, &model.ShellNavigationSlot{
			Priority:     util.Ptr(slotIndex),
			Slot:         fmt.Sprintf("slot-%d", slotIndex),
			Description:  slot.Description,
			AuthRequired: util.Ptr(slot.AuthRequired),
			Module: &model.ShellNavigationSlotModule{
				Path:          slot.Module.Path,
				ExposedModule: slot.Module.ExposedModule,
				ModuleName:    slot.Module.ModuleName,
				RemoteEntry:   fmt.Sprintf("%s/%s", a.BaseUrl, a.RemoteEntry),
			},
		})

		slotIndex += 1
	}

	// assign a priority to all native apps, so we sort and make sure that
	// verathread apps always take the first slots
	for _, a := range data {
		p := 100

		if strings.HasPrefix(a.Package, "vth:azarc") {
			p = 0
		}

		prioritizedApps = append(prioritizedApps, &priorityWrapper{
			priority: p,
			app:      a,
		})
	}

	// sort apps by priority
	slices.SortFunc(prioritizedApps, func(a, b *priorityWrapper) int {
		return cmp.Compare(a.priority, b.priority)
	})

	for _, pa := range prioritizedApps {
		a := pa.app
		log.Info().Msgf("processing app %s", a.Package)
		if a.Navigation != nil {
			// sort nav items by title
			slices.SortFunc(a.Navigation, func(a, b *apptypes.Navigation) int {
				return strings.Compare(a.Title, b.Title)
			})

			for _, navigation := range a.Navigation {
				e := &model.ShellNavigation{}
				e.MapFromEntity(navigation, a.Available)

				switch navigation.Category {
				case app.CategoryDashboard:
					dashboardCategory.Entries = append(dashboardCategory.Entries, e)
				case app.CategoryApp:
					appCategory.Entries = append(appCategory.Entries, e)
				case app.CategorySetting:
					settingsCategory.Entries = append(settingsCategory.Entries, e)
				}
			}
		}

		// register azarc slots first
		if a.Slot1 != nil {
			addSlotFn(a.Slot1, a)
		}
		if a.Slot2 != nil {
			addSlotFn(a.Slot2, a)
		}
		if a.Slot3 != nil {
			addSlotFn(a.Slot3, a)
		}
	}

	result.Categories = append(result.Categories, dashboardCategory)
	result.Categories = append(result.Categories, appCategory)
	result.Categories = append(result.Categories, settingsCategory)
	result.Slots = slots

	return result
}

func MapRegisterSlotToEntity(req *app.RegisterAppSlot) *apptypes.NavigationSlot {
	return &apptypes.NavigationSlot{
		Description:  req.Description,
		AuthRequired: req.AuthRequired,
		Module: &apptypes.NavigationSlotModule{
			Path:          req.Module.Path,
			ExposedModule: req.Module.ExposedModule,
			ModuleName:    req.Module.ModuleName,
		},
	}
}
