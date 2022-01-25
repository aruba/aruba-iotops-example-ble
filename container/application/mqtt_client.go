// Copyright 2021 Hewlett Packard Enterprise (HPE)
//
// Licensed under the MIT License;
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.github.com/aruba-iotops-example-ble/LICENSE
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"crypto/tls"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const mqttVersion = 4

// NewMqttClient connect with MQTT broker.
// It will establish a MQTT connection to publish and subscribe to MQTT broker.
// It has two fields: PubDataCh、SubDataCh.
// PubDataCh: put data into this field, then data will be sent to MQTT broker.
// SubDataCh: get data from MQTT broker, then put data into this field.
func NewMqttClient(url string, userName string, password string,
	clientID string, pubTopic string, subTopic string) *MqttClient {
	return &MqttClient{
		URL:       url, // MQTT broker url
		userName:  userName,
		password:  password,
		ClientID:  clientID,
		PubTopic:  pubTopic,
		SubTopic:  subTopic,
		PubDataCh: make(chan string, 1),
		SubDataCh: make(chan []byte, 1),
	}
}

type MqttClient struct {
	URL       string
	userName  string
	password  string
	ClientID  string
	PubTopic  string
	SubTopic  string
	PubDataCh chan string
	SubDataCh chan []byte
}

// Connect establish a MQTT connection to publish and subscribe to MQTT broker.
func (c *MqttClient) Connect() {
	log.Default().Println("mqtt url : " + c.URL + " ; clientId :" +
		c.ClientID + " ; pubTopic: " + c.PubTopic + " ; subtopic: " + c.SubTopic)

	opts := mqtt.NewClientOptions().
		AddBroker(c.URL).
		SetClientID(c.ClientID).
		SetProtocolVersion(mqttVersion).
		SetTLSConfig(c.newTLSConfig()).
		SetUsername(c.userName).
		SetPassword(c.password).
		SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
			c.SubDataCh <- message.Payload()
		})
	opts.OnConnect = func(client mqtt.Client) {
		log.Default().Println("Mqtt Connected")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Default().Printf("Connect lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		log.Default().Println("mqtt connected fail! ")

		if token.Error() != nil {
			log.Default().Println(token.Error().Error())
		}

		// If connect failed, will retry after 1 second
		<-time.After(1 * time.Second)

		go c.Connect()

		return
	}

	c.subscribe(client)

	go c.publish(client)
}

func (c *MqttClient) GetPubDataCh() chan string {
	return c.PubDataCh
}

func (c *MqttClient) GetSubDataCh() chan []byte {
	return c.SubDataCh
}

func (c *MqttClient) publish(client mqtt.Client) {
	for data := range c.PubDataCh {
		token := client.Publish(c.PubTopic, 0, false, data)
		token.Wait()
	}
}

func (c *MqttClient) subscribe(client mqtt.Client) {
	token := client.Subscribe(c.SubTopic, 1, nil)
	token.Wait()
}

func (c *MqttClient) newTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
}
