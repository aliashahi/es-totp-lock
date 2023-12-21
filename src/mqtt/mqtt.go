package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"es-project/src/webserver"
	"fmt"
	"log"
	"math/rand"
	"os"
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

func Boot() {
	client := createMqttClient()
	subscribe(client)
}

func createMqttClient() mqtt.Client {
	connectAddress := fmt.Sprintf("%s://%s:%d", protocol, broker, port)
	// rand.Seed(time.Now().UnixNano())
	clientID := fmt.Sprintf("go-client-%d", rand.Int())

	webserver.Logger("connect address: %s", connectAddress)
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
		webserver.Logger("connecting to mqtt server")
		token := client.Connect()
		if token.WaitTimeout(3*time.Second) && token.Error() != nil {
			log.Println(token.Error())
		}

		if client.IsConnected() {
			webserver.Logger("connected to mqtt server")
			break
		}
		webserver.Logger("connecting failed, retrying. . .")
		time.Sleep(2 * time.Second)
	}

	return client
}

func Publish(client mqtt.Client, payload string) {
	qos := 0
	if token := client.Publish(pub_topic, byte(qos), false, payload); token.Wait() && token.Error() != nil {
		webserver.Logger("publish failed, topic: %s, payload: %s\n", pub_topic, payload)
		Publish(client, payload)
	} else {
		webserver.Logger("publish success, topic: %s, payload: %s\n", pub_topic, payload)
	}
}

func subscribe(client mqtt.Client) {
	qos := 0
	client.Subscribe(sub_topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		webserver.Logger("Received `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
		u, err := webserver.GetUserByPasscode(fmt.Sprint(msg.Payload()))
		if err != nil {
			webserver.Logger("authentication failed")
			webserver.Logger("%s fail - tried with code %s", time.Now().Format(time.RFC3339), msg.Payload())
			Publish(client, "0")
			return
		}

		webserver.Logger("%s success - user %s entered the room", time.Now().Format(time.RFC3339), u.Username)
		Publish(client, fmt.Sprintf("Hello %s", u.Username[:10]))
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
