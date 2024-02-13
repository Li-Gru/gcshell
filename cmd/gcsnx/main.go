package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gcshell/internal"
	"io"
	"net/http"
	"os"
	"os/user"
)

func getParams() map[string]string {
	var config (internal.Config)
	params := map[string]string{}
	config.Init()
	jsonData, err := json.Marshal(config)
	req, err := http.NewRequest("POST", config.ExtenderUrl, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &params)
	return params
}

func main() {
	config := (&internal.Config{}).Init()
	params := getParams()

	// p := url_extender_marshal_json(&config)
	// json.NewDecoder(p).Decode(&params)

	snx := &internal.SNX{
		SnxPath: config.SnxPath,
		Params:  params,
		Debug:   false,
	}
	snx.GenerateSNXInfo()

	currentUser, err := user.Current()
	internal.Iferr(err)
	if currentUser.Uid == "0" {
		var fingerprint string
		if len(config.FingerPrint) > 0 {
			fingerprint = config.FingerPrint
		} else {
			fingerprint = snx.Params["server_fingerprint"]
		}

		file := fmt.Sprintf("/etc/snx/%s.db", "root")
		data := []byte(fmt.Sprintf("[%s]\n%s\n", snx.Params["server_cn"], fingerprint))
		err := os.WriteFile(file, data, 0644)
		internal.Iferr(err)
	}

	fmt.Println(snx.Params)
	snx.CallSNX()
}
