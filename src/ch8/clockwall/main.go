/*
Exercise 8.1
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type ServerResponse struct {
	city string
	time string
	done bool
}

func main() {

	timeMap := make(map[string]string)
	cityList := []string{}

	channel := make(chan ServerResponse)

	serverCount := len(os.Args) - 1

	for _, server := range os.Args[1:] {
		strs := strings.Split(server, "=")
		city := strs[0]
		addr := strs[1]
		cityList = append(cityList, city)
		go readServer(city, addr, channel)
	}

	done := 0

	for response := range channel {
		if response.done {
			done++
			if done == serverCount {
				break
			}
		} else {
			timeMap[response.city] = response.time
			drawTime(cityList, timeMap)
		}
	}
}

func readServer(city string, addr string, channel chan<- ServerResponse) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		conn.Close()
		channel <- ServerResponse{
			city: city,
			time: "",
			done: true,
		}
	}()
	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			break
		}
		channel <- ServerResponse{
			city: city,
			time: scanner.Text(),
			done: false,
		}
	}
}

func drawTime(cityList []string, timeMap map[string]string) {
	fmt.Print("\r")
	for _, city := range cityList {
		time := timeMap[city]
		fmt.Printf("%s: %s | ", city, time)
	}
}
