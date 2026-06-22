package metartafparser

// MetarType distinguishes between routine METAR and special SPECI reports.
type MetarType string

const (
	MetarTypeMETAR MetarType = "METAR"
	MetarTypeSPECI MetarType = "SPECI"
)

// CloudQuantity indicates the amount of sky cover (e.g., FEW, SCT, BKN, OVC).
type CloudQuantity string

const (
	CloudQuantitySKC CloudQuantity = "SKC"
	CloudQuantityFEW CloudQuantity = "FEW"
	CloudQuantityBKN CloudQuantity = "BKN"
	CloudQuantitySCT CloudQuantity = "SCT"
	CloudQuantityOVC CloudQuantity = "OVC"
	CloudQuantityNSC CloudQuantity = "NSC"
)

// CloudType classifies cloud types (e.g., Cumulonimbus, Stratus, etc.).
type CloudType string

const (
	CloudTypeCB  CloudType = "CB"
	CloudTypeTCU CloudType = "TCU"
	CloudTypeCI  CloudType = "CI"
	CloudTypeCC  CloudType = "CC"
	CloudTypeCS  CloudType = "CS"
	CloudTypeAC  CloudType = "AC"
	CloudTypeST  CloudType = "ST"
	CloudTypeCU  CloudType = "CU"
	CloudTypeAS  CloudType = "AS"
	CloudTypeNS  CloudType = "NS"
	CloudTypeSC  CloudType = "SC"
)

// Intensity qualifies the intensity of a weather phenomenon (light, heavy, in the vicinity).
type Intensity string

const (
	IntensityLight      Intensity = "Light"
	IntensityHeavy      Intensity = "Heavy"
	IntensityModerate   Intensity = "Moderate"
	IntensityInVicinity Intensity = "VC"
)

// Descriptive qualifies a weather phenomenon with descriptive terms (e.g., showers, thunderstorm).
type Descriptive string

const (
	DescriptiveShowers      Descriptive = "SH"
	DescriptiveShallow      Descriptive = "MI"
	DescriptivePatches      Descriptive = "BC"
	DescriptivePartial      Descriptive = "PR"
	DescriptiveDrifting     Descriptive = "DR"
	DescriptiveThunderstorm Descriptive = "TS"
	DescriptiveBlowing      Descriptive = "BL"
	DescriptiveFreezing     Descriptive = "FZ"
)

// Phenomenon enumerates weather phenomena (rain, snow, fog, etc.).
type Phenomenon string

const (
	PhenomenonRain                 Phenomenon = "RA"
	PhenomenonDrizzle              Phenomenon = "DZ"
	PhenomenonSnow                 Phenomenon = "SN"
	PhenomenonSnowGrains           Phenomenon = "SG"
	PhenomenonIcePellets           Phenomenon = "PL"
	PhenomenonIceCrystals          Phenomenon = "IC"
	PhenomenonHail                 Phenomenon = "GR"
	PhenomenonSmallHail            Phenomenon = "GS"
	PhenomenonUnknownPrecip        Phenomenon = "UP"
	PhenomenonFog                  Phenomenon = "FG"
	PhenomenonVolcanicAsh          Phenomenon = "VA"
	PhenomenonMist                 Phenomenon = "BR"
	PhenomenonHaze                 Phenomenon = "HZ"
	PhenomenonWidespreadDust       Phenomenon = "DU"
	PhenomenonSmoke                Phenomenon = "FU"
	PhenomenonSand                 Phenomenon = "SA"
	PhenomenonSpray                Phenomenon = "PY"
	PhenomenonSquall               Phenomenon = "SQ"
	PhenomenonSandWhirls           Phenomenon = "PO"
	PhenomenonThunderstorm         Phenomenon = "TS"
	PhenomenonDuststorm            Phenomenon = "DS"
	PhenomenonSandstorm            Phenomenon = "SS"
	PhenomenonFunnelCloud          Phenomenon = "FC"
	PhenomenonTornado              Phenomenon = "FC"
	PhenomenonNoSignificantWeather Phenomenon = "NSW"
)

// TimeIndicator specifies whether a time is exact (AT), from (FM), or until (TL).
type TimeIndicator string

const (
	TimeIndicatorAT TimeIndicator = "AT"
	TimeIndicatorFM TimeIndicator = "FM"
	TimeIndicatorTL TimeIndicator = "TL"
)

// WeatherChangeType describes how weather conditions change (BECMG, TEMPO, FM, etc.).
type WeatherChangeType string

const (
	WeatherChangeTypeFM    WeatherChangeType = "FM"
	WeatherChangeTypeBECMG WeatherChangeType = "BECMG"
	WeatherChangeTypeTEMPO WeatherChangeType = "TEMPO"
	WeatherChangeTypeINTER WeatherChangeType = "INTER"
	WeatherChangeTypePROB  WeatherChangeType = "PROB"
)

// Direction represents a cardinal or intercardinal wind direction, or VRB (variable).
type Direction string

