# SafeHTTP Client
SafeHTTP provides users with a standard http client with sensible defaults that can be used to make http requests to untrusted URLs. This is useful, for instance,
for dispatching outbound webhook requests with URLs that are provided by customers. SafeHTTP Client inspects the final ip address, instead of simply the provided URL,
protecting you against the threat of SSRF even if the attacker uses a DNS record that looks harmless, but that points to localhost or other forbidden addresses.

Much of the code was adapted from [Andrew Ayer's public domain code](https://www.agwa.name/blog/post/preventing_server_side_request_forgery_in_golang).
## Installing
```
go get github.com/AlchemyTelcoSolutions/safehttp
```

## Usage
```go
// create client
client := safehttp.NewClient(c.opts)
// set forbidden endpoint
req, err := http.NewRequest("POST", "http://localhost:5000", nil)
if err != nil {
    // crash and burn
}
// do actual request, should get a safehttp error
res, err := client.Do(req)
if err != nil {
    // safehttp won't let the call to actually happen
}
```
