.PHONY: build test lint vet fuzz

build:
	go build ./...

test:
	go test ./...

fuzz-shortest:
	go test -fuzz=^FuzzParseMetar$$ -fuzztime=3s -count=1 .
	go test -fuzz=^FuzzParseTAF$$ -fuzztime=3s -count=1 .
	go test -fuzz=^FuzzParseMetarDated$$ -fuzztime=3s -count=1 .
	go test -fuzz=^FuzzParseTAFDated$$ -fuzztime=3s -count=1 .
	go test -fuzz=^FuzzParseTAFAsForecast$$ -fuzztime=3s -count=1 .
	go test -fuzz=^FuzzGetCompositeForecastForDate$$ -fuzztime=3s -count=1 .

fuzz-shorter:
	go test -fuzz=^FuzzParseMetar$$ -fuzztime=10s -count=1 .
	go test -fuzz=^FuzzParseTAF$$ -fuzztime=10s -count=1 .
	go test -fuzz=^FuzzParseMetarDated$$ -fuzztime=10s -count=1 .
	go test -fuzz=^FuzzParseTAFDated$$ -fuzztime=10s -count=1 .
	go test -fuzz=^FuzzParseTAFAsForecast$$ -fuzztime=10s -count=1 .
	go test -fuzz=^FuzzGetCompositeForecastForDate$$ -fuzztime=10s -count=1 .

fuzz-short:
	go test -fuzz=^FuzzParseMetar$$ -fuzztime=30s -count=1 .
	go test -fuzz=^FuzzParseTAF$$ -fuzztime=30s -count=1 .
	go test -fuzz=^FuzzParseMetarDated$$ -fuzztime=30s -count=1 .
	go test -fuzz=^FuzzParseTAFDated$$ -fuzztime=30s -count=1 .
	go test -fuzz=^FuzzParseTAFAsForecast$$ -fuzztime=30s -count=1 .
	go test -fuzz=^FuzzGetCompositeForecastForDate$$ -fuzztime=30s -count=1 .

fuzz-long:
	go test -fuzz=^FuzzParseMetar$$ -fuzztime=10m -count=1 .
	go test -fuzz=^FuzzParseTAF$$ -fuzztime=10m -count=1 .
	go test -fuzz=^FuzzParseMetarDated$$ -fuzztime=10m -count=1 .
	go test -fuzz=^FuzzParseTAFDated$$ -fuzztime=10m -count=1 .
	go test -fuzz=^FuzzParseTAFAsForecast$$ -fuzztime=10m -count=1 .
	go test -fuzz=^FuzzGetCompositeForecastForDate$$ -fuzztime=10m -count=1 .

fuzz-longer:
	go test -fuzz=^FuzzParseMetar$$ -fuzztime=1h -count=1 .
	go test -fuzz=^FuzzParseTAF$$ -fuzztime=1h -count=1 .
	go test -fuzz=^FuzzParseMetarDated$$ -fuzztime=1h -count=1 .
	go test -fuzz=^FuzzParseTAFDated$$ -fuzztime=1h -count=1 .
	go test -fuzz=^FuzzParseTAFAsForecast$$ -fuzztime=1h -count=1 .
	go test -fuzz=^FuzzGetCompositeForecastForDate$$ -fuzztime=1h -count=1 .

fuzz-longest:
	go test -fuzz=^FuzzParseMetar$$ -fuzztime=8h -count=1 .
	go test -fuzz=^FuzzParseTAF$$ -fuzztime=8h -count=1 .
	go test -fuzz=^FuzzParseMetarDated$$ -fuzztime=8h -count=1 .
	go test -fuzz=^FuzzParseTAFDated$$ -fuzztime=8h -count=1 .
	go test -fuzz=^FuzzParseTAFAsForecast$$ -fuzztime=8h -count=1 .
	go test -fuzz=^FuzzGetCompositeForecastForDate$$ -fuzztime=8h -count=1 .

lint:
	golangci-lint run

vet:
	go vet ./...
