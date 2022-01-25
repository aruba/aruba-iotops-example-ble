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

func TestExampleApp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "SE request"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := SEServerMock()

			// http client
			log.Default().Println("Request Ble data and transfer it to a third party server")

			BleRequestURL := server.URL + "/api/v2/ble/stream/packets"
			httpClient := NewHTTPClient(BleRequestURL, "", http.MethodGet)
			httpClient.Connect(context.Background())

			mqttDataCh := make(chan string, 1)

			// bleClient process ble data
			go NewBleClient().ProcessBleData(httpClient.GetDataCh(), mqttDataCh)

			iBeaconData := &IBeaconData{}

			go func() {
				result := <-mqttDataCh
				_ = json.Unmarshal([]byte(result), iBeaconData)
			}()

			<-time.After(20 * time.Millisecond)

			if strings.ToUpper(iBeaconData.UUID) != "F7826DA64FA24E988024BC5B71E0893E" {
				t.Error("Get iBeacon data failed")
			}
		})
	}
}

func SEServerMock() *httptest.Server {
	log.Default().Println("start HTTP server. send Ble data to client.")
	// mock data
	var frameType BleFrameType = 3

	bleDataMock, _ := hex.DecodeString("0201041AFF4C000215F7826DA64FA24E988024BC5B71E0893E00000000C5")

	// ble structure: 0201041AFF 4C00 02 15 4152554EF94A3B869470706978210A00(UUID) 0000(Major) 0000(Minor) C8(power)
	bleData := &BleData{
		Mac:       "50:31:ad:02:5c:93",
		Data:      bleDataMock,
		Rssi:      -57,
		FrameType: &frameType,
		ApMac:     "11:22:33:44:55:66",
	}
	bleDataJSON, _ := json.Marshal(bleData)
	testData := "data:" + string(bleDataJSON) + "\n"

	// HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.EscapedPath() != "/api/v2/ble/stream/packets" {
			_, _ = fmt.Fprintf(writer, "Reqeust path error")
		}
		if request.Method != http.MethodGet {
			_, _ = fmt.Fprintf(writer, "Request method error")
		}

		flusher, _ := writer.(http.Flusher)

		_, _ = fmt.Fprint(writer, testData)
		flusher.Flush()

		time.Sleep(1 * time.Second)
	}))

	return server
}
