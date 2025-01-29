package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(to, token string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")

	verificationURL := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)

	subject := "Verify Your Email"
	body := fmt.Sprintf("Click the link below to verify your email:\n\n%s", verificationURL)
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", smtpFrom, to, subject, body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpFrom, []string{to}, []byte(message)); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
	} else {
		fmt.Printf("Verification email sent to %s\n", to)
	}
}
