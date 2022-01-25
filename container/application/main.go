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
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

// IoT Operations example ble app is a demonstration of how to interact with HPE IoT Operations infrastructure service.
// It will be used to instruct any partner using HPE IoT Operations service to develop their own app.
// This app use HiveMQ as cloud server. you can interact with this app through HiveMQ web page
// (http://www.hivemq.com/demos/websocket-client/).
// HiveMQ is only used to demonstrate this app, and you should replace it with your own service link in your own app.
func main() {
	// apiKey is a token to access HPE IoT Operations infrastructure service,
	// this property will be set by IoT Operations service during installing.
	// you only need to get this url from environment then access IoT Operations API through this apikey.
	// the value will be used in HTTP request (Header.Set("apikey", c.APIKey)).
	apiKey := os.Getenv("APIKEY")

	// apiGwURL is to access HPE IoT Operations infrastructure service,
	// this property will be set by IoT Operations service during installing.
	// you only need to get this url from environment then access IoT Operations API through HTTP request.
	// Example: http://apiGwUrl/$(api-method)
	apiGwURL := os.Getenv("APIGW_URL")

	// clientID is used to identify your MQTT connections.
	// The value should not be the same as the value in the MQTT web page
	// (MQTT web page : http://www.hivemq.com/demos/websocket-client/).
	clientID := strings.ReplaceAll(uuid.New().String(), "-", "")

	// MQTT url
	serverURL := "wss://test.mosquitto.org:8091/mqtt"

	// MQTT username/password
	userName := "rw"
	password := "readwrite"

	// MQTT publish data topic.
	// default topic name is "app2broker_topic".
	// IoT Operations data will be sent into this topic,
	// you can subscribe this topic in the MQTT web page.
	pubTopic := os.Getenv("APP_TO_BROKER_TOPIC")
	if pubTopic == "" {
		pubTopic = "app2broker_topic"
	}

	// MQTT subscribe data topic.
	// default topic name is "broker2app_topic"
	// you can send data into this topic through MQTT web page,
	// Example app will accept data from that topic.
	subTopic := os.Getenv("BROKER_TO_APP_TOPIC")
	if subTopic == "" {
		subTopic = "broker2app_topic"
	}

	log.Println("Example app start")

	// mqtt client
	mqttClient := NewMqttClient(serverURL, userName, password, clientID, pubTopic, subTopic)
	mqttClient.Connect()

	// bleAPIURL: example app will get data from HPE IoT Operations infrastructure services through this API url
	bleAPIURL := "http://" + apiGwURL + "/api/v2/ble/stream/packets"
	httpClient := NewHTTPClient(bleAPIURL, apiKey, http.MethodGet)
	httpClient.Connect(context.Background())

	// bleClient: process ble data
	bleClient := NewBleClient()
	bleClient.ProcessBleData(httpClient.GetDataCh(), mqttClient.GetPubDataCh())
}
