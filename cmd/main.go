package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/kodekoding/go-agent/entity"
	"github.com/kodekoding/go-agent/lib"
	"github.com/tokopedia/enterpriseapp-audit/lib/client"
)

var (
	secret       string
	key          string
	appName      string
	port         string
	clientSpawn  uint
	timeout      uint
	serverHost   string
	maxInFlight  int
	timeInFlight int
	ctx          = context.Background()
	willBeSent   map[string]string
)

func main() {
	flag.StringVar(&secret, "secret", "", "secret of audit-trail account")
	flag.StringVar(&key, "key", "", "key of audit-trail account")
	flag.StringVar(&appName, "name", "", "your app name")
	flag.StringVar(&port, "port", "8321", "setup UDP port")
	flag.UintVar(&clientSpawn, "client-spawn", 100, "setup audit-trail client spawn")
	flag.UintVar(&timeout, "timeout", 10, "setup your timeout")
	flag.IntVar(&maxInFlight, "max_in_flight", 1000, "limit request that should be sent to server")
	flag.IntVar(&timeInFlight, "time_in_flight", 60, "limit time request that should be sent to server (in second)")
	flag.StringVar(&serverHost, "host", "", "setup audit-trail server host with port")
	flag.Parse()

	// Validation section
	if secret == "" || key == "" {
		log.Fatalln("Please provide audit-trail Secret or Key")
	}
	if appName == "" {
		log.Fatalln("please provide your App Name")
	}
	if serverHost == "" {
		log.Fatalln("Please provide audit-trail server host")
	}

	udpHost := ":" + port

	s, err := net.ResolveUDPAddr("udp4", udpHost)
	if err != nil {
		log.Println("error when resolve UDP address" + err.Error())
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		log.Println("error when resolve Listen UDP" + err.Error())
		return
	}

	defer connection.Close()
	atConfig := &entity.AuditTrailConfig{
		Key:         key,
		Secret:      secret,
		AppName:     appName,
		Host:        serverHost,
		ClientSpawn: clientSpawn,
		Timeout:     timeout,
	}

	initAuditTrail(atConfig)

	log.Println("Client Agent for service: " + appName + " ready, running on " + udpHost)

	maxBuffer := 60 * 1024 //60kB
	buffer := make([]byte, maxBuffer)
	rand.Seed(time.Now().Unix())

	dataList := lib.NewQueue()
	mutex := new(sync.Mutex)
	limitCounter := 0
	go watch(dataList, limitCounter, mutex)
	var activityData entity.ActivityLog

	for {
		if limitCounter == maxInFlight {
			processData(dataList, limitCounter, mutex)
		}
		limitCounter++
		n, addr, _ := connection.ReadFromUDP(buffer)
		receiveData := buffer[0:n]
		_ = json.Unmarshal(receiveData, &activityData)

		dataList.Insert(activityData)

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			log.Println("Exiting UDP server!")
			return
		}

		_, err = connection.WriteToUDP([]byte("success receive data"), addr)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func watch(dataList *lib.List, counter int, mutex *sync.Mutex) {
	ticker := time.NewTicker(time.Duration(timeInFlight) * time.Second)

	for {
		select {
		case <-ticker.C:
			processData(dataList, counter, mutex)
		}
	}
}

func processData(data *lib.List, counter int, mutex *sync.Mutex) {
	// reset the limit maxInFlight
	mutex.Lock()
	counter = 0

	data.Reverse()
	list := data.GetHead()
	i := 0
	for list != nil {
		value := lib.GetValue(list)
		valByte, _ := json.Marshal(value)
		err := json.Unmarshal(valByte, &willBeSent)
		if err != nil {
			log.Println("error when unmarshal to struct: " + err.Error())
		}
		trx := client.Start(ctx, willBeSent["element_id"])
		if willBeSent["payload"] != "" {
			trx.RecordPayload(willBeSent["payload"])
		}

		if willBeSent["payload"] == willBeSent["new_data"] {
			// if payload has exactly_same/equal data with new_data, then empty new_data field
			willBeSent["new_data"] = ""
		}
		trx.RecordEvent("log_data", willBeSent)
		trx.End()
		list = lib.GetNext(list)
		data.SetHead(nil)

		i++
	}
	mutex.Unlock()
	if i > 0 {
		log.Printf("\nsuccess sent %d data\n", i)
	}
}

func initAuditTrail(cfg *entity.AuditTrailConfig) {
	if _, err := client.NewClient(
		client.ConfigAppName(cfg.AppName),
		client.ConfigHostname(cfg.Host),
		client.ConfigKey(cfg.Key),
		client.ConfigSecret(cfg.Secret),
		client.ConfigSpawn(int(cfg.ClientSpawn)),
		client.ConfigTimeout(time.Duration(cfg.Timeout)*time.Second),
	); err != nil {
		log.Fatalln("Cannot init audit-trail client: ", err.Error())
	}

	log.Println("audit-trail is up")

}
