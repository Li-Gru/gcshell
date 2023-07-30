package internal

import (
	"crypto/tls"
	"github.com/ThalesIgnite/crypto11"
)

type Token struct {
	Context *crypto11.Context
}

func (t *Token) Init(library *string, serial *string, pin *string) *Token {
	config := &crypto11.Config{
		Path:        *library,
		TokenSerial: *serial,
		Pin:         *pin,
	}
	context, err := crypto11.Configure(config)
	Iferr(err)
	t.Context = context
	return t
}

func (t *Token) GetCertificate(selector string) *tls.Certificate {
	certificates, err := t.Context.FindAllPairedCertificates()
	Iferr(err)
	if len(certificates) == 0 {
		return nil
	}

	if selector == "" {
		return &certificates[0]
	}

	for _, certificate := range certificates {
		if contains(certificate.Leaf.EmailAddresses, selector) {
			return &certificate
		}
	}
	return nil
}
