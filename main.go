package main

import (
	"crypto/subtle"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"time"

	log "github.com/golang/glog"
	rpio "github.com/stianeikeland/go-rpio/v4"
	"golang.org/x/crypto/acme/autocert"
)

var domain = flag.String("domain", "", "the domain to get tls certs for")
var username = flag.String("username", "", "username to auth requests")
var password = flag.String("password", "", "password to auth requests")

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func openDoor() error {
	err := rpio.Open()
	if err != nil {
		return fmt.Errorf("rpio.Open: %w", err)
	}

	pin := rpio.Pin(10)
	pin.Output()
	pin.High()
	time.Sleep(3 * time.Second)
	pin.Low()

	rpio.Close()

	return nil
}

func main() {
	flag.Parse()

	if *domain == "" {
		log.Fatal("domain flag cannot be empty")
	}
	if *username == "" {
		log.Fatal("username flag cannot be empty")
	}
	if *password == "" {
		log.Fatal("password flag cannot be empty")
	}

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(*domain),
		Cache:      autocert.DirCache("certs"),
	}

	http.HandleFunc("/", BasicAuth(func(w http.ResponseWriter, r *http.Request) {
		go func() {
			err := openDoor()
			if err != nil {
				log.Errorf("openDoor: %+v", err)
			}
		}()
	}, *username, *password, "please enter username and password"))

	server := &http.Server{
		Addr: ":443", // 443 is required by letsencrypt(?)
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil)) // 80 is required by letsencrypt(?)

	log.Fatal(server.ListenAndServeTLS("", "")) // key and cert are coming from Let's Encrypt
}
