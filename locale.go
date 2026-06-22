package metartafparser

import (
	"fmt"
	"regexp"
)

// Locale maps localization keys to translated strings or nested maps.
type Locale map[string]any

func localeGet(path string, lang Locale) *string {
	val := resolve(lang, path)
	if val == nil {
		return nil
	}
	s, ok := val.(string)
	if !ok {
		return nil
	}
	return &s
}

var formatRe = regexp.MustCompile(`\{(\d+)\}`)

func formatMsg(message *string, args ...any) *string {
	if message == nil {
		return nil
	}
	for _, arg := range args {
		if arg == nil {
			return nil
		}
	}
	result := formatRe.ReplaceAllStringFunc(*message, func(match string) string {
		idxStr := match[1 : len(match)-1]
		var idx int
		_, _ = fmt.Sscanf(idxStr, "%d", &idx)
		if idx < len(args) {
			return fmt.Sprintf("%v", args[idx])
		}
		return match
	})
	return &result
}

// DefaultLocale returns a Locale pre-populated with English translations for all weather report fields.
func DefaultLocale() Locale {
	return Locale{
		"CloudQuantity":          buildCloudQuantityLocale(),
		"CloudType":              buildCloudTypeLocale(),
		"Converter":              buildConverterLocale(),
		"DepositBrakingCapacity": buildDepositBrakingCapacityLocale(),
		"DepositCoverage":        buildDepositCoverageLocale(),
		"DepositThickness":       buildDepositThicknessLocale(),
		"DepositType":            buildDepositTypeLocale(),
		"Descriptive":            buildDescriptiveLocale(),
		"Error":                  Locale{"prefix": "An error occurred. Error code n°"},
		"ErrorCode":              Locale{"AirportNotFound": "The airport was not found for this message.", "InvalidMessage": "The entered message is invalid."},
		"Indicator":              Locale{"M": "less than", "P": "greater than"},
		"Intensity":              Locale{"Light": "Light", "Moderate": "Moderate", "Heavy": "Heavy", "VC": "In the vicinity"},
		"MetarFacade":            Locale{"InvalidIcao": "Icao code is invalid."},
		"Phenomenon":             buildPhenomenonLocale(),
		"Remark":                 buildRemarkLocale(),
		"TimeIndicator":          Locale{"AT": "at", "FM": "From", "TL": "until"},
		"ToString":               buildToStringLocale(),
		"WeatherChangeType":      Locale{"BECMG": "Becoming", "FM": "From", "PROB": "Probability", "TEMPO": "Temporary"},
	}
}

func buildCloudQuantityLocale() map[string]any {
	return Locale{
		"BKN": "broken",
		"FEW": "few",
		"NSC": "no significant clouds.",
		"OVC": "overcast",
		"SCT": "scattered",
		"SKC": "sky clear",
	}
}

func buildCloudTypeLocale() map[string]any {
	return Locale{
		"AC": "Altocumulus", "AS": "Altostratus", "CB": "Cumulonimbus",
		"CC": "CirroCumulus", "CI": "Cirrus", "CS": "Cirrostratus",
		"CU": "Cumulus", "NS": "Nimbostratus", "SC": "Stratocumulus",
		"ST": "Stratus", "TCU": "Towering cumulus",
	}
}

func buildConverterLocale() map[string]any {
	return Locale{
		"D": "decreasing", "E": "East", "ENE": "East North East",
		"ESE": "East South East", "N": "North", "NE": "North East",
		"NNE": "North North East", "NNW": "North North West",
		"NSC": "no significant change", "NW": "North West", "S": "South",
		"SE": "South East", "SSE": "South South East",
		"SSW": "South South West", "SW": "South West",
		"U": "up rising", "VRB": "Variable", "W": "West",
		"WNW": "West North West", "WSW": "West South West",
	}
}

func buildDepositBrakingCapacityLocale() map[string]any {
	return Locale{
		"GOOD": "good", "MEDIUM": "medium", "MEDIUM_GOOD": "medium/good",
		"MEDIUM_POOR": "poor/medium", "NOT_REPORTED": "not reported",
		"POOR": "poor", "UNRELIABLE": "figures unreliable",
	}
}

