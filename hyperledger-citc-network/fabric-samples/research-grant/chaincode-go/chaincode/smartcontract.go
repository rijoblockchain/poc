package chaincode

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"time"
	"strings"
	//"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// START CONSTANTS

var GrantorMSP = "GrantorMSP"
var AwardeeMSP = "AwardeeMSP"
var SubawardeeMSP = "SubawardeeMSP"

// END CONSTANTS

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}


// Grant describes details of research grant
type Grant struct {
	ID              string 		`json:"ID"`
	Amount          float64 	`json:"amount"`
	Awardee         []Awardee   `json:"awardee"`
	Benefit         []Benefit	`json:"benefit"`
	Cashed_Out      float64     `json:"cashed_out"`
	Description     string      `json:"description"`
	End_Date		string	    `json:"end_date"`
	Grantor			string      `json:"grantor"`
	Grantor_ID		string      `json:"grantor_id"`
	Notes			string      `json:"notes"`
	Paid_Amount		float64		`json:"paid_amount"`
	Payment			[]Payment	`json:"payment"`
	Payment_Type	string 		`json:"payment_type"`
	Progress		[]Progress	`json:"progress"`
	Progress_Freq	string 		`json:"progress_freq"`
	Start_Date		string  	`json:"start_date"`
	Status			string 		`json:"status"`
	Sub 			float64	 	`json:"sub"`
}

// Awardee describes details of Awardee and Subawardee
type Awardee struct {
	Account_Number  		string 		`json:"account_number"`
	Contact         		string		`json:"contact"`
	ID						string		`json:"id"`
	Name        			string		`json:"name"`
	Principal_Investigator  string		`json:"principal_investigator"`
	Organization			string		`json:"organization"`
	Awardee_Type        	string		`json:"awardee_type"`
}

// Benefit describes details of availed benefits for the research
type Benefit struct {
	Benefit    string	`json:"benefit"`
	Amount     float64	`json:"amount"`
}

// Payment describes details of payments
type Payment struct {
	ID              string 		`json:"ID"`
	Awardee_ID      string 	    `json:"awardee_id"`
	Date			string      `json:"date"`
	Item         	[]Benefit   `json:"item"`
	Notes			string      `json:"notes"`
	Status			string 		`json:"status"`
	Total			float64		`json:"total"`
}

// Progress describes details of research developments
type Progress struct {   
	Notes			string      `json:"notes"`
	Percentage		string 		`json:"percentage"`
}

type AmountResponse struct {
	Cashed_Out    		 float64	`json:"cashedOut"`
	Requested_Amount     float64	`json:"requestedAmount"`
}

