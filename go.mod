module github.com/scatkit/jupswap

go 1.23.1

require (
	github.com/gagliardetto/binary v0.8.0
	github.com/shopspring/decimal v1.4.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/google/uuid v1.5.0 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/oapi-codegen/runtime v1.1.1
	github.com/scatkit/pumpdexer v0.0.0-20250101140745-b2f8fd8ca090 //inderct
	github.com/streamingfast/logging v0.0.0-20230608130331-f22c91403091 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
)

//replace github.com/scatkit/pumpdexer [v0.0.0-20250101140745-b2f8fd8ca090] => ../pumpdexer
replace github.com/scatkit/pumpdexer v0.0.0-20250101140745-b2f8fd8ca090 => ../pumpdexer
