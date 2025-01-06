package main

import (
	"context"
  "time"
	//"fmt"

	jup "github.com/scatkit/jupswap/jupcore"
	"github.com/scatkit/jupswap/jup_client"
	"github.com/scatkit/pumpdexer/solana"
	//"github.com/davecgh/go-spew/spew"
)

const URL = "https://quote-api.jup.ag/v6"

func main() {
	jupiterClient, err := jup.NewClient(URL)
	if err != nil {
		panic(err)
	}

	slippageBps := 250

	quoteResp, err := jupiterClient.GetQuoteWithResponse(context.Background(), &jup.GetQuoteParams{
		InputMint:   "So11111111111111111111111111111111111111112",
		OutputMint:  "WENWENvqqNya429ubCdR81ZmD69brwQaaBYY6p3LCpk",
		Amount:      100000,
		SlippageBps: &slippageBps,
	})

	if err != nil {
		panic(err)
	}
	if quoteResp.JSON200 == nil {
		panic("invalid GetQuoteResponse response")
	}
	quote := quoteResp.JSON200

	// Setting prioritization fees to `Auto`
	prioritizationFeeLamports := jup.SwapRequest_PrioritizationFeeLamports{}
	if err := prioritizationFeeLamports.UnmarshalJSON([]byte(`"auto"`)); err != nil {
		panic(err)
	}
  

	// When enabled, it will do a swap simulation to get the compute unit used and set it in ComputeBudget's compute unit limit.
	// This will increase latency slightly since there will be one extra RPC call to simulate this. Default is false.
	var DCUL = true
	userPubKey := solana.MustPubkeyFromBase58("CJMTJWF97jd3dspsN5qhPp4EpKBHMTnkRvDkpSHUWSGJ")

	swapResp, err := jupiterClient.PostSwapWithResponse(context.TODO(), jup.SwapRequest{
		PrioritizationFeeLamports: &prioritizationFeeLamports,
		QuoteResponse:             *quote,
		UserPublicKey:             userPubKey.String(),
		DynamicComputeUnitLimit:   &DCUL,
	})

	if err != nil {
		panic(err)
	}

	if swapResp.JSON200 == nil {
		panic("invalid PostSwapWithResponse{} response")
	}

	swap := swapResp.JSON200
  //spew.Dump(swap)
  
  wallet, err := jup_client.NewWalletFromPrivateKeyBase58("5rg7jXrAYXoAYt1ARV1RzuRFCsH948MyjMjKVG8Kiw7pdZZ7QBjuJnEfufvukPJ5hLyRHUXkPBuc9mP7AS35i5yC")
  if err != nil{
    panic(err)
  }
  jClient,err:= jup_client.NewClient(wallet, "https://api.mainnet-beta.solana.com") 
  if err != nil{
    panic(err)
  }
  
  signedTx, err := jClient.SendTransactionOnChain(context.Background(), swap.SwapTransaction)
  if err != nil{
    panic(err)
  }
  
  time.Sleep(20 * time.Second)
  
  _, err = jClient.CheckSignature(context.Background(), signedTx)
  if err != nil{
    panic(err)
  }
}
