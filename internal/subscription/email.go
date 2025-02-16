package subscription

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

func SendApprovalEmail(userEmail, expiryDate, transactionID string) error {
	log.Printf("SENDING receipt email to: %s (Expiry: %s)\n", userEmail, expiryDate)

	receiptPath, err := GenerateReceipt(transactionID, userEmail, expiryDate)
	if err != nil {
		log.Printf("Failed to generate receipt: %v\n", err)
		return fmt.Errorf("failed to generate receipt: %w", err)
	}

	log.Printf("Receipt generated: %s\n", receiptPath)

	from := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	to := []string{userEmail}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Read PDF file
	fileData, err := os.ReadFile(receiptPath)
	if err != nil {
		log.Printf("Failed to read PDF receipt: %v\n", err)
		return fmt.Errorf("failed to read PDF receipt: %w", err)
	}

	// Encode file to base64
	encoded := base64.StdEncoding.EncodeToString(fileData)

	// email headers
	subject := "YOU ARE NOW PREMIUM USER - Your Receipt"
	boundary := "my-boundary"

	var emailBody bytes.Buffer
	writer := multipart.NewWriter(&emailBody)
	writer.SetBoundary(boundary)

	// text message
	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Type", "text/plain; charset=\"utf-8\"")
	headers.Set("Content-Transfer-Encoding", "7bit")

	part, err := writer.CreatePart(headers)
	if err != nil {
		return err
	}
	part.Write([]byte(fmt.Sprintf("Your subscription has been approved.\n\nYour receipt is attached.\n\nExpiry Date: %s", expiryDate)))

	// Attach the PDF
	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Type", "application/pdf; name="+filepath.Base(receiptPath))
	fileHeader.Set("Content-Disposition", "attachment; filename="+filepath.Base(receiptPath))
	fileHeader.Set("Content-Transfer-Encoding", "base64")

	part, err = writer.CreatePart(fileHeader)
	if err != nil {
		return err
	}
	part.Write([]byte(encoded))

	writer.Close()

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n%s",
		from, strings.Join(to, ","), subject, boundary, emailBody.String())

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		log.Printf("Failed to send email to %s: %v\n", userEmail, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email with receipt sent to %s\n", userEmail)
	return nil
}

func SendRejectionEmail(userEmail string) error {
	log.Printf("Sending rejection email to %s...", userEmail)

	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	to := []string{userEmail}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	message := fmt.Sprintf("Subject: Payment Rejected\n\nYour payment has been rejected.")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("✅ Rejection email sent to %s\n", userEmail)
	return nil
}
