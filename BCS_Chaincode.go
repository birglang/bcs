package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type BCSChaincode struct {
	contractapi.Contract
}

type ControllerReputation struct {
	ID         string  `json:"id"`
	Reputation float64 `json:"reputation"`
}

// Initialize Ledger
func (b *BCSChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	controllers := []ControllerReputation{
		{ID: "controller_1", Reputation: 1.0},
		{ID: "controller_2", Reputation: 1.0},
	}

	for _, controller := range controllers {
		controllerJSON, err := json.Marshal(controller)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(controller.ID, controllerJSON)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get Reputation
func (b *BCSChaincode) GetReputation(ctx contractapi.TransactionContextInterface, id string) (float64, error) {
	controllerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return 0, fmt.Errorf("failed to get controller: %v", err)
	}
	if controllerJSON == nil {
		return 0, fmt.Errorf("controller not found")
	}

	var controller ControllerReputation
	err = json.Unmarshal(controllerJSON, &controller)
	if err != nil {
		return 0, err
	}

	return controller.Reputation, nil
}

// Update Reputation
func (b *BCSChaincode) UpdateReputation(ctx contractapi.TransactionContextInterface, id string, delta float64) error {
	controllerJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to get controller: %v", err)
	}
	if controllerJSON == nil {
		return fmt.Errorf("controller not found")
	}

	var controller ControllerReputation
	err = json.Unmarshal(controllerJSON, &controller)
	if err != nil {
		return err
	}

	controller.Reputation += delta
	if controller.Reputation < 0 {
		controller.Reputation = 0
	} else if controller.Reputation > 1 {
		controller.Reputation = 1
	}

	updatedControllerJSON, err := json.Marshal(controller)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedControllerJSON)
}
