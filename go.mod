module github.com/gissleh/litxap-service

go 1.22.2

require (
	github.com/fwew/fwew-lib/v5 v5.24.1
	github.com/gissleh/litxap v1.4.7
	github.com/stretchr/testify v1.10.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//for testing edits to dependencies, make sure you go get after you comment them back.
//replace github.com/gissleh/litxap/litxap => ../litxap
//replace github.com/gissleh/litxap/litxaputil => ../litxap/litxaputil
//replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
