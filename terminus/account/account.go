package account

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/wad"
)

// Account represents a terminus account.
type Account struct {
	// AccountName is the name of the account
	AccountName string `json:"account_name"`
	// AccountUUID is the uuid for the account
	AccountUUID string `json:"account_uuid"`
	// Wad is the wad used by the account
	Wad wad.Wad `json:"wad"`
	// SharedSeeds is a list of shared seeds nicluded with the account
	SharedSeeds []beacon.SharedSeed `json:"shared_seeds"`
	// Beacons is a list of beacons associated with the account
	Beacons []beacon.Beacon `json:"beacons"`
}

// AddBeacon adds a beacon to the account.
func (a *Account) AddBeacon(newBeacon beacon.Beacon) {
	a.Beacons = append(a.Beacons, newBeacon)
}

// Payer is whether or not account is payer, always true in terminus.
func (a *Account) Payer() bool {
	return true
}

// Payee is whether or not account is payee, always true in terminus.
func (a *Account) Payee() bool {
	return true
}

// Ready determines if an account is ready for use. Right now this always
// returns true since the account's only use files right now.
func (a *Account) Ready() bool {
	return true
}

// RemoveBeacon removes a beacon from the Account.
func (a *Account) RemoveBeacon(toRemove beacon.Beacon) {
	for index, tmpBeacon := range a.Beacons {
		if toRemove.ToBech32Str() == tmpBeacon.ToBech32Str() {
			a.Beacons = append(a.Beacons[:index], a.Beacons[index+1:]...)
		}
	}
}

// AddSharedSeed adds a shared seed from the account.
func (a *Account) AddSharedSeed(ss beacon.SharedSeed) {
	a.SharedSeeds = append(a.SharedSeeds, ss)
}

// RemoveSharedSeed removes a shared seed from the account.
func (a *Account) RemoveSharedSeed(ss beacon.SharedSeed) {
	for index, tmpSeed := range a.SharedSeeds {
		if ss.Equal(tmpSeed) {
			a.SharedSeeds = append(a.SharedSeeds[:index], a.SharedSeeds[index+1:]...)
		}
	}
}

// UnmarshalJSON decodes an account from a json payload.
func (a *Account) UnmarshalJSON(b []byte) (err error) {
	type AccountJSON struct {
		AccountName string  `json:"account_name"`
		AccountUUID string  `json:"account_uuid"`
		Wad         wad.Wad `json:"wad"`
		// pending TODO
		SharedSeeds []string `json:"shared_seeds"`
		Beacons     []string `json:"beacons"`
	}
	var aj AccountJSON
	err = json.Unmarshal(b, &aj)
	if err != nil {
		return err
	}
	a.AccountName = aj.AccountName
	a.AccountUUID = aj.AccountUUID
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

// MarshalJSON encodes an Account as a json payload.
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
