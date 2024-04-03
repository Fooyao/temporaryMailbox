package main

import (
	"fmt"
	"github.com/alash3al/go-smtpsrv"
	"strings"
)

type mailContent struct {
	from    string
	to      string
	title   string
	content string
}

func handler(c *smtpsrv.Context) error {
	UserMail := c.To().String()
	UserMail = strings.Trim(UserMail, "<>")
	st := strings.Split(UserMail, "@")
	domain := st[1]

	// 检查域名是否在允许的域名列表中
	if !contains(allowedDomains, domain) {
		return fmt.Errorf("invalid domain")
	}

	// msg, _ := mail.ReadMessage(c)
	msg, _ := c.Parse()
	content := mailContent{
		from:    strings.Trim(c.From().String(), "<>"),
		title:   msg.Subject,
		content: msg.TextBody,
	}
	if mailBox[UserMail] == nil {
		mailBox[UserMail] = make([]mailContent, 0)
	}
	mailBox[UserMail] = append(mailBox[UserMail], content)
	return nil
}
func startSMTPServer(allowedDomains []string) {
	cfg := smtpsrv.ServerConfig{
		BannerDomain:    allowedDomains[0],
		ListenAddr:      ":25",
		MaxMessageBytes: 5 * 1024,
		Handler:         handler,
	}
	fmt.Println(smtpsrv.ListenAndServe(&cfg))
}
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
