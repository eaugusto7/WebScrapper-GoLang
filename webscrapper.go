package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func CreateFile(texto, nameFile string) {
	file, err := os.OpenFile(nameFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer file.Close()
	file.WriteString(texto)
}

func ReadFile(nameFileCards, priceFileCards string) {
	readFileNameCards, err := os.Open(nameFileCards)

	if err != nil {
		fmt.Println(err)
	}

	readFilePriceCards, err := os.Open(priceFileCards)

	if err != nil {
		fmt.Println(err)
	}

	fileScannerName := bufio.NewScanner(readFileNameCards)
	fileScannerName.Split(bufio.ScanLines)

	fileScannerPrice := bufio.NewScanner(readFilePriceCards)
	fileScannerPrice.Split(bufio.ScanLines)

	var fileLinesName []string
	var fileLinesPrice []string

	for fileScannerName.Scan() {
		fileLinesName = append(fileLinesName, fileScannerName.Text())
	}

	for fileScannerPrice.Scan() {
		fileLinesPrice = append(fileLinesPrice, fileScannerPrice.Text())
	}

	readFileNameCards.Close()
	readFilePriceCards.Close()

	for i := 0; i <= len(fileLinesName)-1; i++ {
		fmt.Println(fileLinesName[i] + " | " + fileLinesPrice[i])
	}
}

func getAllCards(nameFileCards, priceFileCards, linkFileCards string, interator int) int {
	collector := colly.NewCollector()
	collector.UserAgent = "Go program"

	collector.OnHTML("h2", func(e *colly.HTMLElement) {
		if e.Attr("class") == "MuiTypography-root jss76 jss77 MuiTypography-h6" {
			CreateFile(e.Text+"\n", nameFileCards)
		}
	})

	collector.OnHTML("div[class]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "jss79" {
			CreateFile(e.Text+"\n", priceFileCards)
		}

		if e.Attr("class") == "MuiCardContent-root jss62" {
			if strings.Contains(e.Text, "Esgotado") {
				CreateFile("Esgotado\n", priceFileCards)
			}
		}
	})

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), "/placa-de-video") {
			CreateFile("https://www.pichau.com.br"+e.Attr("href")+"\n", linkFileCards)
		}
	})

	collector.Visit("https://www.pichau.com.br/hardware/placa-de-video?page=" + strconv.Itoa(interator))
	return interator + 1
}

func main() {
	nameFileCards := "namefileCards.txt"
	priceFileCards := "pricefileCards.txt"
	linkFileCards := "linkfileCards.txt"

	interator := 1

	_, err := os.Stat(nameFileCards)
	if err == nil {
		e := os.Remove(nameFileCards)
		if e != nil {
			log.Fatal(e)
		}
	}

	_, err = os.Stat(priceFileCards)
	if err == nil {
		e := os.Remove(priceFileCards)
		if e != nil {
			log.Fatal(e)
		}
	}

	_, err = os.Stat(linkFileCards)
	if err == nil {
		e := os.Remove(linkFileCards)
		if e != nil {
			log.Fatal(e)
		}
	}

	for interator < 50 {
		interator = getAllCards(nameFileCards, priceFileCards, linkFileCards, interator)
	}
}
