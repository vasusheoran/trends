package cards

import (
	"encoding/csv"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/services/ticker/ma"
	"github.com/vsheoran/trends/test"
	"github.com/vsheoran/trends/utils"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestNewCard(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const ticker = "test"

	records, err := readInputCSV("test/input/4-12-24.csv")
	if err != nil {
		t.Fatal(err)
	}

	data, err := parseRecords(logger, records)
	if err != nil {
		t.Fatal(err)
	}

	cardSvc := NewCard(logger)

	for i, expected := range data {
		err = cardSvc.Add(ticker, expected.Date, expected.W, expected.X, expected.Y, expected.Z)
		if err != nil {
			t.Fatal(err)
		}

		actualData := cardSvc.Get(ticker)

		if i > 623 {
			validateResult(t, expected, actualData[0])
		}
	}
}

func getCardService(logger log.Logger) *card {

	emaData := map[string]*ma.EMAData{
		"M5": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 6.0,
			Values: []float64{},
			EMA:    []float64{},
		},
		"AS5": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 6.0,
			Values: []float64{},
			EMA:    []float64{},
		},
		"O21": {
			Window: 5,
			Delay:  20,
			Decay:  2.0 / 21.0,
			Values: []float64{},
			EMA:    []float64{},
		},
		"BN21": {
			Window: 5,
			Delay:  0,
			Decay:  2.0 / 21.0,
			Values: []float64{},
			EMA:    []float64{},
		},
	}

	maData := map[string]*ma.MAData{
		"AR10": {
			Values:    []float64{},
			WindowSum: []float64{},
			Window:    10,
		},
		"AR50": {
			Values:    []float64{},
			WindowSum: []float64{},
			Window:    50,
			Offset:    0,
		},
	}

	return &card{
		logger: logger,
		ticker: make(map[string]*tickerData),
		ema:    ma.NewExponentialMovingAverageV2(logger, emaData),
		ma:     ma.NewMovingAverageV2(logger, maData),
	}
}

func validateResult(t *testing.T, expected, actualData models.Ticker) {
	if expected.AD > 0 && strings.Contains(expected.Date, "25-Oct-2005,12-Apr-2013,6-May-2014,6-Dec-2017") {
		assert.True(t, test.IsValueWithinTolerance(actualData.AD, expected.AD, 0.1), fmt.Sprintf("actualAD: %f, expected: %f, diff: %f, date: %s", actualData.AD, expected.AD, math.Abs(actualData.AD-expected.AD), actualData.Date))
	}

	if expected.M > 0.0 {
		assert.True(t, test.IsValueWithinTolerance(actualData.M, expected.M, 0.9), fmt.Sprintf("actualM: %f, expected: %f, diff: %f, date: %s", actualData.M, expected.M, math.Abs(actualData.M-expected.M), actualData.Date))
	}

	if expected.AS > 0.0 {
		assert.True(t, test.IsValueWithinTolerance(actualData.AS, expected.AS, 0.5), fmt.Sprintf("actualAS: %f, expected: %f, diff: %f, date: %s", actualData.AS, expected.AS, math.Abs(actualData.AS-expected.AS), actualData.Date))
	}

	if expected.O > 0 {
		assert.True(t, test.IsValueWithinTolerance(actualData.O, expected.O, 0.5), fmt.Sprintf("actualO: %f, expected: %f, diff: %f, date: %s", actualData.O, expected.O, math.Abs(actualData.O-expected.O), actualData.Date))
	}
	if expected.BN > 0 {
		assert.True(t, test.IsValueWithinTolerance(actualData.BN, expected.BN, 0.5), fmt.Sprintf("actualBN: %f, expected: %f, diff: %f, date: %s", actualData.BN, expected.BN, math.Abs(actualData.BN-expected.BN), actualData.Date))
	}

	if expected.BP > 0 {
		assert.True(t, test.IsValueWithinTolerance(actualData.BP, expected.BP, 0.5), fmt.Sprintf("actualBP: %f, expected: %f, diff: %f, date: %s", actualData.BP, expected.BP, math.Abs(actualData.BP-expected.BP), actualData.Date))
	}

	if expected.AR > 0 {
		assert.True(t, test.IsValueWithinTolerance(actualData.AR, expected.AR, 0.5), fmt.Sprintf("actualAR: %f, expected: %f, diff: %f, date: %s", actualData.AR, expected.AR, math.Abs(actualData.AR-expected.AR), actualData.Date))
	}

	if expected.C > 0.1 {
		assert.True(t, test.IsValueWithinTolerance(actualData.C, expected.C, 0.1), fmt.Sprintf("actualC: %f, expected: %f, diff: %f, date: %s", actualData.C, expected.C, math.Abs(actualData.C-expected.C), actualData.Date))
	}

	if expected.E > 0 {
		assert.True(t, test.IsValueWithinTolerance(actualData.E, expected.E, 0.5), fmt.Sprintf("actualE: %f, expected: %f, diff: %f, date: %s", actualData.E, expected.E, math.Abs(actualData.E-expected.E), actualData.Date))
	}

	if expected.D > 0 {
		assert.True(t, test.IsValueWithinTolerance(actualData.D, expected.D, 0.5), fmt.Sprintf("actualD: %f, expected: %f, diff: %f, date: %s", actualData.D, expected.D, math.Abs(actualData.D-expected.D), actualData.Date))
	}

	if expected.CW > 0.1 {
		assert.True(t, test.IsValueWithinTolerance(actualData.CW, expected.CW, 0.5), fmt.Sprintf("actualCW: %f, expected: %f, diff: %f, date: %s", actualData.CW, expected.CW, math.Abs(actualData.CW-expected.CW), actualData.Date))
	}
}

