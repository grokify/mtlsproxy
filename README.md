# mTLS Proxy

This is a simple mTLS reverse proxy that handles mTLS with the client.

## Clone the Server

```
% git clone https://github.com/grokify/mtlsproxy
% cd mtlsproxy
```

## Generate Key and Certificate

Create test keys and certificates for server and client usage in non-interactive mode with 1 day expiration.

```
% openssl req -x509 -newkey rsa:4096 -keyout server_key.pem -out server_cert.pem -sha256 -days 1 -nodes -subj "/C=US/ST=California/L=Silicon Valley/O=Local/OU=Host/CN=localhost"
% openssl req -x509 -newkey rsa:4096 -keyout client_key.pem -out client_cert.pem -sha256 -days 1 -nodes -subj "/C=US/ST=California/L=Silicon Valley/O=Local/OU=Host/CN=client"
```

Ref: https://stackoverflow.com/a/10176685/1908967

## Configure and Start Server

```
% export MTLSP_SERVER_KEY_PATH=server_key.pem
% export MTLSP_SERVER_CERT_PATH=server_cert.pem
% export MTLSP_CLIENT_CA_PATHS=client_cert.pem
% export MTLSP_UPSTREAM_URL=http://example.com
% export MTLSP_PORT=8080
% go run main.go
2024/03/26 07:36:57 listen:  [::]:8080
```

## Make mTLS Request using cURL

```
% curl --cert client_cert.pem --key client_key.pem --cacert server_crt.pem  https://localhost:8080
```

Returns page from http://example.com.

## References

1. https://pkg.go.dev/net/http/httputil#ReverseProxy
1. https://gist.github.com/JalfResi/6287706
1. https://github.com/habibiefaried/mtls-tcp-proxy/tree/main
1. https://github.com/picatz/mtls-proxy
1. https://stackoverflow.com/questions/76684798/proxy-server-in-go-using-existing-net-conn
1. https://stackoverflow.com/questions/35390726/confirm-tls-certificate-while-performing-reverseproxy-in-golang
1. https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/
1. https://medium.com/trendyol-tech/golang-ile-custom-reverse-proxy-yapmak-7a4198fe86fc
1. https://stackoverflow.com/questions/63899700/how-to-stop-showing-target-url-in-reverseproxy-in-golang-using-newsinglehostreve
1. https://stackoverflow.com/questions/54385164/golang-reverse-proxy-multiple-target-urls-without-appending-subpaths
1. https://stackoverflow.com/questions/23164547/golang-reverseproxy-not-working
1. https://stackoverflow.com/questions/50694429/curl-with-client-certificate-authentication
1. https://stackoverflow.com/questions/70207857/reverse-proxy-using-go-to-cloud-run-instance
1. https://www.integralist.co.uk/posts/golang-reverse-proxy/
1. https://github.com/picatz/mtls-proxy/blob/main/pkg/proxy/server.go
1. https://www.reddit.com/r/golang/comments/zdsgon/octoproxy_simple_tcptls_proxy_support_mutual/
1. https://smallstep.com/hello-mtls/doc/combined/go/nginx-proxy
1. https://smallstep.com/hello-mtls/doc/client/curl
1. https://medium.com/@brucifi/working-with-tls-client-certificates-f67437a9aeb9
1. https://downey.io/notes/dev/curl-using-mutual-tls/