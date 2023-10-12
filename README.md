# Netlify CMS Git Gateway

Deploys https://github.com/netlify/git-gateway to Vercel.

## Usage

You can use this package in your own Vercel projects by including it in a Go Serverless function like so:

```go
package handler

import (
	"net/http"

	handler "github.com/tom-sherman/git-gateway-vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handler.Handler(w, r)
}
```

## Development

```
vercel env pull
vercel dev
```
