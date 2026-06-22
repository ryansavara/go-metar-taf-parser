package metartafparser

import (
	"math"
	"time"
)

// ParseOptions configures parsing behavior.
type ParseOptions struct {
	// Locale for localized string output. Use DefaultLocale() for English.
	Locale Locale
}

// ParseMetar parses a raw METAR string and returns the structured result.
func ParseMetar(rawMetar string, opts *ParseOptions) (*Metar, error) {
	locale := DefaultLocale()
	if opts != nil && opts.Locale != nil {
		locale = opts.Locale
	}
	parser := newMetarParser(locale)
	return parser.Parse(rawMetar)
}

// ParseMetarDated parses a raw METAR string and hydrates dates using the issued time.
func ParseMetarDated(rawMetar string, issued time.Time, opts *ParseOptions) (*MetarDated, error) {
	locale := DefaultLocale()
	if opts != nil && opts.Locale != nil {
		locale = opts.Locale
	}
	parser := newMetarParser(locale)
	metar, err := parser.Parse(rawMetar)
	if err != nil {
		return nil, err
	}
	return metarDatesHydrator(metar, issued), nil
}

// ParseTAF parses a raw TAF string and returns the structured result.
func ParseTAF(rawTAF string, opts *ParseOptions) (*TAF, error) {
	locale := DefaultLocale()
	if opts != nil && opts.Locale != nil {
		locale = opts.Locale
	}
	parser := newTAFParser(locale)
	return parser.Parse(rawTAF)
}

// ParseTAFDated parses a raw TAF string and hydrates dates using the issued time.
func ParseTAFDated(rawTAF string, issued time.Time, opts *ParseOptions) (*TAFDated, error) {
	locale := DefaultLocale()
	if opts != nil && opts.Locale != nil {
		locale = opts.Locale
	}
	parser := newTAFParser(locale)
	taf, err := parser.Parse(rawTAF)
	if err != nil {
		return nil, err
	}
	return tafDatesHydrator(taf, issued), nil
}

// ParseTAFAsForecast parses a raw TAF string and returns a forecast container.
func ParseTAFAsForecast(rawTAF string, issued time.Time, opts *ParseOptions) (*ForecastContainer, error) {
	taf, err := ParseTAFDated(rawTAF, issued, opts)
	if err != nil {
		return nil, err
	}
	return getForecastFromTAF(taf), nil
}

// ============================================================
// Date Hydration
// ============================================================

func determineReportDate(date time.Time, day, hour *int, minute ...int) time.Time {
	if day == nil || hour == nil {
		return date
	}
	minVal := 0
	if len(minute) > 0 {
		minVal = minute[0]
	}

	candidates := []time.Time{
		addMonthsUTC(date, -1),
		date,
		addMonthsUTC(date, 1),
	}

	best := candidates[0]
	bestDiff := math.MaxInt64

	for _, c := range candidates {
		d := setDateComponents(c, *day, *hour, minVal)
		diff := int(math.Abs(float64(d.Unix() - date.Unix())))
		if diff < bestDiff {
			bestDiff = diff
			best = d
		}
	}
	return best
}

func setDateComponents(date time.Time, day, hour, minute int) time.Time {
	y, m, _ := date.Date()
	loc := date.Location()
	return time.Date(y, m, day, hour, minute, 0, 0, loc)
}

