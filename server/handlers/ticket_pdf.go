package handlers

import (
	"bytes"
	"net/http"
	"on-air/models"
	"on-air/repository"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
)

type TicketPDF struct {
	DB *gorm.DB
}

func generate_output(ticket models.Ticket) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTextColor(0, 0, 0)
	ticketWidth := 175.0
	ticketHeight := 50.0
	ticketX := 5.0
	ticketY := 5.0
	c1 := 10.0
	c2 := 35.0
	c3 := 20.0
	c4 := 30.0
	c5 := 15.0
	c6 := 20.0
	font := "Times"
	titleFont := "Arial"
	titleFontSize := 7.0
	itemFontSize := 11.0
	qrLocation := 0.0
	lightSeparator := 7.0
	for i, passenger := range ticket.Passengers {
		if i%4 == 0 {
			pdf.AddPage()
			qrLocation = 15.0
			ticketX = 5.0
			ticketY = 10.0
		}

		pdf.Rect(ticketX, ticketY, ticketWidth, ticketHeight, "D")

		pdf.SetFont(font, "B", 16)
		pdf.Cell(0, 10, "ON-AIR Travels")

		pdf.Ln(lightSeparator + 5)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c1, 10, "Name:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c2, 10, passenger.FirstName+" "+passenger.LastName)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c3, 10, "National Code:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c4, 10, passenger.NationalCode)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c5, 10, "Gender:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c6, 10, passenger.Gender)

		pdf.Ln(lightSeparator)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c1, 10, "Airline:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c2, 10, ticket.Flight.Airline)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c3, 10, "Flight No:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c4, 10, ticket.Flight.Number)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c5, 10, "AirPlane:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c6, 10, ticket.Flight.Airplane)

		pdf.Ln(lightSeparator)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c1, 10, "From:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c2, 10, ticket.Flight.FromCity.Name+"/"+ticket.Flight.FromCity.Country.Name)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c3, 10, "Departure:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c4+10, 10, ticket.Flight.StartedAt.Format(time.RFC822))

		pdf.Ln(lightSeparator)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c1, 10, "To:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c2, 10, ticket.Flight.ToCity.Name+"/"+ticket.Flight.ToCity.Country.Name)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c3, 10, "Arrival:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c4+10, 10, ticket.Flight.EndedAt.Format(time.RFC822))

		pdf.Ln(lightSeparator)

		pdf.SetFont(titleFont, "", titleFontSize)
		pdf.Cell(c1, 10, "Price:")

		pdf.SetFont(font, "BI", itemFontSize)
		pdf.Cell(c2, 10, strconv.Itoa(ticket.UnitPrice)+" Rials")

		id := strconv.FormatUint(uint64(ticket.ID), 10)
		qrCode, err := qrcode.New("https://www.onair.org/trace/"+id+"/"+passenger.NationalCode, qrcode.Medium)
		if err != nil {
			return nil, err
		}

		png, err := qrCode.PNG(256)
		if err != nil {
			return nil, err
		}

		filename := "qr" + strconv.Itoa(i) + ".png"
		pdf.RegisterImageOptionsReader(filename, gofpdf.ImageOptions{ImageType: "png"}, bytes.NewReader(png))
		pdf.ImageOptions(filename, 137, qrLocation, 40, 40, false, gofpdf.ImageOptions{}, 0, "")

		pdf.Ln(15)

		ticketY += ticketHeight + 5
		qrLocation += ticketHeight + 5
	}
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
	// return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	// err := pdf.OutputFileAndClose("ticket.pdf")
	// if err != nil {
	// 	panic(err)
	// }
}

func (t *TicketPDF) Get(ctx echo.Context) error {
	// userID, _ := strconv.Atoi(ctx.Get("id").(string))
	userID := 2
	ticketID, err := strconv.Atoi(ctx.QueryParam("ticket_id"))
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	ticket, err := repository.GetTicket(t.DB, userID, ticketID)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	result, err := generate_output(ticket)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	ctx.Response().Header().Set("Content-Type", "application/pdf")
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename=myfile.pdf")
	ctx.Response().Header().Set("Content-Length", strconv.Itoa(len(result)))

	return ctx.Blob(http.StatusOK, "application/pdf", result)
	// return ctx.JSON(http.StatusOK, result)
}
