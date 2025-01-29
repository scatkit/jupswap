package jupcore
import(
  "encoding/json"
  "net/http"
  "context"
  "bytes"
  "io"
  "strings"
  //"github.com/davecgh/go-spew/spew"
  "fmt"
  "net/url"
) 

type SwapRequest_PrioritizationFeeLamports struct{
  union json.RawMessage
}

func (pFee *SwapRequest_PrioritizationFeeLamports) UnmarshalJSON(b []byte) error{
  err := pFee.union.UnmarshalJSON(b)
  return err
}

func (pFee SwapRequest_PrioritizationFeeLamports) MarshalJSON() ([]byte, error) {
	b, err := pFee.union.MarshalJSON()
	return b, err
}

// SwapRequest_ComputeUnitPriceMicroLamports The compute unit price to prioritize the transaction, the additional fee will be `computeUnitLimit (1400000) * computeUnitPriceMicroLamports`. If `auto` is used, Jupiter will automatically set a priority fee and it will be capped at 5,000,000 lamports / 0.005 SOL.
type SwapRequest_ComputeUnitPriceMicroLamports struct {
  union json.RawMessage
}

func (t SwapRequest_ComputeUnitPriceMicroLamports) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *SwapRequest_ComputeUnitPriceMicroLamports) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

type SwapResponse struct{
  SwapTransaction           string   `json:"swapTransaction"`
  LastValidBlockHeight      float32  `json:"lastValidBlockHeight"`
  PrioritizationFeeLamports *float32 `json:"prioritizationFeeLamprots"`
  DynamicSlippageReport     *struct{
    AmplificationRatio           *string                                        `json:"amplificationRatio,omitempty"`
		CategoryName                 *SwapResponseDynamicSlippageReportCategoryName `json:"categoryName,omitempty"`
		HeuristicMaxSlippageBps      *int                                           `json:"heuristicMaxSlippageBps,omitempty"`
		OtherAmount                  *int                                           `json:"otherAmount,omitempty"`
		SimulatedIncurredSlippageBps *int                                           `json:"simulatedIncurredSlippageBps,omitempty"`
		SlippageBps                  *int                                           `json:"slippageBps,omitempty"`
  }
} 

type SwapResponseDynamicSlippageReportCategoryName string

const (
  Bluechip SwapResponseDynamicSlippageReportCategoryName = "bluechip"
  Lst      SwapResponseDynamicSlippageReportCategoryName = "lst"
	Stable   SwapResponseDynamicSlippageReportCategoryName = "stable"
	Verified SwapResponseDynamicSlippageReportCategoryName = "verified"
)

type PostSwapResponse struct{
  Body          []byte
  HTTPResponse  *http.Response
  JSON200       *SwapResponse
}

