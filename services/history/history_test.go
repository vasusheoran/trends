package history

import (
	"github.com/vsheoran/trends/utils"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	logger := utils.InitializeDefaultLogger()

	testCases := []struct {
		input   string
		output  time.Time
		wantErr bool
	}{
		{"2-Jan-2006", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), false},
		{"02-Jan-2006", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), false},
		{"2-Jan-06", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), false},
		{"02-Jan-06", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), false},
		{"2-01-2006", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), false},
		{"02-1-2006", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), false},
		{"2-01-06", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), false},
		{"02-1-06", time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC), false},
		{"invalid_date", time.Time{}, true},
		{"2023-12-25", time.Time{}, true}, // Format not supported
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			svc := history{
				logger: logger,
			}
			parsedDate, err := svc.parseDate(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !parsedDate.Equal(tc.output) {
					t.Errorf("Expected %v, got %v", tc.output, parsedDate)
				}
			}
		})
	}
}
