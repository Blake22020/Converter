package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type result struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка .env")
	}

	currencies := []string{
		"WST", "ANG", "FKP", "MKD", "BOB", "BYN", "ILS", "LYD", "PKR", "RSD",
		"SBD", "TWD", "AWG", "BMD", "HRK", "ISK", "MZN", "UGX", "ARS", "LSL",
		"NAD", "SSP", "AZN", "BZD", "JMD", "KES", "MWK", "BGN", "BIF", "CNY",
		"IMP", "NOK", "SYP", "XCG", "BTN", "AED", "BDT", "MDL", "COP", "GBP",
		"GMD", "SLL", "BWP", "GYD", "JOD", "KWD", "BSD", "AMD", "DJF", "MYR",
		"NZD", "UZS", "BRL", "LRD", "MGA", "MVR", "NGN", "STN", "TOP", "VND",
		"EGP", "KMF", "KRW", "SEK", "CLP", "FJD", "IDR", "SAR", "YER", "AOA",
		"MNT", "RUB", "TJS", "TTD", "PAB", "CZK", "HKD", "KZT", "TND", "XAF",
		"CUP", "KGS", "RON", "UAH", "CVE", "PEN", "SHP", "DZD", "HTG", "HUF",
		"MXN", "OMR", "TVD", "CDF", "BND", "GNF", "GTQ", "SOS", "FOK", "MUR",
		"NIO", "ZMW", "ZWL", "USD", "CAD", "HNL", "INR", "TMT", "ZAR", "CHF",
		"GEL", "JEP", "LAK", "TZS", "BBD", "IQD", "XDR", "AFN", "MMK", "QAR",
		"SCR", "SDG", "SRD", "GGP", "GIP", "LBP", "MRU", "NPR", "PHP", "PLN",
		"SZL", "BAM", "CRC", "ERN", "KHR", "XCD", "AUD", "VUV", "XOF", "XPF",
		"SGD", "VES", "DKK", "ETB", "KID", "KYD", "PYG", "UYU", "BHD", "DOP",
		"GHS", "IRR", "PGK", "RWF", "SLE", "THB", "MOP", "ALL", "EUR", "JPY",
		"LKR", "MAD", "TRY",
	}

	var apikey = os.Getenv("API_KEY")

	var base string
	for {
		fmt.Print("Введите базовую валюту: ")
		_, err := fmt.Scanln(&base)
		if err != nil {
			fmt.Println("\nОшибка ввода, попробуйте ещё раз")
			continue
		}

		base = strings.ToUpper(base)

		is_valid := false
		for _, v := range currencies {
			if base == v {
				is_valid = true
			}
		}
		if !is_valid {
			fmt.Println("\nНет такой валюты =(")
			continue
		}

		break
	}

	data, err := http.Get("https://v6.exchangerate-api.com/v6/" + apikey + "/latest/" + base)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Ошибка открытия тела =(")
		}
	}(data.Body)

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal("Ошибка создания тела =(")
	}

	var r result
	if err := json.Unmarshal(body, &r); err != nil {
		log.Fatal("Ошибка открытия json =(")
	}

	var values = r.ConversionRates

	var get string
	var val float64
	var ok bool

	history, err := os.Create("history.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer history.Close()

	for {
		fmt.Println("Напишите help что бы увидеть все доступные команды")
		fmt.Print("Введите название валюты: ")
		_, err := fmt.Scanln(&get)
		if err != nil {
			fmt.Println("Попробуйте ещё раз =)")
			continue
		}
		get = strings.ToUpper(get)
		switch get {

		case "STOP":
			return
		case "HELP":
			fmt.Println("help: Выводит все доступные команды")
			fmt.Println("list: Выводит все доступные валюты")
			fmt.Println("stop: Остонавливает программу")
			continue
		case "LIST":
			for key, _ := range values {
				fmt.Println(key)
			}
			continue
		}
		val, ok = values[get]
		if !ok {
			fmt.Println("Попробуйте ещё раз =(")
			continue
		}
		fmt.Println(val)
		history.WriteString(get + "\t" + strconv.FormatFloat(val, 'f', -1, 64) + "\n")

	}
}
