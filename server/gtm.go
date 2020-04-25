package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type gtmSourceCodeCache struct {
	lastUpdate int64
	src        []byte
	headers    http.Header
}

var srcGtmCache = gtmSourceCodeCache{lastUpdate: 0}

func googleTagManagerHandle(w http.ResponseWriter, r *http.Request) {
	var GtmContainerID string
	var GtmDatalayerVar string
	var GtmCookies []string

	var sourceCodeToReturn []byte
	var statusCodeToReturn int = 200
	var headersToReturn http.Header
	var usedCache bool

	if innerID, ok := r.URL.Query()[`id`]; ok {
		GtmContainerID = innerID[0]
	} else {
		fmt.Println(`No get-parameter 'id' set.`)

		w.Write([]byte(`No get-parameter 'id' set.`))

		return
	}

	if innerID, ok := r.URL.Query()[`l`]; ok {
		GtmDatalayerVar = `&l=` + innerID[0]
	}

	if !isInSlice(AllowedGtmIds, r.URL.Query()[`id`][0]) {
		fmt.Println(`Tried to open disallowed GTM ID: ` + r.URL.Query()[`id`][0])

		w.Write([]byte(`ID (` + r.URL.Query()[`id`][0] + `) needs to be whitelisted.`))

		return
	}

	// Picking gtm_* Cookies
	for _, cookie := range r.Cookies() {
		name := []rune(cookie.Name)

		if string(name[0:4]) == `gtm_` {
			GtmCookies = append(GtmCookies, cookie.Name+`=`+cookie.Value)
		}
	}

	if len(GtmCookies) > 0 || srcGtmCache.lastUpdate > (time.Now().Unix()-GtmCacheTime) {
		sourceCodeToReturn = srcGtmCache.src
		headersToReturn = srcGtmCache.headers
		usedCache = true
	} else {
		client := &http.Client{}

		req, err := http.NewRequest(`GET`, `https://www.googletagmanager.com/gtm.js?id=GTM-`+GtmContainerID+GtmDatalayerVar, nil)
		if err != nil {
			fmt.Println(`Experienced problems on requesting gtm.js from google. Aborting.`)

			return
		}

		req.Header.Set(`User-Agent`, `GoGtmGaProxy `+os.Getenv(`APP_VERSION`)+`; github.com/blaumedia/go-gtm-ga-proxy`)

		// Redirect gtm_* cookies to GTM for preview mode
		req.Header.Set(`Cookie`, strings.Join(GtmCookies, `; `))

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(`Experienced problems on requesting gtm.js from google. Aborting.`)

			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(`Experienced problems on requesting gtm.js from google. Aborting.`)

			return
		}

		JsSubdirectoryWithoutLeadingSlash := []rune(JsSubdirectory)
		JsSubdirectoryWithoutLeadingSlash = JsSubdirectoryWithoutLeadingSlash[1:]

		re := regexp.MustCompile(`googletagmanager.com`)
		body = re.ReplaceAll([]byte(body), []byte(EndpointURL))

		re = regexp.MustCompile(`\/gtm.js`)
		body = re.ReplaceAll([]byte(body), []byte(string(JsSubdirectoryWithoutLeadingSlash)+GtmFilename))

		re = regexp.MustCompile(`www.google-analytics.com`)
		body = re.ReplaceAll([]byte(body), []byte(EndpointURL))

		re = regexp.MustCompile(`analytics.js`)
		body = re.ReplaceAll([]byte(body), []byte(string(JsSubdirectoryWithoutLeadingSlash)+GaFilename))

		re = regexp.MustCompile(`u\/analytics_debug.js`)
		body = re.ReplaceAll([]byte(body), []byte(string(JsSubdirectoryWithoutLeadingSlash)+GaDebugFilename))

		if resp.StatusCode == 200 && len(GtmCookies) == 0 {
			srcGtmCache.headers = resp.Header
			srcGtmCache.src = body
			srcGtmCache.lastUpdate = time.Now().Unix()
		}

		headersToReturn = resp.Header
		statusCodeToReturn = resp.StatusCode
		sourceCodeToReturn = body
		usedCache = false
	}

	setResponseHeaders(w, headersToReturn)

	if usedCache {
		w.Header().Add(`X-Cache-Hit`, `true`)
	} else {
		w.Header().Add(`X-Cache-Hit`, `false`)
	}

	w.WriteHeader(statusCodeToReturn)

	w.Write([]byte(sourceCodeToReturn))
}
