package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/kodekoding/go-agent/entity"
)

var counter = 0

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a host:port string")
		return
	}
	CONNECT := arguments[1]

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text)
		times, err := strconv.Atoi(text[0 : len(text)-1])
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
		sendMultipleData(c, times)
		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Exiting UDP client!")
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Reply: %s\n", string(buffer[0:n]))
	}
}

func sendMultipleData(c *net.UDPConn, times int) {
	for i := 0; i < times; i++ {
		counter++
		data := &entity.ActivityLog{
			ElementID:      fmt.Sprintf("123%d", counter),
			NewData:        fmt.Sprintf(`{"id": %d, "data": "perubahan%d"}`, i, i+1),
			OldData:        fmt.Sprintf(`{"id": %d, "data": "perubahan%d"}`, i, i),
			DisplayMessage: `{"action": "merubah data"}`,
			UriPath:        "/shop/shopcore",
			ActivityName:   "Perubahan Data",
			TribeName:      "Enterprise IT",
			ElementName:    "Shop Core - Update",
		}
		dataByte, _ := json.Marshal(data)
		fmt.Printf("%s", dataByte)
		fmt.Println()
		fmt.Println()
		c.Write(dataByte)
	}
}
