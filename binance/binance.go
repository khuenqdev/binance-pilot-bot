package binance

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "time"

    "github.com/google/go-querystring/query"
)

type Client struct {
    Client *http.Client

    baseURL *url.URL

    apiKey    string
    apiSecret string

    // Used for debugging
    RequestUrlString  string
    RequestBodyString string

    // Services
    MarketService MarketServiceInterface
}

func NewClient(baseUrl, apiKey, apiSecret string) *Client {
    httpClient := http.DefaultClient
    parsedBaseUrl, _ := url.Parse(baseUrl)

    c := &Client{
        Client:    httpClient,
        baseURL:   parsedBaseUrl,
        apiKey:    apiKey,
        apiSecret: apiSecret,
    }

    c.MarketService = &MarketService{client: c}

    return c
}

func (c *Client) NewRequest(method, urlStr string, body, options interface{}) (*http.Request, error) {
    rel, err := url.Parse(urlStr)
    if err != nil {
        return nil, err
    }

    // Make the full url based on the relative path
    u := c.baseURL.ResolveReference(rel)

    // Add custom options
    if options != nil {
        optionsQuery, err := query.Values(options)
        if err != nil {
            return nil, err
        }

        for k, values := range u.Query() {
            for _, v := range values {
                optionsQuery.Add(k, v)
            }
        }

        optionsQuery.Add("timestamp", fmt.Sprintf("%d", time.Now().Nanosecond()))
        optionsQuery.Add("recvWindow", defaultRecvWindow)

        u.RawQuery = optionsQuery.Encode()
    }

    // A bit of JSON ceremony
    var js []byte = nil

    if body != nil {
        js, err = json.Marshal(body)
        if err != nil {
            return nil, err
        }
    }

    c.RequestUrlString = fmt.Sprintf("[%s] %s", method, u.String())
    c.RequestBodyString = string(js)
    req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(js))

    if err != nil {
        return nil, err
    }

    req.Header.Add("Content-Type", requestContentTypeJson)
    req.Header.Add("Accept", requestContentTypeJson)
    req.Header.Add("User-Agent", UserAgent)
    req.Header.Add("X-MBX-APIKEY", c.apiKey)

    return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
    resp, err := c.Client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    err = CheckResponseError(resp)
    if err != nil {
        fmt.Println("Request Url:", c.RequestUrlString)
        fmt.Println("Request body:", c.RequestBodyString)
        return err
    }

    if http.StatusNoContent == resp.StatusCode {
        return nil
    }

    if v != nil {
        decoder := json.NewDecoder(resp.Body)
        err := decoder.Decode(&v)
        if err != nil {
            return err
        }
    }

    return nil
}

func (c *Client) CreateAndDo(method, path string, data, options, resource interface{}) error {
    req, err := c.NewRequest(method, path, data, options)
    if err != nil {
        return err
    }

    err = c.Do(req, resource)
    if err != nil {
        return err
    }

    return nil
}

// Get performs a GET request for the given path and saves the result in the
// given resource.
func (c *Client) Get(path string, resource, options interface{}) error {
    return c.CreateAndDo("GET", path, nil, options, resource)
}

// Post performs a POST request for the given path and saves the result in the
// given resource.
func (c *Client) Post(path string, data, resource interface{}, options interface{}) error {
    return c.CreateAndDo("POST", path, data, options, resource)
}

// Put performs a PUT request for the given path and saves the result in the
// given resource.
func (c *Client) Put(path string, data, resource interface{}, options interface{}) error {
    return c.CreateAndDo("PUT", path, data, options, resource)
}

// Delete performs a DELETE request for the given path
func (c *Client) Delete(path string, data interface{}, resource interface{}, options interface{}) error {
    return c.CreateAndDo("DELETE", path, data, options, resource)
}

func GenerateSignature(apiSecret string, encodedQueryParams string) string {
    if "" == apiSecret {
        return ""
    }

    mac := hmac.New(sha256.New, []byte(apiSecret))
    mac.Write([]byte(encodedQueryParams))

    return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func CheckResponseError(r *http.Response) error {
    if r.StatusCode >= 200 && r.StatusCode < 300 {
        return nil
    }

    type BinanceError struct {
        Code int    `json:"code,omitempty"`
        Msg  string `json:"msg,omitempty"`
    }

    var binanceError BinanceError

    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return err
    }

    if len(bodyBytes) > 0 {
        err := json.Unmarshal(bodyBytes, &binanceError)

        if nil != err {
            return ResponseDecodingError{
                Body:    bodyBytes,
                Message: err.Error(),
                Status:  r.StatusCode,
            }
        }
    }

    return wrapSpecificError(r, ResponseError{
        Status:    r.StatusCode,
        Message:   binanceError.Msg,
        ErrorCode: binanceError.Code,
    })
}

func wrapSpecificError(r *http.Response, err ResponseError) error {
    if err.Status == 429 {
        f, _ := strconv.ParseFloat(r.Header.Get("retry-after"), 64)
        return RateLimitError{
            ResponseError: err,
            RetryAfter:    int(f),
        }
    }

    if err.Status == 406 {
        err.Message = "Not acceptable"
    }
    return err
}
