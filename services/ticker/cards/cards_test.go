package cards

import (
	"encoding/csv"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/vsheoran/trends/services/ticker/cards/models"
	"github.com/vsheoran/trends/services/ticker/ma"
	"github.com/vsheoran/trends/trendstest"
	"github.com/vsheoran/trends/utils"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	ceCount = 0
	brCount = 0
)

func TestNewCard(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	const symbol = "test"

	files := []string{
		//"test/input/9-12-24.csv",
		//"test/input/21-12-24.csv",
		"test/input/Nifty-CI.csv",
	}

	for _, file := range files {

		records, err := readInputCSV(file)
		if err != nil {
			t.Fatal(err)
		}

		data, err := parseRecords(logger, records)
		if err != nil {
			t.Fatal(err)
		}

		cardSvc := getCardService(logger)

		for i, expected := range data {

			ticker := models.Ticker{
				Date: expected.Date,
				Time: time.Now(),
				Name: symbol,
				W:    expected.W,
				X:    expected.X,
				Y:    expected.Y,
				Z:    expected.Z,
			}
			if i == 5500 {
				logger.Log(
					"index", i,
				)
			}

			err = cardSvc.Add(ticker)
			if err != nil {
				t.Fatal(err)
			}

			result := cardSvc.Get(symbol)
			validateResult(t, logger, i, expected, result[0])

			if i == 5500 {
				logger.Log(
					"CI", fmt.Sprintf("%.2f", result[0].CI),
					"CD", fmt.Sprintf("%.2f", result[0].CD),
					"CH", fmt.Sprintf("%.2f", result[0].CH),
				)
			}

		}
	}
}

func getCardService(logger log.Logger) *card {

	emaData := map[string]*ma.EMAConfig{
		"M5": {
			Window:   5,
			Delay:    0,
			Decay:    2.0 / 6.0,
			Capacity: 20,
		},
		"AS5": {
			Window:   5,
			Delay:    0,
			Decay:    2.0 / 6.0,
			Capacity: 20,
		},
		"O21": {
			Window:   5,
			Delay:    20,
			Decay:    2.0 / 21.0,
			Capacity: 50,
		},
		"BN21": {
			Window:   5,
			Delay:    0,
			Decay:    2.0 / 21.0,
			Capacity: 50,
		},
		"CD5": {
			Window:   5,
			Delay:    0,
			Decay:    2.0 / 6.0,
			Capacity: 20,
		},
	}

	maData := map[string]*ma.MAConfig{
		"AR10": {
			Window:   10,
			Capacity: 20,
		},
		"AR50": {
			Window:   50,
			Offset:   0,
			Capacity: 100,
		},
	}

	return &card{
		logger: logger,
		ticker: make(map[string]*tickerData),
		ema:    ma.NewExponentialMovingAverageV2(logger, emaData),
		ma:     ma.NewMovingAverageV2(logger, maData),
	}
}

