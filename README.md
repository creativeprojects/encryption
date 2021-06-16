# encryption-handler

Proof of concept of a poor man's encryption over http (without using SSL/TLS).
* The traffic is encrypted using AES-256 (GCM).
* The key is generated with scrypt using a passphrase and a salt.
* Unless specified, each message is encrypted with a new unique nonce.
* Both salt and nonce are sent (in clear) as http header

## Server side

```go
package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/creativeprojects/encryption"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeString(w, "This message is going to be encrypted if your client supports it!")
	})

	port := 3001
	passphrase := []byte("I have been eating too much chocolate")

	// generate a random salt
	salt := make([]byte, 41)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err)
	}

	// create the encryption handler around your usual handler
	handler, err := encryption.NewHandler(passphrase, salt, mux)
	if err != nil {
		panic(err)
	}

	log.Printf("Service is listening on port %d...", port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func writeString(w http.ResponseWriter, payload string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf8")
	_, _ = io.WriteString(w, payload+"\n")
}


```

See the standard response:

```
% curl -v http://localhost:3001
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 3001 (#0)
> GET / HTTP/1.1
> Host: localhost:3001
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf8
< Date: Wed, 16 Jun 2021 16:20:52 GMT
< Content-Length: 66
<
This message is going to be encrypted if your client supports it!
* Connection #0 to host localhost left intact
* Closing connection 0
```

See the encrypted response:

```
% curl -v -H "Accept-Encoding: aesgcm" http://localhost:3001
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 3001 (#0)
> GET / HTTP/1.1
> Host: localhost:3001
> User-Agent: curl/7.64.1
> Accept: */*
> Accept-Encoding: aesgcm
>
< HTTP/1.1 200 OK
< Content-Encoding: aesgcm
< Content-Type: text/plain; charset=utf8
< Vary: Accept-Encoding
< X-Content-Encryption: 2e79168e7d19a7538892583e611020e4fa3c58cb0b45e9de1281d2d12d6cfbbf76b3e098a074c35f25|37a94bfa35ed7e69c23b99f6
< Date: Wed, 16 Jun 2021 16:19:46 GMT
< Content-Length: 82
<
Warning: Binary output can mess up your terminal. Use "--output -" to tell
Warning: curl to output it to your terminal anyway, or consider "--output
Warning: <FILE>" to save to a file.
* Failed writing body (0 != 82)
* Closing connection 0
```

## Client side

```go
package main

import (
	"log"
	"net/http"

	"github.com/creativeprojects/encryption"
)

const serverURL = "http://localhost:3001/"

func main() {
	log.Printf("Calling %s", serverURL)

	req, err := http.NewRequest(http.MethodGet, serverURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// add the accept-encoding header saying we accept our encrypted content
	encryption.AddAcceptEncryptionHeader(req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// decrypt the response
	passphrase := []byte("I have been eating too much chocolate")
	message, err := encryption.DecryptResponse(resp, passphrase)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Answer: %q\n", string(message))
}

```