package email

import (
	"bytes"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func getCookedEmail(email *Email) *gomail.Message {
	m := gomail.NewMessage()

	// 生成 Message-ID
	id := uuid.New().String()
	messageDomain := "artalk.local" // Fallback domain

	if email.FromAddr != "" {
		if atIndex := strings.LastIndex(email.FromAddr, "@"); atIndex != -1 {
			messageDomain = email.FromAddr[atIndex+1:]
		}
	}
	messageID := "<" + id + "@" + messageDomain + ">"
	m.SetHeader("Message-ID", messageID)

	// 发送人
	m.SetHeader("From", m.FormatAddress(email.FromAddr, email.FromName))
	// 接收人
	m.SetHeader("To", email.ToAddr)
	// 抄送人
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	// 主题
	m.SetHeader("Subject", email.Subject)
	// 内容
	m.SetBody("text/html", email.Body)
	// 附件
	//m.Attach("./file.png")

	return m
}

func getEmailMineTxt(email *Email) string {
	emailBuffer := bytes.NewBuffer([]byte{})
	getCookedEmail(email).WriteTo(emailBuffer)
	return string(emailBuffer.Bytes()[:])
}
