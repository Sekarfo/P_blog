package subscription

import (
	"fmt"
	"log"

	"github.com/jung-kurt/gofpdf"
)

// GenerateReceipt creates a PDF receipt and returns the filename
func GenerateReceipt(transactionID, customerName, expiryDate string) (string, error) {
	fileName := fmt.Sprintf("receipt_%s.pdf", transactionID)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Subscription Payment Receipt")
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)

	pdf.Cell(40, 10, fmt.Sprintf("Transaction ID: %s", transactionID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Customer Name: %s", customerName))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Expiry Date: %s", expiryDate))
	pdf.Ln(10)

	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		log.Printf("Failed to generate receipt: %v", err)
		return "", err
	}
	return fileName, nil
}
