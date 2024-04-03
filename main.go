package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
	"time"
)

var mailBox = make(map[string][]mailContent)
var allowedDomains []string
var mainDomain string

func main() {
	domainsEnv := os.Getenv("ALLOWED_DOMAINS")
	mainDomain = os.Getenv("Main_DOMAIN")
	allowedDomains = strings.Split(domainsEnv, ",")
	go scheduleDailyMidnightTask(clearMailBox)
	go startHTTPServer(allowedDomains)
	startSMTPServer(allowedDomains)
}

func scheduleDailyMidnightTask(task func()) {
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	time.AfterFunc(nextMidnight.Sub(now), func() {
		task()
		scheduleDailyMidnightTask(task)
	})
}

func clearMailBox() {
	mailBox = make(map[string][]mailContent)
	fmt.Println("Mailbox cleared at:", time.Now().Format("2006-01-02 15:04:05"))
}
