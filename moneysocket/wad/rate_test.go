package wad

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

const (
	BASE  = "BASE"
	QUOTE = "QUOTE"
)

func TestRateParity(t *testing.T) {
	rateLtcbtc := NewRate("LTC", "BTC", 0.00420700)

	rateBtcUsd := NewRate("BTC", "USD", 10723.12)
	Equal(t, rateBtcUsd.ToString(), "10723.12000000 BTCUSD")

	rateUsdBtc := rateBtcUsd.Invert()
	Equal(t, rateUsdBtc.ToString(), "0.00009326 USDBTC")

	rateBtcCad := NewRate("BTC", "CAD", 14011.28)
	Equal(t, rateBtcCad.ToString(), "14011.28000000 BTCCAD")

	rateCadBtc := rateBtcCad.Invert()
	Equal(t, rateCadBtc.ToString(), "0.00007137 CADBTC")

	rateEggplantCad := NewRate("EGGPLANT", "CAD", 2.49)
	Equal(t, rateEggplantCad.ToString(), "2.49000000 EGGPLANTCAD")

	rateCadEggplant := rateEggplantCad.Invert()
	Equal(t, rateCadEggplant.ToString(), "0.40160643 CADEGGPLANT")

	// derived 1
	rateBtcEggplant := DeriveRate("BTC", "EGGPLANT", [2]Rate{rateEggplantCad, rateBtcCad})
	Equal(t, rateBtcEggplant.ToString(), "5627.02008032 BTCEGGPLANT")

	// derived 2
	rateBtcEggplant = DeriveRate("BTC", "EGGPLANT", [2]Rate{rateBtcCad, rateEggplantCad})
	Equal(t, rateBtcEggplant.ToString(), "5627.02008032 BTCEGGPLANT")

	// derived 3
	rateBtcEggplant = DeriveRate("BTC", "EGGPLANT", [2]Rate{rateBtcCad, rateCadEggplant})
	Equal(t, rateBtcEggplant.ToString(), "5627.02008032 BTCEGGPLANT")

	// derived 4
	rateBtcEggplant = DeriveRate("BTC", "EGGPLANT", [2]Rate{rateBtcCad, rateCadEggplant})
	Equal(t, rateBtcEggplant.ToString(), "5627.02008032 BTCEGGPLANT")

	// derived ltccad
	rateLtcCad := DeriveRate("LTC", "CAD", [2]Rate{rateLtcbtc, rateBtcCad})
	Equal(t, rateLtcCad.ToString(), "58.94545496 LTCCAD")

	b := BitcoinWad(1234)
	u := UsdWad(12.99, rateBtcUsd)
	c := CadWad(30, rateBtcCad)
	e := CustomWad(3.14, rateBtcEggplant, "EGGPLANT", "All", 3, "Eggplant", "🍆")

	Equal(t, b.FmtShort(), "₿ 1.234 sat")
	Equal(t, b.FmtLong(), "₿ 1.234 sat")

	Equal(t, u.FmtShort(), "$ 12.99 USD")
	Equal(t, u.FmtLong(), "$ 12.99 USD (121140.116 sat)")

	Equal(t, c.FmtShort(), "$ 30.00 CAD")
	Equal(t, c.FmtLong(), "$ 30.00 CAD (214113.200 sat)")

	Equal(t, e.FmtShort(), "🍆 3.140 EGGPLANT")
	Equal(t, e.FmtLong(), "🍆 3.140 EGGPLANT (55802.182 sat)")
}

func TestQuoteIncludes(t *testing.T) {
	r := NewRate(BASE, QUOTE, 0.099912)
	True(t, r.Includes(BASE))
	True(t, r.Includes(QUOTE))
	False(t, r.Includes("other"))
}

func TestQuoteOtherPanic(t *testing.T) {
	r := NewRate(BASE, QUOTE, 0.099912)
	Equal(t, r.Other(BASE), QUOTE)
	Equal(t, r.Other(QUOTE), BASE)
	Panics(t, func() {
		r.Other("whatever")
	})
}
