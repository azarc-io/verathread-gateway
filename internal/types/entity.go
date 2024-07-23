package apptypes

import (
	"encoding/json"
	"github.com/azarc-io/verathread-gateway/internal/gql/graph/common/model"
	"time"
)

type (
	App struct {
		ID                      string            `json:"id" bson:"_id,omitempty"`
		Name                    string            `json:"name" bson:"name,omitempty"`
		Package                 string            `json:"package" bson:"package,omitempty"`
		Version                 string            `json:"version" bson:"version,omitempty"`
		APIURL                  string            `json:"apiURL" bson:"apiURL,omitempty"`
		WebURL                  string            `json:"webURL" bson:"webURL,omitempty" yaml:"WebURL"`
		RemoteEntry             string            `json:"remoteEntry,omitempty" bson:"remoteEntry"`
		Proxy                   bool              `json:"proxy" bson:"proxy,omitempty" yaml:"proxy"`
		Navigation              []*Navigation     `json:"navigation" bson:"navigation,omitempty"`
		Slot1                   *NavigationSlot   `json:"slot1,omitempty" bson:"slot1,omitempty" yaml:"slot1"`
		Slot2                   *NavigationSlot   `json:"slot2,omitempty" bson:"slot2,omitempty" yaml:"slot2"`
		Slot3                   *NavigationSlot   `json:"slot3,omitempty" bson:"slot3,omitempty" yaml:"slot3"`
		CreatedAt               time.Time         `json:"createdAt" bson:"createdAt,omitempty"`
		UpdatedAt               time.Time         `json:"updatedAt" bson:"updatedAt,omitempty"`
		Adopted                 bool              `json:"adopted" bson:"adopted,omitempty"`
		Available               bool              `json:"available" bson:"available,omitempty"`
		RemoteEntryRewriteRegEx map[string]string `json:"remoteEntryRewriteRegEx,omitempty" bson:"remoteEntryRewriteRegEx,omitempty"`
	}

	Navigation struct {
		ID           string                    `json:"id" bson:"id,omitempty" yaml:"id"`
		Title        string                    `json:"title" bson:"title,omitempty" yaml:"title"`
		SubTitle     string                    `json:"subTitle,omitempty" bson:"subTitle,omitempty" yaml:"subTitle"`
		AuthRequired bool                      `json:"authRequired,omitempty" bson:"authRequired,omitempty" yaml:"authRequired"`
		Hidden       bool                      `json:"hidden,omitempty" bson:"hidden,omitempty" yaml:"hidden"`
		Category     model.RegisterAppCategory `json:"category" bson:"category,omitempty" yaml:"category"`
		Children     []*Navigation             `json:"children,omitempty" bson:"children,omitempty" yaml:"children"`
		RemoteEntry  string                    `json:"remoteEntry" bson:"remoteEntry,omitempty" yaml:"remoteEntry"`
		Module       *NavigationModule         `json:"module,omitempty" bson:"module,omitempty" yaml:"module"`
		Icon         string                    `json:"icon,omitempty" bson:"icon" yaml:"icon"`
	}

	NavigationChild struct {
		Title        string             `json:"title,omitempty" bson:"title"`
		SubTitle     string             `json:"subTitle,omitempty" bson:"subTitle"`
		AuthRequired bool               `json:"authRequired,omitempty" bson:"authRequired"`
		Path         string             `json:"path,omitempty" bson:"path"`
		Children     []*NavigationChild `json:"children,omitempty" bson:"children"`
		Icon         string             `json:"icon,omitempty" bson:"icon"`
		Module       *NavigationModule  `json:"module,omitempty" bson:"module"`
	}

	NavigationSlot struct {
		Description  string                `json:"description,omitempty"`
		AuthRequired bool                  `json:"authRequired,omitempty"`
		Module       *NavigationSlotModule `json:"module,omitempty"`
	}

	NavigationModule struct {
		Path          string `json:"path,omitempty" bson:"path"`
		ExposedModule string `json:"exposedModule,omitempty" bson:"exposedModule"`
		ModuleName    string `json:"moduleName,omitempty" bson:"moduleName"`
		Outlet        string `json:"outlet,omitempty" bson:"outlet"`
	}

	NavigationSlotModule struct {
		Path          string `json:"path,omitempty" bson:"path"`
		ExposedModule string `json:"exposedModule,omitempty" bson:"exposedModule"`
		ModuleName    string `json:"moduleName,omitempty" bson:"moduleName"`
	}
)

func (a App) MarshalBinary() (data []byte, err error) {
	return json.Marshal(a)
}

func (a *App) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}
