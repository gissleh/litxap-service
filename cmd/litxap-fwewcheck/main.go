package main

import (
	"log"

	litxapfwew "github.com/gissleh/litxap-fwew"
)

func main() {
	entries, err := litxapfwew.Global().LookupEntries("tìfmetok")
	if err != nil {
		log.Fatalln("Failed to lookup tìfmetok:", err)
		return
	}

	log.Println("Tìfmetok gave", len(entries), "entries")
}
