package util

import (
	"github.com/btcsuite/btcutil/bech32"
)

// Bech32EncodeBytes encodes a bech32 string
func Bech32EncodeBytes(toEncode []byte, hrp string) (encodedBytes string, err error) {
	convertedBytes, err := bech32.ConvertBits(toEncode, 8, 5, true)
	if err != nil {
		panic(err)
	}
	encodedBytes, err = bech32.Encode(hrp, convertedBytes)
	if err != nil {
		panic(err)
	}
	return encodedBytes, nil
}

// Bech32DecodeBytes decodes a bech32 string
func Bech32DecodeBytes(rawString string) (hrp string, decodedKey []byte, err error) {
	hrp, decodedKey, err = decodeBech32(rawString)
	if err != nil {
		return "", nil, err
	}
	deconverted, err := bech32.ConvertBits(decodedKey, 5, 8, false)
	if err != nil {
		return "", nil, err
	}
	return hrp, deconverted, nil
}
