package terminus

import (
	"strconv"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

type TerminusDirectory struct {
	config                *config.Config
	AccountBySharedSeed   map[string]account.AccountDb
	SharedSeedsByAccount  map[string][]beacon.SharedSeed
	Accounts              map[string]account.AccountDb
	AccountsByPaymentHash map[string]account.AccountDb
}

func NewTerminusDirectory(config *config.Config) *TerminusDirectory {
	return &TerminusDirectory{
		config:                config,
		AccountBySharedSeed:   make(map[string]account.AccountDb),
		SharedSeedsByAccount:  make(map[string][]beacon.SharedSeed),
		Accounts:              make(map[string]account.AccountDb),
		AccountsByPaymentHash: make(map[string]account.AccountDb),
	}
}

// python verison is an iterator
func (t *TerminusDirectory) GetAccounts() (accounts []account.AccountDb) {
	for _, v := range t.Accounts {
		accounts = append(accounts, v)
	}
	return accounts
}

func (t *TerminusDirectory) GetAccountList() []account.AccountDb {
	return t.GetAccounts()
}

// generate an account name from an autoincrementing int
func (t *TerminusDirectory) GenerateAccountName() string {
	for i := 0; i < 1000; i++ {
		acct := t.LookupByName(strconv.Itoa(i))
		if acct == nil {
			return strconv.Itoa(i)
		}
	}
	panic("more than 1,000 accounts exist")
}

// get list of acount names
func (t *TerminusDirectory) GetAccountNameSet() (accounts []string) {
	for _, account := range t.Accounts {
		accounts = append(accounts, account.Details.AccountName)
	}
	return accounts
}

func (t *TerminusDirectory) LookupByName(name string) *account.AccountDb {
	if val, ok := t.Accounts[name]; ok {
		return &val
	}
	return nil
}

func (t *TerminusDirectory) LookupBySeed(seed beacon.SharedSeed) account.AccountDb {
	return t.AccountBySharedSeed[seed.ToString()]
}

func (t *TerminusDirectory) LookupByPaymentHash(hash string) {
	panic("method not yet implemented")
}

func (t *TerminusDirectory) ReindexAccount(acct account.AccountDb) {
	t.AddAccount(acct)
}

func (t *TerminusDirectory) AddAccount(acct account.AccountDb) {
	details := acct.Details
	acct.ConnectionAttempts = make(map[string]error)
	t.Accounts[details.AccountName] = acct
	sharedSeeds := details.SharedSeeds
	for _, sharedSeed := range sharedSeeds {
		if _, ok := t.SharedSeedsByAccount[details.AccountName]; !ok {
			t.SharedSeedsByAccount[details.AccountName] = []beacon.SharedSeed{}
		}
		t.SharedSeedsByAccount[details.AccountName] = append(t.SharedSeedsByAccount[details.AccountName], sharedSeed)
		t.AccountBySharedSeed[sharedSeed.ToString()] = acct
	}
	//for paymentHash, _ := range details.getPending(){
	//	// todo
	//}
}
