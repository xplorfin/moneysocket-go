package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/moneysocket/wad"
	"github.com/xplorfin/ozzo-validators/rules"
)

// represents the data store for a single account
type AccountDb struct {
	Details Account
	config  *config.Config
	// map of beacon string to most recent error when trying to connect
	ConnectionAttempts map[string]error
}

func GetPersistedAccounts(config *config.Config) (accts []AccountDb) {
	persistedDbs := GetAccountDbs(config)
	accts = append(accts, persistedDbs...)
	return accts
}

func NewAccountDb(accountName string, config *config.Config) (adb AccountDb) {
	adb = AccountDb{
		ConnectionAttempts: make(map[string]error),
		config:             config,
		Details: Account{
			AccountName: accountName,
			AccountUuid: uuid.New().String(),
			Wad:         wad.BitcoinWad(0),
		},
	}
	err := adb.makeDbIfNotExists()
	if err != nil {
		panic(err)
	}
	adb.Details, err = adb.readDetails()
	if err != nil {
		panic(err)
	}

	// use the account db

	return adb
}

// get a list of accoutdbs from our config file
// in python this is an iter, but I can only assume we're not using this
// for a large number of accounts (hopefully)
func GetAccountDbs(configuration *config.Config) (adList []AccountDb) {
	fmt.Println(configuration.GetAccountPersistDir())
	err := filepath.Walk(configuration.GetAccountPersistDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			adb, err := AdbFromFile(path)
			if err != nil {
				return err
			}
			adList = append(adList, AccountDb{Details: adb, config: configuration})
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	return adList
}

// read an account db out of a file
func AdbFromFile(filename string) (adb Account, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return adb, err
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return adb, err
	}

	err = json.Unmarshal(byteValue, &adb)
	if err != nil {
		return adb, err
	}
	return adb, nil
}

// read account details from file
func (a AccountDb) readDetails() (adb Account, err error) {
	return AdbFromFile(a.filename())
}

func (a AccountDb) Persist() error {
	file, err := os.Create(a.filename())
	if err != nil {
		return err
	}
	defer file.Close()

	accountDetails, err := json.Marshal(&a.Details)
	if err != nil {
		return err
	}
	_, err = file.Write(accountDetails)
	return err
}

// create the account db file if it doesn't already exist
func (a AccountDb) makeDbIfNotExists() error {
	if rules.IsFile(a.filename()) {
		return nil
	}
	return a.Persist()
}

func (a AccountDb) Depersist() error {
	return os.Remove(a.filename())
}

func (a AccountDb) filename() string {
	return fmt.Sprintf("%s/%s.json", a.config.GetAccountPersistDir(), a.Details.AccountName)
}

// add a connection attempt for the current account
func (a *AccountDb) AddConnectionAttempt(attemptedBeacon beacon.Beacon, err error) {
	a.ConnectionAttempts[attemptedBeacon.ToBech32Str()] = err
}

// get array of beacons which have disconnected
func (a *AccountDb) GetDisconnectedBeacons() (beacons []beacon.Beacon) {
	for _, detailBeacon := range a.Details.Beacons {
		beaconStr := detailBeacon.ToBech32Str()
		if val, ok := a.ConnectionAttempts[beaconStr]; ok {
			if val != nil {
				beacons = append(beacons, detailBeacon)
			}
		} else {
			continue
		}
	}
	return beacons
}

func (a AccountDb) GetSummaryString(locations []location.Location) (summaryStr string) {
	summaryStr += fmt.Sprintf("\n%s: wad: %s\n", a.Details.AccountName, a.Details.Wad.FmtLong())
	for _, detailBeacon := range a.Details.Beacons {
		beaconStr := detailBeacon.ToBech32Str()
		summaryStr += fmt.Sprintf("\n\t\toutgoing beacon: %s", detailBeacon.ToBech32Str())
		if val, ok := a.ConnectionAttempts[beaconStr]; ok {
			summaryStr += fmt.Sprintf("\n\t\t\tconnection attempts: %s", val)
		} else {
			summaryStr += "\n\t\t\tconnection attempts: none"
		}
	}

	for _, sharedSeed := range a.Details.SharedSeeds {
		seedBeacon := beacon.NewBeaconFromSeed(sharedSeed)
		for _, loc := range locations {
			seedBeacon.AddLocation(loc)
		}
		summaryStr += fmt.Sprintf("\n\t\tincoming shared seed %s", sharedSeed.ToString())
		summaryStr += fmt.Sprintf("\n\t\t\t incoming beacon: %s", seedBeacon.ToBech32Str())
	}
	return summaryStr
}
