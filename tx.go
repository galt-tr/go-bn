package bn

import (
	"context"
	"encoding/json"
	imodels "github.com/libsv/go-bn/internal/models"
	"github.com/libsv/go-bn/models"
	"github.com/libsv/go-bt/v2"
	"log"
)

// TransactionClient interfaces interaction with the transaction sub commands on a bitcoin node.
type TransactionClient interface {
	AddToConfiscationTransactionWhitelist(ctx context.Context, funds []models.ConfiscationTransactionDetails) (*models.AddToConfiscationTransactionWhitelistResponse, error)
	AddToConsensusBlacklist(ctx context.Context, funds []models.Fund) (*models.BlacklistResponse, error)
	RemoveFromPolicyBlacklist(ctx context.Context, funds []models.Fund) (*models.BlacklistResponse, error)
	CreateRawTransaction(ctx context.Context, utxos bt.UTXOs, params models.ParamsCreateRawTransaction) (*bt.Tx, error)
	FundRawTransaction(ctx context.Context, tx *bt.Tx,
		opts *models.OptsFundRawTransaction) (*models.FundRawTransaction, error)
	RawTransaction(ctx context.Context, txID string) (*bt.Tx, error)
	SignRawTransaction(ctx context.Context, tx *bt.Tx,
		opts *models.OptsSignRawTransaction) (*models.SignedRawTransaction, error)
	SendRawTransaction(ctx context.Context, tx *bt.Tx, opts *models.OptsSendRawTransaction) (string, error)
	SendRawTransactions(ctx context.Context,
		params ...models.ParamsSendRawTransactions) (*models.SendRawTransactionsResponse, error)
}

// NewTransactionClient returns a client only capable of interfacing with the transaction sub commands
// on a bitcoin node.
func NewTransactionClient(oo ...BitcoinClientOptFunc) TransactionClient {
	return NewNodeClient(oo...)
}

func (c *client) CreateRawTransaction(ctx context.Context, utxos bt.UTXOs,
	params models.ParamsCreateRawTransaction) (*bt.Tx, error) {
	params.SetIsMainnet(c.isMainnet)
	var resp string
	if err := c.rpc.Do(ctx, "createrawtransaction", &resp, c.argsFor(&params, utxos.NodeJSON())...); err != nil {
		return nil, err
	}
	return bt.NewTxFromString(resp)
}

func (c *client) FundRawTransaction(ctx context.Context, tx *bt.Tx,
	opts *models.OptsFundRawTransaction) (*models.FundRawTransaction, error) {
	resp := imodels.InternalFundRawTransaction{FundRawTransaction: &models.FundRawTransaction{}}
	return resp.FundRawTransaction, c.rpc.Do(ctx, "fundrawtransaction", &resp, c.argsFor(opts, tx.String())...)
}

func (c *client) RawTransaction(ctx context.Context, txID string) (*bt.Tx, error) {
	var resp bt.Tx
	return &resp, c.rpc.Do(ctx, "getrawtransaction", &resp, txID, true)
}

func (c *client) SignRawTransaction(ctx context.Context, tx *bt.Tx,
	opts *models.OptsSignRawTransaction) (*models.SignedRawTransaction, error) {
	var resp imodels.InternalSignRawTransaction
	return resp.SignedRawTransaction, c.rpc.Do(ctx, "signrawtransaction", &resp, c.argsFor(opts, tx.String())...)
}

func (c *client) SendRawTransaction(ctx context.Context, tx *bt.Tx,
	opts *models.OptsSendRawTransaction) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "sendrawtransaction", &resp, c.argsFor(opts, tx.String())...)
}

func (c *client) SendRawTransactions(ctx context.Context,
	params ...models.ParamsSendRawTransactions) (*models.SendRawTransactionsResponse, error) {
	var resp models.SendRawTransactionsResponse
	return &resp, c.rpc.Do(ctx, "sendrawtransactions", &resp, params)
}

func (c *client) AddToConsensusBlacklist(ctx context.Context, funds []models.Fund) (*models.BlacklistResponse, error) {
	var resp models.BlacklistResponse
	req := models.BlacklistArgs{Funds: funds}
	return &resp, c.rpc.Do(ctx, "addToConsensusBlacklist", &resp, req)
}

func (c *client) RemoveFromPolicyBlacklist(ctx context.Context, funds []models.Fund) (*models.BlacklistResponse, error) {
	var resp models.BlacklistResponse
	req := models.BlacklistArgs{Funds: funds}
	return &resp, c.rpc.Do(ctx, "removeFromPolicyBlacklist", &resp, req)
}

func (c *client) AddToConfiscationTransactionWhitelist(ctx context.Context, confiscationTransactions []models.ConfiscationTransactionDetails) (*models.AddToConfiscationTransactionWhitelistResponse, error) {
	var resp models.AddToConfiscationTransactionWhitelistResponse
	req := models.AddToConfiscationTxIdWhitelistArgs{
		confiscationTransactions,
	}
	a, _ := json.Marshal(req)
	log.Printf("%s", a)
	return &resp, c.rpc.Do(ctx, "addToConfiscationTxidWhitelist", &resp, req)
}
