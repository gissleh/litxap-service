package namedict

import (
	"fmt"
	"github.com/gissleh/litxap/litxaputil"
	"log"
	"strings"

	"github.com/gissleh/litxap"
	"github.com/gissleh/litxap-service/adapters/fwewdict"
)

type nameDict struct {
	table      map[string][]string
	definition string
}

func (n *nameDict) LookupEntries(word string) ([]litxap.Entry, error) {
	entryStrs, ok := n.table[word]
	if !ok {
		return nil, litxap.ErrEntryNotFound
	}

	entries := make([]litxap.Entry, len(entryStrs))
	for i, entryStr := range entryStrs {
		entries[i] = *litxap.ParseEntry(entryStr + ": " + n.definition)
	}

	return entries, nil
}

func New(names ...string) litxap.Dictionary {
	table := make(map[string][]string)

	for _, name := range names {
		beforeTrim := len(name)
		name := strings.TrimPrefix(name, "-")
		noStress := ""
		if len(name) != beforeTrim {
			noStress = "no_stress"
		}

		name = strings.Replace(strings.ToLower(name), "-", ".", -1)
		key := strings.Replace(strings.Replace(name, "*", "", -1), ".", "", -1)

		table[key] = append(table[key], fmt.Sprintf("%s: %s", name, noStress))
		for _, suffix := range suffixes {
			key := key + suffix
			table[key] = append(table[key], fmt.Sprintf("%s: -%s %s", name, suffix, noStress))
		}

		if possibleLoanWord := strings.TrimSuffix(name, "ì"); possibleLoanWord != name {
			key := strings.Replace(strings.Replace(possibleLoanWord, "*", "", -1), ".", "", -1)
			table[key] = append(table[key], fmt.Sprintf("%s: %s", name, noStress))
			for _, suffix := range loanWordSuffixes {
				key := key + suffix
				table[key] = append(table[key], fmt.Sprintf("%s: -%s %s", name, suffix, noStress))
			}
		}
	}

	return &nameDict{table: table, definition: "Custom Name"}
}

func FromFwewMultiWordParts() litxap.Dictionary {
	names := make([]string, 0, 16)
	for _, ipa := range fwewdict.FindMultis() {
		wordOptions, stressOptions := litxaputil.RomanizeIPA(ipa)
		for i := range wordOptions {
			words := wordOptions[i]
			stress := stressOptions[i]

			if len(words) == 1 && len(stress) == 1 {
				syllables := words[0]
				stressIndex := -1
				for _, index := range stress {
					if index != -1 {
						stressIndex = index
						break
					}
				}

				if stressIndex == -1 {
					syllables[0] = "-" + syllables[0]
				} else {
					for i := range syllables {
						if i == stressIndex {
							syllables[i] = "*" + syllables[i]
						}
					}
				}

				names = append(names, strings.Join(syllables, "-"))
			}
		}
	}

	res := New(names...)
	res.(*nameDict).definition = "Part of multi-part word"
	return res
}

var suffixes = []string{
	"l", "ìl",
	"t", "ti", "it",
	"r", "ur", "ru",
	"ri", "ìri",
	"yä", "ä", "ye",
}

var loanWordSuffixes = []string{
	"ìl",
	"it",
	"ur",
	"ìri",
	"ä",
}

func init() {
	adp, err := fwewdict.Adpositions()
	if err != nil {
		log.Println("Failed to get adpositions from fwew:", err)
		return
	}

	suffixes = append(suffixes, adp...)
}
