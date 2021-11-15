package binance

import (
    "fmt"
)

type ResponseError struct {
    Status    int
    Message   string
    ErrorCode int
}

func (e ResponseError) Error() string {
    if e.Message != "" {
        return e.Message
    }

    return getErrorFromStatusCode(e.Status)
}

type ResponseDecodingError struct {
    Body    []byte
    Message string
    Status  int
}

func (e ResponseDecodingError) Error() string {
    return e.Message
}

type RateLimitError struct {
    ResponseError
    RetryAfter int
}

func getErrorFromStatusCode(statusCode int) string {
    errStatusCodeMap := map[int]string{
        403: "The WAF Limit (Web Application Firewall) has been violated",
        404: "The particular resource does not exist or could not be found",
        405: "Request method is not allowed",
        418: "An IP has been auto-banned for continuing to send requests after receiving 429 codes",
        429: "Too many requests, rate limit has been reached",
    }

    if errorMsg, ok := errStatusCodeMap[statusCode]; ok {
        return errorMsg
    }

    return fmt.Sprintf("Unknown Error (Status: %d)", statusCode)
}
