package responsecredit

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing response credit

type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

// Issue creates a new response credit and stores it in the world state
func (c *Contract) Issue(ctx TransactionContextInterface, issuer string, creditNumber string, issueDateTime string) (*ResponseCredit, error) {
	credit := ResponseCredit{CreditNumber: creditNumber, Issuer: issuer, IssueDateTime: issueDateTime, Owner: issuer}
	credit.SetIssued()

	err := ctx.GetCreditList().AddCredit(&credit)

	if err != nil {
		return nil, err
	}

	return &credit, nil
}

// Buy updates a response credit to be in trading status and sets the new owner
func (c *Contract) Buy(ctx TransactionContextInterface, issuer string, creditNumber string, currentOwner string, newOwner string, price int, purchaseDateTime string) (*ResponseCredit, error) {
	credit, err := ctx.GetCreditList().GetCredit(issuer, creditNumber)

	if err != nil {
		return nil, err
	}

	if credit.Owner != currentOwner {
		return nil, fmt.Errorf("Credit %s:%s is not owned by %s", issuer, creditNumber, currentOwner)
	}

	if credit.IsIssued() {
		credit.SetTrading()
	}

	if !credit.IsTrading() {
		return nil, fmt.Errorf("Credit %s:%s is not trading. Current state = %s", issuer, creditNumber, credit.GetState())
	}

	credit.Owner = newOwner

	err = ctx.GetCreditList().UpdateCredit(credit)

	if err != nil {
		return nil, err
	}

	return credit, nil
}

// Redeem updates a response credit status to be redeemed
func (c *Contract) Redeem(ctx TransactionContextInterface, issuer string, creditNumber string, redeemingOwner string, redeenDateTime string) (*ResponseCredit, error) {
	credit, err := ctx.GetCreditList().GetCredit(issuer, creditNumber)

	if err != nil {
		return nil, err
	}

	if credit.Owner != redeemingOwner {
		return nil, fmt.Errorf("Credit %s:%s is not owned by %s", issuer, creditNumber, redeemingOwner)
	}

	if credit.IsRedeemed() {
		return nil, fmt.Errorf("Credit %s:%s is already redeemed", issuer, creditNumber)
	}

	credit.Owner = credit.Issuer
	credit.SetRedeemed()

	err = ctx.GetCreditList().UpdateCredit(credit)

	if err != nil {
		return nil, err
	}

	return credit, nil
}