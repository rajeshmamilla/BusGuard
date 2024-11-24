package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// BusService represents the structure of a bus service
type BusService struct {
	ServiceNo  int       `json:"serviceNo"`
	VehNo      string    `json:"vehNo"`
	BusID      int       `json:"busId"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	Date       string    `json:"date"`
	Time       string    `json:"time"`
	Status     string    `json:"status"`
	PtAlloted  int       `json:"ptAlloted"`
}

// SmartContract provides functions for managing bus services
type SmartContract struct {
	contractapi.Contract
}

// AddBusService adds a new bus service to the ledger
func (s *SmartContract) AddBusService(ctx contractapi.TransactionContextInterface, serviceNo int, vehNo string, busID int, from string, to string, date string, time string, status string, ptAlloted int) error {
	busService := BusService{
		ServiceNo: serviceNo,
		VehNo:     vehNo,
		BusID:     busID,
		From:      from,
		To:        to,
		Date:      date,
		Time:      time,
		Status:    status,
		PtAlloted: ptAlloted,
	}

	busServiceKey := "BUS" + strconv.Itoa(serviceNo)

	busServiceBytes, err := json.Marshal(busService)
	if err != nil {
		return fmt.Errorf("failed to marshal bus service: %v", err)
	}

	return ctx.GetStub().PutState(busServiceKey, busServiceBytes)
}

// QueryBusByVehNo retrieves a bus service by its vehicle number
func (s *SmartContract) QueryBusByVehNo(ctx contractapi.TransactionContextInterface, vehNo string) (*BusService, error) {
	queryString := fmt.Sprintf(`{"selector":{"vehNo":"%s"}}`, vehNo)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query bus by vehicle number: %v", err)
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext() {
		return nil, fmt.Errorf("no bus service found with vehicle number: %s", vehNo)
	}

	queryResponse, err := resultsIterator.Next()
	if err != nil {
		return nil, err
	}

	var busService BusService
	err = json.Unmarshal(queryResponse.Value, &busService)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal query response: %v", err)
	}

	return &busService, nil
}

// QueryAllBus retrieves all bus services from the ledger
func (s *SmartContract) QueryAllBus(ctx contractapi.TransactionContextInterface) ([]BusService, error) {
	queryString := `{"selector":{}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query all buses: %v", err)
	}
	defer resultsIterator.Close()

	var busServices []BusService
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var busService BusService
		err = json.Unmarshal(queryResponse.Value, &busService)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query response: %v", err)
		}
		busServices = append(busServices, busService)
	}

	return busServices, nil
}

// QueryBusByFromToDate retrieves bus services matching from, to, and date
func (s *SmartContract) QueryBusByFromToDate(ctx contractapi.TransactionContextInterface, from, to, date string) ([]BusService, error) {
	queryString := fmt.Sprintf(`{"selector":{"from":"%s","to":"%s","date":"%s"}}`, from, to, date)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query buses by from, to, and date: %v", err)
	}
	defer resultsIterator.Close()

	var busServices []BusService
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var busService BusService
		err = json.Unmarshal(queryResponse.Value, &busService)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal query response: %v", err)
		}
		busServices = append(busServices, busService)
	}

	return busServices, nil
}

// UpdateBusStatus updates the status of a bus service
func (s *SmartContract) UpdateBusStatus(ctx contractapi.TransactionContextInterface, serviceNo int, status string) error {
	busServiceKey := "BUS" + strconv.Itoa(serviceNo)

	busServiceBytes, err := ctx.GetStub().GetState(busServiceKey)
	if err != nil {
		return fmt.Errorf("failed to read bus service: %v", err)
	}
	if busServiceBytes == nil {
		return fmt.Errorf("bus service with Service No %d does not exist", serviceNo)
	}

	var busService BusService
	err = json.Unmarshal(busServiceBytes, &busService)
	if err != nil {
		return fmt.Errorf("failed to unmarshal bus service: %v", err)
	}

	busService.Status = status
	updatedBusServiceBytes, err := json.Marshal(busService)
	if err != nil {
		return fmt.Errorf("failed to marshal updated bus service: %v", err)
	}

	return ctx.GetStub().PutState(busServiceKey, updatedBusServiceBytes)
}

// Main function
func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating Bus Service chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Bus Service chaincode: %s", err.Error())
	}
}
