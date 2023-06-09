package meli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"real-state-finder/pkg/entities"
	"real-state-finder/pkg/storage"
	"sort"
	"strconv"
	"strings"
	"time"
)

const apiUrl = "https://api.mercadolibre.com"
const methodCountry = "classified_locations/countries/$COUNTRY_ID"
const methodState = "classified_locations/states/$STATE_ID"
const methodCity = "classified_locations/cities/$CITY_ID"
const authorizationKey = "Authorization"
const realStateFile = "real_state.json"
const droppedDueToUsdFile = "dropped_usd.json"
const droppedDueToAmbientsFile = "dropped_ambients.json"
const droppedDueToPriceFile = "dropped_price.json"
const droppedDueToTotalAreaFile = "dropped_total_area.json"
const droppedDueToWrongNeighborhoodMF = "dropped_wrong_neighborhood_mf.json"

// Consts for SearchAPI
const (
	methodSearch             = "sites/MLA/search"
	searchQueryParams        = "category=$CATEGORY_ID&state=$STATE_ID&OPERATION=$OPERATION&PROPERTY_TYPE=$PROPERTY_TYPE&limit=$LIMIT"
	categoryId               = "MLA1459"              // Inmuebles
	stateId                  = "TUxBUENBUGw3M2E1"     // CABA
	operationType            = "242073"               // Alquiler
	propertyTypes            = "242062,242069,242060" // Departamento,Ph,Casas
	defaultSearchResultLimit = 2
)

type Api interface {
	GetCountry(countryId string) (*entities.Country, error)
	GetState(stateId string) (*entities.State, error)
	GetRealState(offset int) ([]entities.RealState, error)
	CmdSearch(offset int, filterNeighborhood string) error
	CmdRead() error
	CmdGenerateHtml() error
	CmdInitCache() error
}

type api struct {
	accessToken       string
	searchResultLimit int
	maxOffset         int
	minPrice          int
	maxPrice          int
	minAmbients       int
	minTotalArea      int
	storage           storage.Storage
}

func NewApi(a string, l int, m int, minP int, maxP int, minA int, minT int, s storage.Storage) Api {
	return &api{a, l, m, minP, maxP, minA, minT, s}
}

func (a *api) GetCountry(countryId string) (*entities.Country, error) {
	var country *entities.Country
	bArr, err := a.doRequest(fmt.Sprintf("%s/%s", apiUrl, strings.Replace(methodCountry, "$COUNTRY_ID", countryId, 1)))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bArr, &country)
	if err != nil {
		return nil, err
	}

	return country, err
}

func (a *api) GetState(stateId string) (*entities.State, error) {
	var state *entities.State
	bArr, err := a.doRequest(fmt.Sprintf("%s/%s", apiUrl, strings.Replace(methodState, "$STATE_ID", stateId, 1)))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bArr, &state)
	if err != nil {
		return nil, err
	}

	return state, err
}

func (a *api) GetRealState(offset int) ([]entities.RealState, error) {
	result := make([]entities.RealState, 0)

	totalResults, err := a.getTotalResults()
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("Total results in query %d", totalResults))

	loopCount := totalResults / a.searchResultLimit

	// Let's limit the loop count to 5 for now
	if loopCount > a.maxOffset {
		loopCount = a.maxOffset
	}

	for i := 0; i < loopCount; i++ {
		realStateList, err := a.doGetRealStateResults(a.searchResultLimit * i)
		if err != nil {
			return nil, err
		}
		result = append(result, realStateList...)

		fmt.Println(fmt.Sprintf("Done for offset %d. Wait for 5 seconds", i))
		time.Sleep(5 * time.Second)
	}

	return result, nil
}

func (a *api) getTotalResults() (int, error) {
	var realStateResults *entities.RealStateResults
	searchUrl := a.getSearchUrl(0)
	fmt.Println(fmt.Sprintf("Search URL to get Totals: %s", searchUrl))
	bArr, err := a.doRequest(searchUrl)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(bArr, &realStateResults)
	if err != nil {
		return 0, err
	}

	return int(realStateResults.Paging.Total), nil
}

func (a *api) doGetRealStateResults(offset int) ([]entities.RealState, error) {
	var realStateResults *entities.RealStateResults
	searchUrl := a.getSearchUrl(offset)
	fmt.Println(fmt.Sprintf("Search URL: %s", searchUrl))
	bArr, err := a.doRequest(searchUrl)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bArr, &realStateResults)
	if err != nil {
		return nil, err
	}

	return realStateResults.Results, nil
}

func (a *api) getSearchUrl(offset int) string {
	limit := defaultSearchResultLimit
	if a.searchResultLimit > defaultSearchResultLimit {
		limit = a.searchResultLimit
	}

	searchUrl := strings.Replace(searchQueryParams, "$CATEGORY_ID", categoryId, 1)
	searchUrl = strings.Replace(searchUrl, "$STATE_ID", stateId, 1)
	searchUrl = strings.Replace(searchUrl, "$OPERATION", operationType, 1)
	searchUrl = strings.Replace(searchUrl, "$PROPERTY_TYPE", propertyTypes, 1)
	searchUrl = strings.Replace(searchUrl, "$LIMIT", fmt.Sprintf("%d", limit), 1)

	if offset > 0 {
		searchUrl = fmt.Sprintf("%s&offset=%d", searchUrl, offset)
	}

	fmt.Println(fmt.Sprintf("Search Query Params: %s", searchUrl))
	return fmt.Sprintf("%s/%s?%s", apiUrl, methodSearch, searchUrl)
}