func validateResult(t *testing.T, logger log.Logger, index int, expected, actualData models.Ticker) {
	if expected.W > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.W, expected.W, 0.1), fmt.Sprintf("actualW: %f, expected: %f, diff: %f, date: %s", actualData.W, expected.W, math.Abs(actualData.W-expected.W), actualData.Date))
	}

	if expected.X > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.X, expected.X, 0.1), fmt.Sprintf("actualX: %f, expected: %f, diff: %f, date: %s", actualData.X, expected.X, math.Abs(actualData.X-expected.X), actualData.Date))
	}

	if expected.Y > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.Y, expected.Y, 0.1), fmt.Sprintf("actualY: %f, expected: %f, diff: %f, date: %s", actualData.Y, expected.Y, math.Abs(actualData.Y-expected.Y), actualData.Date))
	}

	if expected.Z > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.Z, expected.Z, 0.1), fmt.Sprintf("actualZ: %f, expected: %f, diff: %f, date: %s", actualData.Z, expected.Z, math.Abs(actualData.Z-expected.Z), actualData.Date))
	}

	if expected.AD > 0 && strings.Contains(expected.Date, "25-Oct-2005,12-Apr-2013,6-May-2014,6-Dec-2017") {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.AD, expected.AD, 0.99), fmt.Sprintf("actualAD: %f, expected: %f, diff: %f, date: %s", actualData.AD, expected.AD, math.Abs(actualData.AD-expected.AD), actualData.Date))
	}

	if expected.M > 0.0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.M, expected.M, 0.99), fmt.Sprintf("actualM: %f, expected: %f, diff: %f, date: %s", actualData.M, expected.M, math.Abs(actualData.M-expected.M), actualData.Date))
	}

	if expected.AS > 0.0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.AS, expected.AS, 0.99), fmt.Sprintf("actualAS: %f, expected: %f, diff: %f, date: %s", actualData.AS, expected.AS, math.Abs(actualData.AS-expected.AS), actualData.Date))
	}

	if expected.O > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.O, expected.O, 0.99), fmt.Sprintf("actualO: %f, expected: %f, diff: %f, date: %s", actualData.O, expected.O, math.Abs(actualData.O-expected.O), actualData.Date))
	}
	if expected.BN > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.BN, expected.BN, 0.99), fmt.Sprintf("actualBN: %f, expected: %f, diff: %f, date: %s", actualData.BN, expected.BN, math.Abs(actualData.BN-expected.BN), actualData.Date))
	}

	if expected.BP > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.BP, expected.BP, 0.99), fmt.Sprintf("actualBP: %f, expected: %f, diff: %f, date: %s", actualData.BP, expected.BP, math.Abs(actualData.BP-expected.BP), actualData.Date))
	}

	if expected.AR > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.AR, expected.AR, 0.99), fmt.Sprintf("actualAR: %f, expected: %f, diff: %f, date: %s", actualData.AR, expected.AR, math.Abs(actualData.AR-expected.AR), actualData.Date))
	}

	if expected.C > 0.1 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.C, expected.C, 0.99), fmt.Sprintf("actualC: %f, expected: %f, diff: %f, date: %s", actualData.C, expected.C, math.Abs(actualData.C-expected.C), actualData.Date))
	}

	if expected.E > 0 {
	}

	if expected.D > 0 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.D, expected.D, 0.99), fmt.Sprintf("actualD: %f, expected: %f, diff: %f, date: %s", actualData.D, expected.D, math.Abs(actualData.D-expected.D), actualData.Date))
	}

	//if expected.EB > 0 {
	//	assert.True(t, trendstest.IsValueWithinTolerance(actualData.EB, expected.EB, 0.5), fmt.Sprintf("actualEB: %f, expected: %f, diff: %f, date: %s", actualData.EB, expected.EB, math.Abs(actualData.EB-expected.EB), actualData.Date))
	//}

	if expected.CW > 0.1 && index > 623 {
		assert.True(t, trendstest.IsValueWithinTolerance(actualData.CW, expected.CW, 0.99), fmt.Sprintf("actualCW: %f, expected: %f, diff: %f, date: %s", actualData.CW, expected.CW, math.Abs(actualData.CW-expected.CW), actualData.Date))
	}
	//&& actualData.Date == "1-Nov-2024"
	//if actualData.Date == "31-Oct-2024" {
	//	ceCount++
	//	assert.True(t, test.IsValueWithinTolerance(actualData.CE, expected.CE, 0.01), fmt.Sprintf("actualCE: %f, expected: %f, diff: %f, date: %s", actualData.CE, expected.CE, math.Abs(actualData.CE-expected.CE), actualData.Date))
	//}

	//if expected.BR > 0.0 && index > 3594 {
	//	brCount++
	//	assert.True(t, trendstest.IsValueWithinTolerance(actualData.BR, expected.BR, 0.99), fmt.Sprintf("actualBR: %f, expected: %f, diff: %f, date: %s", actualData.BR, expected.BR, math.Abs(actualData.BR-expected.BR), actualData.Date))
	//}

	//if expected.CE > 0.0 && index > 3594 {
	//	assert.True(t, trendstest.IsValueWithinTolerance(actualData.CE, expected.CE, 0.99), fmt.Sprintf("actualCE: %f, expected: %f, diff: %f, date: %s", actualData.CE, expected.CE, math.Abs(actualData.CE-expected.CE), actualData.Date))
	//}
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
	var ticker models.Ticker

	ticker.Date = row[headers.Date]
	date, err := parseDate(ticker.Date)
	if err == nil {
		ticker.Time = date
	}

	data, err := strconv.ParseFloat(row[headers.W], 64)
	if err == nil {
		ticker.W = data
	}
	data, err = strconv.ParseFloat(row[headers.X], 64)
	if err == nil {
		ticker.X = data
	}
	data, err = strconv.ParseFloat(row[headers.Y], 64)
	if err == nil {
		ticker.Y = data
	}
	data, err = strconv.ParseFloat(row[headers.Z], 64)
	if err == nil {
		ticker.Z = data
	}
	data, err = strconv.ParseFloat(row[headers.AD], 64)
	if err == nil {
		ticker.AD = data
	}
	data, err = strconv.ParseFloat(row[headers.AR], 64)
	if err == nil {
		ticker.AR = data
	}
	data, err = strconv.ParseFloat(row[headers.AS], 64)
	if err == nil {
		ticker.AS = data
	}
	data, err = strconv.ParseFloat(row[headers.BN], 64)
	if err == nil {
		ticker.BN = data
	}
	data, err = strconv.ParseFloat(row[headers.BP], 64)
	if err == nil {
		ticker.BP = data
	}
	data, err = strconv.ParseFloat(row[headers.CW], 64)
	if err == nil {
		ticker.CW = data
	}
	data, err = strconv.ParseFloat(row[headers.BR], 64)
	if err == nil {
		ticker.BR = data
	}
	data, err = strconv.ParseFloat(row[headers.CC], 64)
	if err == nil {
		ticker.CC = data
	}
	data, err = strconv.ParseFloat(row[headers.CE], 64)
	if err == nil {
		ticker.CE = data
	}
	data, err = strconv.ParseFloat(row[headers.CH], 64)
	if err == nil {
		ticker.CH = data
	}
	data, err = strconv.ParseFloat(row[headers.E], 64)
	if err == nil {
		ticker.E = data
	}
	data, err = strconv.ParseFloat(row[headers.C], 64)
	if err == nil {
		ticker.C = data
	}
	data, err = strconv.ParseFloat(row[headers.D], 64)
	if err == nil {
		ticker.D = data
	}
	data, err = strconv.ParseFloat(row[headers.O], 64)
	if err == nil {
		ticker.O = data
	}
	data, err = strconv.ParseFloat(row[headers.M], 64)
	if err == nil {
		ticker.M = data
	}
	data, err = strconv.ParseFloat(row[headers.CD], 64)
	if err == nil {
		ticker.CD = data
	}
	data, err = strconv.ParseFloat(row[headers.DK], 64)
	if err == nil {
		ticker.DK = data
	}
	data, err = strconv.ParseFloat(row[headers.EC], 64)
	if err == nil {
		ticker.EC = data
	}
	data, err = strconv.ParseFloat(row[headers.EB], 64)
	if err == nil {
		ticker.EB = data
	}
	data, err = strconv.ParseFloat(row[headers.CH], 64)
	if err == nil {
		ticker.CH = data
	}

	return ticker
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
		case "CH":
			index.CH = i
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
		case "CI":
			index.EB = i
		case "CJ":
			index.EB = i
		case "CK":
			index.EB = i
		case "CL":
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
	CH int `json:"CH"`

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
