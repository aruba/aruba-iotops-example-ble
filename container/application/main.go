// Copyright 2021 Hewlett Packard Enterprise (HPE)
//
// Licensed under the MIT License;
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.github.com/aruba-iotops-example-ble/LICENSE
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
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

	// mqtt client
	mqttClient := NewMqttClient()
	mqttClient.Connect()

	// bleAPIURL: example app will get data from HPE IoT Operations infrastructure services through this API url
	bleAPIURL := "http://" + apiGwURL + "/api/v3/ble/stream/packets"
	httpClient := NewHTTPClient(bleAPIURL, apiKey, http.MethodGet)
	httpClient.Connect(context.Background())

	ProcessBleData(httpClient.GetDataCh(), mqttClient.GetPubDataCh())
}

const minBleDataLen = 30

type IBeaconData struct {
	DeviceClass string
	UUID        string
	Major       string
	Minor       string
	Power       string
}

// ProcessBleData get data from HPE IoT Operations infrastructure service.
// then decode and decorate and put data into data channel,
// data channel will be consumed by MQTT client.
func ProcessBleData(httpDataCh <-chan *BleData, mqttDataCh chan<- string) {
	for bleData := range httpDataCh {
		payload := bleData.Result.Payload
		if len(payload) < minBleDataLen {
			continue
		}

		// below is to convert iBeacon byte data to iBeacon string data.
		// you need to overwrite this code when you decode your device data.
		// note: field "data" is hexadecimal byte array.
		// If you want to get string data. please process it with method hex.EncodeToString([]byte)
		iBeaconData := &IBeaconData{
			DeviceClass: "iBeacon",
			UUID:        hex.EncodeToString(payload[9:25]),
			Major:       hex.EncodeToString(payload[25:27]),
			Minor:       hex.EncodeToString(payload[27:29]),
			Power:       hex.EncodeToString(payload[29:30]),
		}
		iBeacon, _ := json.Marshal(iBeaconData)

		log.Println("iBeacon uuid: " + iBeaconData.UUID)

		mqttDataCh <- string(iBeacon)
	}
}
