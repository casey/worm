package app

import "appengine"
import "appengine/datastore"
import "net/http"
import "regexp"

import . "flotilla"

var put_re = regexp.MustCompile(`^/(?P<key>[a-zA-Z0-9_.-]+)/(?P<value>[a-zA-Z0-9_.-]{0,64})$`)
var get_re = regexp.MustCompile(`^/(?P<key>[a-zA-Z0-9_.-]+)$`)

func init() {
  Handle("/").Put(put).Get(get).Options(options)
}

func options(r *http.Request) {
  Status(http.StatusOK)
}

func put(r *http.Request) {
  c := appengine.NewContext(r)
  components := Components(put_re, r.URL.Path)
  Ensure(components != nil, http.StatusForbidden)
  key := components["key"]
  value := components["value"]

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
    Text(http.StatusCreated, value)
  } else if *stored == value {
    Text(http.StatusOK, value)
  } else {
    Status(http.StatusForbidden)
  }
}

func get(r *http.Request) {
  c := appengine.NewContext(r)
  components := Components(get_re, r.URL.Path)
  Ensure(components != nil, http.StatusForbidden)
  key := components["key"]
  value, e := getValue(c, key)
  Check(e)
  Ensure(value != nil, http.StatusNotFound)
  Body(http.StatusOK, *value, "text/plain; charset=utf-8")
}
