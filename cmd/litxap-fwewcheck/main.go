package main

import (
	"github.com/gissleh/litxap-service/adapters/fwewdict"
	"log"
)

func main() {
	entries, err := fwewdict.Global().LookupEntries("tìfmetok")
	if err != nil {
		log.Fatalln("Failed to lookup tìfmetok:", err)
		return
	}

	log.Println("Tìfmetok gave", len(entries), "entries")
}
