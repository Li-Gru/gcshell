package main

import (
	"fmt"
	"gcshell/internal"
	"os"
	"os/user"
)

func main() {
	config := (&internal.Config{}).Init()
	token, _ := (&internal.Token{}).Init(&config.Pkcs11Lib, &config.Pkcs11Id, &config.Pkcs11Pin)
	defer token.Context.Close()

	certificate, _ := token.GetCertificate(config.Selector)
	if certificate == nil {
		panic("no certs")
	}

	extender := (&internal.Extender{}).Init(*certificate)
	params, _ := extender.GetParams(config.SnxGateway, config.SnxPrefix, config.SnxRealm)
	err := token.Context.Close()
	if err != nil {
		panic(err)
	}

	snx := &internal.SNX{
		SnxPath: config.SnxPath,
		Params:  *params,
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
	snx.CallSNX()
}
