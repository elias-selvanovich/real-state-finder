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

func main() {

	accessTokenFlag := flag.String("access-token", "", "access token for meli api")
	searchResultLimitFlag := flag.Int("search-result-limit", 2, "Limit of the search result results")
	commandFlag := flag.String("command", "search", "Type of command. Available options: search / print")
	maxOffsetFlag := flag.Int("max-offset", 5, "Maximum offset to search")

	flag.Parse()

	accessToken := *accessTokenFlag
	searchResultLimit := *searchResultLimitFlag
	command := *commandFlag
	maxOffset := *maxOffsetFlag

	if accessToken == "" {
		panic("access token can't be empty")
	}

	printHeader()

	api := meli.NewApi(accessToken, searchResultLimit, maxOffset)

	switch command {
	case CommandSearch:
		err := api.CmdSearch(0)
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
