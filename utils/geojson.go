package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// ----------------------------------------------
// GeoJSON types
// ----------------------------------------------

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

// ----------------------------------------------
// GeoJSON creation functions
// ----------------------------------------------

func NewFeatureCollection(features []Feature) FeatureCollection {
	return FeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}
}

func NewPointFeature(lat, lon float64, properties map[string]interface{}) Feature {
	return Feature{
		Type: "Feature",
		Geometry: Geometry{
			Type:        "Point",
			Coordinates: []float64{lon, lat},
		},
		Properties: properties,
	}
}

func NewLineStringFeature(coordinates [][]float64, properties map[string]interface{}) Feature {
	return Feature{
		Type: "Feature",
		Geometry: Geometry{
			Type:        "LineString",
			Coordinates: coordinates,
		},
		Properties: properties,
	}
}

// ----------------------------------------------
// Export functions
// ----------------------------------------------

func ExportToFile(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("error formatting JSON: %w", err)
	}

	if filename != "" {
		if err := os.WriteFile(filename, jsonData, 0644); err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}
		fmt.Printf("Successfully exported to %s\n", filename)
	} else {
		fmt.Println(string(jsonData))
	}

	return nil
}
