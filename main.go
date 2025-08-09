package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type result struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func main() {
	var key = "5b1e831d2eee621c6c93ee67"
	data, err := http.Get("https://v6.exchangerate-api.com/v6/" + key + "/latest/USD")
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(data.Body)

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		panic(err)
	}

	var r result
	if err := json.Unmarshal(body, &r); err != nil {
		panic(err)
	}

	var values = r.ConversionRates
	var get string
	var val float64
	var ok bool
	for {
		fmt.Println("Напишите help что бы увидеть все доступные команды")
		fmt.Print("Введите название валюты: ")
		_, err := fmt.Scanln(&get)
		if err != nil {
			fmt.Println("Попробуйте ещё раз")
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
			fmt.Println("Попробуйте ещё раз")
			continue
		}
		fmt.Println(val)
	}

}
