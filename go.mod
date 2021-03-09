module github.com/xplorfin/moneysocket-go

go 1.15

require (
	github.com/Flaque/filet v0.0.0-20201012163910-45f684403088
	github.com/brianvoe/gofakeit/v5 v5.11.2
	github.com/btcsuite/btcutil v1.0.2
	github.com/buger/jsonparser v1.1.1
	github.com/dustin/go-humanize v1.0.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/kr/pretty v0.2.1
	github.com/kr/text v0.2.0 // indirect
	github.com/lightningnetwork/lnd v0.11.1-beta.rc5
	github.com/magefile/mage v1.11.0 // indirect
	github.com/mergermarket/go-pkcs7 v0.0.0-20170926155232-153b18ea13c9
	github.com/mvo5/goconfigparser v0.0.0-20201015074339-50f22f44deb5
	github.com/posener/wstest v1.2.0
	github.com/prometheus/common v0.18.0
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	github.com/xplorfin/lndmock v0.4.0
	github.com/xplorfin/netutils v0.23.0
	github.com/xplorfin/ozzo-validators v0.23.0
	github.com/xplorfin/tlsutils v0.14.0
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210304124612-50617c2ba197 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/btcsuite/btcd v0.21.0-beta => github.com/xplorfin/btcd v0.21.0-hotfix
