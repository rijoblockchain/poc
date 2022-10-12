/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"chaincode-go/chaincode"
)

func main() {
	researchChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating research-grant chaincode: %v", err)
	}

	if err := researchChaincode.Start(); err != nil {
		log.Panicf("Error starting research-grant chaincode: %v", err)
	}
}
