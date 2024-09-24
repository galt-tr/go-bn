package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/libsv/go-bn/models"

	"github.com/libsv/go-bn"
)

func main() {
	c := bn.NewNodeClient(
		bn.WithHost("http://localhost:8333"),
		bn.WithCreds("galt", "galt"),
	)
	ctx := context.Background()

	funds := []models.Fund{
		{
			TxOut: models.TxOut{
				TxId: "bfed856c469f4b115a56fad10486a6082ffa7e845f058db542973bca6fefeaff",
				Vout: 0,
			},
			EnforceAtHeight: []models.Enforce{
				{
					Start: 13022,
					Stop:  123100,
				},
			},
			PolicyExpiresWithConsensus: false,
		},
	}

	resp, err := c.AddToConsensusBlacklist(ctx, funds)
	if err != nil {
		panic(err)
	}

	bb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bb))
}
