package wad

import (
	"fmt"
	"strconv"

	"github.com/btcsuite/btcutil"
)

// Wad is a unit of currencies
type Wad struct {
	// MSats is the millisatoshi account in the wad
	MSats float64 `json:"MSats"`
	// AssetStable is wether or not the asset is a stable coin (usd, etc)
	AssetStable bool `json:"asset_stable"`
	// AssetUnit is the conversion rate of the asset
	AssetUnit float64 `json:"asset_unit"`
	// currency is the Currency of the Wad
	currency Currency
}

const (
	// MSatsPerSat is a constant amount of satoshis
	MSatsPerSat = 1000
	// SatsPerBtc is the number of satoshi's in a bitcoin
	SatsPerBtc = btcutil.SatoshiPerBitcoin
	// MSatsPerBtc is the number of satoshis in a bitcoin
	MSatsPerBtc = MSatsPerSat * SatsPerBtc
	// AllCountries is the constant for all nations (international currency)
	AllCountries = "All"
)

// NewWad creates a wad
func NewWad(msats float64, assetStable bool, assetUnit float64, code string) Wad {
	if msats < 0 {
		panic(fmt.Errorf("MSats must be positive, got %f", msats))
	}
	wad := Wad{
		MSats:       msats,
		AssetStable: assetStable,
		AssetUnit:   assetUnit,
		currency: Currency{
			Code:      code,
			Decimals:  -1,
			Symbol:    "₿",
			IsoNum:    -1,
			Name:      "Bitcoin",
			Countries: AllCountries,
		},
	}

	if !wad.AssetStable {
		wad.AssetUnit = msats / MSatsPerBtc
	}

	// if code in fiat, update the currency to reflect fiat params
	if currency, ok := Fiat[code]; ok {
		wad.UpdateCurrency(currency)
	}

	// TODO update params here
	if currency, ok := CryptoCurrency[code]; ok {
		wad.UpdateCurrency(currency)
	}

	return wad
}

// BitcoinWad creates a bitcoin wad
func BitcoinWad(msats float64) Wad {
	return NewWad(msats, false, msats, "BTC")
}

// formatFloat formats a float to a precision
func formatFloat(num float64, prc int) string {
	return fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
}

// FmtShort formats a Wad
func (w Wad) FmtShort() string {
	if !w.AssetStable {
		return fmt.Sprintf("%s %.3f sat", w.currency.Symbol, w.MSats/MSatsPerSat)
	}
	symb := fmt.Sprintf("%s ", w.currency.Symbol)
	var asset string
	if w.currency.Decimals >= 0 {
		asset = formatFloat(w.AssetUnit, w.currency.Decimals)
	} else {
		asset = fmt.Sprintf("%f", w.AssetUnit)
	}
	return fmt.Sprintf("%s%s %s", symb, asset, w.currency.Code)
}

// FmtLong formats a wad
func (w Wad) FmtLong() string {
	if !w.AssetStable {
		return w.FmtShort()
	}
	return fmt.Sprintf("%s (%.3f sat)", w.FmtShort(), w.MSats/MSatsPerSat)
}

// UpdateCurrency changes the Currency of a wad
func (w *Wad) UpdateCurrency(currency Currency) {
	w.currency = currency
}

// UsdWad is a wad denominated in USD
func UsdWad(usd float64, rate Rate) Wad {
	btc, code := rate.Convert(usd, "USD")
	if code != "BTC" {
		panic(fmt.Errorf("expected Code to be BTC, got: %s", code))
	}
	msats := btc * MSatsPerBtc
	return NewWad(msats, true, usd, "USD")
}

// CadWad is a wad denominated in CAD
func CadWad(cad float64, rate Rate) Wad {
	btc, code := rate.Convert(cad, "CAD")
	if code != "BTC" {
		panic(fmt.Errorf("expected Code to be BTC, got: %s", code))
	}
	msats := btc * MSatsPerBtc
	return NewWad(msats, true, cad, "CAD")
}

// CustomWad creates a custom wad at a given Rate
func CustomWad(units float64, rate Rate, code string, countries string, decimals int, name string, symbol string) Wad {
	btc, btcCode := rate.Convert(units, code)
	if btcCode != "BTC" {
		panic(fmt.Errorf("expected Code to be BTC, got: %s", code))
	}
	msats := btc * MSatsPerBtc
	return Wad{
		MSats:       msats,
		AssetStable: true,
		AssetUnit:   units,
		currency: Currency{
			Code:      code,
			Decimals:  decimals,
			Name:      name,
			Symbol:    symbol,
			Countries: countries,
		},
	}
}
