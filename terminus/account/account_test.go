package account

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/filet"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

func TestGetPersistedAccounts(t *testing.T) {
	const testIteratons = 5
	configuration := config.NewConfig()
	configuration.AccountPersistDir = filet.TmpDir(t, "")

	testAccounts := make(map[string]DB)
	for i := 0; i < testIteratons; i++ {
		accountName := gofakeit.BeerAlcohol()
		adb := NewAccountDb(accountName, configuration)
		testAccounts[accountName] = adb
	}

	dbs := GetPersistedAccounts(configuration)
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
