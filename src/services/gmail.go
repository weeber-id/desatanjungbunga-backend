package services

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
)

// Email structure
type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

// Send Email
func (e *Email) Send() error {
	config := variables.GmailConfig
	e.From = config.Email

	msg := "From: " + e.From + "\n" +
		"To: " + e.To + "\n" +
		"Subject: " + e.Subject + "\n\n" +
		e.Body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", e.From, config.Password, "smtp.gmail.com"),
		e.From,
		[]string{e.To},
		[]byte(msg),
	)
	if err != nil {
		log.Printf("Smtp error to %s : %v \n", e.To, err)
		return err
	}
	log.Printf("Smtp SUCCESS send mail to %s\n", e.To)
	return nil
}

// SendNewPasswordForReset via email
func (e *Email) SendNewPasswordForReset(name, username, password string) error {
	e.Subject = "Pemberitahuan Kantor Desa TanjungBunga - Reset Password"
	e.Body = fmt.Sprintf("Halo %s, Akun anda telah direset oleh super admin. Silahkan login kembali dengan menggunakan \n\nUsername: %s\nPassword: %s\n\n", name, username, password)

	return e.Send()
}
