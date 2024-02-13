package internal

import (
	"os"
)

type Config struct {
	Selector    string
	Pkcs11Id    string
	Pkcs11Lib   string
	Pkcs11Pin   string
	SnxGateway  string
	SnxPrefix   string
	SnxRealm    string
	SnxPath     string
	FingerPrint string
	ExtenderUrl string
}

func defaultString(key string, val string) string {
	if key != "" || len(key) > 0 {
		return key
	}
	return val
}
func (c *Config) fromArg() {
	// TODO
}

func (c *Config) fromEnv() {
	c.Selector = os.Getenv("CSHELL_PKCS11_SELECTOR")
	c.Pkcs11Id = os.Getenv("CSHELL_PKCS11_ID")
	c.Pkcs11Lib = os.Getenv("CSHELL_PKCS11_LIB")
	c.Pkcs11Pin = os.Getenv("CSHELL_PKCS11_PIN")
	c.SnxGateway = os.Getenv("CSHELL_SNX_GATEWAY")
	c.SnxPrefix = os.Getenv("CSHELL_SNX_PREFIX")
	c.SnxRealm = os.Getenv("CSHELL_SNX_REALM")
	c.SnxPath = os.Getenv("CSHELL_SNX_PATH")
	c.FingerPrint = os.Getenv("CSHELL_SNX_FINGERPRINT")
	c.ExtenderUrl = os.Getenv("CSHELL_EXTENDER_URL")
}

func (c *Config) defaults() {
	c.Selector = defaultString(c.Selector, "")
	c.SnxPrefix = defaultString(c.SnxPrefix, "/")
	c.SnxPath = defaultString(c.SnxPath, "/usr/bin/snx")
}

func (c *Config) Init() *Config {
	c.fromArg()
	c.fromEnv()
	c.defaults()
	return c
}
