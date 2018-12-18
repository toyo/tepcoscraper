package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/sclevine/agouti"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println(os.Args[0] + ` -id [userid] -password [password]`)
		return
	}

	var (
		id       = flag.String(`id`, ``, `ID`)
		password = flag.String(`password`, ``, `PASSWORD`)
	)

	flag.Parse()
	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalf("ブラウザ(Selenium Webdriver)が見つかりません: %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatalf("ブラウザが開けません: %v", err)
	}

	if err := page.Navigate("https://www.kakeibo.tepco.co.jp/dk/aut/login/"); err != nil {
		log.Fatalf("ページが表示できません: %v", err)
	}

	page.FindByName(`id`).SendKeys(*id)
	page.FindByName(`password`).SendKeys(*password)
	page.FindByID(`idLogin`).Click()

	page.FindByID(`idNotEmptyImg_contents01.jpg`).Click()
	page.FindByID(`bt_time_view.jpg`).Click()

	datepattern := regexp.MustCompile(`20\d\d\/\d\d\/\d\d`)

	for i := 0; i < 30; i++ {
		html, err := page.HTML()
		if err != nil {
			log.Fatal(err)
		}

		lines := strings.Split(html, "\n")

		var date, data string
		for _, line := range lines {
			if idx := strings.Index(line, `var items = [["日次",`); idx != -1 {
				data = strings.SplitN(line[idx+23:], `]`, 2)[0]
			}
			if strings.Index(line, `の電気使用量`) != -1 {
				date = datepattern.FindString(line)
			}
		}
		fmt.Println(date + `,` + data)

		page.FindByID(`doPrevious`).Click()
	}
}
