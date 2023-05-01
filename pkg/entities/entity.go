package entities

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/leekchan/accounting"
)

const Rooms = "ROOMS"
const CoveredArea = "COVERED_AREA"
const TotalArea = "TOTAL_AREA"
const CurrencyDollar = "USD"

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
	Id         string      `json:"id"`
	Title      string      `json:"title"`
	CurrencyId string      `json:"currency_id"`
	Price      float64     `json:"price"`
	Condition  string      `json:"condition"`
	Location   *Location   `json:"location"`
	Permalink  string      `json:"permalink"`
	Attributes []Attribute `json:"attributes"`
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
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Value  string           `json:"value_name"`
	Values []AttributeValue `json:"values"`
}

type AttributeValue struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
}