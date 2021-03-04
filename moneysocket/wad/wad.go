package wad

import (
	"fmt"
	"strconv"

	"github.com/btcsuite/btcutil"
)

// todo add currency in here
type Wad struct {
	Msats       float64 `json:"Msats"`
	AssetStable bool    `json:"asset_stable"`
	AssetUnit   float64 `json:"asset_unit"`
	currency    Currency
}

const (
	MsatsPerSat  = 1000
	SatsPerBtc   = btcutil.SatoshiPerBitcoin
	MsatsPerBtc  = MsatsPerSat * SatsPerBtc
	AllCountries = "All"
)

func NewWad(msats float64, assetStable bool, assetUnit float64, code string) Wad {
	if msats < 0 {
		panic(fmt.Errorf("Msats must be positive, got %f", msats))
	}
	wad := Wad{
		Msats:       msats,
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
		wad.AssetUnit = msats / MsatsPerBtc
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

func BitcoinWad(msats float64) Wad {
	return NewWad(msats, false, msats, "BTC")
}

func formatFloat(num float64, prc int) string {
	return fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
}

func (w Wad) FmtShort() string {
	if !w.AssetStable {
		return fmt.Sprintf("%s %.3f sat", w.currency.Symbol, w.Msats/MsatsPerSat)
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

func (w Wad) FmtLong() string {
	if !w.AssetStable {
		return w.FmtShort()
	}
	return fmt.Sprintf("%s (%.3f sat)", w.FmtShort(), w.Msats/MsatsPerSat)
}

func (w *Wad) UpdateCurrency(currency Currency) {
	w.currency = currency
}

func UsdWad(usd float64, rate Rate) Wad {
	btc, code := rate.Convert(usd, "USD")
	if code != "BTC" {
		panic(fmt.Errorf("expected Code to be BTC, got: %s", code))
	}
	msats := btc * MsatsPerBtc
	return NewWad(msats, true, usd, "USD")
}

func CadWad(cad float64, rate Rate) Wad {
	btc, code := rate.Convert(cad, "CAD")
	if code != "BTC" {
		panic(fmt.Errorf("expected Code to be BTC, got: %s", code))
	}
	msats := btc * MsatsPerBtc
	return NewWad(msats, true, cad, "CAD")
}

func CustomWad(units float64, rate Rate, code string, countries string, decimals int, name string, symbol string) Wad {
	btc, btcCode := rate.Convert(units, code)
	if btcCode != "BTC" {
		panic(fmt.Errorf("expected Code to be BTC, got: %s", code))
	}
	msats := btc * MsatsPerBtc
	return Wad{
		Msats:       msats,
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