func buildDepositCoverageLocale() map[string]any {
	return Locale{
		"FROM_11_TO_25": "from 11% to 25%", "FROM_26_TO_50": "from 26% to 50%",
		"FROM_51_TO_100": "from 51% to 100%", "LESS_10": "less than 10%",
		"NOT_REPORTED": "not reported",
	}
}

func buildDepositThicknessLocale() map[string]any {
	return Locale{
		"CLOSED": "closed", "LESS_1_MM": "less than 1 mm",
		"NOT_REPORTED": "not reported", "THICKNESS_10": "10 cm",
		"THICKNESS_15": "15 cm", "THICKNESS_20": "20 cm",
		"THICKNESS_25": "25 cm", "THICKNESS_30": "30 cm",
		"THICKNESS_35": "35 cm", "THICKNESS_40": "40 cm or more",
	}
}

func buildDepositTypeLocale() map[string]any {
	return Locale{
		"CLEAR_DRY": "clear and dry", "COMPACTED_SNOW": "compacted or rolled snow",
		"DAMP": "damp", "DRY_SNOW": "dry snow", "FROZEN_RIDGES": "frozen ruts or ridges",
		"ICE": "ice", "NOT_REPORTED": "not reported",
		"RIME_FROST_COVERED": "rime or frost covered", "SLUSH": "slush",
		"WET_SNOW": "wet snow", "WET_WATER_PATCHES": "wet or water patches",
	}
}

func buildDescriptiveLocale() map[string]any {
	return Locale{
		"BC": "patches", "BL": "blowing", "DR": "low drifting",
		"FZ": "freezing", "MI": "shallow", "PR": "partial",
		"SH": "showers of", "TS": "thunderstorm",
	}
}

func buildPhenomenonLocale() map[string]any {
	return Locale{
		"BR": "mist", "DS": "duststorm", "DU": "widespread dust",
		"DZ": "drizzle", "FC": "funnel cloud", "FG": "fog",
		"FU": "smoke", "GR": "hail", "GS": "small hail and/or snow pellets",
		"HZ": "haze", "IC": "ice crystals", "PL": "ice pellets",
		"PO": "dust or sand whirls", "PY": "spray", "RA": "rain",
		"SA": "sand", "SG": "snow grains", "SN": "snow",
		"SQ": "squall", "SS": "sandstorm", "TS": "thunderstorm",
		"UP": "unknown precipitation", "VA": "volcanic ash",
		"NSW": "no significant weather",
	}
}

