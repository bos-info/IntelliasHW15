package main

import (
	"testing"
	"time"
)

func TestFindTrains(t *testing.T) {

	tests := map[string]struct {
		arrStation string
		depStation string
		criteria   string
		want       Trains
		wantErr    error
	}{
		"allCorrect": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "price",
			want: Trains{{1177, 1902, 1929, 164.65, time.Date(0, time.January, 1, 10, 25, 00, 0, time.UTC), time.Date(0, time.January, 1, 16, 36, 00, 0, time.UTC)},
				{1178, 1902, 1929, 164.65, time.Date(0, time.January, 1, 10, 25, 0, 0, time.UTC), time.Date(0, time.January, 1, 16, 36, 00, 0, time.UTC)},
				{1141, 1902, 1929, 176.77, time.Date(0, time.January, 1, 10, 25, 00, 0, time.UTC), time.Date(0, time.January, 1, 16, 48, 00, 0, time.UTC)}},
			wantErr: nil,
		},
	}

	for _, test := range tests { //Обрабатывает каждое значение testData в сегменте.
		got, gotErr := FindTrains(test.depStation, test.arrStation, test.criteria)
		if len(got) != 3 || gotErr != nil {
			t.Errorf("ХЕРНЯ")
		}
	}

}
