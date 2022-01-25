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

type BleFrameType int32

type BleData struct {
	// ble data
	Data      []byte
	FrameType *BleFrameType
	// device mac address
	Mac string
	// Received Signal Strength Indication
	Rssi int32
	// AP mac address
	ApMac string
}

type IBeaconData struct {
	DeviceClass string
	UUID        string
	Major       string
	Minor       string
	Power       string
}
