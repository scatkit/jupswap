package jupcore
import(
  "net/http"
  "fmt"
  "context"
  "strings"
)

// Abstraction over http client
type HTTPClient interface{
  Do(req *http.Request) (*http.Response, error)
  CloseIdleConnections()
}

type RequestEditorFunction func(ctx context.Context, req *http.Request) error

type jupClient struct{
  Endpoint        string
  HTTPClient      HTTPClient
  RequestEditors  []RequestEditorFunction
}

type ClientOption func(*jupClient) error

func NewClient(endpointURL string, opts ...ClientOption) (*jupClient, error){
  jupCl := &jupClient{
    Endpoint:   endpointURL,
    HTTPClient: &http.Client{},
  }
  
  if !strings.HasSuffix(jupCl.Endpoint, "/"){
    jupCl.Endpoint += "/"
  }

  if opts != nil{
    for _, option := range opts{
      if err := option(jupCl); err != nil{
        return nil, fmt.Errorf("Could not apply option(%v): %w", option, err)
      }
    }
  }
  
  return jupCl, nil
}

func (cl *jupClient) applyEditors(ctx context.Context, req *http.Request, extraEditors []RequestEditorFunction) error{
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

