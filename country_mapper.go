package country_mapper

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
)

//go:embed files/country_info.csv
var defaultCountryInfo []byte

type CountryInfoClient struct {
	Data []*CountryInfo
}

func (c *CountryInfoClient) MapByName(name string) *CountryInfo {
	for _, row := range c.Data {
		// check Name field
		if strings.ToLower(row.Name) == strings.ToLower(name) {
			return row
		}

		// check AlternateNames field
		if stringInSlice(strings.ToLower(name), row.AlternateNamesLower()) {
			return row
		}
	}
	return nil
}

func (c *CountryInfoClient) MapByAlpha2(alpha2 string) *CountryInfo {
	for _, row := range c.Data {
		if strings.ToLower(row.Alpha2) == strings.ToLower(alpha2) {
			return row
		}
	}
	return nil
}

func (c *CountryInfoClient) MapByAlpha3(alpha3 string) *CountryInfo {
	for _, row := range c.Data {
		if strings.ToLower(row.Alpha3) == strings.ToLower(alpha3) {
			return row
		}
	}
	return nil
}

func (c *CountryInfoClient) MapByCurrency(currency string) []*CountryInfo {
	rowList := []*CountryInfo{}
	for _, row := range c.Data {
		if stringInSlice(strings.ToLower(currency), row.CurrencyLower()) {
			rowList = append(rowList, row)
		}
	}
	return rowList
}

func (c *CountryInfoClient) MapByCallingCode(callingCode string) []*CountryInfo {
	rowList := []*CountryInfo{}
	for _, row := range c.Data {
		if stringInSlice(strings.ToLower(callingCode), row.CallingCodeLower()) {
			rowList = append(rowList, row)
		}
	}
	return rowList
}

func (c *CountryInfoClient) MapByRegion(region string) []*CountryInfo {
	rowList := []*CountryInfo{}
	for _, row := range c.Data {
		if strings.ToLower(row.Region) == strings.ToLower(region) {
			rowList = append(rowList, row)
		}
	}
	return rowList
}

func (c *CountryInfoClient) MapBySubregion(subregion string) []*CountryInfo {
	rowList := []*CountryInfo{}
	for _, row := range c.Data {
		if strings.ToLower(row.Subregion) == strings.ToLower(subregion) {
			rowList = append(rowList, row)
		}
	}
	return rowList
}

type CountryInfo struct {
	Name           string
	AlternateNames []string
	Alpha2         string
	Alpha3         string
	Capital        string
	Currency       []string
	CallingCode    []string
	Region         string
	Subregion      string
}

func (c *CountryInfo) AlternateNamesLower() []string {
	updated := []string{}
	for _, alternateName := range c.AlternateNames {
		updated = append(updated, strings.ToLower(alternateName))
	}
	return updated
}

func (c *CountryInfo) CurrencyLower() []string {
	updated := []string{}
	for _, currency := range c.Currency {
		updated = append(updated, strings.ToLower(currency))
	}
	return updated
}

func (c *CountryInfo) CallingCodeLower() []string {
	updated := []string{}
	for _, callingCode := range c.CallingCode {
		updated = append(updated, strings.ToLower(callingCode))
	}
	return updated
}

func readCSVFromURL(fileURL string) ([][]string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// readCSVFromBytes reads csv data from a byte slice
func readCSVFromBytes(b []byte) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(b))
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Load reads the country info data from a csv file and returns a CountryInfoClient
//
// Pass in an optional url if you would like to use your own downloadable csv file for country's data.
// This is useful if you prefer to host the data file yourself or if you have modified some of the fields
// for your specific use case.
func Load(specifiedURL ...string) (*CountryInfoClient, error) {
	if len(specifiedURL) > 1 {
		// For some reason the original code allowed multiple urls to be passed
		// in, but only used the first one.
		//
		// For now, we'll just return an error if more than one url is passed
		// in. Previously, the code would just ignore all but the first url.
		return nil, fmt.Errorf("only one url can be passed in")
	}

	var data [][]string
	var err error

	// use user specified url for csv file if provided, else use default
	if len(specifiedURL) > 0 {
		data, err = readCSVFromURL(specifiedURL[0])
	} else {
		data, err = readCSVFromBytes(defaultCountryInfo)
	}

	if err != nil {
		return nil, err
	}

	recordList := []*CountryInfo{}
	for idx, row := range data {
		// skip header
		if idx == 0 {
			continue
		}

		// get name
		name := strings.Split(row[0], ",")[:1][0]

		// use commonly used & altSpellings names as AlternateNames
		alternateNames := strings.Split(row[0], ",")[1:]
		alternateNames = append(alternateNames, strings.Split(row[8], ",")...)

		record := &CountryInfo{
			Name:           name,
			AlternateNames: alternateNames,
			Alpha2:         row[2],
			Alpha3:         row[4],
			Capital:        row[7],
			Currency:       strings.Split(row[5], ","),
			CallingCode:    strings.Split(row[6], ","),
			Region:         row[10],
			Subregion:      row[11],
		}

		recordList = append(recordList, record)
	}

	return &CountryInfoClient{Data: recordList}, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
