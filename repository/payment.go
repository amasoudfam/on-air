package repository

import (
	"errors"
	"log"
	"on-air/config"
	"on-air/models"
	"strconv"

	pasargad "github.com/pepco-api/golang-rest-sdk"

	"gorm.io/gorm"
)

func PayTicket(db *gorm.DB, ipg *config.IPG, ticketID uint) (string, error) {
	var dbticket models.Ticket

	result := db.First(&dbticket, "ID = ?", ticketID)
	if result.RowsAffected == 0 {
		return "", errors.New("Ticket not found")
	}

	dbPayment := models.Payment{
		TicketID: ticketID,
		Amount:   dbticket.UnitPrice * dbticket.Count,
		Status:   "Requested",
	}

	result = db.Create(&dbPayment)
	if err := result.Error; err != nil {
		return "", err
	}

	//TODO: Third party call to get bank address with token to redirect to bank
	pasargadApi := pasargadApi(ipg)

	request := pasargad.CreatePaymentRequest{
		Amount:        int64(dbPayment.Amount),
		InvoiceNumber: strconv.Itoa(int(dbPayment.ID)),
		InvoiceDate:   dbPayment.CreatedAt.String(),
	}

	response, err := pasargadApi.Redirect(request)
	if err != nil {
		return "", err
	}

	return response, nil
}

var notFountPaymentError = errors.New("Payment not found")

func VerifyPayment(db *gorm.DB, ipg *config.IPG, paymentID uint) (string, error) {
	var dbPayment models.Payment

	result := db.First(&dbPayment, "ID = ?", paymentID)
	if result.RowsAffected == 0 {
		return "", notFountPaymentError
	}

	pasargadApi := pasargadApi(ipg)

	checkRequest := pasargad.CreateCheckTransactionRequest{
		InvoiceNumber: strconv.Itoa(int(dbPayment.ID)),
		InvoiceDate:   dbPayment.CreatedAt.String(),
	}

	checkResponse, err := pasargadApi.CheckTransaction(checkRequest)
	if err != nil {
		return "", err
	}

	if checkResponse.IsSuccess != true && checkResponse.Amount != int64(dbPayment.Amount) {
		RefundPayment(ipg, dbPayment)
		return "", errors.New("Transaction not correct!")
	}

	//TODO : add TransactionReferenceID to payment
	verifyRequest := pasargad.CreateVerifyPaymentRequest{
		InvoiceNumber: strconv.Itoa(int(dbPayment.ID)),
		InvoiceDate:   dbPayment.CreatedAt.String(),
	}

	verifyResponse, err := pasargadApi.VerifyPayment(verifyRequest)
	if err != nil {
		return "", err
	}

	if verifyResponse.IsSuccess {
		dbPayment.Status = "Verified"
		result = db.Save(dbPayment)

		if result.RowsAffected == 0 {
			return "", notFountPaymentError
		}

	} else {
		RefundPayment(ipg, dbPayment)
	}

	ChangeTicketStatus(db, dbPayment.TicketID, "Payed")

	return dbPayment.Status, nil
}

func RefundPayment(ipg *config.IPG, dbPayment models.Payment) {

	pasargadApi := pasargadApi(ipg)

	request := pasargad.CreateRefundRequest{
		InvoiceNumber: strconv.Itoa(int(dbPayment.ID)),
		InvoiceDate:   dbPayment.CreatedAt.String(),
	}

	//TODO : rertry multiple times
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
