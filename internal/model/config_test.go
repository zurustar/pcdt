// Package model contains tests for the config model.
package model

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Feature: countdown-timer, Property 10: 設定の保存・読み込みラウンドトリップ
// **Validates: Requirements 6.4**
//
// *任意の*有効な設定値に対して、保存してから読み込むと、元の設定値と等しい値が復元される。
func TestAppConfig_Property10_JSONRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Any valid config values, when serialized to JSON and deserialized back, produce equal values
	properties.Property("config JSON round trip preserves all values", prop.ForAll(
		func(alwaysOnTop bool, windowWidth float32, windowHeight float32, lastMinutes int, lastSeconds int) bool {
			original := &AppConfig{
				AlwaysOnTop:  alwaysOnTop,
				WindowWidth:  windowWidth,
				WindowHeight: windowHeight,
				LastMinutes:  lastMinutes,
				LastSeconds:  lastSeconds,
			}

			// Serialize to JSON
			jsonData, err := original.ToJSON()
			if err != nil {
				return false
			}

			// Deserialize from JSON
			loaded, err := FromJSON(jsonData)
			if err != nil {
				return false
			}

			// Verify all fields are equal
			return original.AlwaysOnTop == loaded.AlwaysOnTop &&
				original.WindowWidth == loaded.WindowWidth &&
				original.WindowHeight == loaded.WindowHeight &&
				original.LastMinutes == loaded.LastMinutes &&
				original.LastSeconds == loaded.LastSeconds
		},
		gen.Bool(),                  // alwaysOnTop
		gen.Float32Range(200, 1000), // windowWidth: reasonable window sizes
		gen.Float32Range(150, 800),  // windowHeight: reasonable window sizes
		gen.IntRange(0, 99),         // lastMinutes: 0-99 as per input constraints
		gen.IntRange(0, 59),         // lastSeconds: 0-59 as per input constraints
	))

	properties.TestingRun(t)
}