// doRequest calls the api with the uri passed as parameter and return the response byte array
func (a *api) doRequest(uri string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(authorizationKey, fmt.Sprintf("Bearer %s", a.accessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf("Error querying search api. Status %d", resp.StatusCode))
	}

	bArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bArr, nil
}

func (a *api) CmdSearch(offset int, filterNeighborhood string) error {
	var droppedUsdList, droppedAmbientsList, droppedPriceList, droppedTotalAreaList, droppedFilterNeighbor, realStateList []entities.RealState

	a.storage.ResetNew()

	realStates, err := a.GetRealState(0) // Start the recursive function from 0
	fmt.Println(fmt.Sprintf("Found %d results", len(realStates)))
	if err != nil {
		return err
	}

	for _, rs := range realStates {
		ambients, err := strconv.Atoi(rs.GetAttributeValue(entities.Rooms))
		if err != nil {
			fmt.Println("ambients is not int")
		}

		totalArea := int(rs.GetValueStruct(entities.TotalArea))

		// If currency of the item is in dollars we don't care about it
		if rs.CurrencyId == entities.CurrencyDollar {
			droppedUsdList = append(droppedUsdList, rs)
			continue
		}

		// If price is less than min price or greater than max price then drop it
		if int(rs.Price) < a.minPrice || int(rs.Price) > a.maxPrice {
			droppedPriceList = append(droppedPriceList, rs)
			continue
		}

		if ambients < a.minAmbients {
			droppedAmbientsList = append(droppedAmbientsList, rs)
			continue
		}

		if totalArea < a.minTotalArea {
			droppedTotalAreaList = append(droppedTotalAreaList, rs)
			continue
		}

		if filterNeighborhood != "" {
			if rs.Location.Neighborhood.Name != filterNeighborhood {
				droppedFilterNeighbor = append(droppedFilterNeighbor, rs)
				continue
			}
		}

		fmt.Println(fmt.Sprintf("About to store real state with id %s", rs.Id))
		if !a.storage.Exists(rs.Id) {
			fmt.Println(fmt.Sprintf("Real state with id %s didn't exist, storing it", rs.Id))
			rs.IsNew = true
			rs.CreatedDate = time.Now()
			a.storage.Save(rs)
			rs.Print()
		}

		realStateList = append(realStateList, rs)
		fmt.Println(fmt.Sprintf("Didn't store because real state with id %s already exists in storage", rs.Id))

	}

	fmt.Println()
	fmt.Println()
	fmt.Println(fmt.Sprintf("Dropped %d due to USD, %d due to less than desired ambients, %d due to price, %d due to TotalArea out of %d found.", len(droppedUsdList), len(droppedAmbientsList), len(droppedPriceList), len(droppedTotalAreaList), len(realStates)))
	fmt.Println()
	fmt.Println()
	fmt.Println(fmt.Sprintf("Found %d appartments out of %d", len(realStateList), len(realStates)))
	fmt.Println()
	fmt.Println()

	fmt.Println(fmt.Print("Saving files..."))
	saveToFile(realStateFile, a.storage.GetList())
	saveToFile(droppedDueToUsdFile, droppedUsdList)
	saveToFile(droppedDueToAmbientsFile, droppedAmbientsList)
	saveToFile(droppedDueToPriceFile, droppedPriceList)
	saveToFile(droppedDueToTotalAreaFile, droppedTotalAreaList)
	saveToFile(droppedDueToWrongNeighborhoodMF, droppedFilterNeighbor)

	return nil
}

func saveToFile(filename string, rs []entities.RealState) error {
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

func (a *api) CmdRead() error {
	realState, err := readFromFile(realStateFile)
	if err != nil {
		return err
	}

	for _, rs := range realState {
		rs.Print()
	}

	return nil
}

func readFromFile(filename string) ([]entities.RealState, error) {
	var realState []entities.RealState
	bArr, err := ioutil.ReadFile(realStateFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bArr, &realState)
	if err != nil {
		return nil, err
	}

	return realState, nil
}

func (a *api) CmdGenerateHtml() error {
	simpleRealState := make([]entities.SimpleRealState, 0)
	realStateList, err := readFromFile(realStateFile)
	if err != nil {
		return err
	}

	sort.Sort(entities.ByPrice(realStateList))

	for _, rs := range realStateList {
		simpleRealState = append(simpleRealState, rs.ToSimpleRealState())
	}

	htmlData := entities.HtmlRepresentation{
		TotalCount: len(simpleRealState),
		RealState:  simpleRealState,
		Timestamp:  time.Now().Format("Monday, 02 Jan 2006 15:04:05"),
	}

	allFiles := []string{"content.tmpl", "footer.tmpl", "header.tmpl", "page.tmpl"}

	var allPaths []string

	for _, tmpl := range allFiles {
		allPaths = append(allPaths, "./templates/"+tmpl)
	}

	templates := template.Must(template.New("").ParseFiles(allPaths...))

	var processed bytes.Buffer
	templates.ExecuteTemplate(&processed, "page", htmlData)

	f, _ := os.Create("./index.html")
	w := bufio.NewWriter(f)
	w.WriteString(string(processed.Bytes()))
	w.Flush()
	return nil
}

func (a *api) CmdInitCache() error {
	return nil
}
