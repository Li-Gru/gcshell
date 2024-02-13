package internal

import (
	"crypto/tls"
	"errors"
	"github.com/ThalesIgnite/crypto11"
)

type Token struct {
	Context *crypto11.Context
}

func (t *Token) Init(library *string, serial *string, pin *string) (*Token, error) {
	config := &crypto11.Config{
		Path:        *library,
		TokenSerial: *serial,
		Pin:         *pin,
	}
	context, err := crypto11.Configure(config)
	if err != nil {
		return nil, err
	}
	t.Context = context
	return t, nil
}

func (t *Token) GetCertificate(selector string) (*tls.Certificate, error) {
	var ErrEmptyCert = errors.New("empty cert found")
	certificates, err := t.Context.FindAllPairedCertificates()
	if err != nil {
		return nil, err
	}
	if len(certificates) == 0 {
		return nil, ErrEmptyCert
	}

	if selector == "" {
		return &certificates[0], nil
	}

	for _, certificate := range certificates {
		if contains(certificate.Leaf.EmailAddresses, selector) {
			return &certificate, nil
		}
	}
	return nil, nil
}
