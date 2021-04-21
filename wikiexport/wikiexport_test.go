package wikiexport

import (
	"testing"
	"time"
)

func TestGetLastUpdate(t *testing.T) {
	lastUpdate, err := GetLastUpdate("16.04.2021_дневная")
	if err != nil {
		t.Fatal(err)
	}

	expected := time.Date(2021, 4, 16, 21, 4, 0, 0, time.Local)

	if lastUpdate != expected {
		t.Fatalf("%v != %v", lastUpdate, expected)
	}
}
