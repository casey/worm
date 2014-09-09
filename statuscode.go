package app

import "net/http"

type statusCode int

func (status statusCode) informational() bool { return status >= 100 && status < 200 }
func (status statusCode) successful() bool    { return status >= 200 && status < 300 }
func (status statusCode) redirection() bool   { return status >= 300 && status < 400 }
func (status statusCode) badRequest() bool    { return status >= 400 && status < 500 }
func (status statusCode) serverError() bool   { return status >= 500 && status < 600 }

func (status statusCode) mustNotIncludeMessageBody(method string) bool {
  return status.informational() ||
    status == http.StatusNoContent ||
    status == http.StatusResetContent ||
    status == http.StatusNotModified ||
    status == http.StatusOK && method == "HEAD"
}

func (status statusCode) text() string {
  if text := http.StatusText(status.number()); text != "" {
    return text
  }

  return "Mystery Status Code"
}

func (status statusCode) number() int {
  return int(status)
}
