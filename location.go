package location

import (
	"fmt"
	"strings"

	"github.com/ecletus/admin"
	"github.com/ecletus/core"
	"github.com/ecletus/core/resource"
	"github.com/ecletus/core/utils"
)

type Plugin struct {
}

var GoogleAPIKey string

type LocationConfig struct {
	GoogleAPIKey string
}

func (LocationConfig) ConfigureQorMeta(meta resource.Metaor) {
}

// Location is a struct, you could embedded it into your model to get the Location feature for your model
type Location struct {
	Address   string
	City      string
	Region    string
	Country   string
	Zip       string
	Latitude  float64
	Longitude float64
}

type locationInterface interface {
	GetLocation() *Location
}

// GetLocation get location from your model
func (location *Location) GetLocation() *Location {
	return location
}

// ConfigureResource configure qor locale for Admin
func (*Location) ConfigureResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		Admin := res.GetAdmin()
		res.UseTheme("location")

		Admin.RegisterFuncMap("replace_suffix", func(str, suffix, newSuffix string) string {
			return fmt.Sprint(strings.TrimSuffix(str, suffix), newSuffix)
		})

		if field, ok := res.ModelStruct.FieldsByName["Location"]; ok {
			labelName := field.Name
			if customName, ok := utils.ParseTagOption(field.Tag.Get("location"))["NAME"]; ok {
				labelName = customName
			}

			if res.GetMeta(field.Name) == nil {
				res.Meta(&admin.Meta{Name: field.Name, Label: labelName, Type: "location", Config: &LocationConfig{GoogleAPIKey: GoogleAPIKey}, Valuer: func(resource interface{}, ctx *core.Context) interface{} {
					return resource.(locationInterface).GetLocation()
				}})
				res.IndexAttrs(res.IndexAttrs(), "-"+field.Name, "-Latitude", "-Longitude")
				res.NewAttrs(res.NewAttrs(), "Location", "-Address", "-City", "-Region", "-Country", "-Zip", "-Latitude", "-Longitude")
				res.EditAttrs(res.EditAttrs(), "Location", "-Address", "-City", "-Region", "-Country", "-Zip", "-Latitude", "-Longitude")
				res.ShowAttrs(res.ShowAttrs(), "Location", "-"+field.Name, false)
			}
		}
	}
}
