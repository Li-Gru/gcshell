package main

import (
	"fmt"
	"gcshell/internal"
	"os"
	"os/user"
)


func main()  {
	config := (&internal.Config{}).Init()
	token := (&internal.Token{}).Init(&config.Pkcs11Lib, &config.Pkcs11Id, &config.Pkcs11Pin)
	defer token.Context.Close()

	certificate := token.GetCertificate(config.Selector)
	if certificate == nil { panic("no certs") }

	extender := (&internal.Extender{}).Init(*certificate)
	params := extender.GetParams(config.SnxGateway, config.SnxPrefix, config.SnxRealm)
	token.Context.Close()

	snx := &internal.SNX{
		SnxPath: config.SnxPath,
		Params:  *params,
		Debug:   false,
	}
	snx.GenerateSNXInfo()

	currentUser, err := user.Current()
	internal.Iferr(err)
	if currentUser.Uid == "0" {
		file := fmt.Sprintf("/etc/snx/%s.db", "root")
		data := []byte(fmt.Sprintf("[%s]\n%s\n", snx.Params["server_cn"], snx.Params["server_fingerprint"]))
		err := os.WriteFile(file, data, 0644)
		internal.Iferr(err)
	}


	fmt.Println(snx.Params)
	snx.CallSNX()
}