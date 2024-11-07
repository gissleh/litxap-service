package fwewdict

import "strings"

/* To help deduce phonemes */
var romanization2 = map[string]string{
	// Vowels
	"a": "a", "i": "i", "ɪ": "ì",
	"o": "o", "ɛ": "e", "u": "u",
	"æ": "ä", "õ": "õ", //võvä' only
	// Diphthongs
	"aw": "aw", "ɛj": "ey",
	"aj": "ay", "ɛw": "ew",
	// Psuedovowels
	"ṛ": "rr", "ḷ": "ll",
	// Consonents
	"t": "t", "p": "p", "ʔ": "'",
	"n": "n", "k": "k", "l": "l",
	"s": "s", "ɾ": "r", "j": "y",
	"t͡s": "ts", "t'": "tx", "m": "m",
	"v": "v", "w": "w", "h": "h",
	"ŋ": "ng", "z": "z", "k'": "kx",
	"p'": "px", "f": "f", "r": "r",
	// Reef dialect
	"b": "b", "d": "d", "g": "g",
	"ʃ": "sh", "tʃ": "ch", "ʊ": "ù",
	// mistakes and rarities
	"ʒ": "ch", "": "", " ": ""}

func nth_rune(word string, n int) (output string) {
	r := []rune(word)
	if n < 0 { // negative index
		n = len(r) + n
	}
	if n >= len(r) {
		return ""
	}
	return string(r[n])
}

func has(word string, character string) (output bool) {
	r := []rune(word)
	if len(character) == 0 {
		return false
	}
	c := []rune(character)[0]
	for i := 0; i < len(r); i++ {
		if c == r[i] {
			return true
		}
	}
	return false
}

// Helper function to get phonetic transcriptions of secondary pronunciations
func RomanizeIPA(IPA string) []string {
	// now Romanize the IPA
	IPA = strings.ReplaceAll(IPA, "ʊ", "u")
	IPA = strings.ReplaceAll(IPA, "õ", "o") // vonvä' as võvä' only
	word := strings.Split(IPA, " ")

	results := []string{}

	// Make sure it's not the same word with different stresses
	if len(word) > 2 {
		word[0] = strings.ReplaceAll(word[0], "[", "")
		word[0] = strings.ReplaceAll(word[0], "]", "")
		word[2] = strings.ReplaceAll(word[2], "[", "")
		word[2] = strings.ReplaceAll(word[2], "]", "")

		if strings.ReplaceAll(word[0], "ˈ", "") == strings.ReplaceAll(word[2], "ˈ", "") {
			word = []string{strings.ReplaceAll(word[0], "ˈ", "")}
		}
	}

	// get the last one only
	for j := 0; j < len(word); j++ {
		breakdown := ""

		word[j] = strings.ReplaceAll(word[j], "[", "")
		word[j] = strings.ReplaceAll(word[j], "]", "")
		// "or" means there's more than one IPA in this word, and we only want one
		if word[j] == "or" {
			breakdown = ""
			continue
		}

		syllables := strings.Split(word[j], ".")

		/* Onset */
		for k := 0; k < len(syllables); k++ {
			stressed := strings.Contains(syllables[k], "ˈ")

			syllable := strings.ReplaceAll(syllables[k], "·", "")
			syllable = strings.ReplaceAll(syllable, "ˈ", "")
			syllable = strings.ReplaceAll(syllable, "ˌ", "")

			if stressed {
				breakdown += "__"
			}

			// tsy
			if strings.HasPrefix(syllable, "tʃ") {
				breakdown += "ch"
				syllable = strings.TrimPrefix(syllable, "tʃ")
			} else if len(syllable) >= 4 && syllable[0:4] == "t͡s" {
				// ts
				breakdown += "ts"
				//tsp
				if has("ptk", nth_rune(syllable, 3)) {
					if nth_rune(syllable, 4) == "'" {
						// ts + ejective onset
						breakdown += romanization2[syllable[4:6]]
						syllable = syllable[6:]
					} else {
						// ts + unvoiced plosive
						breakdown += romanization2[string(syllable[4])]
						syllable = syllable[5:]
					}
				} else if has("lɾmnŋwj", nth_rune(syllable, 3)) {
					// ts + other consonent
					breakdown += romanization2[nth_rune(syllable, 3)]
					syllable = syllable[4+len(nth_rune(syllable, 3)):]
				} else {
					// ts without a cluster
					syllable = syllable[4:]
				}
			} else if has("fs", nth_rune(syllable, 0)) {
				//
				breakdown += nth_rune(syllable, 0)
				if has("ptk", nth_rune(syllable, 1)) {
					if nth_rune(syllable, 2) == "'" {
						// f/s + ejective onset
						breakdown += romanization2[syllable[1:3]]
						syllable = syllable[3:]
					} else {
						// f/s + unvoiced plosive
						breakdown += romanization2[string(syllable[1])]
						syllable = syllable[2:]
					}
				} else if has("lɾmnŋwj", nth_rune(syllable, 1)) {
					// f/s + other consonent
					breakdown += romanization2[nth_rune(syllable, 1)]
					syllable = syllable[1+len(nth_rune(syllable, 1)):]
				} else {
					// f/s without a cluster
					syllable = syllable[1:]
				}
			} else if has("ptk", nth_rune(syllable, 0)) {
				if nth_rune(syllable, 1) == "'" {
					// ejective
					breakdown += romanization2[syllable[0:2]]
					syllable = syllable[2:]
				} else {
					// unvoiced plosive
					breakdown += romanization2[string(syllable[0])]
					syllable = syllable[1:]
				}
			} else if has("ʔlɾhmnŋvwjzbdg", nth_rune(syllable, 0)) {
				// other normal onset
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = syllable[len(nth_rune(syllable, 0)):]
			} else if has("ʃʒ", nth_rune(syllable, 0)) {
				// one sound representd as a cluster
				if nth_rune(syllable, 0) == "ʃ" {
					breakdown += "sh"
				}
				syllable = syllable[len(nth_rune(syllable, 0)):]
			}

			/*
			 * Nucleus
			 */
			if len(syllable) > 1 && has("jw", nth_rune(syllable, 1)) {
				//diphthong
				breakdown += romanization2[syllable[0:len(nth_rune(syllable, 0))+1]]
				syllable = string([]rune(syllable)[2:])
			} else if len(syllable) > 1 && has("lr", nth_rune(syllable, 0)) {
				breakdown += romanization2[syllable[0:3]]
				continue
			} else {
				//vowel
				breakdown += romanization2[nth_rune(syllable, 0)]
				syllable = string([]rune(syllable)[1:])
			}

			/*
			 * Coda
			 */
			if len(syllable) > 0 {
				if nth_rune(syllable, 0) == "s" {
					breakdown += "sss" //oìsss only
				} else {
					if syllable == "k̚" {
						breakdown += "k"
					} else if syllable == "p̚" {
						breakdown += "p"
					} else if syllable == "t̚" {
						breakdown += "t"
					} else if syllable == "ʔ̚" {
						breakdown += "'"
					} else {
						if syllable[0] == 'k' && len(syllable) > 1 {
							breakdown += "kx"
						} else {
							breakdown += romanization2[syllable]
						}
					}
				}
			}
			if stressed {
				breakdown += "__"
			}
		}
		results = append(results, breakdown)
	}

	return results
}
