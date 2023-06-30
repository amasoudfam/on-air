package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const baseAddress = ""

type FlightInfoResponse struct {
	Price int
}

func GetInfo(flightNumber string) (*FlightInfoResponse, error) {

	//TODO :flightNumber
	response, err := http.Get(baseAddress + "" + flightNumber)

	if err != nil {
		return &FlightInfoResponse{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &FlightInfoResponse{}, err
	}

	var flightInfoResponse FlightInfoResponse
	json.Unmarshal(responseData, &flightInfoResponse)

	return &flightInfoResponse, nil
}

type ReserveResponse struct {
	Reserved bool
}

func Reserve(flightNumber string, passengersCount int) (*ReserveResponse, error) {

	//TODO :flightNumber , paspassengersCount
	response, err := http.Get(baseAddress + "")

	if err != nil {
		return &ReserveResponse{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &ReserveResponse{}, err
	}

	var reserveResponse ReserveResponse
	json.Unmarshal(responseData, &reserveResponse)

	return &reserveResponse, nil
}

type RefundResponse struct {
	Result bool
}

func Refund(flightNumber string, passengersCount int) (*RefundResponse, error) {

	//TODO :flightNumber , paspassengersCount
	response, err := http.Get(baseAddress + "")

	if err != nil {
		return &RefundResponse{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &RefundResponse{}, err
	}

	var refundResponse RefundResponse
	json.Unmarshal(responseData, &refundResponse)

	return &refundResponse, nil
}
