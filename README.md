# IoT Operations Example Ble App

This is an example application for Hewlett Packard Enterprise (HPE) IoT Operations (IoTOps) project that guides you to develop your own applications.
<br> This example app covers the basic app functions: 
<br> 1) classifying BLE devices with iBeacon Device Class using a Lua script
<br> 2) taking data from IoT Operations collector and transferring the data to your own cloud service.
<br> In this app, we choose HiveMQ as the cloud service as an example.

Before developing your own app, please read the documents related to HPE IoTOps, You need to know how the app will interact with HPE IoTOps. More information please visit [https://help.central.arubanetworks.com/latest/documentation/online_help/content/allowlist/iot.htm](https://help.central.arubanetworks.com/latest/documentation/online_help/content/allowlist/iot.htm).

Contents
* [Topology](#topology)
* [Project structure](#project-structure)
  * [Lua](#lua)
  * [Container](#container)
    * [Access IoT Operations infrastructure service](#access-iot-operations-infrastructure-service)
* [Build and Application Onboard](#build-and-application-onboard)
  * [Container Build](#container-build)
  * [AppBundle Configuration](#appbundle-configuration)
  * [Installation](#installation)
  * [Data Display](#data-display)
* [License](#license)

## Topology
![topology](./resource/topology.jpg)

## Project structure
```
aruba-iotops-example-ble
    |-- container
    |   |-- application
    |   |-- Dockerfile
    |   |-- Makefile
    |-- lua
    |-- resource
    |-- README.md
    |-- VERSION
    |-- LICENSE
```

### lua
The `lua` directory contains a Lua script for BLE device data packets processing.
IoTOps collectors use the Lua scripts to parse BLE device packet data and classify BLE device into Device Classes.
Application developers need to define a `function decode(address, addressType, advType, elements)` function as the entrance of device data processing.
In this example we use a Lua script to decode iBeacon packets. Different devices with different data types and structures should have different scripts.

### container
The `container` directory holds the source code of a containerized program.
IoTOps runs third party applications in a containerized environment.
Application developers can define their own data processing logic and integrate their private solutions using within this container.

In this example, the application receives the subscribed device packets from collector's API Gateway, then forwards the packet to HiveMQ topic.

#### Access IoT Operations infrastructure service
Container programs can access collector through HTTP API as follows:
```go
func (c *httpClient) Connect() {
	req, _ := http.NewRequest("GET", "http://$(apiGwUrl)/api/v2/ble/stream/packets", nil)
	req.Header.Set("apikey", "")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer func() {
		resp.Body.Close()
	}()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(data))
	}
}
```

NOTE: More IoTOps API please visit [https://app.swaggerhub.com/apis/davix/aruba_iot_gateway_container_api/2.0](https://app.swaggerhub.com/apis/davix/aruba_iot_gateway_container_api/2.0)

## Build and Application Onboard
### Container Build
```
cd $(BASE_DIR)/container/
make docker
```

### AppBundle Configuration
You should provide appBundle config info as below to let us help you config,
- [required] app name: 
  for example: iot_operations_example_ble_app
- [required] app icon: 
  please design your own app icon
- [required] Summary: 
  app function summary
- [required] app description: 
  please give us a brief introduction of your app
- [optional] Author's website: 
- [required] app categories:   
  for example: example app、iBeacon
- [optional] Lua script: 
  If you use Lua script to classify device, you have to send your script file to us
  - [required]subscription list:
    every item have two value: match type、match value. you should tell us how to classify your device data. 
    for example: iBeacon data: 0201041AFF 4C0002 15F7826DA64FA24E988024BC5B71E0893E00000000C5. match type:VENDOR / match value: 4C0002
- [optional] container image path: 
  If you only use Lua scripts to interact with HPE IoT Operations platform, you don't need to supply container image path 
  - [optional]container environment: 
    if your container need to config environment variables, please give us the keys and values. format like: key:value
- [optional] Container Resource Usage CPU Min: 
  for example: 0.5
- [optional] Container Resource Usage CPU Max: 
  for example: 1
- [optional] Container Resource Usage Mem Min: 
  for example: 512Mi
- [optional] Container Resource Usage Mem Max: 
  for example: 1024Mi
- [optional] Container Subscription list: 
  container should config subscription list. for example: device class name (It already is set in Lua script)
  Match Type: DEVICE_CLASS  MATCH_VALUE: iBeacon
- [optional] Container Required Environment Variables
  If your app need to config some environment variables, please let us know.Every one variable we have below fields can be set
    - [required]Name: variable name
    - [optional]Description 
    - [optional]Default value
    - [optional]Validation: regex pattern that can be used to validate the environment variable
    - [required]Type: have two values: OPTIONAL、REQUIRED; you need to choose one of them
    - [optional]Placeholder
    - [optional]Supported Values: add supported values for this key 
- [optional]Default App Permission Allowed External URLS: 
  if you need to interact with app through url, you should config external url.
  for example: test.mosquitto.org
- [optional] Default App Permission Allowed Device Class Classifications: 
  device class name will be config
- [optional] Default App Permission Device Class Permissions
  device class name will be config

### Installation
After config AppBundle, you will see your app in your Aruba central. You only need to config the MQTT topic.
After the installation is complete, you can access [http://www.hivemq.com/demos/websocket-client/](http://www.hivemq.com/demos/websocket-client/) and subscribe MQTT data topic to get data.

### Data Display
Access MQTT web page [http://www.hivemq.com/demos/websocket-client/](http://www.hivemq.com/demos/websocket-client/), subscribe data topic as shown in the image below:
<br>
![appBundle web page](./resource/mqttclient.jpg)
<br>
On this web page, we will have the following operate step:
```
step 1: on Connection part, type "test.mosquitto.org" into Host field, then click "Connect" button  (It will establish a connection with MQTT broker)
step 2: on Subscriptions part, click "Add New Topic Subscription" button, then it will pop up a window, you should type your own public topic name into Topic field, then click "Subscribe" button. The topic name is set when the app is installed. If you haven't set this topic, we also have default value: "app2broker_topic".
step 3: on the "Messages" part, you will see the data from your app.
```

## License
[MIT LICENSE](./LICENSE)