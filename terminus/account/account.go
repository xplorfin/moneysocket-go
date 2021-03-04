package account

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/wad"
)

// represents a terminus account
type Account struct {
	AccountName string  `json:"account_name"`
	AccountUuid string  `json:"account_uuid"`
	Wad         wad.Wad `json:"wad"`
	// pending TODO
	SharedSeeds []beacon.SharedSeed `json:"shared_seeds"`
	Beacons     []beacon.Beacon     `json:"beacons"`
}

func (a *Account) AddBeacon(newBeacon beacon.Beacon) {
	a.Beacons = append(a.Beacons, newBeacon)
}

// wether or not account is payer, always true in terminus
func (a *Account) Payer() bool {
	return true
}

// wether or not account is payee, always true in terminus
func (a *Account) Payee() bool {
	return true
}

func (a *Account) Ready() bool {
	return true
}

func (a *Account) RemoveBeacon(toRemove beacon.Beacon) {
	for index, tmpBeacon := range a.Beacons {
		if toRemove.ToBech32Str() == tmpBeacon.ToBech32Str() {
			a.Beacons = append(a.Beacons[:index], a.Beacons[index+1:]...)
		}
	}
}

func (a *Account) AddSharedSeed(ss beacon.SharedSeed) {
	a.SharedSeeds = append(a.SharedSeeds, ss)
}

func (a *Account) RemoveSharedSeed(ss beacon.SharedSeed) {
	for index, tmpSeed := range a.SharedSeeds {
		if ss.Equal(tmpSeed) {
			a.SharedSeeds = append(a.SharedSeeds[:index], a.SharedSeeds[index+1:]...)
		}
	}
}

func (a *Account) UnmarshalJSON(b []byte) (err error) {
	type AccountJson struct {
		AccountName string  `json:"account_name"`
		AccountUuid string  `json:"account_uuid"`
		Wad         wad.Wad `json:"wad"`
		// pending TODO
		SharedSeeds []string `json:"shared_seeds"`
		Beacons     []string `json:"beacons"`
	}
	var aj AccountJson
	err = json.Unmarshal(b, &aj)
	if err != nil {
		return err
	}
	a.AccountName = aj.AccountName
	a.AccountUuid = aj.AccountUuid
	a.Wad = aj.Wad
	for _, seed := range aj.SharedSeeds {
		ss, err := beacon.HexToSharedSeed(seed)
		if err != nil {
			return err
		}
		a.SharedSeeds = append(a.SharedSeeds, ss)
	}

	for _, b := range aj.Beacons {
		bcn, err := beacon.DecodeFromBech32Str(b)
		if err != nil {
			return err
		}
		a.Beacons = append(a.Beacons, bcn)
	}

	return err
}

func (a *Account) MarshalJSON() ([]byte, error) {
	type Alias Account
	// overwrite struct fields with converted counterparts
	var convertedBeacons []string
	for _, tmpBeacon := range a.Beacons {
		convertedBeacons = append(convertedBeacons, tmpBeacon.ToBech32Str())
	}
	var convertedSeeds []string
	for _, seed := range a.SharedSeeds {
		convertedSeeds = append(convertedSeeds, seed.ToString())
	}

	return json.Marshal(&struct {
		Beacons     []string `json:"beacons"`
		SharedSeeds []string `json:"shared_seeds"`
		*Alias
	}{
		Beacons:     convertedBeacons,
		SharedSeeds: convertedSeeds,
		Alias:       (*Alias)(a),
	})
}
