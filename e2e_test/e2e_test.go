package e2e_test

import (
	"testing"

	"github.com/celo-org/celo-blockchain/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// This statement is commented out but left here since its very useful for
	// debugging problems and its non trivial to construct.
	//
	// log.Root().SetHandler(log.LvlFilterHandler(log.LvlWarn, log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
}

func TestNewNetwork(t *testing.T) {
	accounts := test.Accounts(3)
	gc := test.GenesisConfig(accounts)
	network, err := test.NewNetwork(accounts, gc)
	require.NoError(t, err)
	assert.Equal(t, 3, len(network))
}

// This test starts a network submits a transaction and waits for the whole
// network to process the transaction.
// func TestSendCelo(t *testing.T) {
// 	accounts := test.Accounts(3)
// 	gc := test.GenesisConfig(accounts)
// 	network, err := test.NewNetwork(accounts, gc)
// 	require.NoError(t, err)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
// 	defer cancel()

// 	// Send 1 celo from the dev account attached to node 0 to the dev account
// 	// attached to node 1.
// 	tx, err := network[0].SendCelo(ctx, network[1].DevAddress, 1)
// 	require.NoError(t, err)
// 	// Wait for the whole network to process the transaction.
// 	err = network.AwaitTransactions(ctx, tx)
// 	require.NoError(t, err)
// 	println("-------------------------------------------------------------------")
// 	println("-------------------------------------------------------------------")
// 	println("-------------------------------------------------------------------")
// 	println("-------------------------------------------------------------------")
// 	println("-------------------------------------------------------------------")
// 	println("-------------------------------------------------------------------")
// 	println("-------------------------------------------------------------------")
// 	network.Shutdown()

// 	// time.Sleep(5 * time.Minute)

// }