const (
	DirectionVRB Direction = "VRB"
	DirectionE   Direction = "E"
	DirectionENE Direction = "ENE"
	DirectionESE Direction = "ESE"
	DirectionN   Direction = "N"
	DirectionNE  Direction = "NE"
	DirectionNNE Direction = "NNE"
	DirectionNNW Direction = "NNW"
	DirectionNW  Direction = "NW"
	DirectionS   Direction = "S"
	DirectionSE  Direction = "SE"
	DirectionSSE Direction = "SSE"
	DirectionSSW Direction = "SSW"
	DirectionSW  Direction = "SW"
	DirectionW   Direction = "W"
	DirectionWNW Direction = "WNW"
	DirectionWSW Direction = "WSW"
)

// DistanceUnit specifies whether a distance is in meters or statute miles.
type DistanceUnit string

const (
	DistanceUnitMeters       DistanceUnit = "m"
	DistanceUnitStatuteMiles DistanceUnit = "SM"
)

// SpeedUnit specifies the unit for wind speed (knots, m/s, km/h).
type SpeedUnit string

const (
	SpeedUnitKnot              SpeedUnit = "KT"
	SpeedUnitMetersPerSecond   SpeedUnit = "MPS"
	SpeedUnitKilometersPerHour SpeedUnit = "KM/H"
)

// ValueIndicator marks a value as greater than (P) or less than (M) the reported number.
type ValueIndicator string

const (
	ValueIndicatorGreaterThan ValueIndicator = "P"
	ValueIndicatorLessThan    ValueIndicator = "M"
)

// RunwayInfoTrend indicates whether runway visual range is increasing, decreasing, or unchanged.
type RunwayInfoTrend string

const (
	RunwayInfoTrendUprising            RunwayInfoTrend = "U"
	RunwayInfoTrendDecreasing          RunwayInfoTrend = "D"
	RunwayInfoTrendNoSignificantChange RunwayInfoTrend = "N"
)

// RunwayInfoUnit specifies the unit for runway visual range (feet or meters).
type RunwayInfoUnit string

const (
	RunwayInfoUnitFeet   RunwayInfoUnit = "FT"
	RunwayInfoUnitMeters RunwayInfoUnit = "m"
)

// IcingIntensity rates the severity and type of in-flight icing on a 0-9 scale.
type IcingIntensity string

const (
	IcingIntensityNone                     IcingIntensity = "0"
	IcingIntensityLight                    IcingIntensity = "1"
	IcingIntensityLightRimeIcingCloud      IcingIntensity = "2"
	IcingIntensityLightClearIcingPrecip    IcingIntensity = "3"
	IcingIntensityModerateMixedIcing       IcingIntensity = "4"
	IcingIntensityModerateRimeIcingCloud   IcingIntensity = "5"
	IcingIntensityModerateClearIcingPrecip IcingIntensity = "6"
	IcingIntensitySevereMixedIcing         IcingIntensity = "7"
	IcingIntensitySevereRimeIcingCloud     IcingIntensity = "8"
	IcingIntensitySevereClearIcingPrecip   IcingIntensity = "9"
)

// TurbulenceIntensity rates turbulence severity and frequency on a 0-9 scale plus extreme.
type TurbulenceIntensity string

const (
	TurbulenceIntensityNone                 TurbulenceIntensity = "0"
	TurbulenceIntensityLight                TurbulenceIntensity = "1"
	TurbulenceIntensityModerateClearAirOcc  TurbulenceIntensity = "2"
	TurbulenceIntensityModerateClearAirFreq TurbulenceIntensity = "3"
	TurbulenceIntensityModerateCloudOcc     TurbulenceIntensity = "4"
	TurbulenceIntensityModerateCloudFreq    TurbulenceIntensity = "5"
	TurbulenceIntensitySevereClearAirOcc    TurbulenceIntensity = "6"
	TurbulenceIntensitySevereClearAirFreq   TurbulenceIntensity = "7"
	TurbulenceIntensitySevereCloudOcc       TurbulenceIntensity = "8"
	TurbulenceIntensitySevereCloudFreq      TurbulenceIntensity = "9"
	TurbulenceIntensityExtreme              TurbulenceIntensity = "X"
)

// DepositType describes the type of surface deposit on a runway.
type DepositType string

const (
	DepositTypeNotReported      DepositType = "/"
	DepositTypeClearDry         DepositType = "0"
	DepositTypeDamp             DepositType = "1"
	DepositTypeWetWaterPatches  DepositType = "2"
	DepositTypeRimeFrostCovered DepositType = "3"
	DepositTypeDrySnow          DepositType = "4"
	DepositTypeWetSnow          DepositType = "5"
	DepositTypeSlush            DepositType = "6"
	DepositTypeIce              DepositType = "7"
	DepositTypeCompactedSnow    DepositType = "8"
	DepositTypeFrozenRidges     DepositType = "9"
)

// DepositCoverage indicates the percentage of runway area covered by a deposit.
type DepositCoverage string

