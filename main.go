package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	for range 60 {
		mainImpl()
	}
}

func mainImpl() {
	url := "http://srv.msk01.gigacorp.local/_stats"
	regex, _ := regexp.Compile(`^\d{1,19}(?:,\d{1,19}){6}$`)

	for range 3 {

		resp, err := http.Get(url)

		if err != nil {
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}

		bodyStr := string(body)

		if !regex.MatchString(bodyStr) {
			continue
		}

		values := strings.Split(bodyStr, ",")

		getValue := func(index int) float64 {
			number, _ := strconv.ParseFloat(strings.TrimSpace(values[index]), 64)
			return number
		}

		loadAverage := getValue(0)
		totalStorage := getValue(3)
		occupiedStorage := getValue(4)
		bandwidth := getValue(5)
		throughput := getValue(6)

		if loadAverage > 30 {
			fmt.Printf("Load Average is too high: %v\n", loadAverage)
		}

		currMemoryUsage := getValue(2) / getValue(1)

		if currMemoryUsage > 0.8 {
			fmt.Printf("Memory usage too high: %v%%\n", math.Floor(currMemoryUsage*100))
		}

		if occupiedStorage/totalStorage > 0.9 {
			fmt.Printf("Free disk space is too low: %v Mb left\n", math.Floor((totalStorage-occupiedStorage)/math.Pow(2, 20)))
		}

		if throughput/bandwidth > 0.9 {
			fmt.Printf("Network bandwidth usage high: %v Mbit/s available\n", math.Floor((bandwidth-throughput)/1_000_000))
		}

		return
	}

	fmt.Print("Unable to fetch server statistic")
}