type PaymentStatus struct {
	Grant_ID		string		`json:"grant_id"`
	Payment			[]Payment	`json:"payment"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("Research grant ledger is initiated")

	return nil
}

// Create a new Grant
func (s *SmartContract) InitiateGrant(ctx contractapi.TransactionContextInterface) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to initiate grant", clientMSPID)
	}

	// Get new transaction definition details from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return false, fmt.Errorf("error getting transient: %v", err)
	}

	// Private records get passed in transient field, instead of func args
	transientGrantJSON, ok := transientMap["grant"]
	if !ok {
		//log error to stdout
		return false, fmt.Errorf("grant not found in the transient map input")
	}

	var grant Grant

	json.Unmarshal([]byte(transientGrantJSON), &grant)

	id := grant.ID
	grant.Grantor = strings.Replace(clientMSPID, "MSP", "", 1)
	grant.Grantor_ID = userId
	grant.Status = "Not Assigned"

	var benefitAmount float64
	for _, benefit := range grant.Benefit {
		benefitAmount += benefit.Amount
	}

	if benefitAmount != grant.Amount {
		return false, fmt.Errorf("Total Benefit %.2f doesn't match with the Grant Amount %.2f", benefitAmount, grant.Amount)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grantExists, err := ctx.GetStub().GetState(requestCompositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if grantExists != nil {
		return false, fmt.Errorf("the grant %s exists", id)
	}

	grantJSON, err := json.Marshal(grant)
	if err != nil {
		return false,  fmt.Errorf("error marshaling json: %v", err)
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	
	return true, nil
}

// ReadGrant returns the grant stored in the world state with given id.
func (s *SmartContract) ReadGrant(ctx contractapi.TransactionContextInterface, id string) (*Grant, error) {
	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grantJSON, err := ctx.GetStub().GetState(requestCompositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if grantJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var grant Grant
	err = json.Unmarshal(grantJSON, &grant)
	if err != nil {
		return nil, err
	}

	return &grant, nil
}

// Assign grant to awardee - Grantor
func (s *SmartContract) AssignGrant(ctx contractapi.TransactionContextInterface) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to assign grant", clientMSPID)
	}

	type assignTransientInput struct {
		Grant_ID		string		`json:"grant_id"`
		Awardee         []Awardee   `json:"awardee"`
		Status			string		`json:"status"`
	}

	// Get new transaction definition details from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return false, fmt.Errorf("error getting transient: %v", err)
	}

	// Private records get passed in transient field, instead of func args
	transientGrantJSON, ok := transientMap["assign_grant"]
	if !ok {
		//log error to stdout
		return false, fmt.Errorf("assign_grant not found in the transient map input")
	}

	var assignGrantInput assignTransientInput

	err = json.Unmarshal(transientGrantJSON, &assignGrantInput)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	assignGrantInput.Status = "Pending"
	for i := 0; i < len(assignGrantInput.Awardee); i++ {
        if len(assignGrantInput.Awardee[i].Name) == 0 {
			return false, fmt.Errorf("Name field must be a non-empty string")
		}
		if len(assignGrantInput.Awardee[i].ID) == 0 {
			return false, fmt.Errorf("ID field must be a non-empty string")
		}
		if len(assignGrantInput.Awardee[i].Contact) == 0 {
			return false, fmt.Errorf("Contact field must be a non-empty string")
		}
		if len(assignGrantInput.Awardee[i].Principal_Investigator) == 0 {
			return false, fmt.Errorf("Principal_Investigator field must be a non-empty string")
		}
		if len(assignGrantInput.Awardee[i].Account_Number) == 0 {
			return false, fmt.Errorf("Account_Number field must be a non-empty string")
		}
		if len(assignGrantInput.Awardee[i].Organization) == 0 {
			return false, fmt.Errorf("Organization field must be a non-empty string")
		}
		assignGrantInput.Awardee[i].Organization = strings.Title(strings.ToLower(strings.Replace(assignGrantInput.Awardee[i].Organization, "MSP", "", 1)))
		if len(assignGrantInput.Awardee[i].Awardee_Type) == 0 {
			return false, fmt.Errorf("Awardee_Type field must be a non-empty string")
		}
    }


	grant, err := s.ReadGrant(ctx, assignGrantInput.Grant_ID)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", assignGrantInput.Grant_ID)
	}

	if grant.Status != "Not Assigned" {
		return false, fmt.Errorf("Grant %s is in %s status", grant.ID, grant.Status)	
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("User %s from org %v is not authorized to assign grant for this grant %s",userId, clientMSPID, grant.ID)
	}

	assignGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        assignGrantInput.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		grant.Payment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			assignGrantInput.Status,
		Sub: 			grant.Sub,
	}

	grantJSONasBytes, err := json.Marshal(assignGrant)
	if err != nil {
		return false, fmt.Errorf("failed to marshal grant into JSON: %v", err)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{grant.ID})
	err = ctx.GetStub().PutState(requestCompositeKey, grantJSONasBytes)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil

}

// Awardee accept grant
func (s *SmartContract) AcceptGrant(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != AwardeeMSP {
		return false, fmt.Errorf("User from org %v is not authorized to accept grant", clientMSPID)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grant, err := s.ReadGrant(ctx, id)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", id)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if grant.Status != "Pending" {
		return false, fmt.Errorf("Grant %s is in %s status. It should be in assigned to an awardee", grant.ID, grant.Status)	
	}

	
	if !checkAwardee(grant.Awardee, userId) {
		return false, fmt.Errorf("Awardee %s is not assigned in the Grant %s", userId, grant.ID)	
	}

	if checkSubAwardee(grant.Awardee, userId) {
		return false, fmt.Errorf("Awardee %s is not allowed to accept this Grant %s", userId, grant.ID)	
	}

	approveGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		grant.Payment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			"Approved",
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(approveGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Awardee reject grant
func (s *SmartContract) RejectGrant(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != AwardeeMSP {
		return false, fmt.Errorf("User from org %v is not authorized to reject grant", clientMSPID)
	}
	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grant, err := s.ReadGrant(ctx, id)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", id)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if !checkAwardee(grant.Awardee, userId) {
		return false, fmt.Errorf("Awardee %s is not assigned in the Grant %s", userId, grant.ID)	
	}

	if grant.Status != "Pending" {
		return false, fmt.Errorf("Grant %s is not in Pending status", grant.ID)	
	}


	rejectGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		grant.Payment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			"Rejected",
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(rejectGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Awardee reject grant
func (s *SmartContract) RevokeGrant(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to revoke grant", clientMSPID)
	}
	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grant, err := s.ReadGrant(ctx, id)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", id)
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("Grantor %s is not allowed to revoke the Grant %s", userId, grant.ID)	
	}

	revokeGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		grant.Payment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			"Revoked",
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(revokeGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Update Grant
func (s *SmartContract) UpdateGrant(ctx contractapi.TransactionContextInterface) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to update grant", clientMSPID)
	}
	// Get new transaction definition details from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return false, fmt.Errorf("error getting transient: %v", err)
	}

	// Private records get passed in transient field, instead of func args
	transientGrantJSON, ok := transientMap["update_grant"]
	if !ok {
		//log error to stdout
		return false, fmt.Errorf("grant not found in the transient map input")
	}

	var updatedGrant Grant

	json.Unmarshal([]byte(transientGrantJSON), &updatedGrant)

	id := updatedGrant.ID

	grant, err := s.ReadGrant(ctx, id)
	if err != nil {
		return false, err
	}
	if grant == nil {
		return false, fmt.Errorf("the grant %s does not exist", id)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", updatedGrant.ID)	
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("Grantor %s is not allowed to update the Grant %s", userId, updatedGrant.ID)	
	}

	var benefitAmount float64
	for _, benefit := range updatedGrant.Benefit {
		benefitAmount += benefit.Amount
	}

	if benefitAmount != updatedGrant.Amount {
		return false, fmt.Errorf("Total Benefit %.2f doesn't match with the Grant Amount %.2f", benefitAmount, updatedGrant.Amount)
	}

	updateGrant := Grant{
		ID:             grant.ID,
		Amount:         updatedGrant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        updatedGrant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		grant.Payment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})

	grantJSON, err := json.Marshal(updateGrant)
	if err != nil {
		return false,  fmt.Errorf("error marshaling json: %v", err)
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put grant into ledger: %v", err)
	}
	return true, nil
}

// GetAllGrants returns all grants found in world state
func (s *SmartContract) GetAllGrants(ctx contractapi.TransactionContextInterface) ([]*Grant, error) {
	//requestPartialCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("grant", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var grants []*Grant
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var grant Grant
		err = json.Unmarshal(queryResponse.Value, &grant)
		if err != nil {
			return nil, err
		}
		grants = append(grants, &grant)
	}

	return grants, nil
}


// Request Reimbursement by awardee
func (s *SmartContract) RequestReimbursement(ctx contractapi.TransactionContextInterface) (string, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	

	dt := time.Now()
	formattedTime := dt.Format("01-02-2006 15:04:05")
	type reimbursementTransientInput struct {
		ID				string		`json:"ID"`
		Grant_ID		string		`json:"grant_id"`
		Awardee_ID      string   	`json:"awardee_id"`
		Date		    string   	`json:"date"`
		Notes			string		`json:"notes"`
		Item			[]Benefit	`json:"item"`
	}
	// Get new transaction definition details from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "", fmt.Errorf("error getting transient: %v", err)
	}

	// Private records get passed in transient field, instead of func args
	transientReimbursementJSON, ok := transientMap["request_reimbursement"]
	if !ok {
		//log error to stdout
		return "", fmt.Errorf("grant not found in the transient map input")
	}

	var reimbursementInput reimbursementTransientInput

	err = json.Unmarshal(transientReimbursementJSON, &reimbursementInput)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	id := reimbursementInput.Grant_ID

	grant, err := s.ReadGrant(ctx, id)
	if err != nil {
		return "", err
	}
	if grant == nil {
		return "", fmt.Errorf("the asset %s does not exist", id)
	}
	if grant.Status == "Revoked" {
		return "", fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if !checkAwardee(grant.Awardee, reimbursementInput.Awardee_ID) && !checkSubAwardee(grant.Awardee, reimbursementInput.Awardee_ID) {
		return "", fmt.Errorf("Awardee %s is not assigned in the Grant %s", reimbursementInput.Awardee_ID, grant.ID)	
	}

	if grant.Status != "Approved" {
		return "", fmt.Errorf("Grant %s is not approved by the Awardee %s", grant.ID, userId)	
	}
	
	if clientMSPID != AwardeeMSP && clientMSPID != SubawardeeMSP {

		return "", fmt.Errorf("User from org %v is not authorized to request reimbursement", clientMSPID)
	}

	if checkPayment(grant.Payment, reimbursementInput.ID) {
		return "", fmt.Errorf("Payment with ID %s already exists in the Grant %s", reimbursementInput.ID, grant.ID)	
	}

	if reimbursementInput.Awardee_ID != userId {
		return "", fmt.Errorf("User %s is not allowed to request reimbursement for %s", userId, reimbursementInput.Awardee_ID)
	}
	
	var flag bool
	for _, awardee := range grant.Awardee {
			if reimbursementInput.Awardee_ID == awardee.ID {
				flag = true
				break
			}
		} 

	if !flag {
		return "", fmt.Errorf("the awardee %s does not exist in the Grant %s", reimbursementInput.Awardee_ID, id)
	}

	var paid_amount float64
	var itemMap  = make(map[string]float64)
	for _, item := range reimbursementInput.Item {
		itemMap[item.Benefit] = item.Amount
		paid_amount += item.Amount
	}

	var awardee_type string
	for _, awardee := range grant.Awardee {
		if awardee.ID == reimbursementInput.Awardee_ID {
			awardee_type = awardee.Awardee_Type
		}
	}

	var benefitMap  = make(map[string]float64)
	for _, benefit := range grant.Benefit {
		benefitMap[benefit.Benefit] = benefit.Amount
	}
	fmt.Println(benefitMap)

	var benefitAmountMapMain  = make(map[string]float64)
	var benefitAmountMapSub  = make(map[string]float64)
	for _, benefit := range grant.Benefit {
		benefitAmountMapMain[benefit.Benefit] = 0.0
		benefitAmountMapSub[benefit.Benefit] = 0.0
	}

	var payment_amount float64
	for _, payment := range grant.Payment {
		if payment.Status != "Requested" && payment.Status != "Accepted" && payment.Status != "Pending-redeem" && payment.Status != "Accept_redeem" {
			continue
		}
		for _, item := range payment.Item {
			if checkAwardee(grant.Awardee, payment.Awardee_ID) {
				benefitAmountMapMain[item.Benefit] = benefitAmountMapMain[item.Benefit] + item.Amount
			} else {
				benefitAmountMapSub[item.Benefit] = benefitAmountMapSub[item.Benefit] + item.Amount
			}
			
		}
		payment_amount += payment.Total
	}


	flag = false
	var message string
	var totalBenefitAmount float64
	var totalBenefitAmountForItem float64
	for key, value := range itemMap {
		totalBenefitAmount += value	
		totalBenefitAmountForItem = benefitAmountMapMain[key] + benefitAmountMapSub[key] + value	
		if awardee_type == "Main" {
			benefitAmountMapMain[key] = benefitAmountMapMain[key] + value 
			flag = checkBenefitAmount(totalBenefitAmountForItem, benefitMap[key], benefitAmountMapMain[key], "Main", 100.00)
			if !flag{
				balanceAmount := 0.0
				if benefitMap[key]-totalBenefitAmountForItem+value < 0 {
					balanceAmount = 0
				}else {
					balanceAmount = benefitMap[key]-totalBenefitAmountForItem+value
				}
				message = fmt.Sprintf("Requested value of %.2f exceeds the allocated amount of %.2f. Remaining amount available for %s benefit is %.2f", value, benefitMap[key], key, balanceAmount)
				break
			}
			} else if awardee_type == "Sub" {
			benefitAmountMapSub[key] = benefitAmountMapSub[key] + value
			flag = checkBenefitAmount(totalBenefitAmountForItem, benefitMap[key], benefitAmountMapSub[key], "Sub", grant.Sub)
			if !flag{
				balanceAmount := 0.0
				if (grant.Sub*benefitMap[key]/100)-totalBenefitAmountForItem+value < 0 {
					balanceAmount = 0
				}else {
					balanceAmount = (grant.Sub*benefitMap[key]/100)-totalBenefitAmountForItem+value
				}
				message = fmt.Sprintf("Requested value of %.2f exceeds the allocated %.2f percentage of the total amount (%.2f). Remaining amount available for %s benefit for subawardee is %.2f.", value, grant.Sub, benefitMap[key], key, balanceAmount)
				break
			}
		}
	}

	
	if payment_amount + totalBenefitAmount > grant.Amount {
		return "", fmt.Errorf("Already requested and paid amount adds upto %.2f. The requested total amount of %.2f exceeds Grant's amount of %.2f", payment_amount, paid_amount, grant.Amount)
	}


	if !flag {
		return "", fmt.Errorf("%s", message)
	}

	payment := Payment{
		ID:				reimbursementInput.ID,
		Awardee_ID:     reimbursementInput.Awardee_ID,
		Date:			formattedTime,
		Item:         	reimbursementInput.Item,
		Notes:			reimbursementInput.Notes,
		Status:			"Requested",
		Total:			totalBenefitAmount,
	}

	grant.Payment = append(grant.Payment, payment)

	grantJSON, err := json.Marshal(grant)
	if err != nil {
		return "", err
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return "", fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	
	return fmt.Sprintf("Reimbursement Request for the Payment %s is successful", payment.ID), nil
}

// Grantor accept reimbursement
func (s *SmartContract) AcceptReimbursement(ctx contractapi.TransactionContextInterface, grant_id string, payment_id string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to accept reimbursement", clientMSPID)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{grant_id})
	grant, err := s.ReadGrant(ctx, grant_id)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", grant_id)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("User %s from org %v is not authorized to accept reimbursement for this grant %s",userId, clientMSPID, grant.ID)
	}

	if !checkPayment(grant.Payment, payment_id) {
		return false, fmt.Errorf("Payment %s doesn't exist in the Grant %s", payment_id, grant.ID)	
	}

	paid_amount := grant.Paid_Amount
	var updatedPayment []Payment
	for _, payment := range grant.Payment {
		if payment.ID == payment_id {
			if payment.Status != "Requested" {
				return false, fmt.Errorf("Payment %s is Not Requested", payment_id)	
			}
			payment.Status = "Accepted"
			paid_amount += payment.Total
		}
		updatedPayment = append(updatedPayment, payment)
	}

	updatedGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	paid_amount,
		Payment:		updatedPayment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(updatedGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Grantor reject reimbursement
func (s *SmartContract) RejectReimbursement(ctx contractapi.TransactionContextInterface, grant_id string, payment_id string, msg string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to reject reimbursement", clientMSPID)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{grant_id})
	grant, err := s.ReadGrant(ctx, grant_id)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", grant_id)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("User %s from org %v is not authorized to reject reimbursement for this grant %s",userId, clientMSPID, grant.ID)
	}

	if !checkPayment(grant.Payment, payment_id) {
		return false, fmt.Errorf("Payment %s doesn't exist in the Grant %s", payment_id, grant.ID)	
	}

	var updatedPayment []Payment
	for _, payment := range grant.Payment {
		if payment.ID == payment_id {
			if payment.Status != "Requested" {
				return false, fmt.Errorf("Payment %s is Not Requested", payment_id)	
			}
			payment.Status = msg
		}
		updatedPayment = append(updatedPayment, payment)
	}

	updatedGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		updatedPayment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(updatedGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Awardee redeem tokens
func (s *SmartContract) RedeemTokens(ctx contractapi.TransactionContextInterface, grant_id string, payment_id string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != AwardeeMSP && clientMSPID != SubawardeeMSP {
		return false, fmt.Errorf("User from org %v is not authorized to redeem tokens", clientMSPID)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{grant_id})
	grant, err := s.ReadGrant(ctx, grant_id)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", grant_id)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if !checkAwardee(grant.Awardee, userId) && !checkSubAwardee(grant.Awardee, userId) {
		return false, fmt.Errorf("Awardee %s is not assigned in the Grant %s", userId, grant.ID)	
	}

	if !checkPayment(grant.Payment, payment_id) {
		return false, fmt.Errorf("Payment %s doesn't exist in the Grant %s", payment_id, grant.ID)	
	}

	var updatedPayment []Payment
	for _, payment := range grant.Payment {
		if payment.ID == payment_id {
			if payment.Awardee_ID != userId {
				return false, fmt.Errorf("Awardee %s is not allowed to redeem tokens for this payment %s", userId, payment_id)	
			}
			if payment.Status != "Accepted" {
				return false, fmt.Errorf("Reimbursement for payment %s has Not Accepted by the Grantor", payment_id)	
			}
			payment.Status = "Pending-redeem"
		}
		updatedPayment = append(updatedPayment, payment)
	}

	updatedGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		updatedPayment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(updatedGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Grantor accept redeem
func (s *SmartContract) AcceptRedeem(ctx contractapi.TransactionContextInterface, grant_id string, payment_id string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]
	
	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to accept redeem", clientMSPID)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{grant_id})
	grant, err := s.ReadGrant(ctx, grant_id)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", grant_id)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("Grantor %s is not allowed to accept redeem for the Grant %s", userId, grant.ID)	
	}

	if !checkPayment(grant.Payment, payment_id) {
		return false, fmt.Errorf("Payment %s doesn't exist in the Grant %s", payment_id, grant.ID)	
	}

	var updatedPayment []Payment
	var cashedOut float64
	var paid_amount float64
	for _, payment := range grant.Payment {
		if payment.ID == payment_id {
			if payment.Status != "Pending-redeem" {
				return false, fmt.Errorf("Payment %s is not in Pending-redeem status", payment_id)	
			}
			payment.Status = "Accept_redeem"
			cashedOut = grant.Cashed_Out + payment.Total
			paid_amount = payment.Total

		}
		updatedPayment = append(updatedPayment, payment)
	}

	updatedGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     cashedOut,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount - paid_amount,
		Payment:		updatedPayment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(updatedGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Grantor reject redeem
func (s *SmartContract) RejectRedeem(ctx contractapi.TransactionContextInterface, grant_id string, payment_id string, msg string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to reject redeem", clientMSPID)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{grant_id})
	grant, err := s.ReadGrant(ctx, grant_id)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", grant_id)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("Grantor %s is not allowed to accept redeem for the Grant %s", userId, grant.ID)	
	}

	if !checkPayment(grant.Payment, payment_id) {
		return false, fmt.Errorf("Payment %s doesn't exist in the Grant %s", payment_id, grant.ID)	
	}

	var updatedPayment []Payment
	var paid_amount float64
	for _, payment := range grant.Payment {
		if payment.ID == payment_id {
			if payment.Status != "Pending-redeem" {
				return false, fmt.Errorf("Payment %s is not in Pending-redeem status", payment_id)	
			}
			payment.Status = msg
			paid_amount = payment.Total
		}
		updatedPayment = append(updatedPayment, payment)
	}

	updatedGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount - paid_amount,
		Payment:		updatedPayment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(updatedGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Grantor add Awardees
func (s *SmartContract) AddAwardee(ctx contractapi.TransactionContextInterface) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to add awardee", clientMSPID)
	}


	type awardeeTransientInput struct {
		Grant_ID		string		`json:"grant_id"`
		Awardee         Awardee   	`json:"awardee"`
	}

	// Get new transaction definition details from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return false, fmt.Errorf("error getting transient: %v", err)
	}

	// Private records get passed in transient field, instead of func args
	transientGrantJSON, ok := transientMap["add_awardee"]
	if !ok {
		//log error to stdout
		return false, fmt.Errorf("add_awardee not found in the transient map input")
	}

	var awardeeInput awardeeTransientInput
	err = json.Unmarshal(transientGrantJSON, &awardeeInput)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	awardeeInput.Awardee.Awardee_Type = "Main"

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{awardeeInput.Grant_ID})
	grant, err := s.ReadGrant(ctx, awardeeInput.Grant_ID)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", awardeeInput.Grant_ID)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("Grantor %s is not allowed to add awardee in the Grant %s", userId, grant.ID)	
	}

	if checkAwardee(grant.Awardee, awardeeInput.Awardee.ID) {
		return false, fmt.Errorf("Awardee %s is already exists in the Grant %s", awardeeInput.Awardee.ID, grant.ID)	
	}

	updatedGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        append(grant.Awardee, awardeeInput.Awardee),
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		grant.Payment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(updatedGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil
}

// Awardee add SubAwardees
func (s *SmartContract) AddSubawardee(ctx contractapi.TransactionContextInterface) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != AwardeeMSP {
		return false, fmt.Errorf("User from org %v is not authorized to add subawardee", clientMSPID)
	}


	type subAwardeeTransientInput struct {
		Grant_ID		string		`json:"grant_id"`
		Awardee_ID		string 		`json:"awardee_id"`
		Awardee         Awardee   	`json:"awardee"`
	}

	// Get new transaction definition details from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return false, fmt.Errorf("error getting transient: %v", err)
	}

	// Private records get passed in transient field, instead of func args
	transientGrantJSON, ok := transientMap["add_subawardee"]
	if !ok {
		//log error to stdout
		return false, fmt.Errorf("add_subawardee not found in the transient map input")
	}

	var subAwardeeInput subAwardeeTransientInput
	err = json.Unmarshal(transientGrantJSON, &subAwardeeInput)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	subAwardeeInput.Awardee.Awardee_Type = "Sub"

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{subAwardeeInput.Grant_ID})
	grant, err := s.ReadGrant(ctx, subAwardeeInput.Grant_ID)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", subAwardeeInput.Grant_ID)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if !checkAwardee(grant.Awardee, userId) {
		return false, fmt.Errorf("Awardee %s is not assigned in the Grant %s", userId, grant.ID)	
	}

	if checkSubAwardee(grant.Awardee, subAwardeeInput.Awardee.ID) {
		return false, fmt.Errorf("Awardee %s is already exists in the Grant %s", subAwardeeInput.Awardee.ID, grant.ID)	
	}

	if grant.Status != "Approved" {
		return false, fmt.Errorf("Grant %s is not approved by the Awardee %s", grant.ID, userId)	
	}

	updatedGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        append(grant.Awardee, subAwardeeInput.Awardee),
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		grant.Payment,
		Payment_Type:	grant.Payment_Type,
		Progress:		grant.Progress,
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(updatedGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil

}

// Awardee add Progress
func (s *SmartContract) AddProgress(ctx contractapi.TransactionContextInterface) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != AwardeeMSP && clientMSPID != SubawardeeMSP {
		return false, fmt.Errorf("User from org %v is not authorized to add progress", clientMSPID)
	}


	type progressTransientInput struct {
		Grant_ID		string			`json:"grant_id"`
		Progress		Progress 		`json:"progress"`	
	}

	// Get new transaction definition details from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return false, fmt.Errorf("error getting transient: %v", err)
	}

	// Private records get passed in transient field, instead of func args
	transientGrantJSON, ok := transientMap["add_progress"]
	if !ok {
		//log error to stdout
		return false, fmt.Errorf("add_progress not found in the transient map input")
	}

	var progressInput progressTransientInput
	err = json.Unmarshal(transientGrantJSON, &progressInput)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{progressInput.Grant_ID})
	grant, err := s.ReadGrant(ctx, progressInput.Grant_ID)
	if err != nil {
		return false, fmt.Errorf("Grant %s does not exist", progressInput.Grant_ID)
	}

	if grant.Status == "Revoked" {
		return false, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if !checkAwardee(grant.Awardee, userId) && !checkSubAwardee(grant.Awardee, userId) {
		return false, fmt.Errorf("Awardee %s is not assigned in the Grant %s", userId, grant.ID)	
	}

	if grant.Status != "Approved" {
		return false, fmt.Errorf("Grant %s is not approved by the Awardee %s", grant.ID, userId)	
	}

	updatedGrant := Grant{
		ID:             grant.ID,
		Amount:         grant.Amount,
		Awardee:        grant.Awardee,
		Benefit:        grant.Benefit,
		Cashed_Out:     grant.Cashed_Out,
		Description:    grant.Description,
		End_Date:		grant.End_Date,
		Grantor:		grant.Grantor,
		Grantor_ID:		grant.Grantor_ID,
		Notes:			grant.Notes,
		Paid_Amount:	grant.Paid_Amount,
		Payment:		grant.Payment,
		Payment_Type:	grant.Payment_Type,
		Progress:		append(grant.Progress, progressInput.Progress),
		Progress_Freq:	grant.Progress_Freq,
		Start_Date:		grant.Start_Date,
		Status:			grant.Status,
		Sub: 			grant.Sub,
	}

	grantJSON, err := json.Marshal(updatedGrant)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(requestCompositeKey, grantJSON)

	if err != nil {
		return false, fmt.Errorf("failed to put transaction definition into ledger: %v", err)
	}
	return true, nil

}

// Get Wallet with Specified Status
func (s *SmartContract) GetWallet(ctx contractapi.TransactionContextInterface, grant_id string, awardee_id string, status string) (float64, error) {
	
	grant, err := s.ReadGrant(ctx, grant_id)
	if err != nil {
		return 0.0, fmt.Errorf("Grant %s does not exist", grant_id)
	}

	if grant.Status == "Revoked" {
		return 0.0, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	var totalAmount float64
	for _, payment := range grant.Payment{
		if payment.Awardee_ID == awardee_id && payment.Status == status{
			totalAmount += payment.Total
		}
	}

	return totalAmount, nil
}

// Get Wallet with Specified Status
func (s *SmartContract) MyWallet(ctx contractapi.TransactionContextInterface, grant_id string) (*AmountResponse, error) {

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	grant, err := s.ReadGrant(ctx, grant_id)
	if err != nil {
		return nil, fmt.Errorf("Grant %s does not exist", grant_id)
	}

	if grant.Status == "Revoked" {
		return nil, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	if !checkAwardee(grant.Awardee, userId) && !checkSubAwardee(grant.Awardee, userId) {
		return nil, fmt.Errorf("Awardee %s is not assigned in the Grant %s", userId, grant.ID)	
	}

	var cashedOut float64
	var requestedAmount float64
	for _, payment := range grant.Payment {
		if payment.Awardee_ID == userId {
			if payment.Status == "Accepted" || payment.Status == "Pending-redeem" {
				requestedAmount += payment.Total
			} else if payment.Status == "Accept_redeem" {
				cashedOut += payment.Total
			}
		}
	}


	response := AmountResponse {
		Cashed_Out:			cashedOut,
		Requested_Amount:	requestedAmount,
	}
	return &response, nil
}

// GetAllGrants for specific user returns all grants for grantors and awardee found in world state
func (s *SmartContract) GetAllGrantsUser(ctx contractapi.TransactionContextInterface) ([]Grant, error) {
	//requestPartialCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("grant", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var grants []Grant
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var grant Grant
		err = json.Unmarshal(queryResponse.Value, &grant)
		if err != nil {
			return nil, err
		}

		if checkAwardee(grant.Awardee, userId) || grant.Grantor_ID == userId || checkSubAwardee(grant.Awardee, userId) {
			grants = append(grants, grant)
		}
		
	}

	if len(grants)  == 0 {
		return []Grant{}, nil
	} else {
		return grants, nil
	}
	
}

// GetAllApprovedGrants for specific user returns all approved grants for awardee found in world state
func (s *SmartContract) GetAllApprovedGrants(ctx contractapi.TransactionContextInterface) ([]Grant, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	//requestPartialCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("grant", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var grants []Grant
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var grant Grant
		err = json.Unmarshal(queryResponse.Value, &grant)
		if err != nil {
			return nil, err
		}

		if (grant.Grantor_ID == userId || checkAwardee(grant.Awardee, userId) || checkSubAwardee(grant.Awardee, userId)) && grant.Status == "Approved" {
			grants = append(grants, grant)
		}
		
	}

	if len(grants)  == 0 {
		return []Grant{}, nil
	} else {
		return grants, nil
	}
	
}

// GetGrantsByStatus returns all grants with specific status
func (s *SmartContract) GetGrantsByStatus(ctx contractapi.TransactionContextInterface, status string) ([]Grant, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	//requestPartialCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("grant", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var grants []Grant
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var grant Grant
		err = json.Unmarshal(queryResponse.Value, &grant)
		if err != nil {
			return nil, err
		}

		if grant.Status == status && (checkAwardee(grant.Awardee, userId) || checkSubAwardee(grant.Awardee, userId) || grant.Grantor_ID == userId) {
			grants = append(grants, grant)
		}
		
	}

	if len(grants)  == 0 {
		return []Grant{}, nil
	} else {
		return grants, nil
	}
	
}

// GetGrantBenefits returns the benefits assigned for a grant with given id.
func (s *SmartContract) GetGrantBenefits(ctx contractapi.TransactionContextInterface, id string) ([]Benefit, error) {
	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grantJSON, err := ctx.GetStub().GetState(requestCompositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if grantJSON == nil {
		return nil, fmt.Errorf("the grant %s does not exist", id)
	}

	var grant Grant
	var benefit []Benefit
	err = json.Unmarshal(grantJSON, &grant)
	if err != nil {
		return nil, err
	}

	benefitJSON, err := json.Marshal(grant.Benefit)
	if err != nil {
		return nil,  fmt.Errorf("error marshaling json: %v", err)
	}

	err = json.Unmarshal(benefitJSON, &benefit)
	if err != nil {
		return nil, err
	}

	return benefit, nil
}

// GetPayments returns all the payments in a grant with given id.
func (s *SmartContract) GetPayments(ctx contractapi.TransactionContextInterface, id string) ([]Payment, error) {
	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grantJSON, err := ctx.GetStub().GetState(requestCompositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if grantJSON == nil {
		return nil, fmt.Errorf("the Grant %s does not exist", id)
	}

	var grant Grant
	var payment []Payment
	err = json.Unmarshal(grantJSON, &grant)
	if err != nil {
		return nil, err
	}

	paymentJSON, err := json.Marshal(grant.Payment)
	if err != nil {
		return nil,  fmt.Errorf("error marshaling json: %v", err)
	}

	err = json.Unmarshal(paymentJSON, &payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPaymentByStatus returns all the payments with given status in a grant for the user.
func (s *SmartContract) GetPaymentByStatus(ctx contractapi.TransactionContextInterface, status string) ([]PaymentStatus, error) {
    // convert string to json string

    jsonString := strings.ReplaceAll(status, "'", "\"")

    var statusStringArray []string
    json.Unmarshal([]byte(jsonString), &statusStringArray)

	var response PaymentStatus
	var responses []PaymentStatus
	var payments []Payment
	

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	// clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed getting the client's MSPID: %v", err)
	// }
	// if clientMSPID != GrantorMSP {
	// 	return nil, fmt.Errorf("User from org %v is not authorized to read grant payments", clientMSPID)
	// }

	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("grant", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var grant Grant
		err = json.Unmarshal(queryResponse.Value, &grant)
		if err != nil {
			return nil, err
		}

		for _, payment := range grant.Payment {
			for _, paymentStatus := range statusStringArray {
				if payment.Status == paymentStatus && (checkAwardee(grant.Awardee, userId) || checkSubAwardee(grant.Awardee, userId) || grant.Grantor_ID == userId) {
					payments = append(payments, payment)
				}
			}
			
		}

		if len(payments) != 0 {
			response = PaymentStatus {
				Grant_ID:	grant.ID,
				Payment:	payments,	
			}
	
			responses = append(responses, response)
		}
		
		
	}

	
	if len(responses)  == 0 {
		var nullResponse []PaymentStatus
		nullResponse = append(nullResponse, PaymentStatus{
		Grant_ID: "",
		Payment:  make([]Payment, 0),
	})
		return nullResponse, nil
	} else {
		return responses, nil
	}
	
	
}

// GetPaymentByStatusForAllGrants returns all the payments with given status in a grant.
func (s *SmartContract) GetPaymentByStatusForAllGrants(ctx contractapi.TransactionContextInterface, status []string) ([]PaymentStatus, error) {
	
	var response PaymentStatus
	var responses []PaymentStatus
	var payments []Payment

	// clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed getting the client's MSPID: %v", err)
	// }
	// if clientMSPID != GrantorMSP {
	// 	return nil, fmt.Errorf("User from org %v is not authorized to read grant payments", clientMSPID)
	// }

	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("grant", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var grant Grant
		err = json.Unmarshal(queryResponse.Value, &grant)
		if err != nil {
			return nil, err
		}

		for _, payment := range grant.Payment {
			for _, paymentStatus := range status {
				if payment.Status == paymentStatus {
					payments = append(payments, payment)
				}
			}
			
		}

		if len(payments) != 0 {
			response = PaymentStatus {
				Grant_ID:	grant.ID,
				Payment:	payments,	
			}
	
			responses = append(responses, response)
		}
		
	}

	if len(responses)  == 0 {
		var nullResponse []PaymentStatus
		nullResponse = append(nullResponse, PaymentStatus{
		Grant_ID: "",
		Payment:  make([]Payment, 0),
	})
		return nullResponse, nil
	} else {
		return responses, nil
	}
}

// GetPaymentByAwardee returns all the payments with given awardee id in a grant.
func (s *SmartContract) GetPaymentByAwardee(ctx contractapi.TransactionContextInterface, id string, awardeeId string) ([]Payment, error) {
	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grantJSON, err := ctx.GetStub().GetState(requestCompositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if grantJSON == nil {
		return nil, fmt.Errorf("the Grant %s does not exist", id)
	}

	var grant Grant
	err = json.Unmarshal(grantJSON, &grant)
	if err != nil {
		return nil, err
	}

	if !checkAwardee(grant.Awardee, awardeeId) && !checkSubAwardee(grant.Awardee, awardeeId) {
		return nil, fmt.Errorf("Awardee %s is not assigned in the Grant %s", awardeeId, grant.ID)	
	}

	var payments []Payment
	for _, payment := range grant.Payment {
		if payment.Awardee_ID == awardeeId {
			payments = append(payments, payment)
		}
	}

	if len(payments)  == 0 {
		return []Payment{}, nil
	} else {
		return payments, nil
	}
	
}

// GetProgress returns all the progresses in a grant with given id.
func (s *SmartContract) GetProgress(ctx contractapi.TransactionContextInterface, id string) ([]Progress, error) {
	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	grantJSON, err := ctx.GetStub().GetState(requestCompositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if grantJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var grant Grant
	err = json.Unmarshal(grantJSON, &grant)
	if err != nil {
		return nil, err
	}

	return grant.Progress, nil
}

// Get Remaining Amount
func (s *SmartContract) GetRemainingAmount(ctx contractapi.TransactionContextInterface, grant_id string) (float64, error) {
	
	grant, err := s.ReadGrant(ctx, grant_id)
	if err != nil {
		return 0.0, fmt.Errorf("Grant %s does not exist", grant_id)
	}

	if grant.Status == "Revoked" {
		return 0.0, fmt.Errorf("Grant %s is revoked", grant.ID)	
	}

	var remainingAmount float64
	remainingAmount = grant.Amount - grant.Cashed_Out

	return remainingAmount, nil
}


// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteGrant(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's ID: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		return false, fmt.Errorf("error: %v", err)
	}
	userId := strings.Split(string(data), ",")[0][9:]

	clientMSPID, err:= ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	if clientMSPID != GrantorMSP {
		return false, fmt.Errorf("User from org %v is not authorized to initiate grant", clientMSPID)
	}
	
	grant, err := s.ReadGrant(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if grant == nil {
		return false, fmt.Errorf("the grant %s doesn't exist", id)
	}

	if grant.Grantor_ID != userId {
		return false, fmt.Errorf("User %s from org %v is not authorized to assign grant for this grant %s",userId, clientMSPID, grant.ID)
	}

	requestCompositeKey, _ := ctx.GetStub().CreateCompositeKey("grant", []string{id})
	err = ctx.GetStub().DelState(requestCompositeKey)
	if err != nil {
		return false, fmt.Errorf("Deleting Grant failed", err)
	}

	return true, nil
}

func checkBenefitAmount(itemAmount float64, benefitAmount float64, totalBenefit float64, awardeeType string, percentage float64) (bool) {
	switch awardeeType {
		case "Main":
			if itemAmount <= benefitAmount && totalBenefit <= benefitAmount {
				return true
			} else {
				return false
			}
		case "Sub":
			if itemAmount <= benefitAmount && totalBenefit <= (percentage*benefitAmount)/100 {
				return true
			} else {
				return false
			}

	}

	return false
}


func checkAwardee(awardees []Awardee, userId string) (bool) {
	flag := false
	for _, awardee := range awardees {
		if awardee.ID == userId && awardee.Awardee_Type == "Main" {
			flag = true
			break
		}
	}
	return flag
}

func checkSubAwardee(awardees []Awardee, userId string) (bool) {
	flag := false
	for _, awardee := range awardees {
		if awardee.ID == userId && awardee.Awardee_Type == "Sub" {
			flag = true
			break
		}
	}
	return flag
}

func checkPayment(payments []Payment, payment_id string) (bool) {
	flag := false
	for _, payment := range payments {
		if payment.ID == payment_id {
			flag = true
			break
		}
	}
	return flag
}




