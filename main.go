package main

import (
	"github.com/Depal/repscr/telegram"
	"github.com/Depal/repscr/utils"
	"github.com/Depal/repscr/wikiexport"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	tickFrequency = 2
)

func configure() (err error) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	log.Info("Загружаю файл конфигурации...")
	viper.SetConfigFile("config.yaml")
	viper.SetDefault("LastMessageId", 0)
	viper.SetDefault("LastSentAt", time.Time{})
	viper.SetDefault("LastDocInfoUpdate", time.Time{})
	viper.SetDefault("RepostNext", true)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
	//log.SetLevel(log.DebugLevel)
}

func sendReport(offset int, silent bool) (messageId int, err error) {
	log.Info("Определяю смену...")
	shiftName := utils.ShiftOffsetToShiftName(time.Now(), offset)

	log.Info("Подключаюсь к Wiki...")
	filePath, err := wikiexport.ExportShiftToPdf(shiftName)
	if err != nil {
		return 0, err
	}

	log.Info("Отправляю PDF в Telegram...")
	messageId, err = telegram.SendFileToReportChannel(filePath, "", silent)
	if err != nil {
		return 0, err
	}

	log.Info("Удаляю PDF с машины...")
	err = os.Remove(filePath)
	if err != nil {
		return 0, err
	}

	return messageId, nil
}

func reportWasUpdated(offset int) (wasReportUpdated bool, err error) {
	lastDocInfoUpdate := viper.GetTime("LastDocInfoUpdate")

	shiftName := utils.ShiftOffsetToShiftName(time.Now(), offset)
	currentDocInfoUpdate, err := wikiexport.GetLastUpdate(shiftName)
	if err != nil {
		return false, err
	}

	if currentDocInfoUpdate.After(lastDocInfoUpdate) {
		viper.Set("LastDocInfoUpdate", currentDocInfoUpdate)
		return true, nil
	} else {
		return false, nil
	}
}

func processReport(offset int, repost bool) (err error) {
	lastMessageId := viper.GetInt("LastMessageId")

	messageId, err := sendReport(offset, true)
	if err != nil {
		return err
	}

	log.Info("Сохраняю результаты выгрузки...")
	viper.Set("LastMessageId", messageId)
	viper.Set("LastSentAt", time.Now())
	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	log.Info("Отчёт отправлен")

	if repost {
		log.Info("Т.к. произведено переразмещение отчёта, удаляю предыдущую выгрузку...")
		err = telegram.DeleteMessageFromReportChannel(lastMessageId)
		if err != nil {
			return err
		}
	}

	return nil
}

func processReportUpdateBased() (err error) {
	log.Info("Проверяю наличие обновлений в отчёте...")
	now := time.Now()
	var offset int

	if now.Hour() == utils.SwitchNightToDay || now.Hour() == utils.SwitchDayToNight {
		offset = -1
	} else {
		offset = 0
		viper.Set("RepostNext", true)
	}

	reportWasUpdated, err := reportWasUpdated(offset)
	if err != nil {
		return err
	}

	if !reportWasUpdated {
		log.Info("Обновлений отчёта не обнаружено")
		return nil
	}

	log.Info("! Обнаружено обновление отчёта")

	lastSentAt := viper.GetTime("LastSentAt")
	if utils.NewShiftStarted(now, lastSentAt) {
		log.Info("Обнаружена пересменка. Предыдущий отчёт не будет заменён")
		viper.Set("RepostNext", false)
	} else {
		viper.Set("RepostNext", true)
	}

	err = processReport(offset, viper.GetBool("RepostNext"))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := configure()
	if err != nil {
		log.Fatal(err)
	}

	for true {
		err = processReportUpdateBased()
		if err != nil {
			if err.Error() == "docInfo not found" {
				log.Warn("Отчёт пуст, поэтому не будет отправлен")
			} else {
				log.Error(err)
			}

		}
		log.Info("Работаю в цикле. Пауза до следующего срабатывания...")
		time.Sleep(tickFrequency * time.Minute)
	}
}
