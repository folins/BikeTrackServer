package service

import (
	"bytes"
	"crypto/tls"
	"text/template"

	"github.com/folins/biketrackserver"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailContains struct {
	Name string
	Code int
}

type SMTPService struct {
	SmtpHost    string
	SmtpPort    int
	SenderEmail string
	Password    string
}

func NewSMTPService(email, password, host string, port int) *SMTPService {
	return &SMTPService{
		SenderEmail: email,
		SmtpHost: host,
		SmtpPort: port,
		Password: password,
	}
}

func (_smtp *SMTPService) SendConfirmCode(user biketrackserver.User, code int) error {

	to := user.Email

	var body bytes.Buffer

	data := EmailContains{
		Name: user.Name,
		Code: code,
	}

	tmpl, err := template.ParseFiles("./templates/register_email.html")
	if err != nil {
		logrus.Fatalf("Error occured while parsing template: %s", err.Error())
		return err
	}

	if err := tmpl.Execute(&body, data); err != nil {
		logrus.Fatalf("Error occured while executing template: %s", err.Error())
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", _smtp.SenderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Дякуємо за реєстрацію!")

	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(_smtp.SmtpHost, _smtp.SmtpPort, _smtp.SenderEmail, _smtp.Password)
	logrus.Debugf("host: %s, port: %d, smail: %s, pass: %s", _smtp.SmtpHost, _smtp.SmtpPort, _smtp.SenderEmail, _smtp.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logrus.Fatalf("Error occured when send email: %s", err.Error())
		return err
	}

	return nil
}