const (
	DepositCoverageNone        DepositCoverage = "0"
	DepositCoverageNotReported DepositCoverage = "/"
	DepositCoverageLess10      DepositCoverage = "1"
	DepositCoverageFrom11To25  DepositCoverage = "2"
	DepositCoverageFrom26To50  DepositCoverage = "5"
	DepositCoverageFrom51To100 DepositCoverage = "9"
)

// AltimeterUnit specifies whether altimeter pressure is in inches of mercury or hectopascals.
type AltimeterUnit string

const (
	AltimeterUnitInHg AltimeterUnit = "inHg"
	AltimeterUnitHPa  AltimeterUnit = "hPa"
)

// RemarkType enumerates all possible remark types found in METAR/TAF reports.
type RemarkType string

const (
	RemarkTypeUnknown                         RemarkType = "Unknown"
	RemarkTypeAO1                             RemarkType = "AO1"
	RemarkTypeAO2                             RemarkType = "AO2"
	RemarkTypePRESFR                          RemarkType = "PRESFR"
	RemarkTypePRESRR                          RemarkType = "PRESRR"
	RemarkTypeTORNADO                         RemarkType = "TORNADO"
	RemarkTypeFUNNELCLOUD                     RemarkType = "FUNNELCLOUD"
	RemarkTypeWATERSPOUT                      RemarkType = "WATERSPOUT"
	RemarkTypeVIRGA                           RemarkType = "VIRGA"
	RemarkTypeWindPeak                        RemarkType = "WindPeak"
	RemarkTypeWindShiftFropa                  RemarkType = "WindShiftFropa"
	RemarkTypeWindShift                       RemarkType = "WindShift"
	RemarkTypeTowerVisibility                 RemarkType = "TowerVisibility"
	RemarkTypeSurfaceVisibility               RemarkType = "SurfaceVisibility"
	RemarkTypePrevailingVisibility            RemarkType = "PrevailingVisibility"
	RemarkTypeSecondLocationVisibility        RemarkType = "SecondLocationVisibility"
	RemarkTypeSectorVisibility                RemarkType = "SectorVisibility"
	RemarkTypeTornadicActivityBegEnd          RemarkType = "TornadicActivityBegEnd"
	RemarkTypeTornadicActivityBeg             RemarkType = "TornadicActivityBeg"
	RemarkTypeTornadicActivityEnd             RemarkType = "TornadicActivityEnd"
	RemarkTypePrecipitationBeg                RemarkType = "PrecipitationBeg"
	RemarkTypePrecipitationBegEnd             RemarkType = "PrecipitationBegEnd"
	RemarkTypePrecipitationEnd                RemarkType = "PrecipitationEnd"
	RemarkTypeThunderStormLocationMoving      RemarkType = "ThunderStormLocationMoving"
	RemarkTypeThunderStormLocation            RemarkType = "ThunderStormLocation"
	RemarkTypeSmallHailSize                   RemarkType = "SmallHailSize"
	RemarkTypeHailSize                        RemarkType = "HailSize"
	RemarkTypeSnowPellets                     RemarkType = "SnowPellets"
	RemarkTypeVirgaDirection                  RemarkType = "VirgaDirection"
	RemarkTypeCeilingHeight                   RemarkType = "CeilingHeight"
	RemarkTypeObscuration                     RemarkType = "Obscuration"
	RemarkTypeVariableSkyHeight               RemarkType = "VariableSkyHeight"
	RemarkTypeVariableSky                     RemarkType = "VariableSky"
	RemarkTypeCeilingSecondLocation           RemarkType = "CeilingSecondLocation"
	RemarkTypeSeaLevelPressure                RemarkType = "SeaLevelPressure"
	RemarkTypeSnowIncrease                    RemarkType = "SnowIncrease"
	RemarkTypeHourlyMaximumMinimumTemperature RemarkType = "HourlyMaximumMinimumTemperature"
	RemarkTypeHourlyMaximumTemperature        RemarkType = "HourlyMaximumTemperature"
	RemarkTypeHourlyMinimumTemperature        RemarkType = "HourlyMinimumTemperature"
	RemarkTypeHourlyPrecipitationAmount       RemarkType = "HourlyPrecipitationAmount"
	RemarkTypeHourlyTemperatureDewPoint       RemarkType = "HourlyTemperatureDewPoint"
	RemarkTypeHourlyPressure                  RemarkType = "HourlyPressure"
	RemarkTypeIceAccretion                    RemarkType = "IceAccretion"
	RemarkTypePrecipitationAmount36Hour       RemarkType = "PrecipitationAmount36Hour"
	RemarkTypePrecipitationAmount24Hour       RemarkType = "PrecipitationAmount24Hour"
	RemarkTypeSnowDepth                       RemarkType = "SnowDepth"
	RemarkTypeSunshineDuration                RemarkType = "SunshineDuration"
	RemarkTypeWaterEquivalentSnow             RemarkType = "WaterEquivalentSnow"
	RemarkTypeNextForecastBy                  RemarkType = "NextForecastBy"
)
