package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/kodekoding/go-agent/lib"
)

var (
	secretKey string
	appName   string
	port      string
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	flag.StringVar(&secretKey, "secret", "", "secret key of audit-trail account")
	flag.StringVar(&appName, "name", "", "your app name")
	flag.StringVar(&port, "port", "8321", "setup your UDP port")
	flag.Parse()

	PORT := ":" + port

	if secretKey == "" {
		log.Fatalln("Please provide your Secret Key")
	}
	if appName == "" {
		log.Fatalln("please provide your App Name")
	}

	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()

	fmt.Println("Client Agent for service: " + appName + " with secret: " + secretKey)
	fmt.Println("UDP server is runing on " + PORT)

	maxBuffer := 60 * 1024 //60kB
	buffer := make([]byte, maxBuffer)
	rand.Seed(time.Now().Unix())

	dataList := lib.NewQueue()
	go watch(dataList)
	for {
		n, addr, _ := connection.ReadFromUDP(buffer)
		receiveData := buffer[0 : n-1]
		dataList.Insert(string(receiveData))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		_, err = connection.WriteToUDP([]byte("success receive data"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func watch(dataList *lib.List) {
	ticker := time.NewTicker(time.Duration(10) * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("before sent")
			dataList.Reverse()
			list := dataList.GetHead()
			i := 0
			for list != nil {
				fmt.Printf("%+v", lib.GetValue(list))
				list = lib.GetNext(list)
				dataList.SetHead(nil)
				fmt.Println()
				i++
			}
			fmt.Printf("%d data receive\n", i)
		}
	}
}
