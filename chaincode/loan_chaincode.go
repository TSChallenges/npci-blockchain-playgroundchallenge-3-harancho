package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type LoanContract struct {
	contractapi.Contract
}

type Loan struct {
	LoanID        string    `json:"loanID"`
	ApplicantName string    `json:"applicantName"`
	LoanAmount    float64   `json:"loanAmount"`
	TermMonths    int       `json:"termMonths"`
	InterestRate  float64   `json:"interestRate"`
	Outstanding   float64   `json:"outstanding"`
	Status        string    `json:"status"`
	Repayments    []float64 `json:"repayments"`
}

// TODO: Implement ApplyForLoan
func (c *LoanContract) ApplyForLoan(ctx contractapi.TransactionContextInterface, loanID, applicantName string, loanAmount float64, termMonths int, interestRate float64) error {

	data, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve loanID: %v", err)
	}
	if data != nil {
		fmt.Println("loanID already exists")
		return fmt.Errorf("loanID already exists")
	}

	if len(loanID) == 0 {
		fmt.Println("loan Id cannot be null")
		return errors.New("loan Id cannot be null")
	}

	if len(applicantName) == 0 {
		fmt.Println("applicantName cannot be null")
		return errors.New("applicantName cannot be null")
	}

	if loanAmount <= 0 {
		fmt.Println("loanAmount cannot be 0 or less")
		return errors.New("loanAmount cannot be 0 or less")
	}

	if termMonths <= 0 {
		fmt.Println("termMonths cannot be 0 or less")
		return errors.New("termMonths cannot be 0 or less")
	}

	if interestRate <= 0 {
		fmt.Println("interestRate cannot be 0 or less")
		return errors.New("interestRate cannot be 0 or less")
	}

	var loandPayload = Loan{
		LoanID:        loanID,
		ApplicantName: applicantName,
		LoanAmount:    loanAmount,
		TermMonths:    termMonths,
		InterestRate:  interestRate,
		Outstanding:   loanAmount,
		Status:        "Pending",
	}

	loanPayloadInByte, err := json.Marshal(loandPayload)
	if err != nil {
		fmt.Printf("error in marshalling data %v \n", err)
		return err
	}

	err = ctx.GetStub().PutState(loanID, loanPayloadInByte)
	if err != nil {
		return fmt.Errorf("Failed to store loan request data: %v", err)
	}
	return nil
}

// TODO: Implement ApproveLoan
func (c *LoanContract) ApproveLoan(ctx contractapi.TransactionContextInterface, loanID string, status string) error {

	data, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve loanID: %v", err)
	}
	if data == nil {
		return fmt.Errorf("loanID not found")
	}

	var loandPayload Loan
	err = json.Unmarshal(data, &loandPayload)
	if err != nil {
		fmt.Printf("error in unmarshalling loan payload %v \n", err)
		return err
	}

	if loandPayload.Status != "Pending" {
		fmt.Printf("cannot approve non Pending Loans \n")
		return errors.New("cannot approve non Pending Loans")
	}

	loandPayload.Status = status

	loanPayloadInByte, err := json.Marshal(loandPayload)
	if err != nil {
		fmt.Printf("error in marshalling data %v \n", err)
		return err
	}

	err = ctx.GetStub().PutState(loanID, loanPayloadInByte)
	if err != nil {
		return fmt.Errorf("Failed to update loan status data: %v", err)
	}

	return nil
}

// TODO: Implement MakeRepayment
func (c *LoanContract) MakeRepayment(ctx contractapi.TransactionContextInterface, loanID string, repaymentAmount float64) error {

	data, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve loanID: %v", err)
	}
	if data == nil {
		return fmt.Errorf("loanID not found")
	}

	var loandPayload Loan
	err = json.Unmarshal(data, &loandPayload)
	if err != nil {
		fmt.Printf("error in unmarshalling loan payload %v \n", err)
		return err
	}

	if loandPayload.Status == "Pending" {
		fmt.Printf("cannot make repayment for Pending Loans \n")
		return errors.New("cannot make repayment for Pending Loans")
	}

	loandPayload.Outstanding = loandPayload.Outstanding - repaymentAmount
	loandPayload.Repayments = append(loandPayload.Repayments, repaymentAmount)

	loanPayloadInByte, err := json.Marshal(loandPayload)
	if err != nil {
		fmt.Printf("error in marshalling data %v \n", err)
		return err
	}

	err = ctx.GetStub().PutState(loanID, loanPayloadInByte)
	if err != nil {
		return fmt.Errorf("Failed to make repayment: %v", err)
	}

	return nil
}

// TODO: Implement CheckLoanBalance
func (c *LoanContract) CheckLoanBalance(ctx contractapi.TransactionContextInterface, loanID string) (*Loan, error) {

	data, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		return &Loan{}, fmt.Errorf("Failed to retrieve loanID: %v", err)
	}
	if data == nil {
		return &Loan{}, fmt.Errorf("loanID not found")
	}

	var loandPayload Loan
	err = json.Unmarshal(data, &loandPayload)
	if err != nil {
		fmt.Printf("error in unmarshalling loan payload %v \n", err)
		return &Loan{}, err
	}

	return &loandPayload, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(LoanContract))
	if err != nil {
		fmt.Printf("Error creating loan chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting loan chaincode: %s", err)
	}
}
