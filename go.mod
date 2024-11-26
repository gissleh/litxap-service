module github.com/gissleh/litxap-service

go 1.22.2

require (
	github.com/fwew/fwew-lib/v5 v5.22.3-0.20241112155007-5f43c7781335
	github.com/stretchr/testify v1.9.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gissleh/litxap v1.0.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/gissleh/litxap => ../litxap
//replace github.com/gissleh/litxap/litxaputil => ../litxap/litxaputil
//for testing on a local machine's fwew-lib
//replace github.com/fwew/fwew-lib/v5 => ../fwew-lib
