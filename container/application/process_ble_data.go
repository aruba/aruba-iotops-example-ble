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
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
)

const minBleDataLen = 30

// NewBleClient is to get data from HPE IoT Operations infrastructure service and decode data.
// It will send HTTP request to get data.
// After get response data, it will decode data and decorate data, then put data into data channel.
func NewBleClient() *BleClient {
	return &BleClient{}
}

type BleClient struct{}

// ProcessBleData get data from HPE IoT Operations infrastructure service.
// then decode and decorate and put data into data channel,
// data channel will be consumed by MQTT client.
func (c *BleClient) ProcessBleData(httpDataCh <-chan []byte, mqttDataCh chan<- string) {
	// the structure of data is : `data:{"key":"value"}`
	// below code will replace `data:` with "", leaving only json structured data.
	for data := range httpDataCh {
		if strings.Contains(string(data), "data:") {
			str := strings.ReplaceAll(string(data), "data:", "")
			str = strings.ReplaceAll(str, "\n", "")

			// below is to convert iBeacon byte data to iBeacon string data.
			// you need to overwrite this code when you decode your device data.
			// note: field "data" is hexadecimal byte array.
			// If you want to get string data. please process it with method hex.EncodeToString([]byte)
			bleData := &BleData{}
			_ = json.Unmarshal([]byte(str), bleData)

			if len(bleData.Data) < minBleDataLen {
				continue
			}

			iBeaconData := &IBeaconData{
				DeviceClass: "iBeacon",
				UUID:        hex.EncodeToString(bleData.Data[9:25]),
				Major:       hex.EncodeToString(bleData.Data[25:27]),
				Minor:       hex.EncodeToString(bleData.Data[27:29]),
				Power:       hex.EncodeToString(bleData.Data[29:30]),
			}
			iBeacon, _ := json.Marshal(iBeaconData)

			log.Default().Println("iBeacon uuid: " + iBeaconData.UUID)

			mqttDataCh <- string(iBeacon)
		}
	}
}
