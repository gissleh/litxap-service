package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gissleh/litxap"
	litxapfwew "github.com/gissleh/litxap-fwew"
	"github.com/gissleh/litxap/litxapfilter"
	"github.com/gissleh/litxap/litxapformats"
)

func main() {
	dict := litxap.MultiDictionary{
		litxapfwew.Global(),
		litxapfwew.MultiWordPartDictionary(),
		&litxap.NumberDictionary{},
	}

	listenAddr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	if listenAddr == ":" {
		listenAddr = ":8081"
	}

	debugAllowOrigin := os.Getenv("LITXAP_ALLOW_ORIGIN")

	log.Println("Starting with address:", listenAddr)

	errCh := make(chan error)
	go func() {
		errCh <- http.ListenAndServe(listenAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			totalStart := time.Now()

			type outputFormats struct {
				IPA             []string `json:"ipa"`
				DiscordMarkdown []string `json:"discordMarkdown"`
				BBCode          []string `json:"bbCode"`
				CompactHtml     []string `json:"compactHtml"`
				IRC             []string `json:"irc"`
			}

			type output struct {
				Line             litxap.Line   `json:"line,omitempty"`
				Lines            []litxap.Line `json:"lines,omitempty"`
				FilteredLines    []litxap.Line `json:"filteredLines,omitempty"`
				Formats          outputFormats `json:"formats"`
				Ambiguities      []any         `json:"ambiguities"`
				UnknownWords     []any         `json:"unknownWords"`
				RunDurationMS    float64       `json:"runDurationMs"`
				TotalDurationMS  float64       `json:"totalDurationMs"`
				FilterDurationMS float64       `json:"filterDurationMs"`
				FormatDurationMS float64       `json:"formatDurationMs"`
			}

			type inputFilters struct {
				DiphthongFromWeakVowel            bool `json:"diphthongFromWeakVowel,omitempty"`
				ReanalyzeDiphthongs               bool `json:"reanalyzeDiphthongs,omitempty"`
				DemoteEjectivesBeforeConsonants   bool `json:"demoteEjectivesBeforeConsonants,omitempty"`
				RemoveRepeatedEjective            bool `json:"removeRepeatedEjective,omitempty"`
				NasalAssimilation                 bool `json:"nasalAssimilation,omitempty"`
				SaeRemover                        bool `json:"saeRemover,omitempty"`
				SpellOeAsWe                       bool `json:"spellOeAsWe,omitempty"`
				ReefUnstressedAeAsE               bool `json:"reefUnstressedAeAsE,omitempty"`
				ReefEjectiveToVoiced              bool `json:"reefEjectiveToVoiced,omitempty"`
				ReefDropGlottalStopsBetweenVowels bool `json:"reefDropGlottalStopsBetweenVowels,omitempty"`
				ReefApplyChSh                     bool `json:"reefApplyChSh,omitempty"`
				ElideUnstressedEWordEndings       bool `json:"elideUnstressedEWordEndings,omitempty"`
				ElideMiSiNiBeforeAy               bool `json:"elideMiSiNiBeforeAy,omitempty"`
				ElideAdvPrefixAndE                bool `json:"elideAdvPrefixAndE,omitempty"`
			}

			var input struct {
				Lines       []string            `json:"lines,omitempty"`
				Multiline   bool                `json:"multiline,omitempty"`
				Selections  map[int]map[int]int `json:"selections,omitempty"`
				CustomWords []string            `json:"customWords,omitempty"`
				Filters     inputFilters        `json:"filters"`
			}
			dict := dict

			w.Header().Set("Content-Type", "application/json")

			if r.URL.Path != "/api/run" {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "Endpoint not found"})
				return
			}

			switch r.Method {
			case "GET":
				q := r.URL.Query()
				input.Lines = q["line"]
				input.Selections = make(map[int]map[int]int)
				input.Multiline = q.Get("multiline") == "true" || len(input.Lines) > 1
				for sel := range strings.SplitSeq(q.Get("selections"), ";") {
					if sel == "" {
						continue
					}

					kv := strings.SplitN(sel, ":", 3)
					if len(kv) < 2 {
						w.WriteHeader(http.StatusBadRequest)
						_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad selection format"})
					}

					var lineIndex int
					var partIndex int
					var selection int
					var err error
					if len(kv) == 2 {
						partIndex, err = strconv.Atoi(kv[0])
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad selection format"})
						}
						selection, err = strconv.Atoi(kv[1])
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad selection format"})
						}
					} else if len(kv) == 3 {
						lineIndex, err = strconv.Atoi(kv[0])
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad selection format"})
						}
						partIndex, err = strconv.Atoi(kv[1])
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad selection format"})
						}
						selection, err = strconv.Atoi(kv[2])
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad selection format"})
						}
					}

					if input.Selections[lineIndex] == nil {
						input.Selections[lineIndex] = make(map[int]int)
					}
					input.Selections[lineIndex][partIndex] = selection
				}
				if names := r.URL.Query().Get("names"); names != "" {
					input.CustomWords = strings.Split(names, ",")
					for i, word := range input.CustomWords {
						input.CustomWords[i] = strings.TrimSpace(word)
					}
				}
				if filters := r.URL.Query().Get("filters"); filters != "" {
					for _, filter := range strings.Split(filters, ",") {
						switch filter {
						case "diphthongFromWeakVowel", "dfwv":
							input.Filters.DiphthongFromWeakVowel = true
						case "reanalyzeDiphthongs", "rd":
							input.Filters.ReanalyzeDiphthongs = true
						case "demoteEjectivesBeforeConsonants", "debc":
							input.Filters.DemoteEjectivesBeforeConsonants = true
						case "removeRepeatedEjective", "rre":
							input.Filters.RemoveRepeatedEjective = true
						case "nasalAssimilation", "na":
							input.Filters.NasalAssimilation = true
						case "saeRemover", "sr":
							input.Filters.SaeRemover = true
						case "spellOeAsWe", "soaw":
							input.Filters.SpellOeAsWe = true
						case "elideUnstressedEWordEndings", "euewe":
							input.Filters.ElideUnstressedEWordEndings = true
						case "elideMiSiNiBeforeAy", "emsnba":
							input.Filters.ElideMiSiNiBeforeAy = true
						case "elideAdvPrefixAndE", "eapae":
							input.Filters.ElideAdvPrefixAndE = true
						case "reefUnstressedAeAsE", "r_uaae":
							input.Filters.ReefUnstressedAeAsE = true
						case "reefEjectiveToVoiced", "r_etv":
							input.Filters.ReefEjectiveToVoiced = true
						case "reefDropGlottalStopsBetweenVowels", "r_dgsbv":
							input.Filters.ReefDropGlottalStopsBetweenVowels = true
						case "reefApplyChSh", "r_chsh":
							input.Filters.ReefApplyChSh = true
						default:
						}
					}
				}
			case "POST":
				err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16384)).Decode(&input)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
					return
				}

				if len(input.Lines) > 1 {
					input.Multiline = true
				}
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
				return
			}

			totalLength := 0
			for _, l := range input.Lines {
				totalLength += len(l)
				if totalLength > 8192 {
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(map[string]string{"error": "line too long"})
					return
				}
			}

			if debugAllowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", debugAllowOrigin)
			}

			if len(input.CustomWords) > 0 {
				dict = litxap.MultiDictionary{dict, litxap.CustomWords(input.CustomWords, "")}
			}

			runStart := time.Now()
			lines, err := litxap.RunLines(input.Lines, dict)
			runDurationMs := time.Since(runStart).Seconds() * 1000
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			/*
				DiscordMarkdown: line.Format(litxapformats.DiscordMarkdown(), input.Selections),
				BBCode:          line.Format(litxapformats.BBCode(), input.Selections),
				CompactHtml:     line.Format(litxapformats.CompactHTML(), input.Selections),
			*/

			res := output{RunDurationMS: runDurationMs}

			filterStart := time.Now()

			var emptyFilters inputFilters
			if input.Filters != emptyFilters {
				filters := make([]litxapfilter.Filter, 0, 16)
				if input.Filters.DiphthongFromWeakVowel {
					filters = append(filters, litxapfilter.DiphthongFromWeakVowel)
				}
				if input.Filters.ReanalyzeDiphthongs {
					filters = append(filters, litxapfilter.ReanalyzeDiphthongs)
				}
				if input.Filters.DemoteEjectivesBeforeConsonants {
					filters = append(filters, litxapfilter.DemoteEjectivesBeforeConsonants)
				}
				if input.Filters.RemoveRepeatedEjective {
					filters = append(filters, litxapfilter.RemoveRepeatedEjective)
				}
				if input.Filters.NasalAssimilation {
					filters = append(filters, litxapfilter.NasalAssimilation)
				}
				if input.Filters.SaeRemover {
					filters = append(filters, litxapfilter.SaeRemover)
				}
				if input.Filters.SpellOeAsWe {
					filters = append(filters, litxapfilter.SpellOeAsWe)
				}
				if input.Filters.ElideUnstressedEWordEndings {
					filters = append(filters, litxapfilter.ElideUnstressedEWordEndings)
				}
				if input.Filters.ElideMiSiNiBeforeAy {
					filters = append(filters, litxapfilter.ElideMiSiNiBeforeAy)
				}
				if input.Filters.ElideAdvPrefixAndE {
					filters = append(filters, litxapfilter.ElideAdvPrefixAndE)
				}
				if input.Filters.ReefUnstressedAeAsE {
					filters = append(filters, litxapfilter.ReefUnstressedAeAsE)
				}
				if input.Filters.ReefEjectiveToVoiced {
					filters = append(filters, litxapfilter.ReefEjectiveToVoiced)
				}
				if input.Filters.ReefDropGlottalStopsBetweenVowels {
					filters = append(filters, litxapfilter.ReefDropGlottalStopsBetweenVowels)
				}
				if input.Filters.ReefApplyChSh {
					filters = append(filters, litxapfilter.ReefApplyChSh)
				}

				res.FilteredLines = make([]litxap.Line, len(lines))
				for i := range lines {
					res.FilteredLines[i] = litxapfilter.ApplyFilters(
						lines[i].WithSelections(input.Selections[i], true),
						filters...,
					)
				}
			}

			if input.Multiline {
				res.Lines = lines
			} else {
				res.Line = lines[0]
			}

			res.FilterDurationMS = time.Since(filterStart).Seconds() * 1000

			formatStart := time.Now()
			for i, line := range lines {
				selections := input.Selections[i]
				if res.FilteredLines != nil && res.FilteredLines[i] != nil {
					line = res.FilteredLines[i]
					selections = nil
				}

				res.Formats.DiscordMarkdown = append(
					res.Formats.DiscordMarkdown,
					line.Format(litxapformats.DiscordMarkdown(), selections),
				)
				res.Formats.BBCode = append(
					res.Formats.BBCode,
					line.Format(litxapformats.BBCode(), selections),
				)
				res.Formats.CompactHtml = append(
					res.Formats.CompactHtml,
					line.Format(litxapformats.CompactHTML(), selections),
				)
				res.Formats.IRC = append(
					res.Formats.IRC,
					line.Format(litxapformats.IRCDefaultColors(), selections),
				)

				for j, part := range line {
					selection, ok := selections[j]
					if !ok {
						selection = -1
					}

					var ij any = j
					if input.Multiline {
						ij = [2]int{i, j}
					}

					_, stress := part.GetSyllables(selection)
					switch stress {
					case litxap.LPSAmbiguousMatches:
						res.Ambiguities = append(res.Ambiguities, ij)
					case litxap.LPSNoMatches:
						res.UnknownWords = append(res.UnknownWords, ij)
					}
				}

				ipa, err := line.IPA(selections, ".")
				if err != nil {
					res.Formats.IPA = append(res.Formats.IPA, "ERROR: "+err.Error())
				} else {
					res.Formats.IPA = append(res.Formats.IPA, ipa)
				}
			}
			res.FormatDurationMS = time.Since(formatStart).Seconds() * 1000
			res.TotalDurationMS = time.Since(totalStart).Seconds() * 1000

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(res)
		}))
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		log.Println("Got signal:", s)
		os.Exit(0)
	case err := <-errCh:
		log.Println("Got error:", err)
		os.Exit(1)
	}
}
