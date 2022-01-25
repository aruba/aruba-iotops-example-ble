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
	"strings"
	"testing"
	"time"
)

func TestMqttClient(t *testing.T) {
	t.Parallel()

	mqttClient := NewMqttClient("wss://test.mosquitto.org:8091/mqtt", "rw", "readwrite",
		"random_client", "iotops_topic", "iotops_topic")

	tests := []struct {
		name       string
		mqttClient *MqttClient
	}{
		{name: "pub and sub", mqttClient: mqttClient},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mqttClient.Connect()

			// pub
			text := "iotops pub message : " + time.Now().String()
			mqttClient.GetPubDataCh() <- text

			// sub
			var data []byte

			go func() {
				data = <-mqttClient.GetSubDataCh()
			}()

			// verify
			<-time.After(500 * time.Millisecond)

			if !strings.Contains(string(data), "iotops pub message") {
				t.Error("mqtt publish data failed")
			}
		})
	}
}
