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
	"bufio"
	"context"
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
		dataCh: make(chan []byte, 1),
	}
}

type HTTPClient struct {
	URL    string
	APIKey string
	Method string // HTTP Method: GET/POST/HEAD/OPTIONS/PUT/PATCH/DELETE/TRACE/CONNECT
	dataCh chan []byte
}

// Connect establish an HTTP connection.
// response data will be put into filed "dataCh".
func (c *HTTPClient) Connect(ctx context.Context) {
	log.Default().Println("Http request, url: " + c.URL + " ; apiKey: " + c.APIKey)

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
		log.Default().Println("HTTP request error!")
		// If connect failed, will retry after 1 second
		<-time.After(1 * time.Second)

		go c.Connect(ctx)

		return
	}

	reader := bufio.NewReader(resp.Body)

	go func() {
		defer func() {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
		}()

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				log.Default().Println(err.Error())

				return
			}

			if len(line) > 0 {
				c.dataCh <- line
			}
		}
	}()
}

// GetDataCh return HTTP client filed "dataCh".
// HTTP response data will be put into this field "dataCh".
// method in "process_ble_data.go" file will consume this data from data channel.
func (c *HTTPClient) GetDataCh() <-chan []byte {
	return c.dataCh
}
