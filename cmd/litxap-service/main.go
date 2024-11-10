package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	fwew_lib "github.com/fwew/fwew-lib/v5"
	"github.com/gissleh/litxap"
	"github.com/gissleh/litxap-service/adapters/fwewdict"
	"github.com/gissleh/litxap-service/adapters/namedict"
)

func main() {
	dict := fwewdict.Global()

	listenAddr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	if listenAddr == ":" {
		listenAddr = ":8081"
	}

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
			result1, _ := fwew_lib.TranslateFromNaviHash(key, true)
			result2, _ := fwew_lib.TranslateFromNaviHash(fullWord.String(), true)
			IPAstring = strings.Split(result2[0][1].IPA, " ")

			if len(result1[0]) < 2 {
				doubles[key] = IPAstring[0]
			}

			for i, multiword := range stringArray {
				res3, _ := fwew_lib.TranslateFromNaviHash(multiword, true)
				if len(res3[0]) < 2 {
					doubles[multiword] = IPAstring[i+1]
				}
			}
		}
	}

	fmt.Println(doubles)

	log.Println("Starting with address:", listenAddr)

	errCh := make(chan error)
	go func() {
		errCh <- http.ListenAndServe(listenAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/run" {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "Endpoint not found"})
				return
			}

			dict := dict
			if names := r.URL.Query().Get("names"); names != "" {
				dict = litxap.MultiDictionary{dict, namedict.New(strings.Split(names, ",")...)}
			}

			w.Header().Set("Content-Type", "application/json")

			q := r.URL.Query()
			line, err := litxap.RunLine(q.Get("line"), dict, doubles)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{"line": line})
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
