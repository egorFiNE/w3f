package main

import (
    "strings"
    "strconv"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type ethBlockResponse struct {
	jsonrpc string `json:"jsonrpc"`
	Result string `json:"result"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <url>\n", os.Args[0])
		os.Exit(1)
	}

	url := os.Args[1]
	response, errString := fetchJson(url)
	if errString != "" {
		fmt.Fprintf(os.Stderr, "Error fetching data: %s\n", errString)
		os.Exit(1)
	}

    blockNumberString := strings.Replace(response.Result, "0x", "", -1)
    blockNumber, err := strconv.ParseUint(blockNumberString, 16, 32)
    if err != nil {
        panic(err)
    }
	fmt.Printf("%d\n", blockNumber);
}

func fetchJson(url string) (*ethBlockResponse, string) {
	body := []byte(`{
		"method": "eth_blockNumber",
		"id": 1,
        "jsonrpc": "2.0",
		"params": []
	}`)

    r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, "failed to make POST request"
	}
    r.Header.Add("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(r)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
	}

	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "failed to read response body"
	}

	var response ethBlockResponse
	err = json.Unmarshal(bodyResponse, &response)
	if err != nil {
		return nil, "failed to unmarshal JSON response"
	}

	return &response, ""
}
