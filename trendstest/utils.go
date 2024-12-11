package trendstest

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type TestData struct {
	Name           string    `json:"-"`
	U              []float64 `json:"u"`
	CI             []float64 `json:"ci"`
	ExpectedUEMA5  []float64 `json:"expected_u_ema_2_6"`
	ExpectedUEMA21 []float64 `json:"expected_u_ema_2_21"`
	ExpectedCPEMA5 []float64 `json:"expected_cp_ema_2_6_delay_50"`

	AD []float64 `json:"AD"`
	AR []float64 `json:"AR"`
	AS []float64 `json:"AS"`
	BN []float64 `json:"BN"`
	BP []float64 `json:"BP"`
	CW []float64 `json:"CW"`
	BR []float64 `json:"BR"`
	CC []float64 `json:"CC"`
	CE []float64 `json:"CE"`
	ED []float64 `json:"ED"`
}

func GetTestFiles(path string) []string {

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	testFiles := []string{}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".json") {
			testFiles = append(testFiles, filepath.Join(path, e.Name()))
		}
	}
	return testFiles
}

func GetTestCasesV2(path string, response interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	return nil
}

func IsValueWithinTolerance(actualValue, expectedValue, tolerance float64) bool {
	diff := math.Abs(actualValue - expectedValue)

	if diff > tolerance {
		return false
	}

	return true
}

func GetTestCases(path string) ([]TestData, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	testFiles := []string{}
	for _, e := range entries {
		if strings.Contains(e.Name(), ".json") {
			testFiles = append(testFiles, filepath.Join(path, e.Name()))
		}
	}

	var testCases []TestData
	var data []byte

	for _, fileName := range testFiles {
		data, err = os.ReadFile(fileName)
		if err != nil {
			return nil, err
		}

		var testData TestData
		err = json.Unmarshal(data, &testData)
		if err != nil {
			return nil, err
		}

		testData.Name = fileName
		testCases = append(testCases, testData)
	}

	return testCases, err
}
