WORM
====

Write once read many key value storage service

[Test instance here.](http://rodarmor-worm.appspot.com)


API
---

KEYs match `/[a-zA-Z0-9.-_]+/`.

* `PUT /KEY` -> Associates KEY with data in request. Puts after the first will be ignored.
* `GET /KEY` ->  Returns the data associated with KEY.

```
> curl -X PUT http://rodarmor-worm.appspot.com/hello --data 'bob'
bob
> curl -X PUT http://rodarmor-worm.appspot.com/hello --data 'frank'
403 Forbidden
> curl -X GET http://rodarmor-worm.appspot.com/hello
bob
```

About
-----

KEYs can be any length, but data is limited to 64 bytes just to avoid too much spam in the test instance datastore.

The sha256 hash of KEYs are used as datastore string IDs, instead of the KEY itself. This allows keys to be arbitrarily long, since they aren't actually stored in the datastore. Also, this lessens the severity of a [potential attack vector](http://ikaisays.com/2011/01/25/app-engine-datastore-tip-monotonically-increasing-values-are-bad/).
