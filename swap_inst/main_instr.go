package main

import (
	"context"
  //"time"
  "fmt"
  "log"
  "encoding/json"
  "encoding/base64"

  //bin "github.com/gagliardetto/binary"
  alt "github.com/scatkit/pumpdexer/programs/address-lookup-table"
  "github.com/scatkit/pumpdexer/rpc"
	jup "github.com/scatkit/jupswap/jupcore"
	//"github.com/scatkit/jupswap/jup_client"
	"github.com/scatkit/pumpdexer/solana"
	"github.com/davecgh/go-spew/spew"
  dec "github.com/shopspring/decimal"

)

const URL = "https://quote-api.jup.ag/v6"

type altOption map[solana.PublicKey]*alt.AddressLookupTableState

type dumbTransactionInstructions struct {
	accounts  []*solana.AccountMeta
	data      []byte
	programID solana.PublicKey
}

func (t *dumbTransactionInstructions) Accounts() []*solana.AccountMeta{
	return t.accounts
}

func (t *dumbTransactionInstructions) ProgramID() solana.PublicKey{
	return t.programID
}

func (t *dumbTransactionInstructions) Data() ([]byte, error) {
	return t.data, nil
}

func formatAmount(amountStr string, decimalNum uint) (formattedAmount int64, err error){
  amountFloat, err  := dec.NewFromString(amountStr)
  if err != nil{
    return 0, fmt.Errorf("invalid float string: %s", amountStr)
  }
  // Base 10^decimalNum 
  smallestUnitOfToken := dec.NewFromInt(10).Pow(dec.NewFromInt(int64(decimalNum)))
  
  return amountFloat.Mul(smallestUnitOfToken).IntPart(), nil
}

func main() {
  userPubkey := solana.MustPubkeyFromBase58("CJMTJWF97jd3dspsN5qhPp4EpKBHMTnkRvDkpSHUWSGJ")
  // maybe not return an error
	jupiterClient, err := jup.NewClient(URL)
	if err != nil {
		panic(err)
	}
  solClient := rpc.New("https://api.mainnet-beta.solana.com")
  
  // Selling case
  var inputTokenMint = "So11111111111111111111111111111111111111112"
  var outputTokenMint = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
  var inputAmount = "0.0001"
  var slippageBps = 150
  
  var decimalPrec uint
  if inputTokenMint == "So11111111111111111111111111111111111111112"{
    decimalPrec = uint(9) // should do this as constant
  } else{
    tokenAccountBalanceResponse, err := solClient.GetTokenSupply(
      context.Background(), 
      solana.MustPubkeyFromBase58(inputTokenMint),
      rpc.CommitmentFinalized)
      if err != nil{
        log.Fatal(err)
      }
      decimalPrec = uint(tokenAccountBalanceResponse.Value.Decimals) // casts uint8 to uint
    } 
  formattedUserInputAmount, err := formatAmount(inputAmount, decimalPrec)
  if err != nil{
    log.Fatal(err)
  }

  // JUP 
  //var asLegacy = true
  quoteResp, err := jupiterClient.GetQuoteWithResponse(context.Background(), &jup.GetQuoteParams{
    InputMint:   inputTokenMint,
    OutputMint:  outputTokenMint, 
    Amount:      formattedUserInputAmount,
    SlippageBps: &slippageBps,
    //AsLegacyTransaction:  &asLegacy,
  })

  if err != nil {
    panic(err)
  }
  if quoteResp.JSON200 == nil {
    panic("invalid GetQuoteResponse response")
  }
  quote := quoteResp.JSON200
  
  dynamicComputeUnitLimit := true 
  prioritizationFeeLamports := jup.SwapRequest_PrioritizationFeeLamports{}
  jtip := map[string]int64{
    "jitoTipLamports": 10000,
  }
  bts,err := json.Marshal(jtip)
  if err != nil{
    panic(err)
  }
  if err = prioritizationFeeLamports.UnmarshalJSON(bts); err != nil {
		panic(err)
	}
  swapInstrResp, err := jupiterClient.PostSwapInstructionsWithResponse(context.TODO(), jup.SwapInstructionsRequest{
    PrioritizationFeeLamports: &prioritizationFeeLamports,
  	QuoteResponse:             *quote,
  	UserPublicKey:             userPubkey.String(),
    //AsLegacyTransaction:       &asLegacy, 
    DynamicComputeUnitLimit: &dynamicComputeUnitLimit,
  })
  
  if err != nil{
    log.Fatal(err)
  }
  
  if swapInstrResp.JSON200 == nil {
  	panic("invalid PostSwapWithResponse{} response")
  }
  swapInstr := swapInstrResp.JSON200
  spew.Dump(swapInstr)
  var output = dumbTransactionInstructions{}
  err = deserializeInstruction(swapInstr.SwapInstruction, &output)
  if err != nil{
    log.Fatal(err)
  }
  //spew.Dump(output)

  // DO NOT DELET:
  
  //accs := []solana.PublicKey{}
  //for _,acc := range swapInstr.AddressLookupTableAddresses{
  //  accs = append(accs, solana.MustPubkeyFromBase58(acc))
  //}
  //
  //fmt.Println(len(accs))
  //res_accs, err := solClient.GetMultipleAccounts(context.Background(), accs...) 
  //
  //altOpt := map[solana.PublicKey]*alt.AddressLookupTableState{}
  //for i,acc := range res_accs.Value{
  //  res, err := alt.DecodeAddressLookupTableState(acc.Data.GetBinary())
  //  if err != nil{
  //    panic(err)
  //  }
  //  altOpt[accs[i]] = res
  //}
  //
  //spew.Dump(altOpt)
}

func deserializeALTs(altAddresses []string,
) (altOption, error){
  accs := []solana.PublicKey{}
  for _,acc := range altAddresses{
    accs = append(accs, solana.MustPubkeyFromBase58(acc))
  }
  
  fmt.Println(len(accs))
  res_accs, err := solClient.GetMultipleAccounts(context.Background(), accs...) 
  
  altOpt := map[solana.PublicKey]*alt.AddressLookupTableState{}
  for i,acc := range res_accs.Value{
    res, err := alt.DecodeAddressLookupTableState(acc.Data.GetBinary())
    if err != nil{
      panic(err)
    }
    altOpt[accs[i]] = res
  }
  
  spew.Dump(altOpt)

  
}

func deserializeInstruction(inst jup.Instruction, out *dumbTransactionInstructions,
) error{ 
  if out == nil{
    return fmt.Errorf("TransactionInstructions cannot be `nil`")
  }
  formattedAccs := make([]*solana.AccountMeta, len(inst.Accounts))
  for i,acc := range inst.Accounts{
    formattedAccs[i] = &solana.AccountMeta{
      PublicKey: solana.MustPubkeyFromBase58(acc.Pubkey),
      IsSigner: acc.IsSigner,
      IsWritable: acc.IsWritable,
    }
  }
  out.accounts = formattedAccs
  dataBytes, err := base64.StdEncoding.DecodeString(inst.Data)
  if err != nil{
    return fmt.Errorf("Failed to decode data (%s) from instruction: %w",inst.Data, err)
  }
  out.data = dataBytes
  out.programID = solana.MustPubkeyFromBase58(inst.ProgramId)
  
  return nil
}

