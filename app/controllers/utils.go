package controllers

import (
	"myapp/app/utils"
	"strconv"
	"strings"

	"gopkg.in/gomail.v1"
)

func parseUintOrDefault(intStr string, _default uint64) uint64 {
	if value, err := strconv.ParseUint(intStr, 0, 64); err != nil {
		return _default
	} else {
		return value
	}
}

func parseIntOrDefault(intStr string, _default int64) int64 {
	if value, err := strconv.ParseInt(intStr, 0, 64); err != nil {
		return _default
	} else {
		return value
	}
}

func EmailProvider(email string) string {
	arrs := strings.Split(email, "@")
	domain := arrs[1]
	rules := map[string]string{
		"163.com":   "smtp.163.com",
		"qq.com":    "smtp.qq.com",
		"gmail.com": "mail.google.com",
	}
	provider, ok := rules[domain]
	if ok {
		return provider
	}
	return "http://mail." + domain
}

func sendMail(subject, content, to string, html bool) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "activation@inaiping.xyz")
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	if html {
		mail.SetBody("text/html", content)
	} else {
		mail.SetBody("text/plain", content)
	}
	provider := gomail.NewMailer("smtp.163.com", "mediatek_cd", "mtkcd100", 587)
	if err := provider.Send(mail); err != nil {
		utils.AssertNoError(err, "Send activation mail failed")
	}
}

// SendMail for activation mail
func SendMail(subject, content, to string) {
	sendMail(subject, content, to, true)
}

func SendMailTextPlain(subject, content, to string) {
	sendMail(subject, content, to, false)
}
