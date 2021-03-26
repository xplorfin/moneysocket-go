package terminus

import (
	"strconv"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
	"github.com/xplorfin/moneysocket-go/terminus/account"
)

type Directory struct {
	config                *config.Config
	AccountBySharedSeed   map[string]account.DB
	SharedSeedsByAccount  map[string][]beacon.SharedSeed
	Accounts              map[string]account.DB
	AccountsByPaymentHash map[string]account.DB
}

func NewTerminusDirectory(config *config.Config) *Directory {
	return &Directory{
		config:                config,
		AccountBySharedSeed:   make(map[string]account.DB),
		SharedSeedsByAccount:  make(map[string][]beacon.SharedSeed),
		Accounts:              make(map[string]account.DB),
		AccountsByPaymentHash: make(map[string]account.DB),
	}
}

// python version is an iterator
func (t *Directory) GetAccounts() (accounts []account.DB) {
	for _, v := range t.Accounts {
		accounts = append(accounts, v)
	}
	return accounts
}

func (t *Directory) GetAccountList() []account.DB {
	return t.GetAccounts()
}

// generate an account name from an autoincrementing int
func (t *Directory) GenerateAccountName() string {
	for i := 0; i < 1000; i++ {
		acct := t.LookupByName(strconv.Itoa(i))
		if acct == nil {
			return strconv.Itoa(i)
		}
	}
	panic("more than 1,000 accounts exist")
}

// get list of acount names
func (t *Directory) GetAccountNameSet() (accounts []string) {
	for _, account := range t.Accounts {
		accounts = append(accounts, account.Details.AccountName)
	}
	return accounts
}

func (t *Directory) LookupByName(name string) *account.DB {
	if val, ok := t.Accounts[name]; ok {
		return &val
	}
	return nil
}

func (t *Directory) LookupBySeed(seed beacon.SharedSeed) account.DB {
	return t.AccountBySharedSeed[seed.ToString()]
}

func (t *Directory) LookupByPaymentHash(hash string) {
	panic("method not yet implemented")
}

func (t *Directory) ReindexAccount(acct account.DB) {
	t.AddAccount(acct)
}

func (t *Directory) AddAccount(acct account.DB) {
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
