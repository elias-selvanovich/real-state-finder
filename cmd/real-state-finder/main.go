package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"real-state-finder/pkg/entities"
	"real-state-finder/pkg/meli"
)

const CommandSearch = "search"
const CommandRead = "read"
const CommandGenerateHtml = "generate-html"
const CommandInitCache = "init-cache"

func main() {

	accessTokenFlag := flag.String("access-token", "", "access token for meli api")
	searchResultLimitFlag := flag.Int("search-result-limit", 2, "Limit of the search result results")
	commandFlag := flag.String("command", "search", "Type of command. Available options: search / print")
	maxOffsetFlag := flag.Int("max-offset", 5, "Maximum offset to search")
	minPriceFlag := flag.Int("min-price", 120000, "Min Price to search")
	maxPriceFlag := flag.Int("max-price", 400000, "Max Price to search")
	minAmbientsFlag := flag.Int("min-ambients", 3, "Min Ambients to search")
	minTotalAreaFlag := flag.Int("min-total-area", 70, "Min Total Area to search")
	filterNeighborhoodFlag := flag.String("filter-neighborhood", "", "Neighborhood to filter")

	flag.Parse()

	accessToken := *accessTokenFlag
	searchResultLimit := *searchResultLimitFlag
	command := *commandFlag
	maxOffset := *maxOffsetFlag
	minPrice := *minPriceFlag
	maxPrice := *maxPriceFlag
	minAmbients := *minAmbientsFlag
	minTotalArea := *minTotalAreaFlag
	filterNeighborhood := *filterNeighborhoodFlag

	if accessToken == "" {
		panic("access token can't be empty")
	}

	printHeader()

	api := meli.NewApi(accessToken, searchResultLimit, maxOffset, minPrice, maxPrice, minAmbients, minTotalArea)

	switch command {
	case CommandSearch:
		err := api.CmdSearch(0, filterNeighborhood)
		if err != nil {
			panic(err)
		}
		break
	case CommandRead:
		err := api.CmdRead()
		if err != nil {
			panic(err)
		}
		break
	case CommandGenerateHtml:
		err := api.CmdGenerateHtml()
		if err != nil {
			panic(err)
		}
		break
	}
}

func printHeader() {
	fmt.Println("================= REAL STATE FINDER =================")
}

func printMessage(msg string) {
	fmt.Println(fmt.Sprintf("> %s", msg))
}

func dump(filename string, rs []entities.RealState) error {
	fmt.Println(fmt.Sprintf("Saving %s", filename))
	bArr, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, bArr, 0644)
	if err != nil {
		return err
	}

	return nil
}