type SwapRequest struct{
  // AllowOptimizedWrappedSolTokenAccount Default is false. Enabling it would reduce use an optimized way to open WSOL that reduce compute unit.
	AllowOptimizedWrappedSolTokenAccount *bool `json:"allowOptimizedWrappedSolTokenAccount,omitempty"`
  
  // AsLegacyTransaction Default is false.
  // Request a legacy transaction rather than the default versioned transaction, 
  // needs to be paired with a quote using asLegacyTransaction otherwise the transaction might be too large
	AsLegacyTransaction *bool `json:"asLegacyTransaction,omitempty"`
  
  // BlockhashSlotsToExpiry Optional. When passed in, Swap object will be returned with your desired slots to epxiry.
	BlockhashSlotsToExpiry *float32 `json:"blockhashSlotsToExpiry,omitempty"`
  
  // ComputeUnitPriceMicroLamports The compute unit price to prioritize the transaction, 
  // the additional fee will be `computeUnitLimit (1400000) * computeUnitPriceMicroLamports`. 
  // If `auto` is used, Jupiter will automatically set a priority fee and it will be capped at 5,000,000 lamports / 0.005 SOL.
	ComputeUnitPriceMicroLamports *SwapRequest_ComputeUnitPriceMicroLamports `json:"computeUnitPriceMicroLamports,omitempty"`
  
  // CorrectLastValidBlockHeight Optional. Default to false. Request Swap object to be returned with the correct blockhash prior to Agave 2.0.
	CorrectLastValidBlockHeight *bool `json:"correctLastValidBlockHeight,omitempty"`

  // DestinationTokenAccount Public key of the token account that will be used to receive the token out of the swap.
  // If not provided, the user's ATA will be used. If provided, we assume that the token account is already initialized.
	DestinationTokenAccount *string `json:"destinationTokenAccount,omitempty"`
  
  // DynamicComputeUnitLimit When enabled, it will do a swap simulation to get the compute unit used and set it in ComputeBudget's compute unit limit. 
  // This will increase latency slightly since there will be one extra RPC call to simulate this. Default is `false`.
	DynamicComputeUnitLimit *bool `json:"dynamicComputeUnitLimit,omitempty"`
	DynamicSlippage         *struct {
		MaxBps *int `json:"maxBps,omitempty"`
		MinBps *int `json:"minBps,omitempty"`
	} `json:"dynamicSlippage,omitempty"`
  
  // Fee token account, same as the output token for ExactIn and as the input token for ExactOut, 
  // it is derived using the seeds = ["referral_ata", referral_account, mint] and the `REFER4ZgmyYx9c6He5XfaTMiGfdLwRnkV4RPp9t9iF3`
  // referral contract (only pass in if you set a feeBps and make sure that the feeAccount has been created).
	FeeAccount *string `json:"feeAccount,omitempty"`

	PrioritizationFeeLamports *SwapRequest_PrioritizationFeeLamports `json:"prioritizationFeeLamports,omitempty"`
  
  // ProgramAuthorityId The program authority id [0;7], load balanced across the available set by default
	ProgramAuthorityId *int          `json:"programAuthorityId,omitempty"`
	QuoteResponse      QuoteResponse `json:"quoteResponse"`

  // When enabled, it will not do any rpc calls check on user's accounts. 
  // Enable it only when you already setup all the accounts needed for the trasaction, like wrapping or unwrapping sol, destination account is already created.
	SkipUserAccountsRpcCalls *bool `json:"skipUserAccountsRpcCalls,omitempty"`  
  
  // Default is true. This enables the usage of shared program accountns.
  // That means no intermediate token accounts or open orders accounts need to be created for the users.
  // But it also means that the likelihood of hot accounts is higher.
	UseSharedAccounts *bool `json:"useSharedAccounts,omitempty"`
  
  // Default is false. This is useful when the instruction before the swap has a transfer that increases the input token amount.
  // Then, the swap will just use the difference between the token ledger token amount and post token amount.
	UseTokenLedger *bool `json:"useTokenLedger,omitempty"`

    // UserPublicKey The user public key.
	UserPublicKey string `json:"userPublicKey"`

	// Default is true. If true, will automatically wrap/unwrap SOL. 
  // If false, it will use wSOL token account.  Will be ignored if `destinationTokenAccount` is set because the `destinationTokenAccount` 
  // may belong to a different user that we have no authority to close.
	WrapAndUnwrapSol *bool `json:"wrapAndUnwrapSol,omitempty"`
}

func (cl *JupClient) PostSwapWithResponse(ctx context.Context, reqBody SwapRequest, reqEditors ...RequestEditorFunction) (*PostSwapResponse, error){
  resp, err := cl.PostSwap(ctx, reqBody, reqEditors...)
  if err != nil{
    return nil, err
  }
  
  bodyBytes, err := io.ReadAll(resp.Body)
  defer func(){resp.Body.Close()}()
  if err != nil{
    return nil, err
  }
  
  response := &PostSwapResponse{
    Body:         bodyBytes,
    HTTPResponse: resp,
  }
  
  switch{
  case strings.Contains(resp.Header.Get("Content-Type"), "json") && resp.StatusCode == 200:
    var dest SwapResponse
    if err := json.Unmarshal(bodyBytes, &dest); err != nil{
      return nil, err
    }
    response.JSON200 = &dest
  }
  
  return response, nil
}

func (cl *JupClient) PostSwap(ctx context.Context, reqBody SwapRequest, reqEditors ...RequestEditorFunction) (*http.Response, error){
  req, err := NewPostSwapRequest(ctx, cl.Endpoint, reqBody)
  if err != nil{
    return nil, err
  }
  
  req = req.WithContext(ctx)
  if err = cl.applyEditors(ctx, req, reqEditors); err != nil{
    return nil, err
  }
  
  resp,err := cl.HTTPClient.Do(req)
  return resp, err
}

func NewPostSwapRequest(ctx context.Context, endpointURL string, reqBody SwapRequest) (*http.Request, error){
  reqBytes, err := json.Marshal(reqBody)
  if err != nil{
    return nil, err
  }
  
  Endpoint, err := url.Parse(endpointURL)
  if err != nil{
    return nil, err
  }
   
  operationPath := fmt.Sprintf("/swap")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := Endpoint.Parse(operationPath)
	if err != nil {
		return nil, err
	}
  request, err := http.NewRequestWithContext(ctx, "POST", queryURL.String(), bytes.NewReader(reqBytes))
  
  request.Header.Set("Content-Type", "application/json")
  request.Header.Set("Accept", "application/json")
  
  return request, nil
}

