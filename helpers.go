package metartafparser

import "strings"

func splitN(s, sep string, n ...int) []string {
	split := strings.Split(s, sep)
	limit := -1
	if len(n) > 0 {
		limit = n[0]
	}
	if limit < 0 || len(split) <= limit {
		return split
	}
	out := make([]string, 0, limit)
	out = append(out, split[:limit]...)
	out = append(out, strings.Join(split[limit:], sep))
	return out
}

func splitFields(s string, n ...int) []string {
	split := strings.Fields(s)
	limit := -1
	if len(n) > 0 {
		limit = n[0]
	}
	if limit < 0 || len(split) <= limit {
		return split
	}
	out := make([]string, 0, limit)
	out = append(out, split[:limit]...)
	out = append(out, strings.Join(split[limit:], " "))
	return out
}

func resolve(obj map[string]any, path string) any {
	parts := strings.Split(path, ".")
	current := any(obj)
	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			if l, ok2 := current.(Locale); ok2 {
				m = map[string]any(l)
				ok = true
			}
		}
		if !ok {
			return nil
		}
		current, ok = m[part]
		if !ok {
			return nil
		}
	}
	return current
}

func isStation(s string) bool {
	return len(s) == 4
}

var allPhenomenonValues = []Phenomenon{
	PhenomenonRain, PhenomenonDrizzle, PhenomenonSnow, PhenomenonSnowGrains,
	PhenomenonIcePellets, PhenomenonIceCrystals, PhenomenonHail, PhenomenonSmallHail,
	PhenomenonUnknownPrecip, PhenomenonFog, PhenomenonVolcanicAsh, PhenomenonMist,
	PhenomenonHaze, PhenomenonWidespreadDust, PhenomenonSmoke, PhenomenonSand,
	PhenomenonSpray, PhenomenonSquall, PhenomenonSandWhirls, PhenomenonThunderstorm,
	PhenomenonDuststorm, PhenomenonSandstorm, PhenomenonFunnelCloud, PhenomenonNoSignificantWeather,
}

var allDescriptiveValues = []Descriptive{
	DescriptiveShowers, DescriptiveShallow, DescriptivePatches, DescriptivePartial,
	DescriptiveDrifting, DescriptiveThunderstorm, DescriptiveBlowing, DescriptiveFreezing,
}

func boolPtr(b bool) *bool                               { return &b }
func intPtr(v int) *int                                  { return &v }
func int64Ptr(i int64) *int64                            { return &i }
func float64Ptr(v float64) *float64                      { return &v }
func valueIndicatorPtr(v ValueIndicator) *ValueIndicator { return &v }
func wctPtr(w WeatherChangeType) *WeatherChangeType      { return &w }

func mergeContainer(dest, src *Container) {
	if src.Wind != nil {
		if dest.Wind == nil {
			dest.Wind = src.Wind
		} else {
			if src.Wind.Degrees != nil {
				dest.Wind.Degrees = src.Wind.Degrees
			}
			if src.Wind.Gust != nil {
				dest.Wind.Gust = src.Wind.Gust
			}
			if src.Wind.MinVariation != nil {
				dest.Wind.MinVariation = src.Wind.MinVariation
			}
			if src.Wind.MaxVariation != nil {
				dest.Wind.MaxVariation = src.Wind.MaxVariation
			}
			if src.Wind.Direction != "" {
				dest.Wind.Direction = src.Wind.Direction
			}
			if src.Wind.Unit != "" {
				dest.Wind.Unit = src.Wind.Unit
			}
			if src.Wind.Speed != 0 {
				dest.Wind.Speed = src.Wind.Speed
			}
		}
	}
	if src.Visibility != nil {
		if dest.Visibility == nil {
			dest.Visibility = src.Visibility
		} else {
			if src.Visibility.Min != nil {
				dest.Visibility.Min = src.Visibility.Min
			}
			if src.Visibility.Indicator != nil {
				dest.Visibility.Indicator = src.Visibility.Indicator
			}
			if src.Visibility.Value != 0 {
				dest.Visibility.Value = src.Visibility.Value
			}
			if src.Visibility.Unit != "" {
				dest.Visibility.Unit = src.Visibility.Unit
			}
			dest.Visibility.Ndv = dest.Visibility.Ndv || src.Visibility.Ndv
		}
	}
	if src.WindShear != nil {
		dest.WindShear = src.WindShear
	}
	if src.VerticalVisibility != nil {
		dest.VerticalVisibility = src.VerticalVisibility
	}
	if src.Cavok != nil {
		dest.Cavok = src.Cavok
	}
	dest.Clouds = append(dest.Clouds, src.Clouds...)
	dest.WeatherConditions = append(dest.WeatherConditions, src.WeatherConditions...)
}
