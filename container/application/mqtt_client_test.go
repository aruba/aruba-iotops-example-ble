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
	"strings"
	"testing"
	"time"
)

func TestMqttClient(t *testing.T) {
	t.Setenv("APP_TO_BROKER_TOPIC", "iotops_topic")
	t.Setenv("BROKER_TO_APP_TOPIC", "iotops_topic")
	mqttClient := NewMqttClient()
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
}
