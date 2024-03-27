package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/grokify/mogo/crypto/tlsutil"
	"github.com/grokify/mogo/net/netutil"
)

const (
	EnvUpstreamURL    = "MTLSP_UPSTREAM_URL"
	EnvClientCAPaths  = "MTLSP_CLIENT_CA_PATHS"
	EnvServerCertPath = "MTLSP_SERVER_CERT_PATH"
	EnvServerKeyPath  = "MTLSP_SERVER_KEY_PATH"
	EnvBindPort       = "MTLSP_PORT"
)

func main() {
	tlsConfig, err := tlsutil.NewTLSConfig(
		os.Getenv(EnvServerCertPath),
		os.Getenv(EnvServerKeyPath),
		[]string{},
		[]string{os.Getenv(EnvClientCAPaths)}, true,
	)
	if err != nil {
		log.Fatal(err)
	}

	origin, err := url.Parse(os.Getenv(EnvUpstreamURL))
	if err != nil {
		log.Fatal(err)
	}

	ln, err := tls.Listen("tcp", "0.0.0.0:"+os.Getenv(EnvBindPort), tlsConfig.Config)
	if err != nil {
		log.Fatalf("failed to create listener: %s", err)
	}

	log.Println("listen: ", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept conn: %s", err)
			continue
		}

		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			log.Printf("failed to cast conn to tls.Conn")
			continue
		}

		go func() {
			tag := fmt.Sprintf("[%s -> %s]", tlsConn.LocalAddr(), tlsConn.RemoteAddr())
			log.Printf("%s accept", tag)

			defer tlsConn.Close()

			// this is required to complete the handshake and populate the connection state
			// we are doing this so we can print the peer certificates prior to reading / writing to the connection
			err := tlsConn.Handshake()
			if err != nil {
				log.Printf("failed to complete handshake: %s", err)
				return
			}

			if len(tlsConn.ConnectionState().PeerCertificates) > 0 {
				log.Printf("%s client common name: %+v", tag, tlsConn.ConnectionState().PeerCertificates[0].Subject.CommonName)
			}

			err = netutil.ModifyConnectionRequest(conn, func(req *http.Request) error {
				req.Host = origin.Host
				req.URL.Host = origin.Host
				req.URL.Scheme = "http"
				return nil
			})
			if err != nil {
				log.Printf("failed to complete request: %s", err)
			}
		}()
	}
}

/*
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.Host = origin.Host
		req.URL.Host = origin.Host
		req.URL.Scheme = "http"
	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("HANDLE_FUNC")
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv(EnvBindPort), nil))
*/
