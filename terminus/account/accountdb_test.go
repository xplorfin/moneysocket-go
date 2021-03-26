package account

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/filet"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

func TestNewAccountDb(t *testing.T) {
	configuration := config.NewConfig()
	configuration.AccountPersistDir = filet.TmpDir(t, "")
	NewAccountDb("test", configuration)
}

func TestListAccountDbs(t *testing.T) {
	const testIteratons = 5
	configuration := config.NewConfig()
	configuration.AccountPersistDir = filet.TmpDir(t, "")

	testAccounts := make(map[string]DB)
	for i := 0; i < testIteratons; i++ {
		accountName := gofakeit.BeerAlcohol()
		adb := NewAccountDb(accountName, configuration)
		testAccounts[accountName] = adb
	}

	dbs := GetAccountDbs(configuration)
	for _, db := range dbs {
		if val, ok := testAccounts[db.Details.AccountName]; ok {
			Equal(t, val.Details.AccountName, db.Details.AccountName)
			Equal(t, val.Details.Wad, db.Details.Wad)
			Equal(t, val.Details.AccountUUID, db.Details.AccountUUID)
			Equal(t, val.Details.SharedSeeds, db.Details.SharedSeeds)
		} else {
			t.Errorf("expected iterated account %s to be in accounts created for testing", db.Details.AccountName)
		}
	}
}

func TestAdbMethods(t *testing.T) {
	configuration := config.NewConfig()
	configuration.AccountPersistDir = filet.TmpDir(t, "")
	accountName := gofakeit.BeerAlcohol()
	adb := NewAccountDb(accountName, configuration)
	// test beacon ops
	testBeacon := beacon.NewBeacon()
	Equal(t, len(adb.Details.Beacons), 0)
	adb.Details.AddBeacon(testBeacon)
	Equal(t, len(adb.Details.Beacons), 1)
	adb.Details.RemoveBeacon(testBeacon)
	Equal(t, len(adb.Details.Beacons), 0)
}
