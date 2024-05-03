package service

import (
	"fmt"
	"github.com/go-mail/mail"
	"os"
	"strconv"
	"time"
)

func SendVerificationCode(code, toEmail string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "PinPals Admin"+"<"+os.Getenv("SMTP_USERNAME")+">")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verification Code for PinPals")
	messageId := "<" + strconv.FormatInt(time.Now().Unix(), 10) + "@dekun.me" + ">"
	m.SetHeader("Message-Id", messageId)
	m.SetBody("text/html", "Your verification code is: <br/> <b>"+code+"</b> <br/> <br/> This code will expire in 10 minutes.")
	d := mail.NewDialer(os.Getenv("SMTP_HOST"), 465, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Send verification code error: ", err)
		return err
	}
	return nil
}
