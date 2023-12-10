package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	protocol  = "ssl"
	broker    = "f25aeeaa.ala.us-east-1.emqxsl.com" // broker address
	sub_topic = "es-project/esp/pub"                // define topic
	pub_topic = "es-project/server/pub"             // define topic
	username  = "test"                              // username for authentication
	password  = "test"                              // password for authentication
	port      = 8883                                // port of MQTT over TLS/SSL

)

func main() {
	client := createMqttClient()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		subscribe(client)
	}() // we use goroutine to run the subscription function
	wg.Wait()
}

func createMqttClient() mqtt.Client {
	connectAddress := fmt.Sprintf("%s://%s:%d", protocol, broker, port)
	// rand.Seed(time.Now().UnixNano())
	clientID := fmt.Sprintf("go-client-%d", rand.Int())

	fmt.Println("connect address: ", connectAddress)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectAddress)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(time.Second * 60)

	// Optional: set server CA
	opts.SetTLSConfig(loadTLSConfig("./emqxsl-ca.crt"))

	client := mqtt.NewClient(opts)
	for {
		fmt.Println("connecting to mqtt server")
		token := client.Connect()
		if token.WaitTimeout(3*time.Second) && token.Error() != nil {
			log.Println(token.Error())
		}

		if client.IsConnected() {
			break
		}
		fmt.Println("connecting failed, retrying. . .")
		time.Sleep(2 * time.Second)
	}

	return client
}

func publish(client mqtt.Client, payload string) {
	qos := 0
	if token := client.Publish(pub_topic, byte(qos), false, payload); token.Wait() && token.Error() != nil {
		fmt.Printf("publish failed, topic: %s, payload: %s\n", pub_topic, payload)
		publish(client, payload)
	} else {
		fmt.Printf("publish success, topic: %s, payload: %s\n", pub_topic, payload)
	}
}

func subscribe(client mqtt.Client) {
	qos := 0
	client.Subscribe(sub_topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
		if string(msg.Payload()) == "1234" {
			log.Println("authentication success")
			publish(client, "1")
		} else {
			log.Println("authentication failed")
			publish(client, "0")
		}
	})
}

func loadTLSConfig(caFile string) *tls.Config {
	// load tls config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	if caFile != "" {
		certpool := x509.NewCertPool()
		ca, err := os.ReadFile(caFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	return &tlsConfig
}
