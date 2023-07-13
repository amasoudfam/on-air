package repository

import (
	"errors"
	"log"
	"on-air/config"
	"on-air/models"
	"on-air/pasargad"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func PayTicket(db *gorm.DB, ipg *config.IPG, ticketID uint) (string, error) {
	var dbticket models.Ticket

	err := db.First(&dbticket, "ID = ?", ticketID).Error

	if err != nil {
		return "", err
	}

	payment := models.Payment{
		TicketID: ticketID,
		Amount:   dbticket.UnitPrice * dbticket.Count,
		Status:   string(models.Requested),
	}

	err = db.Create(&payment).Error

	if err != nil {
		return "", err
	}

	pasargadApi := pasargadApi(ipg)

	request := pasargad.CreatePaymentRequest{
		Amount:        int64(payment.Amount),
		InvoiceNumber: strconv.Itoa(int(payment.ID)),
		InvoiceDate:   payment.CreatedAt.Format("2006/01/02"),
	}

	response, err := pasargadApi.Redirect(request)

	if err != nil {
		return "", err
	}

	return response, nil
}

var notFountPaymentError = errors.New("Payment not found")

func VerifyPayment(db *gorm.DB, ipg *config.IPG, paymentID int, paymentDate time.Time, transactionReferenceID int) (string, error) {
	var dbPayment models.Payment

	err := db.First(&dbPayment, "ID = ?", uint(paymentID)).Error
	if err != nil {
		return "", err
	}

	pasargadApi := pasargadApi(ipg)

	checkRequest := pasargad.CreateCheckTransactionRequest{
		InvoiceNumber:          strconv.Itoa(paymentID),
		InvoiceDate:            paymentDate.Format("2006/01/02"),
		TransactionReferenceID: strconv.Itoa(transactionReferenceID),
	}

	checkResponse, err := pasargadApi.CheckTransaction(checkRequest)
	if err != nil {
		return "", err
	}

	if checkResponse.IsSuccess != true && checkResponse.Amount != int64(dbPayment.Amount) {
		RefundPayment(ipg, dbPayment)
		return "", errors.New("Transaction not correct!")
	}

	verifyRequest := pasargad.CreateVerifyPaymentRequest{
		InvoiceNumber: strconv.Itoa(int(dbPayment.ID)),
		InvoiceDate:   dbPayment.CreatedAt.Format("2006/01/02"),
	}

	verifyResponse, err := pasargadApi.VerifyPayment(verifyRequest)
	if err != nil {
		return "", err
	}

	if verifyResponse.IsSuccess {
		dbPayment.Status = string(models.Verified)
		err = db.Save(dbPayment).Error

		if err != nil {
			RefundPayment(ipg, dbPayment)
			return "", err
		}

	} else {
		RefundPayment(ipg, dbPayment)
	}

	ChangeTicketStatus(db, dbPayment.TicketID, string(models.PaymentPaid))

	return dbPayment.Status, nil
}

func RefundPayment(ipg *config.IPG, dbPayment models.Payment) {

	pasargadApi := pasargadApi(ipg)

	request := pasargad.CreateRefundRequest{
		InvoiceNumber: strconv.Itoa(int(dbPayment.ID)),
		InvoiceDate:   dbPayment.CreatedAt.Format("2006/01/02"),
	}

	_, err := pasargadApi.Refund(request)

	if err != nil {
		log.Fatal(err)
	}

}

func pasargadApi(ipg *config.IPG) (pasrgad *pasargad.PasargadPaymentAPI) {
	return pasargad.PasargadAPI(
		int64(ipg.MerchantCode),
		int64(ipg.TerminalId),
		ipg.RedirectUrl,
		ipg.CertFile,
	)
}

func ChangePaymentStatus(db *gorm.DB, ticketID uint, status string) error {
	var payments []models.Payment

	err := db.Model(&payments).Where("ticket_id = ?", ticketID).Find(&payments).Error

	if err != nil {
		return err
	}

	err = db.Model(&payments).Update("status", status).Error

	if err != nil {
		return err
	}

	return nil
}
