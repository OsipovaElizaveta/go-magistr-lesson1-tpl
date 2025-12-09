package main

import (
	"fmt"
	"io"
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

		getValue := func(index int) int64 {
			number, _ := strconv.ParseInt(strings.TrimSpace(values[index]), 10, 64)
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

		currMemoryUsage := getValue(2) * 100 / getValue(1)

		if currMemoryUsage >= 80 {
			fmt.Printf("Memory usage too high: %v%%\n", currMemoryUsage)
		}

		if occupiedStorage*100/totalStorage >= 90 {
			fmt.Printf("Free disk space is too low: %v Mb left\n", (totalStorage-occupiedStorage)/(1<<20))
		}

		if throughput*100/bandwidth >= 90 {
			fmt.Printf("Network bandwidth usage high: %v Mbit/s available\n", (bandwidth-throughput)/1_000_000)
		}

		return
	}

	fmt.Print("Unable to fetch server statistic")
}
