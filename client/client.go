package main

import (
	"encoding/json"

	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type Loan struct {
	LoanID        string
	ApplicantName string
	LoanAmount    float64
	TermMonths    int
	InterestRate  float64
	Outstanding   float64
	Status        string
	Repayments    []float64
}

func main() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"fabric-samples",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

    channelName := "mychannel"
    if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
        channelName = cname
    }

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

    chaincodeName := "loan"
    if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
        chaincodeName = ccname
    }

	contract := network.GetContract(chaincodeName)

	// TODO: Call ApplyForLoan
	_, err = contract.SubmitTransaction("ApplyForLoan", "loan1", "John Doe", "5000", "12", "5.5")
	if err != nil {
		log.Fatalf("Failed to apply for loan: %v", err)
	}
	fmt.Println("Loan successfully applied")

	_, err = contract.SubmitTransaction("ApproveLoan", "loan1", "Approved")
	if err != nil {
		log.Fatalf("Failed to approve loan: %v", err)
	}
	fmt.Println("Loan status updated to Approved")

	_, err = contract.SubmitTransaction("MakeRepayment", "loan1", "1000")
	if err != nil {
		log.Fatalf("Failed to make repayment for loan: %v", err)
	}
	fmt.Println("Repayment recorded. Outstanding balance updated.")

	// TODO: Call CheckLoanBalance
	result, err := contract.EvaluateTransaction("CheckLoanBalance", "loan1")
	if err != nil {
		log.Fatalf("Failed to get loan balance: %v", err)
	}

	var loan Loan
	err = json.Unmarshal(result, &loan)
	if err != nil {
		log.Fatalf("Failed to unmarshal loan: %v", err)
	}

	fmt.Printf("Outstanding Balance: %f\n", loan.Outstanding)
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"fabric-samples",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := os.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}

