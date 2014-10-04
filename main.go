package app

import "appengine"
import "appengine/datastore"
import "net/http"
import "regexp"
import "strings"

import . "flotilla"

var key_re = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

func init() {
  Handle("/").Put(put).Get(get).Options(options)
}

func options(r *http.Request) {
  Status(StatusOK)
}

func put(r *http.Request) {
  c := appengine.NewContext(r)
  key := strings.TrimPrefix(r.URL.Path, "/")
  Ensure(key_re.MatchString(key), StatusForbidden)
  Ensure(r.ContentLength >= 0, StatusLengthRequired)
  Ensure(r.ContentLength <= 128, StatusRequestEntityTooLarge)
  buffer, e := ReadContent(r)
  Check(e)
  value := string(buffer)
  var stored *string
  Check(datastore.RunInTransaction(c, func(c appengine.Context) error {
    v, e := getValue(c, key)
    stored = v
    Check(e)
    if stored == nil {
      Check(putValue(c, key, value))
    }
    return nil
  }, nil))

  if stored == nil {
    Text(StatusCreated, value)
  } else if *stored == value {
    Text(StatusOK, value)
  } else {
    Status(StatusForbidden)
  }
}

func get(r *http.Request) {
  c := appengine.NewContext(r)
  key := strings.TrimPrefix(r.URL.Path, "/")
  Ensure(key_re.MatchString(key), StatusForbidden)
  value, e := getValue(c, key)
  Check(e)
  Ensure(value != nil, StatusNotFound)
  Body(StatusOK, *value, "text/plain; charset=utf-8")
}
