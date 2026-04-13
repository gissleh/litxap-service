module github.com/gissleh/litxap-service

go 1.24

toolchain go1.24.0

require (
	github.com/gissleh/litxap v1.9.0
	github.com/gissleh/litxap-fwew v1.9.0-fwew5.27.2
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/fwew/fwew-lib/v5 v5.27.2 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
)

//for testing edits to dependencies, make sure you go get after you comment them back.
// replace github.com/gissleh/litxap => ../litxap
// replace github.com/gissleh/litxaputil => ../litxap/litxaputil
// replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
