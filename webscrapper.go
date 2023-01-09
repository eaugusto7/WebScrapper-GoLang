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

func CreateFile(texto, nameFileCards string) {
	fmt.Printf("Writing to a file in Go lang\n")

	file, err := os.OpenFile(nameFileCards, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer file.Close()

	len, err := file.WriteString(texto)

	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	//fmt.Printf("\nFile Name: %s", file.Name())
	fmt.Printf("\nLength: %d bytes", len)
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

func getAllCards(nameFileCards, priceFileCards string, interator int) int {
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
	collector.Visit("https://www.pichau.com.br/hardware/placa-de-video?page=" + strconv.Itoa(interator))

	return interator + 1
}

func main() {

	nameFileCards := "namefileCards.txt"
	priceFileCards := "pricefileCards.txt"

	interator := 1

	_, err := os.Stat(nameFileCards)
	fmt.Println("ERRO: ", os.IsNotExist(err))

	if err == nil {
		e := os.Remove(nameFileCards)
		if e != nil {
			log.Fatal(e)
		}
	}

	_, err = os.Stat(priceFileCards)
	fmt.Println("ERRO: ", os.IsNotExist(err))

	if err == nil {
		e := os.Remove(priceFileCards)
		if e != nil {
			log.Fatal(e)
		}
	}

	for interator < 3 {
		interator = getAllCards(nameFileCards, priceFileCards, interator)
	}

	ReadFile(nameFileCards, priceFileCards)
}
