package internal

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
)

type Extender struct {
	client *http.Client
}

func (e *Extender) Init(certificate tls.Certificate) *Extender {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{certificate},
				Renegotiation:      tls.RenegotiateOnceAsClient,
				InsecureSkipVerify: true,
			},
		},
	}
	e.client = client
	return e
}

func (e *Extender) GetParams(host string, prefix string, realm string) *map[string]string {
	var url string
	url = fmt.Sprintf("https://%s%s/Login/LoginWithCert?selectedRealm=%s", host, prefix, realm)
	resp, err := e.client.Get(url)
	Iferr(err)

	url = fmt.Sprintf("https://%s%s/SNX/extender", host, prefix)
	resp, err = e.client.Get(url)
	Iferr(err)
	bodyBuffer, err := io.ReadAll(resp.Body)
	Iferr(err)

	re := regexp.MustCompile(`/\* .*Extender.* \*/`)
	src := re.FindAllString(string(bodyBuffer), 1)[0]
	params := map[string]string{}
	for _, i := range strings.Split(src, ";") {
		if strings.Contains(i, "Extender.") {
			i = strings.Split(i, "Extender.")[1]
			m := strings.Split(i, " = ")
			k := strings.TrimSpace(m[0])
			v := strings.TrimSpace(m[1])
			params[k] = strings.Trim(v, `"`)

		}
	}
	return &params
}
