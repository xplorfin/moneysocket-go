package wad

import (
	"fmt"
	"time"
)

type Rate struct {
	BaseCode  string
	QuoteCode string
	RateValue float64
	Timestamp time.Time
}

func NewRate(baseCode string, quoteCode string, rate float64) Rate {
	return Rate{
		BaseCode:  baseCode,
		QuoteCode: quoteCode,
		RateValue: rate,
		Timestamp: time.Now(),
	}
}

func (r Rate) ToString() string {
	return fmt.Sprintf("%.8f %s%s", r.RateValue, r.BaseCode, r.QuoteCode)
}

// Includes determines whether or not the Rate.BaseCode is Rate.BaseCode or the Rate.QuoteCode
func (r Rate) Includes(quote string) bool {
	return r.BaseCode == quote || r.QuoteCode == quote
}

// Invert inverts the wad by flipping the currency conversion
func (r Rate) Invert() Rate {
	return Rate{
		BaseCode:  r.QuoteCode,
		QuoteCode: r.BaseCode,
		RateValue: 1.0 / r.RateValue,
		Timestamp: r.Timestamp,
	}
}

// Other determines if base get quote if quote get base
func (r Rate) Other(code string) string {
	if !r.Includes(code) {
		panic(fmt.Errorf("Code must be %s or %s to use other", r.QuoteCode, r.BaseCode))
	}
	if r.BaseCode == code {
		return r.QuoteCode
	}
	return r.BaseCode
}

// Convert converts a value/code into a new currency
func (r Rate) Convert(value float64, valueCode string) (converted float64, code string) {
	if valueCode == r.BaseCode {
		return value * r.RateValue, r.QuoteCode
	} else if valueCode == r.QuoteCode {
		return value / r.RateValue, r.BaseCode
	}
	panic(fmt.Errorf("expected Code to be either %s or %s", r.QuoteCode, r.BaseCode))
}

// DeriveRate does a cross-currency rate conversion
func DeriveRate(baseQuote string, quoteCode string, rates [2]Rate) Rate {
	if !(rates[0].Includes(baseQuote) || rates[1].Includes(baseQuote)) {
		panic(fmt.Errorf("at least one rate must use baseQuote %s", baseQuote))
	}
	if !(rates[0].Includes(quoteCode) || rates[1].Includes(quoteCode)) {
		panic(fmt.Errorf("at least one rate must use quoteCode %s", quoteCode))
	}

	var first, second Rate
	if rates[0].Includes(baseQuote) {
		first = rates[0]
		second = rates[1]
	} else {
		first = rates[1]
		second = rates[0]
	}
	otherCode := first.Other(baseQuote)

	if !second.Includes(otherCode) {
		panic(fmt.Errorf("expected second to include %s", otherCode))
	}

	otherConverted, otherCheck := first.Convert(1.0, baseQuote)
	if otherCheck != otherCode {
		panic(fmt.Errorf("expected %s to equal %s", otherCode, otherCheck))
	}

	quoteConverted, quoteCheck := second.Convert(otherConverted, otherCode)

	if quoteCheck != quoteCode {
		panic(fmt.Errorf("expected %s to equal %s", quoteCheck, quoteCode))
	}

	timestamp := minTime(first.Timestamp, second.Timestamp)
	return Rate{
		BaseCode:  baseQuote,
		QuoteCode: quoteCode,
		RateValue: quoteConverted,
		Timestamp: timestamp,
	}
}

// minTime is a helper function that get the lesser of two times
func minTime(time1 time.Time, time2 time.Time) time.Time {
	if time1.Before(time2) {
		return time1
	}
	return time2
}