func addMonthsUTC(date time.Time, count int) time.Time {
	y, m, d := date.Date()
	newM := int(m) + count
	newY := y + (newM-1)/12
	newM = ((newM - 1) % 12) + 1

	// Handle day overflow
	firstOfMonth := time.Date(newY, time.Month(newM), 1, 0, 0, 0, 0, date.Location())
	daysInMonth := firstOfMonth.AddDate(0, 1, -1).Day()
	if d > daysInMonth {
		d = daysInMonth
	}
	return time.Date(newY, time.Month(newM), d, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
}

// metarDatesHydrator creates a MetarDated from a Metar and issued date.
func metarDatesHydrator(report *Metar, date time.Time) *MetarDated {
	minVal := 0
	if report.Minute != nil {
		minVal = *report.Minute
	}
	return &MetarDated{
		Metar:  *report,
		Issued: determineReportDate(date, report.Day, report.Hour, minVal).UnixMilli(),
	}
}

// tafDatesHydrator creates a TAFDated from a TAF and issued date.
func tafDatesHydrator(report *TAF, date time.Time) *TAFDated {
	minVal := 0
	if report.Minute != nil {
		minVal = *report.Minute
	}
	issued := determineReportDate(date, report.Day, report.Hour, minVal)
	result := &TAFDated{
		TAF:    *report,
		Issued: issued.UnixMilli(),
	}

	start := determineReportDate(issued, &report.Validity.StartDay, &report.Validity.StartHour)
	end := determineReportDate(issued, &report.Validity.EndDay, &report.Validity.EndHour)

	result.ValidityDated = ValidityDated{
		StartDay:  report.Validity.StartDay,
		StartHour: report.Validity.StartHour,
		EndDay:    report.Validity.EndDay,
		EndHour:   report.Validity.EndHour,
		Start:     start.UnixMilli(),
		End:       end.UnixMilli(),
	}

	if report.MaxTemperature != nil {
		t := *report.MaxTemperature
		result.MaxTemperature = &Temperature{
			Temperature: t.Temperature,
			Day:         t.Day,
			Hour:        t.Hour,
		}
	}
	if report.MinTemperature != nil {
		t := *report.MinTemperature
		result.MinTemperature = &Temperature{
			Temperature: t.Temperature,
			Day:         t.Day,
			Hour:        t.Hour,
		}
	}

	// Hydrate trends with dates
	newTrends := make([]TAFTrend, 0, len(report.Trends))
	for _, trend := range report.Trends {
		trend.Remarks = hydrateRemarkDates(trend.Remarks, issued)
		if trend.Validity != nil {
			switch v := trend.Validity.(type) {
			case FMValidity:
				start := determineReportDate(issued, &v.StartDay, &v.StartHour, v.StartMinutes)
				trend.ValidityStartMs = int64Ptr(start.UnixMilli())
			case Validity:
				start := determineReportDate(issued, &v.StartDay, &v.StartHour)
				end := determineReportDate(issued, &v.EndDay, &v.EndHour)
				trend.ValidityStartMs = int64Ptr(start.UnixMilli())
				trend.ValidityEndMs = int64Ptr(end.UnixMilli())
			}
		}
		newTrends = append(newTrends, trend)
	}
	result.Trends = newTrends

	// Hydrate remarks
	result.Remarks = hydrateRemarkDates(report.Remarks, issued)

	return result
}

func hydrateRemarkDates(remarks []Remark, issued time.Time) []Remark {
	result := make([]Remark, 0, len(remarks))
	for _, r := range remarks {
		if r.Type == RemarkTypeNextForecastBy && r.Day != nil && r.Hour != nil {
			d := determineReportDate(issued, r.Day, r.Hour, *r.Minute)
			ms := d.UnixMilli()
			r.Date = &ms
		}
		result = append(result, r)
	}
	return result
}

// ============================================================
// Forecast
// ============================================================

func getForecastFromTAF(taf *TAFDated) *ForecastContainer {
	issuedTime := time.UnixMilli(taf.Issued).UTC()
	start := determineReportDate(issuedTime, &taf.ValidityDated.StartDay, &taf.ValidityDated.StartHour)
	end := determineReportDate(issuedTime, &taf.ValidityDated.EndDay, &taf.ValidityDated.EndHour)

	fc := &ForecastContainer{
		Station:   taf.Station,
		Issued:    taf.Issued,
		Start:     start.UnixMilli(),
		End:       end.UnixMilli(),
		Message:   taf.Message,
		Amendment: taf.Amendment,
		Auto:      taf.Auto,
		Canceled:  taf.Canceled,
		Corrected: taf.Corrected,
		Nil:       taf.Nil,
	}

	if taf.MaxTemperature != nil {
		maxTemp := *taf.MaxTemperature
		maxDate := determineReportDate(issuedTime, &maxTemp.Day, &maxTemp.Hour)
		fc.MaxTemperature = &TemperatureDated{
			Temperature: maxTemp.Temperature,
			Day:         maxTemp.Day,
			Hour:        maxTemp.Hour,
			Date:        maxDate.UnixMilli(),
		}
	}
	if taf.MinTemperature != nil {
		minTemp := *taf.MinTemperature
		minDate := determineReportDate(issuedTime, &minTemp.Day, &minTemp.Hour)
		fc.MinTemperature = &TemperatureDated{
			Temperature: minTemp.Temperature,
			Day:         minTemp.Day,
			Hour:        minTemp.Hour,
			Date:        minDate.UnixMilli(),
		}
	}

	initialForecast := forecastFromContainer(&taf.Container, taf.InitialRaw, nil, taf.ValidityDated.Start, taf.ValidityDated.End, nil)

	trendForecasts := convertTrendsToForecasts(taf.Trends)
	allForecasts := append([]Forecast{initialForecast}, trendForecasts...)

	fc.Forecast = hydrateEndDates(allForecasts, taf.ValidityDated)

	return fc
}

func forecastFromContainer(c *Container, raw string, typ *WeatherChangeType, start, end int64, trend *TAFTrend) Forecast {
	forecast := Forecast{
		Wind:               c.Wind,
		Visibility:         c.Visibility,
		WindShear:          c.WindShear,
		VerticalVisibility: c.VerticalVisibility,
		Cavok:              c.Cavok,
		Remark:             c.Remark,
		Raw:                raw,
		Remarks:            c.Remarks,
		Clouds:             c.Clouds,
		WeatherConditions:  c.WeatherConditions,
		Turbulence:         c.Turbulence,
		Icing:              c.Icing,
		Type:               typ,
		Start:              start,
		End:                end,
	}
	if trend != nil && trend.ValidityEndMs != nil && typ != nil && *typ == WeatherChangeTypeBECMG {
		forecast.By = trend.ValidityEndMs
	}
	return forecast
}

func convertTrendsToForecasts(trends []TAFTrend) []Forecast {
	result := make([]Forecast, 0, len(trends))
	for _, t := range trends {
		start := int64(0)
		if t.ValidityStartMs != nil {
			start = *t.ValidityStartMs
		}
		end := int64(0)
		if t.ValidityEndMs != nil {
			end = *t.ValidityEndMs
		}
		f := forecastFromContainer(&t.Container, t.Raw, wctPtr(t.Type), start, end, &t)
		result = append(result, f)
	}
	return result
}

func isImplicitEnd(f Forecast) bool {
	return f.Type == nil || *f.Type == WeatherChangeTypeFM || *f.Type == WeatherChangeTypeBECMG
}

func hydrateEndDates(forecasts []Forecast, reportValidity ValidityDated) []Forecast {
	findNext := func(index int) *Forecast {
		for i := index; i < len(forecasts); i++ {
			if isImplicitEnd(forecasts[i]) {
				return &forecasts[i]
			}
		}
		return nil
	}

	result := make([]Forecast, 0, len(forecasts))
	var previouslyHydrated *Forecast

	for i := range forecasts {
		f := forecasts[i]

		if !isImplicitEnd(f) {
			// TEMPO/INTER/PROB: keep explicit end date
			result = append(result, f)
			continue
		}

		// FM, BECMG, or initial forecast: end is determined by next FM/BECMG or report validity end
		next := findNext(i + 1)
		if next != nil {
			f.End = next.Start
		} else {
			f.End = reportValidity.End
		}

		f = hydrateWithPreviousContextIfNeeded(f, previouslyHydrated)
		result = append(result, f)
		previouslyHydrated = &result[len(result)-1]
	}

	return result
}

func hydrateWithPreviousContextIfNeeded(forecast Forecast, context *Forecast) Forecast {
	if forecast.Type == nil || *forecast.Type != WeatherChangeTypeBECMG || context == nil {
		return forecast
	}

	// BECMG inherits previous context for fields not explicitly set
	ctx := *context
	ctx.Remark = nil
	ctx.Remarks = nil

	// vertical visibility should not be carried over if clouds exist
	if len(forecast.Clouds) > 0 {
		ctx.VerticalVisibility = nil
	}

	// CAVOK should not propagate if anything other than wind changes
	if len(forecast.Clouds) > 0 || forecast.VerticalVisibility != nil || len(forecast.WeatherConditions) > 0 || forecast.Visibility != nil {
		ctx.Cavok = nil
	}

	// Merge context into forecast (forecast fields take priority)
	if forecast.Wind == nil {
		forecast.Wind = ctx.Wind
	}
	if forecast.Visibility == nil {
		forecast.Visibility = ctx.Visibility
	}
	if forecast.VerticalVisibility == nil {
		forecast.VerticalVisibility = ctx.VerticalVisibility
	}
	if forecast.WindShear == nil {
		forecast.WindShear = ctx.WindShear
	}
	if forecast.Cavok == nil {
		forecast.Cavok = ctx.Cavok
	}
	if forecast.Remark == nil && ctx.Remark != nil {
		forecast.Remark = ctx.Remark
	}
	if len(forecast.Remarks) == 0 && len(ctx.Remarks) > 0 {
		forecast.Remarks = ctx.Remarks
	}
	if len(forecast.Clouds) == 0 && len(ctx.Clouds) > 0 {
		forecast.Clouds = ctx.Clouds
	}
	if len(forecast.WeatherConditions) == 0 && len(ctx.WeatherConditions) > 0 {
		forecast.WeatherConditions = ctx.WeatherConditions
	}

	return forecast
}

// CompositeForecast holds prevailing and supplemental forecasts for a specific timestamp.
type CompositeForecast struct {
	Supplemental []Forecast `json:"supplemental"`
	Prevailing   Forecast   `json:"prevailing"`
}

// GetCompositeForecastForDate finds the prevailing and supplemental forecasts for a given timestamp.
func GetCompositeForecastForDate(date time.Time, fc *ForecastContainer) (*CompositeForecast, error) {
	ts := date.UnixMilli()

	if ts < fc.Start || ts >= fc.End {
		return nil, NewTimestampOutOfBoundsError("Provided timestamp is outside the report validity period")
	}

	var prevailing *Forecast
	var supplemental []Forecast

	for _, f := range fc.Forecast {
		if isImplicitEnd(f) && f.Start <= ts {
			prevailing = &f
		}

		if !isImplicitEnd(f) && f.End > ts && f.Start <= ts {
			supplemental = append(supplemental, f)
		}
	}

	if prevailing == nil {
		return nil, NewUnexpectedParseError("Unable to find trend for date")
	}

	return &CompositeForecast{
		Prevailing:   *prevailing,
		Supplemental: supplemental,
	}, nil
}
