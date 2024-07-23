package apputil

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/azarc-io/verathread-gateway/internal/gql/graph/common/model"
	util2 "github.com/azarc-io/verathread-gateway/internal/gql/graph/util"
	"github.com/rs/zerolog/log"

	apptypes "github.com/azarc-io/verathread-gateway/internal/types"
	"github.com/azarc-io/verathread-next-common/util"
)

type priorityWrapper struct {
	priority int
	app      *apptypes.App
}

const (
	AppPriority       = 0
	SettingsPriority  = 1
	DashboardCategory = 2
)

// MapNavInputToNavEntity maps gql navigation data to entity data
func MapNavInputToNavEntity(an *model.RegisterAppNavigationInput, n *apptypes.Navigation) {
	n.Title = an.Title
	n.SubTitle = an.SubTitle
	n.AuthRequired = an.AuthRequired
	n.Hidden = an.Hidden
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
		MapChildNavInputToNavEntity(nc, child)
		n.Children = append(n.Children, nc)
	}
}

// MapChildNavInputToNavEntity maps child navigation to entity
func MapChildNavInputToNavEntity(n *apptypes.Navigation, an *model.RegisterChildAppNavigationInput) {
	n.Title = an.Title
	n.SubTitle = an.SubTitle
	n.AuthRequired = an.AuthRequired
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
		MapChildNavInputToNavEntity(nc, child)
		n.Children = append(n.Children, nc)
	}
}

// MapAppsToNavigation maps apps to shell configuration data for the gql api
func MapAppsToNavigation(data []*apptypes.App) *model.ShellConfiguration {
	result := &model.ShellConfiguration{
		DefaultRoute: util.Ptr(""),
		Categories:   []*model.ShellNavigationCategory{},
	}

	var (
		appCategory = &model.ShellNavigationCategory{
			Title:    "Apps",
			Priority: AppPriority,
			Category: model.RegisterAppCategoryApp,
			Entries:  make([]*model.ShellNavigation, 0),
		}
		settingsCategory = &model.ShellNavigationCategory{
			Title:    "Settings",
			Priority: SettingsPriority,
			Category: model.RegisterAppCategorySetting,
			Entries:  make([]*model.ShellNavigation, 0),
		}
		dashboardCategory = &model.ShellNavigationCategory{
			Title:    "Dashboards",
			Priority: DashboardCategory,
			Category: model.RegisterAppCategoryDashboard,
			Entries:  make([]*model.ShellNavigation, 0),
		}

		slotIndex = 0
		slots     []*model.ShellNavigationSlot
		addSlotFn func(slot *apptypes.NavigationSlot, a *apptypes.App)

		prioritizedApps = make([]*priorityWrapper, len(data))
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
				RemoteEntry:   fmt.Sprintf("%s/%s", a.WebURL, a.RemoteEntry),
			},
		})

		slotIndex += 1
	}

	// assign a priority to all native apps, so we sort and make sure that
	// verathread apps always take the first slots
	for i, a := range data {
		p := 100

		if strings.HasPrefix(a.Package, "vth:azarc") {
			p = 0
		}

		prioritizedApps[i] = &priorityWrapper{
			priority: p,
			app:      a,
		}
	}

	// sort apps by priority
	slices.SortFunc(prioritizedApps, func(a, b *priorityWrapper) int {
		return cmp.Compare(a.priority, b.priority)
	})

	for _, pa := range prioritizedApps {
		a := pa.app
		if a.Navigation != nil {
			// sort nav items by title
			slices.SortFunc(a.Navigation, func(a, b *apptypes.Navigation) int {
				return strings.Compare(a.Title, b.Title)
			})

			for _, navigation := range a.Navigation {
				e := &model.ShellNavigation{}
				util2.MapFromEntity(e, navigation, a.Available)

				switch navigation.Category {
				case model.RegisterAppCategoryDashboard:
					dashboardCategory.Entries = append(dashboardCategory.Entries, e)
				case model.RegisterAppCategoryApp:
					appCategory.Entries = append(appCategory.Entries, e)
				case model.RegisterAppCategorySetting:
					settingsCategory.Entries = append(settingsCategory.Entries, e)
				}
			}
		} else {
			log.Warn().Msgf("navigation data is empty")
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

	result.Categories = append(result.Categories, dashboardCategory, appCategory, settingsCategory)
	result.Slots = slots

	return result
}

// MapRegisterSlotToEntity maps slot model to slot entity
func MapRegisterSlotToEntity(req *model.RegisterAppSlot) *apptypes.NavigationSlot {
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
