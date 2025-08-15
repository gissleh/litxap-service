module github.com/gissleh/litxap-service

go 1.22.2

require (
	github.com/fwew/fwew-lib/v5 v5.25.1
	github.com/gissleh/litxap v1.7.0
	github.com/stretchr/testify v1.10.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gissleh/litxap-fwew v0.0.0-20250815125607-8b87ab613a4b // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//for testing edits to dependencies, make sure you go get after you comment them back.
// replace github.com/gissleh/litxap => ../litxap
// replace github.com/gissleh/litxaputil => ../litxap/litxaputil
// replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
