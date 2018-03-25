package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

const TIMEOUT = 30

var (
	APIURL     = "https://api.opsgenie.com"
	parameters = make(map[string]string)
)

func init() {
	apiKey := flag.String("apiKey", "", "api key")
	name := flag.String("name", "", "heartbeat name")
	apiURL := flag.String("apiUrl", "", "api url")

	flag.Parse()

	if *apiURL != "" {
		APIURL = *apiURL
	}

	parameters["apiKey"] = *apiKey
	parameters["name"] = *name
}

func getHTTPClient(seconds int) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(seconds))
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(seconds)))
				return conn, nil
			},
		},
	}
	return client
}

func main() {
	var buf, err = json.Marshal(parameters)
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(2)
	}

	requestBody := bytes.NewBuffer(buf)

	requestURL := fmt.Sprintf("%s/v2/heartbeats/%s/ping", APIURL, parameters["name"])

	request, err := http.NewRequest("POST", requestURL, requestBody)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(2)
	}

	request.Header.Add("Authorization", fmt.Sprintf("GenieKey %s", parameters["apiKey"]))

	client := getHTTPClient(TIMEOUT)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Fprintln(os.Stdout, "couldn't send heartbeat to opsgenie", err)
		os.Exit(2)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Couldn't read the response from opsgenie", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode > 199 && resp.StatusCode < 400 {
		fmt.Fprintln(os.Stdout, "OK - successfully sent heartbeat to opsgenie")
		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stdout, "Opsgenie response:"+string(body[:]))
		os.Exit(1)
	}

}
