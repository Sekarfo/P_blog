package subscription

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// Sends an approval email with the subscription expiry date
func SendApprovalEmail(userEmail string, expiryDate string) error {
	log.Printf("Sending approval email to: %s | Expiry Date: %s\n", userEmail, expiryDate)

	from := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	to := []string{userEmail}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	message := []byte(fmt.Sprintf("Subject: Subscription Approved\n\nYour subscription has been approved! It expires on %s.", expiryDate))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v\n", userEmail, err)
		return err
	}

	log.Printf("✅ Email successfully sent to %s\n", userEmail)
	return nil
}

// Sends a rejection email to the user
func SendRejectionEmail(userEmail string) error {
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	to := []string{userEmail}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	message := []byte("Subject: Subscription Rejected\n\nYour subscription request has been rejected. If you believe this is an error, please contact support.")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Printf("Failed to send rejection email to %s: %v\n", userEmail, err)
		return err
	}
	log.Printf("Rejection email sent to %s successfully.\n", userEmail)
	return nil
}