func buildRemarkLocale() map[string]any {
	return Locale{
		"ALQDS": "all quadrants",
		"AO1":   "automated stations without a precipitation discriminator",
		"AO2":   "automated station with a precipitation discriminator",
		"AO2A":  "automated station with a precipitation discriminator (augmented)",
		"BASED": "based",
		"Barometer": []any{
			"Increase, then decrease",
			"Increase, then steady, or increase then Increase more slowly",
			"steady or unsteady increase",
			"Decrease or steady, then increase; or increase then increase more rapidly",
			"Steady",
			"Decrease, then increase",
			"Decrease then steady; or decrease then decrease more slowly",
			"Steady or unsteady decrease",
			"Steady or increase, then decrease; or decrease then decrease more rapidly",
		},
		"Ceiling": Locale{
			"Height": "ceiling varying between {0} and {1} feet",
			"Second": Locale{
				"Location": "ceiling of {0} feet measured by a second sensor located at {1}",
			},
		},
		"DSNT":        "distant",
		"FCST":        "forecast",
		"FUNNELCLOUD": "funnel cloud",
		"HVY":         "heavy",
		"Hail": Locale{
			"0":          "largest hailstones with a diameter of {0} inches",
			"LesserThan": "largest hailstones with a diameter less than {0} inches",
		},
		"Hourly": Locale{
			"Maximum": Locale{
				"Minimum": Locale{
					"Temperature": "24-hour maximum temperature of {0}°C and 24-hour minimum temperature of {1}°C",
				},
				"Temperature": "6-hourly maximum temperature of {0}°C",
			},
			"Minimum": Locale{
				"Temperature": "6-hourly minimum temperature of {0}°C",
			},
			"Temperature": Locale{
				"0": "hourly temperature of {0}°C",
				"Dew": Locale{
					"Point": "hourly temperature of {0}°C and dew point of {1}°C",
				},
			},
		},
		"Ice": Locale{
			"Accretion": Locale{
				"Amount": "{0}/100 of an inch of ice accretion in the past {1} hour(s)",
			},
		},
		"LGT": "light",
		"LTG": "lightning",
		"MOD": "moderate",
		"Next": Locale{
			"Forecast": Locale{
				"By": "next forecast by {0}, {1}:{2}Z",
			},
		},
		"NXT":         "next",
		"ON":          "on",
		"Obscuration": "{0} layer at {1} feet composed of {2}",
		"PRESFR":      "pressure falling rapidly",
		"PRESRR":      "pressure rising rapidly",
		"PeakWind":    "peak wind of {1} knots from {0} degrees at {2}:{3}",
		"Precipitation": Locale{
			"Amount": Locale{
				"24": "{0} inches of precipitation fell in the last 24 hours",
				"3": Locale{
					"6": "{1} inches of precipitation fell in the last {0} hours",
				},
				"Hourly": "{0} inches of precipitation fell in the last hour",
			},
			"Beg": Locale{
				"0":   "{0} {1} beginning at {2}:{3}",
				"End": "{0} {1} beginning at {2}:{3} ending at {4}:{5}",
			},
			"End": "{0} {1} ending at {2}:{3}",
		},
		"Pressure": Locale{
			"Tendency": "of {0} hectopascals in the past 3 hours",
		},
		"SLPNO": "sea level pressure not available",
		"Sea": Locale{
			"Level": Locale{
				"Pressure": "sea level pressure of {0} HPa",
			},
		},
		"Second": Locale{
			"Location": Locale{
				"Visibility": "visibility of {0} SM measured by a second sensor located at {1}",
			},
		},
		"Sector": Locale{
			"Visibility": "visibility of {1} SM in the {0} direction",
		},
		"Snow": Locale{
			"Depth": "snow depth of {0} inches",
			"Increasing": Locale{
				"Rapidly": "snow depth increase of {0} inches in the past hour with a total depth on the ground of {1} inches",
			},
			"Pellets": "{0} snow pellets",
		},
		"Sunshine": Locale{
			"Duration": "{0} minutes of sunshine",
		},
		"Surface": Locale{
			"Visibility": "surface visibility of {0} statute miles",
		},
		"TORNADO": "tornado",
		"Thunderstorm": Locale{
			"Location": Locale{
				"0":      "thunderstorm {0} of the station",
				"Moving": "thunderstorm {0} of the station moving towards {1}",
			},
		},
		"Tornadic": Locale{
			"Activity": Locale{
				"BegEnd":    "{0} beginning at {1}:{2} ending at {3}:{4} {5} SM {6} of the station",
				"Beginning": "{0} beginning at {1}:{2} {3} SM {4} of the station",
				"Ending":    "{0} ending at {1}:{2} {3} SM {4} of the station",
			},
		},
		"Tower": Locale{
			"Visibility": "control tower visibility of {0} statute miles",
		},
		"VIRGA":      "virga",
		"WATERSPOUT": "waterspout",
		"Variable": Locale{
			"Prevailing": Locale{
				"Visibility": "variable prevailing visibility between {0} and {1} SM",
			},
			"Sky": Locale{
				"Condition": Locale{
					"0":      "cloud layer varying between {0} and {1}",
					"Height": "cloud layer at {0} feet varying between {1} and {2}",
				},
			},
		},
		"Virga": Locale{
			"Direction": "virga {0} from the station",
		},
		"Water": Locale{
			"Equivalent": Locale{
				"Snow": Locale{
					"Ground": "water equivalent of {0} inches of snow",
				},
			},
		},
		"WindShift": Locale{
			"0":     "wind shift at {0}:{1}",
			"FROPA": "wind shift accompanied by frontal passage at {0}:{1}",
		},
	}
}

