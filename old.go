package main

//func fromArg() {
//	if len(os.Args) < 2 {
//		log.Info("Использовать так: \"repscr.exe -1\"")
//		log.Info("где -1 - сдвиг смены (0 = текущая)")
//		log.Info()
//		log.Info("Нажмите Enter для выхода...")
//		_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
//		if err != nil {
//			log.Fatal(err)
//		}
//		return
//	}
//
//	offset, err := strconv.Atoi(os.Args[1])
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = sendReport(offset)
//}

//func processReportPeriodical() (err error) {
//	now := time.Now()
//
//	if now.Hour() < 7 || now.Hour() > 22 || (now.Hour() > 10 && now.Hour() < 19) {
//		return nil
//	}
//
//	offset := 0
//	if now.Hour() == 9 || now.Hour() == 21 {
//		offset = -1
//	}
//
//	lastMessageId := viper.GetInt("LastMessageId")
//	lastSentAt := viper.GetTime("LastSentAt")
//	isRepost := now.Sub(lastSentAt).Hours() < 10
//
//	messageId, err := sendReport(offset, isRepost)
//	if err != nil {
//		return err
//	}
//
//	log.Info("Сохраняю результаты выгрузки...")
//	viper.Set("LastMessageId", messageId)
//	viper.Set("LastSentAt", time.Now())
//	err = viper.WriteConfig()
//	if err != nil {
//		return err
//	}
//
//	log.Info("Отчёт отправлен")
//
//	if isRepost {
//		log.Info("Т.к. произведено переразмещение отчёта, удаляю предыдущую выгрузку...")
//		err = telegram.DeleteMessageFromReportChannel(lastMessageId)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
