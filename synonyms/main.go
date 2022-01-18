package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	thesaurus "github.com/yasir2000/go-web-dev-side-project-2/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHugh{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalln("Failed when looking for synonyms for  "+word+"", err)
		}
		if len(syns) == 0 {
			log.Fatalln("Couldnt find any synonyms for "+word+"", err)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
