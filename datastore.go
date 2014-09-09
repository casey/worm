package app

import "crypto/sha256"
import "encoding/hex"
import "appengine"
import "appengine/datastore"

type entity struct {
  Value string `datastore:"value,noindex"`
}

func stringID(key string) string {
  sha := sha256.New()
  sha.Write([]byte(key))
  sum := sha.Sum(nil)
  return hex.EncodeToString(sum)
}

func makeDatastoreKey(c appengine.Context, key string) *datastore.Key {
  return datastore.NewKey(c, "value", stringID(key), 0, nil)
}

func getValue(c appengine.Context, key string) (*string, error) {
  datastoreKey := makeDatastoreKey(c, key)
  entity := entity{}
  e := datastore.Get(c, datastoreKey, &entity)
  c.Errorf("getValue: '%v' -> '%v' %v", key, stringID(key), e)
  if e == datastore.ErrNoSuchEntity {
    return nil, nil
  } else if e != nil {
    return nil, e
  } else {
    return &entity.Value, nil
  }
}

func putValue(c appengine.Context, key string, value string) (*string, error) {
  datastoreKey := makeDatastoreKey(c, key)
  entity := entity{value}
  _, e := datastore.Put(c, datastoreKey, &entity)
  c.Errorf("putValue: '%v' -> '%v' %v", key, stringID(key), e)
  if e == nil {
    return &entity.Value, nil
  } else {
    return nil, e
  }
}
