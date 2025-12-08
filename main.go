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

		matched := regex.MatchString(bodyStr)

		if !matched {
			continue
		}

		var values []float64

		for _, value := range strings.Split(bodyStr, ",") {

			number, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)

			values = append(values, number)
		}

		loadAverage := values[0]
		totalMemory := values[1]
		occupiedMemory := values[2]
		totalStorage := values[3]
		occupiedStorage := values[4]
		bandwidth := values[5]
		throughput := values[6]

		if loadAverage > 30 {
			fmt.Printf("Load Average is too high: %v\n", values[0])
		}

		currMemoryUsage := occupiedMemory / totalMemory

		if currMemoryUsage > 0.8 {
			fmt.Printf("Memory usage too high: %v%%\n", currMemoryUsage*100)
		}

		currStorageUsage := occupiedStorage / totalStorage

		if currStorageUsage > 0.9 {
			fmt.Printf("Free disk space is too low: %v Mb", (totalStorage-occupiedStorage)/math.Pow(2, 20))
		}

		currNetworkLoad := throughput / bandwidth

		if currNetworkLoad > 0.9 {
			fmt.Printf("Network bandwidth usage high: %v Mbit/s available", (bandwidth-throughput)/math.Pow(2, 17))
		}

		return
	}

	fmt.Print("Unable to fetch server statistic")
}
