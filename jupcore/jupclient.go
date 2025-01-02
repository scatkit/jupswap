package jupcore
import(
  "net/http"
  "context"
  "strings"
  //"github.com/davecgh/go-spew/spew"
)

type HttpRequestDoer interface{
  Do(req *http.Request) (*http.Response, error)
}

type RequestEditorFunction func(ctx context.Context, req *http.Request) error

type Client struct{
  EndpointURL string
  
  //Doer for performig http requests
  HTTPClient HttpRequestDoer
  
  // A list of callbacks for modifying requests which are generated before sending over the network.
  RequestEditors []RequestEditorFunction
}

type ClientOption func(*Client) error

func NewClient(serverUrl string, opts ...ClientOption) (*Client, error){
  cl := Client{
    EndpointURL: serverUrl, 
  }
  
  // mutate client and add optional parameters
  for _, opt := range opts{
    if err := opt(&cl); err != nil{
      return nil, err
    }
  }
  
  if !strings.HasSuffix(cl.EndpointURL, "/"){
    cl.EndpointURL += "/"
  }
    
  if cl.HTTPClient == nil{
    cl.HTTPClient = &http.Client{}
  } 
  
  return &cl, nil
}

func (cl *Client) applyEditors(ctx context.Context, req *http.Request, extraEditors []RequestEditorFunction) error{
  for _, editFunc := range cl.RequestEditors{ // Client's editors
    if err := editFunc(ctx, req); err != nil{
      return err
    }
  }
  for _, editFunc := range extraEditors{ // Additonal editors
    if err := editFunc(ctx, req); err != nil{
      return err
    }
  }
  return nil
}

