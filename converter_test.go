package metartafparser

import (
	"errors"
	"strconv"
	"testing"
)

func TestParseFractionFraction(t *testing.T) {
	got, err := parseFraction("1/2")
	if err != nil {
		t.Fatalf("parseFraction(1/2): unexpected error: %v", err)
	}
	assertEqual(t, got, 0.5, "parseFraction(1/2)")
}

func TestParseFractionWholeNumber(t *testing.T) {
	got, err := parseFraction("3")
	if err != nil {
		t.Fatalf("parseFraction(3): unexpected error: %v", err)
	}
	assertEqual(t, got, 3.0, "parseFraction(3)")
}

func TestParseFractionZeroDenominator(t *testing.T) {
	got, err := parseFraction("1/0")
	if err != nil {
		t.Fatalf("parseFraction(1/0): unexpected error: %v", err)
	}
	assertEqual(t, got, 0.0, "parseFraction(1/0)")
}

func TestParseFractionOverflow(t *testing.T) {
	// A large but finite numerator must not overflow to +Inf when scaled by 100.
	_, err := parseFraction("1e308/1")
	if !errors.Is(err, strconv.ErrRange) {
		t.Fatalf("parseFraction(1e308/1): got %v, want %v", err, strconv.ErrRange)
	}
}

func TestConvertFractionalAmountWholePlusFractionOverflow(t *testing.T) {
	_, err := convertFractionalAmount("1.7e308 1e308")
	if !errors.Is(err, strconv.ErrRange) {
		t.Fatalf("convertFractionalAmount: got %v, want %v", err, strconv.ErrRange)
	}
}
