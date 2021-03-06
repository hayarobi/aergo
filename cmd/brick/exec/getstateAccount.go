package exec

import (
	"fmt"
	"math/big"

	"github.com/aergoio/aergo/cmd/brick/context"
	"github.com/aergoio/aergo/contract"
)

func init() {
	registerExec(&getStateAccount{})
}

type getStateAccount struct{}

func (c *getStateAccount) Command() string {
	return "getstate"
}

func (c *getStateAccount) Syntax() string {
	return fmt.Sprintf("%s", context.AccountSymbol)
}

func (c *getStateAccount) Usage() string {
	return fmt.Sprintf("getstate <account_name> `[expected_balance]`")
}

func (c *getStateAccount) Describe() string {
	return "create an account with a given amount of balance"
}

func (c *getStateAccount) Validate(args string) error {
	if context.Get() == nil {
		return fmt.Errorf("load chain first")
	}

	_, _, err := c.parse(args)

	return err
}

func (c *getStateAccount) parse(args string) (string, string, error) {
	splitArgs := context.SplitSpaceAndAccent(args, false)
	if len(splitArgs) < 1 {
		return "", "", fmt.Errorf("need an arguments. usage: %s", c.Usage())
	}

	expectedResult := ""
	if len(splitArgs) == 2 {
		expectedResult = splitArgs[1].Text
	}

	return splitArgs[0].Text, expectedResult, nil
}

func (c *getStateAccount) Run(args string) (string, error) {
	accountName, expectedResult, _ := c.parse(args)

	state, err := context.Get().GetAccountState(accountName)

	if err != nil {
		return "", err
	}
	if expectedResult == "" {
		return fmt.Sprintf("%s = %d", contract.StrToAddress(accountName), new(big.Int).SetBytes(state.GetBalance())), nil
	} else {
		strRet := fmt.Sprintf("%d", new(big.Int).SetBytes(state.GetBalance()))
		if expectedResult == strRet {
			return "state compare successfully", nil
		} else {
			return "", fmt.Errorf("state compre fail. Expected: %s, Actual: %s", expectedResult, strRet)
		}
	}
}
