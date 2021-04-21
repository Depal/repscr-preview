package wikiexport

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"
)

const (
	RootUrl string = "http://[СТЁРТО ОК]"
)

func ExportShiftToPdf(shiftName string) (filePath string, err error) {
	filePath, err = filepath.Abs("./" + shiftName + ".pdf")
	if err != nil {
		return "", err
	}

	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Debugf))
	defer cancel()

	err = chromedp.Run(ctx, login())
	if err != nil {
		return "", err
	}

	var buf []byte
	err = chromedp.Run(ctx, printToPDF(RootUrl+":"+shiftName, &buf))

	log.Info("Сохраняю PDF...")
	if err := ioutil.WriteFile(filePath, buf, 0644); err != nil {
		log.Fatal(err)
	}

	return filePath, nil
}

func GetLastUpdate(shiftName string) (lastUpdate time.Time, err error) {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Debugf))
	defer cancel()

	err = chromedp.Run(ctx, login())
	if err != nil {
		return time.Time{}, err
	}

	var docInfo string
	err = chromedp.Run(ctx,
		chromedp.Navigate(RootUrl+":"+shiftName),
		chromedp.Text(`//div[@class="docInfo"]`, &docInfo, chromedp.NodeVisible, chromedp.BySearch),
	)
	if err != nil {
		return time.Time{}, err
	}
	if docInfo == "" {
		return time.Time{}, errors.New("docInfo not found")
	}

	//log.Infof("Обнаружен docInfo: %v", docInfo)

	//log.Info("Преобразую docInfo в дату...")
	re, err := regexp.Compile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}`)
	if err != nil {
		return time.Time{}, err
	}

	timePart := re.FindString(docInfo)

	lastUpdate, err = time.ParseInLocation("2006/01/02 15:04", timePart, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return lastUpdate, nil
}

func login() chromedp.Tasks {
	log.Info("Вхожу на Wiki...")

	loginField := `//input[@name="u"]`
	passwordField := `//input[@name="p"]`

	return chromedp.Tasks{
		chromedp.Navigate(RootUrl),
		chromedp.WaitVisible(loginField),
		chromedp.WaitVisible(passwordField),
		chromedp.SendKeys(loginField, "[СТЁРТО ОК]"),
		chromedp.SendKeys(passwordField, `[СТЁРТО ОК]`),
		chromedp.Submit(passwordField),
		chromedp.WaitNotPresent(loginField),
	}
}

func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	log.Info("Генерирую PDF...")
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
