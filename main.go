package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Values struct {
	LoadAvg     int
	Ramcap      int
	Ramcons     int
	Diskcap     int
	Diskcons    int
	NetworkBand int
	Network     int
}

func FillStruct(s *Values, array []string) {
	s.LoadAvg, _ = strconv.Atoi(array[0])
	if s.LoadAvg > 30 {
		fmt.Printf("Load Average is too high: %d", s.LoadAvg)
	}
	s.Ramcap, _ = strconv.Atoi(array[1])
	s.Ramcons, _ = strconv.Atoi(array[2])
	if del := float64(s.Ramcons) / float64(s.Ramcap); del > 0.8 {
		del *= 100
		delint := int(del)
		fmt.Printf("Memory usage too high: %d%%", delint)
	}
	s.Diskcap, _ = strconv.Atoi(array[3])
	s.Diskcons, _ = strconv.Atoi(array[4])
	if disk := float64(s.Diskcons) / float64(s.Diskcap); disk > 0.9 {
		disk = (float64(s.Diskcons) - float64(s.Diskcap)) / 1048576
		diskint := int(disk)
		fmt.Printf("Memory usage too high: %d%%", diskint)
	}
	s.NetworkBand, _ = strconv.Atoi(array[5])
	s.Network, _ = strconv.Atoi(array[6])
	if network := float64(s.Network) / float64(s.NetworkBand); network > 0.9 {
		networkint := int(network)
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available", networkint)
	}
}

func main() {
	client := resty.New()
	url := "http://srv.msk01.gigacorp.local/_stats"
	interval := 10 * time.Second
	i := 1
	for {
		value := Values{}
		resp, err := client.R().Get(url)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode() != http.StatusOK || resp.Header().Get("Content-Type") != "text/plain" {
			if i == 3 {
				fmt.Println("Unable to fetch server statistic")
				break
			} else {
				i++
				continue
			}
		}
		respBody := strings.Split(resp.String(), ",")
		FillStruct(&value, respBody)
		time.Sleep(interval)
	}
}
