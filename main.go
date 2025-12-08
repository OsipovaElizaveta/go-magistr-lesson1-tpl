package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	for range 60 {
		mainImpl()
		time.Sleep(3 * time.Second)
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

		debug([]byte(strconv.Itoa(resp.StatusCode)))
		if err != nil {
			debug([]byte(err.Error()))
		}

		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}

		debug(body)

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
			fmt.Printf("Memory usage too high: %v%%\n", math.Floor(currMemoryUsage*100))
		}

		currStorageUsage := occupiedStorage / totalStorage

		if currStorageUsage > 0.9 {
			fmt.Printf("Free disk space is too low: %v Mb left\n", math.Floor((totalStorage-occupiedStorage)/math.Pow(2, 20)))
		}

		currNetworkLoad := throughput / bandwidth

		if currNetworkLoad > 0.9 {
			fmt.Printf("Network bandwidth usage high: %v Mbit/s available\n", math.Floor((bandwidth-throughput)/1_000_000))
		}

		return
	}

	fmt.Print("Unable to fetch server statistic")
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func debug(data []byte) {
	fileName := "debug.out"
	data = append(data, byte('\n'))
	if fileExists(fileName) {

		f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

		defer f.Close()

		f.Write(data)
	} else {
		os.WriteFile(fileName, data, 0644)
	}
}
