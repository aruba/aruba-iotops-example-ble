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
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const iBeaconRawHex = "0201041AFF4C000215F7826DA64FA24E988024BC5B71E0893E00000000C5"

func TestExampleApp(t *testing.T) {
	server := SEServerMock()

	// http client
	log.Println("Request Ble data and transfer it to a third party server")

	BleRequestURL := server.URL + "/api/v3/ble/stream/packets"
	httpClient := NewHTTPClient(BleRequestURL, "", http.MethodGet)
	httpClient.Connect(context.Background())

	mqttDataCh := make(chan string, 1)

	// bleClient process ble data
	go ProcessBleData(httpClient.GetDataCh(), mqttDataCh)

	iBeaconData := &IBeaconData{}

	go func() {
		result := <-mqttDataCh
		_ = json.Unmarshal([]byte(result), iBeaconData)
	}()

	<-time.After(20 * time.Millisecond)

	if strings.ToUpper(iBeaconData.UUID) != iBeaconRawHex[18:50] {
		t.Error("Get iBeacon data failed")
	}
}

func SEServerMock() *httptest.Server {
	log.Println("start HTTP server. send Ble data to client.")
	// mock data
	bleDataMock, _ := hex.DecodeString(iBeaconRawHex)
	payload := base64.StdEncoding.EncodeToString(bleDataMock)

	testData := fmt.Sprintf(`{"result":{"mac":"dc:a6:32:3f:1f:33","apMac":"ff:ff:2c:5d:94:9f","payload":"%s","rssi":-46,"frameType":"BLE_FRAME_TYPE_ADV_IND","radioMac":"ff:11:df:f5:ba:b1","macAddressType":"BLE_MAC_ADDRESS_TYPE_PUBLIC"}}`, payload)

	// HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.EscapedPath() != "/api/v3/ble/stream/packets" {
			_, _ = fmt.Fprintf(writer, "Reqeust path error")
		}
		if request.Method != http.MethodGet {
			_, _ = fmt.Fprintf(writer, "Request method error")
		}

		flusher, _ := writer.(http.Flusher)

		writer.Write(append([]byte(testData), []byte("\n")...))
		flusher.Flush()

		time.Sleep(1 * time.Second)
	}))

	return server
}