func buildToStringLocale() map[string]any {
	return Locale{
		"airport": "airport", "altimeter": "altimeter (hPa)",
		"amendment": "amendment", "auto": "auto", "cavok": "cavok",
		"clouds": "clouds",
		"day":    Locale{"hour": "hour of the day", "month": "day of the month"},
		"deposit": Locale{
			"braking": "braking capacity", "coverage": "coverage",
			"thickness": "thickness", "type": "type of deposit",
		},
		"descriptive": "descriptive",
		"dew":         Locale{"point": "dew point"},
		"end": Locale{
			"day":  Locale{"month": "end day of the month"},
			"hour": Locale{"day": "end hour of the day"},
		},
		"height":    Locale{"feet": "height (ft)", "meter": "height (m)"},
		"indicator": "indicator", "intensity": "intensity",
		"message": "original message", "name": "name", "nosig": "nosig",
		"phenomena": "phenomena", "probability": "probability",
		"quantity": "quantity", "remark": "remarks",
		"report": Locale{"time": "time of report"},
		"runway": Locale{"info": "runways information"},
		"start": Locale{
			"day":    Locale{"month": "starting day of the month"},
			"hour":   Locale{"day": "starting hour of the day"},
			"minute": "starting minute",
		},
		"temperature": Locale{"0": "temperature (°C)", "max": "maximum temperature (°C)", "min": "minimum temperature (°C)"},
		"trend":       "trend", "trends": "trends", "type": "type",
		"vertical":   Locale{"visibility": "vertical visibility (ft)"},
		"visibility": Locale{"main": "main visibility", "max": "maximum visibility", "min": Locale{"0": "minimum visibility", "direction": "minimum visibility direction"}},
		"weather":    Locale{"conditions": "weather conditions"},
		"wind": Locale{
			"direction": Locale{"0": "direction", "degrees": "direction (degrees)"},
			"gusts":     "gusts",
			"max":       Locale{"variation": "maximal wind variation"},
			"min":       Locale{"variation": "minimal wind variation"},
			"speed":     "speed", "unit": "unit",
		},
	}
}

func isKnownRemarkType(s string) bool {
	switch RemarkType(s) {
	case RemarkTypeAO1, RemarkTypeAO2, RemarkTypePRESFR, RemarkTypePRESRR,
		RemarkTypeTORNADO, RemarkTypeFUNNELCLOUD, RemarkTypeWATERSPOUT, RemarkTypeVIRGA,
		RemarkTypeWindPeak, RemarkTypeWindShiftFropa, RemarkTypeWindShift,
		RemarkTypeTowerVisibility, RemarkTypeSurfaceVisibility,
		RemarkTypePrevailingVisibility, RemarkTypeSecondLocationVisibility,
		RemarkTypeSectorVisibility, RemarkTypeTornadicActivityBegEnd,
		RemarkTypeTornadicActivityBeg, RemarkTypeTornadicActivityEnd,
		RemarkTypePrecipitationBeg, RemarkTypePrecipitationBegEnd,
		RemarkTypePrecipitationEnd, RemarkTypeThunderStormLocationMoving,
		RemarkTypeThunderStormLocation, RemarkTypeSmallHailSize, RemarkTypeHailSize,
		RemarkTypeSnowPellets, RemarkTypeVirgaDirection, RemarkTypeCeilingHeight,
		RemarkTypeObscuration, RemarkTypeVariableSkyHeight, RemarkTypeVariableSky,
		RemarkTypeCeilingSecondLocation, RemarkTypeSeaLevelPressure,
		RemarkTypeSnowIncrease, RemarkTypeHourlyMaximumMinimumTemperature,
		RemarkTypeHourlyMaximumTemperature, RemarkTypeHourlyMinimumTemperature,
		RemarkTypeHourlyPrecipitationAmount, RemarkTypeHourlyTemperatureDewPoint,
		RemarkTypeHourlyPressure, RemarkTypeIceAccretion,
		RemarkTypePrecipitationAmount36Hour, RemarkTypePrecipitationAmount24Hour,
		RemarkTypeSnowDepth, RemarkTypeSunshineDuration,
		RemarkTypeWaterEquivalentSnow, RemarkTypeNextForecastBy:
		return true
	}
	return false
}
