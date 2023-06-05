package entities

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/leekchan/accounting"
)

const Rooms = "ROOMS"
const CoveredArea = "COVERED_AREA"
const TotalArea = "TOTAL_AREA"
const CurrencyDollar = "USD"

const thumbnailUrl = "https://http2.mlstatic.com/D_NQ_NP_$THUMBNAIL_ID-W.webp"

type Country struct {
	Name   string  `json:"name"`
	States []State `json:"states"`
}

type State struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Cities []City `json:"cities"`
}

type City struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RealStateResults struct {
	Results []RealState `json:"results"`
	Paging  Paging      `json:"paging"`
}

type Paging struct {
	Total  float64 `json:"total"`
	Offset float64 `json:"offset"`
	Limit  float64 `json:"limit"`
}

type RealState struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	CurrencyId  string      `json:"currency_id"`
	Price       float64     `json:"price"`
	Condition   string      `json:"condition"`
	Location    *Location   `json:"location"`
	Permalink   string      `json:"permalink"`
	ThumbnailId string      `json:"thumbnail_id"`
	Attributes  []Attribute `json:"attributes"`
}

type Location struct {
	AddressLine  string       `json:"address_line"`
	Neighborhood Neighborhood `json:"neighborhood"`
}

type Neighborhood struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Attribute struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Value       string           `json:"value_name"`
	Values      []AttributeValue `json:"values"`
	ValueStruct *ValueStruct     `json:"value_struct"`
}

type ValueStruct struct {
	Number float64 `json:"number"`
	Unit   string  `json:"unit"`
}

type AttributeValue struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ByPrice []RealState

func (a ByPrice) Len() int           { return len(a) }
func (a ByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPrice) Less(i, j int) bool { return a[i].Price > a[j].Price }

func (rs *RealState) GetValueStruct(attrName string) float64 {
	for _, attr := range rs.Attributes {
		if attr.Id == attrName && attr.ValueStruct != nil {
			return attr.ValueStruct.Number
		}
	}
	return 0
}

func (rs *RealState) GetAttributeValue(attrName string) string {
	for _, attr := range rs.Attributes {
		if attr.Id == attrName {
			return attr.Value
		}
	}
	return "Not Found"
}

func (rs *RealState) Print() {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	ambients, err := strconv.Atoi(rs.GetAttributeValue(Rooms))
	if err != nil {
		fmt.Println("ambients is not int")
	}
	coveredArea := rs.GetAttributeValue(CoveredArea)
	totalArea := rs.GetAttributeValue(TotalArea)
	address := "Not Specified"
	neighborhood := "Not Specified"
	if rs.Location != nil {
		address = rs.Location.AddressLine
		neighborhood = rs.Location.Neighborhood.Name
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Title", "Ambients", "Covered Area", "Total Area", "Neighborhood", "Address", "Price"})
	t.AppendRow(table.Row{rs.Title, ambients, coveredArea, totalArea, neighborhood, address, ac.FormatMoney(rs.Price)})
	t.AppendSeparator()

	t.Render()

	fmt.Println(rs.Permalink)

}

func (rs *RealState) ToSimpleRealState() SimpleRealState {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	address := "Not Specified"
	neighborhood := "Not Specified"
	if rs.Location != nil {
		address = rs.Location.AddressLine
		neighborhood = rs.Location.Neighborhood.Name
	}

	return SimpleRealState{
		Title:        rs.Title,
		Ambients:     rs.GetAttributeValue(Rooms),
		CoveredArea:  rs.GetAttributeValue(CoveredArea),
		TotalArea:    rs.GetAttributeValue(TotalArea),
		Address:      address,
		Neighborhood: neighborhood,
		Price:        ac.FormatMoney(rs.Price),
		Permalink:    rs.Permalink,
		Thumbnail:    strings.Replace(thumbnailUrl, "$THUMBNAIL_ID", rs.ThumbnailId, 1),
	}
}

type SimpleRealState struct {
	Title        string
	Ambients     string
	CoveredArea  string
	TotalArea    string
	Neighborhood string
	Address      string
	Price        string
	Permalink    string
	Thumbnail    string
}

type HtmlRepresentation struct {
	TotalCount int
	RealState  []SimpleRealState
	Timestamp  string
}
