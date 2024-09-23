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
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// NewHTTPClient connect with HPE IoT Operations infrastructure service.
// It will establish a connection through Connect(),
// after connected, client will get response data and put data into field "dataCh",
// consumer can invoke GetDataCh() to get field "dataCh" then consume the response data.
func NewHTTPClient(url, apiKey, method string) *HTTPClient {
	return &HTTPClient{
		URL:    url,
		APIKey: apiKey,
		Method: method,
		dataCh: make(chan *BleData, 1),
	}
}

type HTTPClient struct {
	URL    string
	APIKey string
	Method string // HTTP Method: GET/POST/HEAD/OPTIONS/PUT/PATCH/DELETE/TRACE/CONNECT
	dataCh chan *BleData
}

// Connect establish an HTTP connection.
// response data will be put into filed "dataCh".
func (c *HTTPClient) Connect(ctx context.Context) {
	log.Println("Http request, url: " + c.URL + " ; apiKey: " + c.APIKey)

	req, _ := http.NewRequestWithContext(ctx, c.Method, c.URL, nil)
	req.Header.Set("apikey", c.APIKey)

	client := http.DefaultClient

	resp, err := client.Do(req)
	defer func() {
		if err != nil && resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		log.Println("HTTP request error!")
		// If connect failed, will retry after 1 second
		<-time.After(1 * time.Second)

		go c.Connect(ctx)

		return
	}

	scanner := bufio.NewScanner(resp.Body)

	go func() {
		defer func() {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
		}()

		for scanner.Scan() {
			//log.Println(scanner.Text())
			bleData := bleDataFromResult(scanner.Bytes())
			if bleData != nil {
				c.dataCh <- bleData
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println(err.Error())
		}
	}()
}

func bleDataFromResult(aResult []byte) *BleData {
	bleData := new(BleData)
	err := json.Unmarshal(aResult, bleData)
	if err != nil {
		return nil
	}
	return bleData
}

// GetDataCh returns the streaming ble data channel
func (c *HTTPClient) GetDataCh() <-chan *BleData {
	return c.dataCh
}

type BleData struct {
	Result struct {
		Mac            string `json:"mac"`
		ApMac          string `json:"apMac"`
		Payload        []byte `json:"payload"`
		Rssi           int    `json:"rssi"`
		FrameType      string `json:"frameType"`
		RadioMac       string `json:"radioMac"`
		MacAddressType string `json:"macAddressType"`
	} `json:"result"`
}
