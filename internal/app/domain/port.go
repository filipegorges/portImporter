package domain

import (
	"fmt"
	"strings"
)

type Port struct {
	Name        string    `bson:"name" json:"name"`
	City        string    `bson:"city" json:"city"`
	Country     string    `bson:"country" json:"country"`
	Alias       []string  `bson:"alias" json:"alias"`
	Regions     []string  `bson:"regions" json:"regions"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
	Province    string    `bson:"province" json:"province"`
	Timezone    string    `bson:"timezone" json:"timezone"`
	Unlocs      []string  `bson:"unlocs" json:"unlocs"`
	Code        string    `bson:"code,omitempty" json:"code"`
}

func (p Port) ToString() string {
	return fmt.Sprintf(
		"Port{Name: %q, City: %q, Country: %q, Alias: [%s], Regions: [%s], Coordinates: [%s], Province: %q, Timezone: %q, Unlocs: [%s], Code: %q}",
		p.Name,
		p.City,
		p.Country,
		strings.Join(p.Alias, ", "),
		strings.Join(p.Regions, ", "),
		floatSliceToString(p.Coordinates),
		p.Province,
		p.Timezone,
		strings.Join(p.Unlocs, ", "),
		p.Code,
	)
}

func floatSliceToString(floats []float64) string {
	strs := make([]string, len(floats))
	for i, f := range floats {
		strs[i] = fmt.Sprintf("%.6f", f)
	}
	return strings.Join(strs, ", ")
}
