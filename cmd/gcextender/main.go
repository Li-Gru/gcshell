package main

import (
	"encoding/json"
	"gcshell/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func getConfig(w http.ResponseWriter, r *http.Request) {
	var config (internal.Config)
	config.Init()

	// TODO: Merge values of two structs
	json.NewDecoder(r.Body).Decode(&config)

	token, err := (&internal.Token{}).Init(&config.Pkcs11Lib, &config.Pkcs11Id, &config.Pkcs11Pin)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer token.Context.Close()
	certificate, err := token.GetCertificate(config.Selector)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	extender := (&internal.Extender{}).Init(*certificate)
	params, err := extender.GetParams(config.SnxGateway, config.SnxPrefix, config.SnxRealm)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	response, _ := json.Marshal(params)
	w.Write(response)
}
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	r.Post("/extender", getConfig)
	http.ListenAndServe(":3000", r)
}
