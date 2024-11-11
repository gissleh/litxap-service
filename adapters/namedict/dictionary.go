package namedict

import (
	"fmt"
	"log"
	"strings"

	"github.com/gissleh/litxap"
	"github.com/gissleh/litxap-service/adapters/fwewdict"
	"github.com/gissleh/litxap/litxaputil"
)

type nameDict struct {
	table map[string][]string
}

func (n *nameDict) LookupEntries(word string) ([]litxap.Entry, error) {
	entryStrs, ok := n.table[word]
	if !ok {
		return nil, litxap.ErrEntryNotFound
	}

	entries := make([]litxap.Entry, len(entryStrs))
	for i, entryStr := range entryStrs {
		entries[i] = *litxap.ParseEntry(entryStr)
	}

	return entries, nil
}

func New(names ...string) litxap.Dictionary {
	table := make(map[string][]string)

	for _, name := range names {
		name := strings.Replace(name, "-", ".", -1)
		key := strings.Replace(strings.Replace(name, "*", "", -1), ".", "", -1)

		table[key] = append(table[key], fmt.Sprintf("%s: : Custom Name", name))
		for _, suffix := range suffixes {
			key := key + suffix
			table[key] = append(table[key], fmt.Sprintf("%s: -%s: Custom Name", name, suffix))
		}
	}

	return &nameDict{table: table}
}

func Doubles(names map[string]string) litxap.Dictionary {
	table := make(map[string][]string)

	for key, val := range names {
		//name := strings.Replace(val, "-", ".", -1)
		key := strings.Replace(strings.Replace(key, "*", "", -1), ".", "", -1)

		syllables, _ := litxaputil.RomanizeIPA(val)

		table[key] = syllables[0][0]
		/*for _, suffix := range suffixes {
			key := key + suffix
			table[key] = append(table[key], fmt.Sprintf("%s: -%s: Incomplete dict entry", name, suffix))
		}*/
	}

	return &nameDict{table: table}
}

var suffixes = []string{
	"l", "ìl",
	"t", "ti", "it",
	"r", "ur", "ru",
	"ri", "ìri",
	"yä", "ä", "ye",
}

func init() {
	adp, err := fwewdict.Adpositions()
	if err != nil {
		log.Println("Failed to get adpositions from fwew:", err)
		return
	}

	suffixes = append(suffixes, adp...)
}
