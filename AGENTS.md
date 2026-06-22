# AGENTS.md — metar-taf-parser

Single-package Go library (`package metartafparser`) for parsing METAR and TAF aviation weather reports. CLI entrypoint at `cmd/metar-taf-parser`.

## Commands

| Command | What it does |
|---------|-------------|
| `make lint` | `golangci-lint run` (v2 config, `default: all` with ~30 disabled linters) |
| `make vet` | `go vet ./...` |
| `make test` | `go test ./...` |
| `make build` | `go build ./...` |
| `go build ./cmd/metar-taf-parser/` | Build CLI binary |
| `go test -run TestName` | Run a single test (stdlib `testing`, no testify) |
| `make fuzz-shortest` | 3s fuzz tests for all 6 fuzz targets |
| `make fuzz-shorter` | 10s fuzz tests |
| `make fuzz-short` | 30s fuzz tests |
| `make fuzz-long` | 10m fuzz tests |
| `make fuzz-longer` | 1h fuzz tests |
| `make fuzz-longest` | 8h fuzz tests |

Always run `make lint && make test` before commits.

## Documentation

- **`README.md`** and **godoc comments** (doc comments on all exported declarations) must be kept in sync with the codebase.
- When adding, renaming, or removing any exported type, function, const, or field:
  - Update or add its godoc comment in the source `.go` file.
  - Update `README.md` (Data Model, Parse Functions, Enums, or Error Handling sections as applicable).
  - Every exported declaration **must** have a godoc comment — run `go vet ./...` to catch missing ones.
- Review `README.md` for accuracy whenever making significant changes to the public API surface.

## Testing

- Uses stdlib `testing` package with custom generic assert helpers (`assertEqual`, `assertNil`, `assertNotNil`, `assertIntPtr`, `assertFloatPtr`, `assertBoolPtr`, etc.) — **do not import testify**.
- Test files: `*_test.go` at package level (not `_test` external package).
- Tests use `DefaultLocale()`, no external fixtures or services needed.
- All tests are fast unit tests with no hermetic/integration requirements.

## Architecture

- **`tokenizer.go`**: Tokenizer that splits raw reports into tokens.
- **`abstract_parser.go`**: Abstract parser for common METAR/TAF fields.
- **`metar_parser.go`**: METAR-specific parsing logic.
- **`taf_parser.go`**: TAF-specific parsing logic.
- **`commands.go`**: Command pattern — `commonCommander` (wind, visibility, clouds), `metarCommander` (altimeter, temperature, runway), `tafCommander` (icing, turbulence).
- **`metartafparser.go`**: Public API (`ParseMetar`, `ParseTAF`, `ParseMetarDated`, `ParseTAFDated`, `ParseTAFAsForecast`, `GetCompositeForecastForDate`).
- **`types.go`**: All struct types (`Metar`, `TAF`, `Forecast`, `Container`, `Remark`, etc.). `ParseOptions` has only one field: `Locale`.
- **`converter.go`**: Unit conversions (deg→cardinal, visibility, fractional amounts).
- **`locale.go`**: Nested map-based i18n with English default. `DefaultLocale()` returns full English translations.
- **`errors.go`**: Custom error types and sentinel errors.
- **`helpers.go`**: Shared helper utilities.
- **`enums.go`**: Enum type definitions (cloud cover, weather phenomena, etc.).
- **`remark_parser.go`**: Remark entry-point dispatch logic.
- **`remark_commands.go`**: Command pattern for remark sub-parsers.
- **`remark_*.go`**: Remark parsing split across files by category (weather, precip, snow, temp/pressure, visibility, wind).
- **`forecast_test.go`**, **`taf_test.go`**, **`remark_test.go`**, **`metartafparser_test.go`**, **`fuzz_test.go`**, **`cmd/metar-taf-parser/main_test.go`**: Test coverage.

## Key details

- **No external dependencies** — pure stdlib (Go 1.25.7). Module `github.com/ryansavara/metar-taf-parser`.
- Formatting: `gofmt` + `goimports` enforced by golangci-lint.
- CLI auto-detects METAR vs TAF, exits non-zero on parse failure, outputs indented JSON. Supports input via argument or stdin.
- Date hydration converts day/hour/minute timestamps to absolute Unix ms using `determineReportDate` (tries ±1 month, picks closest to issued date).
- BECMG forecasts inherit unset fields from previous context via `hydrateWithPreviousContextIfNeeded`.
