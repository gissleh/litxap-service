package fwewdict

import (
	"bytes"
	"log"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	fwew_lib "github.com/fwew/fwew-lib/v5"
	"github.com/gissleh/litxap"
	"github.com/gissleh/litxap/litxaputil"
)

type fwewDict struct {
	mu sync.Mutex
}

var globalDict = &fwewDict{}

func (d *fwewDict) LookupEntries(word string) ([]litxap.Entry, error) {
	stopped := int32(0)

	d.mu.Lock()
	if len(word) > 16 {
		go func() {
			time.Sleep(time.Millisecond * 250)
			if atomic.CompareAndSwapInt32(&stopped, 0, 1) {
				panic("fwew hung on " + word)
			}
		}()
	}

	res, err := fwew_lib.TranslateFromNaviHash(word, true, false, true)
	if atomic.CompareAndSwapInt32(&stopped, 0, 1) {
		d.mu.Unlock()
	}
	if err != nil {
		return nil, err
	}

	entries := make([]litxap.Entry, 0, len(res))

	for _, matches := range res {
		for _, match := range matches {
			if match.ID == "" {
				continue
			}

			syllables := strings.Split(strings.ReplaceAll(strings.ToLower(match.Syllables), " ", "-"), "-")

			for _, ipa := range strings.Split(match.IPA, "or") {
				ipa = strings.Trim(ipa, " []")
				ipaSyllables := strings.Split(strings.ReplaceAll(ipa, " ", "."), ".")
				if len(ipaSyllables) != len(syllables) {
					continue
				}

				stressIndex := 0
				for i, syllable := range ipaSyllables {
					if strings.HasPrefix(syllable, "ˈ") {
						stressIndex = i
						break
					}
				}

				suffixes := append([]string{}, match.Affixes.Suffix...)

				slices.Reverse(suffixes)

				entry := litxap.Entry{
					Word:        match.Navi,
					Translation: match.EN,
					Syllables:   syllables,
					Stress:      stressIndex,
					InfixPos:    litxaputil.InfixPositionsFromBrackets(match.InfixLocations, syllables),
					Prefixes:    match.Affixes.Prefix,
					Infixes:     match.Affixes.Infix,
					Suffixes:    suffixes,
				}

				entries = append(entries, entry)
			}
		}
	}

	return entries, nil
}

func Adpositions() ([]string, error) {
	list, err := fwew_lib.List([]string{"pos", "has", "adp."}, 0)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(list))
	for _, match := range list {
		res = append(res, strings.TrimSuffix(match.Navi, "+"))
	}

	return res, nil
}

func FindMultis() map[string]string {
	// Calculate the multiword words needed at startup
	// Make sure we have words that must be multiword words
	doubles := map[string]string{}
	multis := fwew_lib.GetMultiwordWords()
	fullWord := bytes.NewBuffer(make([]byte, 0, 16))
	IPAstring := []string{}
	for key, val := range multis {
		for _, stringArray := range val {
			fullWord.Reset()
			fullWord.WriteString(key + " ")
			for i, multiword := range stringArray {
				fullWord.WriteString(multiword)
				if i+1 != len(stringArray) {
					fullWord.WriteString(" ")
				}
			}
			result1, _ := fwew_lib.TranslateFromNaviHash(key, true, false, true)
			result2, _ := fwew_lib.TranslateFromNaviHash(fullWord.String(), true, false, true)
			if len(result2[0]) == 1 {
				log.Println(fullWord.String(), "-- not found")
				continue
			}

			IPAstring = strings.Split(result2[0][1].IPA, " ")

			if len(result1[0]) < 2 {
				doubles[key] = IPAstring[0]
			}

			for i, multiword := range stringArray {
				res3, _ := fwew_lib.TranslateFromNaviHash(multiword, true, false, true)
				if len(res3[0]) < 2 {
					doubles[multiword] = IPAstring[i+1]
				}
			}
		}
	}

	return doubles
}

func Global() litxap.Dictionary {
	return globalDict
}

func init() {
	fwew_lib.StartEverything()
}
