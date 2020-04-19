package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// EndpointURL represents the public URL
var EndpointURL = os.Getenv(`ENDPOINT_URL`)

// JsSubdirectory contains the analytics.js and gtm.js
var JsSubdirectory = os.Getenv(`JS_SUBDIRECTORY`)

// GtmFilename is the new public name for the gtm.js file
var GtmFilename = os.Getenv(`GTM_FILENAME`)

// GaFilename is the new public name for the analytics.js file
var GaFilename = os.Getenv(`GA_FILENAME`)

// GaDebugFilename is the new public name for the analytics_debug.js file
var GaDebugFilename = os.Getenv(`GADEBUG_FILENAME`)

// GaCollectEndpoint is the endpoint for the /collect requests
var GaCollectEndpoint = os.Getenv(`GA_COLLECT_ENDPOINT`)

// GaCollectEndpointRedirect is the endpoint for the /r/collect requests
var GaCollectEndpointRedirect = os.Getenv(`GA_COLLECT_REDIRECT_ENDPOINT`)

// GaCollectEndpointJ is the endpoint for the /j/collect requests
var GaCollectEndpointJ = os.Getenv(`GA_COLLECT_J_ENDPOINT`)

// AllowedGtmIds contains the whitelisted GTM container ids to proxy.
// Will be processed during main() function.
var AllowedGtmIds = strings.Split(os.Getenv(`GTM_IDS`), `,`)

func isInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func javascriptFilesHandle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case JsSubdirectory + GtmFilename:
		googleTagManagerHandle(w, r)
	case JsSubdirectory + GaFilename:
		googleAnalyticsJsHandle(w, r, false)
	case JsSubdirectory + GaDebugFilename:
		googleAnalyticsJsHandle(w, r, true)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Not found`))
	}
	return
}

func collectHandle(w http.ResponseWriter, r *http.Request) {
	googleAnalyticsCollectHandle(w, r)
}

func setResponseHeaders(w http.ResponseWriter, headers http.Header) {
	// Looping through headers from request
	for headerName, headerValue := range headers {
		// Picking and sending relevant headers to client
		if headerName == `Age` || headerName == `Cache-Control` || headerName == `Content-Type` || headerName == `Date` || headerName == `Expires` || headerName == `Last-Modified` {
			w.Header().Set(headerName, headerValue[0])
		}
	}

	w.Header().Set(`Server`, `GoGtmGaProxy 1.0.0`)
}

func main() {
	if EndpointURL == `` || JsSubdirectory == `` || GtmFilename == `` || GaFilename == `` || GaDebugFilename == `` || GaCollectEndpoint == `` || GaCollectEndpointRedirect == `` || GaCollectEndpointJ == `` || len(AllowedGtmIds) < 1 {
		fmt.Println(`ERROR: Seems the environment variables aren't set. Exiting.`)
		os.Exit(1)
	}

	http.HandleFunc(JsSubdirectory, javascriptFilesHandle)

	http.HandleFunc(GaCollectEndpoint, collectHandle)
	http.HandleFunc(GaCollectEndpointRedirect, collectHandle)
	http.HandleFunc(GaCollectEndpointJ, collectHandle)

	if err := http.ListenAndServe(`:8080`, nil); err != nil {
		panic(err)
	}
}