func parseRecords(logger log.Logger, records [][]string) ([]models.Ticker, error) {
	var tickerData []models.Ticker
	var headers *tickerDataIndex
	for i, row := range records {
		if i == 0 {
			headers = parseHeaders(records)
			continue
		}
		data := getTickerData(row, headers)
		tickerData = append(tickerData, data)
	}

	return tickerData, nil
}

func getTickerData(row []string, headers *tickerDataIndex) models.Ticker {
	var tickerData models.Ticker

	tickerData.Date = row[headers.Date]
	date, err := parseDate(tickerData.Date)
	if err == nil {
		tickerData.Time = date
	}

	data, err := strconv.ParseFloat(row[headers.W], 64)
	if err == nil {
		tickerData.W = data
	}
	data, err = strconv.ParseFloat(row[headers.X], 64)
	if err == nil {
		tickerData.X = data
	}
	data, err = strconv.ParseFloat(row[headers.Y], 64)
	if err == nil {
		tickerData.Y = data
	}
	data, err = strconv.ParseFloat(row[headers.Z], 64)
	if err == nil {
		tickerData.Z = data
	}
	data, err = strconv.ParseFloat(row[headers.AD], 64)
	if err == nil {
		tickerData.AD = data
	}
	data, err = strconv.ParseFloat(row[headers.AR], 64)
	if err == nil {
		tickerData.AR = data
	}
	data, err = strconv.ParseFloat(row[headers.AS], 64)
	if err == nil {
		tickerData.AS = data
	}
	data, err = strconv.ParseFloat(row[headers.BN], 64)
	if err == nil {
		tickerData.BN = data
	}
	data, err = strconv.ParseFloat(row[headers.BP], 64)
	if err == nil {
		tickerData.BP = data
	}
	data, err = strconv.ParseFloat(row[headers.CW], 64)
	if err == nil {
		tickerData.CW = data
	}
	data, err = strconv.ParseFloat(row[headers.BR], 64)
	if err == nil {
		tickerData.BR = data
	}
	data, err = strconv.ParseFloat(row[headers.CC], 64)
	if err == nil {
		tickerData.CC = data
	}
	data, err = strconv.ParseFloat(row[headers.CE], 64)
	if err == nil {
		tickerData.CE = data
	}
	data, err = strconv.ParseFloat(row[headers.ED], 64)
	if err == nil {
		tickerData.ED = data
	}
	data, err = strconv.ParseFloat(row[headers.E], 64)
	if err == nil {
		tickerData.E = data
	}
	data, err = strconv.ParseFloat(row[headers.C], 64)
	if err == nil {
		tickerData.C = data
	}
	data, err = strconv.ParseFloat(row[headers.D], 64)
	if err == nil {
		tickerData.D = data
	}
	data, err = strconv.ParseFloat(row[headers.O], 64)
	if err == nil {
		tickerData.O = data
	}
	data, err = strconv.ParseFloat(row[headers.M], 64)
	if err == nil {
		tickerData.M = data
	}
	data, err = strconv.ParseFloat(row[headers.CD], 64)
	if err == nil {
		tickerData.CD = data
	}
	data, err = strconv.ParseFloat(row[headers.DK], 64)
	if err == nil {
		tickerData.DK = data
	}
	data, err = strconv.ParseFloat(row[headers.EC], 64)
	if err == nil {
		tickerData.EC = data
	}
	data, err = strconv.ParseFloat(row[headers.EB], 64)
	if err == nil {
		tickerData.EB = data
	}

	return tickerData
}

func parseHeaders(records [][]string) *tickerDataIndex {
	var index tickerDataIndex
	if records == nil {
		return nil
	}

	for i, val := range records[0] {
		toLowerVal := strings.Trim(val, " ")

		switch toLowerVal {
		case "Date":
			index.Date = i
		case "W":
			index.W = i
		case "X":
			index.X = i
		case "Y":
			index.Y = i
		case "Z":
			index.Z = i
		case "AD":
			index.AD = i
		case "AR":
			index.AR = i
		case "AS":
			index.AS = i
		case "BN":
			index.BN = i
		case "BP":
			index.BP = i
		case "CW":
			index.CW = i
		case "BR":
			index.BR = i
		case "CC":
			index.CC = i
		case "CE":
			index.CE = i
		case "ED":
			index.ED = i
		case "E":
			index.E = i
		case "C":
			index.C = i
		case "D":
			index.D = i
		case "O":
			index.O = i
		case "M":
			index.M = i
		case "CD":
			index.CD = i
		case "DK":
			index.DK = i
		case "EC":
			index.EC = i
		case "EB":
			index.EB = i
		}
	}

	return &index
}

func readInputCSV(filename string) ([][]string, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func GetTestFiles(path string) ([]string, error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	testFiles := []string{}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".csv") {
			testFiles = append(testFiles, filepath.Join(path, e.Name()))
		}
	}
	return testFiles, nil
}

type tickerDataIndex struct {
	Date int `json:"Date"`

	W int `json:"W" description:"Close"`
	X int `json:"X" description:"Open"`
	Y int `json:"Y" description:"High"`
	Z int `json:"Z" description:"Low"`

	AD int `json:"AD" description:""`
	AR int `json:"AR"`
	AS int `json:"AS"`
	BN int `json:"BN"`
	BP int `json:"BP"`
	CW int `json:"CW"`
	BR int `json:"BR"`
	CC int `json:"CC"`
	CE int `json:"CE"`
	ED int `json:"ED"`

	E int `json:"E"`
	C int `json:"C"`
	D int `json:"D"`

	O  int `json:"O"`
	M  int `json:"M"`
	CD int `json:"CD"`
	DK int `json:"DK"`
	EC int `json:"EC"`
	EB int `json:"EB"`
}
