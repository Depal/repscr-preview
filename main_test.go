package main

import (
	"github.com/Depal/repscr/telegram"
	"github.com/Depal/repscr/utils"
	"github.com/Depal/repscr/wikiexport"
	"os"
	"testing"
	"time"
)

func TestSendReport(t *testing.T) {
	shiftName := utils.ShiftOffsetToShiftName(time.Date(2020, 10, 15, 15, 0, 0, 0, time.Local), -1)

	filePath, err := wikiexport.ExportShiftToPdf(shiftName)
	if err != nil {
		t.Fatal(err)
	}

	messageID, err := telegram.SendFileToReportChannel(filePath, "тест", true)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Fatal(err)
	}

	err = telegram.DeleteMessageFromReportChannel(messageID)
	if err != nil {
		t.Fatal(err)
	}
}
