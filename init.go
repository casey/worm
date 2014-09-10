package app

import "appengine"
import "appengine/datastore"
import "fmt"
import "strings"
import "net/http"
import "regexp"

var path_re = regexp.MustCompile(`^(/[a-zA-Z0-9_.-]*)(/[a-zA-Z0-9_.-]{0,64})?$`)

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  status := statusCode(http.StatusInternalServerError)
  body := ""
  headers := make(map[string]string)
  headers["Warranty"] = `THIS IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND EXPRESS OR IMPLIED.`
  headers["Content-Type"] = `text/plain; charset="utf-8"`

  defer func() {
    if e := recover(); e != nil {
      c.Errorf("handler: recovered from panic: %v", e)
    }

    for name, value := range headers {
      w.Header().Set(name, value)
    }

    w.WriteHeader(status.number())

    if status.mustNotIncludeMessageBody(r.Method) {
      fmt.Fprint(w, "\n")
    } else if body == "" {
      fmt.Fprintf(w, "%v %v\n", status.number(), status.text())
    } else {
      fmt.Fprintf(w, "%v", body)
    }
  }()

  ensure := func(condition bool, errorCode int) {
    if !condition {
      status = statusCode(errorCode)
      panic("ensure condition false")
    }
  }

  check := func(e error) {
    if e != nil {
      status = http.StatusInternalServerError
      panic(e)
    }
  }

  put := r.Method == "PUT"
  get := r.Method == "GET"
  ensure(put || get, http.StatusMethodNotAllowed)

  match := path_re.FindStringSubmatch(r.URL.Path)

  ensure(len(match) >= 2, http.StatusForbidden)

  key := strings.TrimPrefix(match[1], "/")
  value := ""

  if put {
    ensure(len(match) == 3, http.StatusForbidden)
    value = strings.TrimPrefix(match[2], "/")
  } 

  check(datastore.RunInTransaction(c, func(c appengine.Context) error {
    pointer, e := getValue(c, key)
    check(e)

    if get {
      ensure(pointer != nil, http.StatusNotFound)
      status = http.StatusOK
    } else if pointer == nil {
      pointer, e = putValue(c, key, value)
      check(e)
      status = http.StatusCreated
    } else {
      if *pointer == value {
        status = http.StatusOK
      } else {
        status = http.StatusForbidden
      }
    }

    value = *pointer
    return nil
  }, nil))

  body = value + "\n"
}
