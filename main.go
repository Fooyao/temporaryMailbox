package main

import "os"

var mailBox = make(map[string][]mailContent)
var allowedDomains []string

func main() {
	domainsEnv := os.Getenv("ALLOWED_DOMAINS")
    	allowedDomains = strings.Split(domainsEnv, ",")
	go startHTTPServer(allowedDomains)
	startSMTPServer(allowedDomains)
}
