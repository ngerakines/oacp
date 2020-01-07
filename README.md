# OAuth Callback Proxy (oacp)

`oacp` is a proxy used to redirect requests by the state query string parameter.

# Operation

First, run the server. By default, local storage (in memory) is used.

    $ oacp server

Next, register state/location pairs using the API.

    $ oacp record --server http://localhost:5000 --state foo --location https://www.google.com/

OR

    $ curl -vv -F "state=foo" -F "location=https://www.google.com/" "http://aocp:aocp@localhost:5000/api/locations"

Last, the callback will redirect as configured.

    $ curl -vv "http://localhost:5000/callback?state=foo&code=what"
    > GET /callback?state=foo&code=what HTTP/1.1
    > Host: localhost:5000
    > User-Agent: curl/7.59.0
    > Accept: */*
    >
    < HTTP/1.1 302 Found
    < Content-Type: text/html; charset=utf-8
    < Location: https://www.google.com/
    < Date: Tue, 07 Jan 2020 17:45:09 GMT
    < Content-Length: 46
    <
    <a href="https://www.google.com/">Found</a>.

**Note**: State tokens are expected to be unique and one-time-use only. Once a state is used, it is removed from storage and cannot be used again. 