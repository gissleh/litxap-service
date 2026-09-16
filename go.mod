module github.com/gissleh/litxap-service

go 1.25.0

require (
	github.com/gissleh/litxap v1.17.1
	github.com/gissleh/litxap-fwew v1.17.1-fwew5.28.1
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/fwew/fwew-lib/v5 v5.28.1 // indirect
	github.com/go-sql-driver/mysql v1.10.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
)

//for testing edits to dependencies, make sure you go get after you comment them back.
// replace github.com/gissleh/litxap => ../litxap
// replace github.com/gissleh/litxaputil => ../litxap/litxaputil
// replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
