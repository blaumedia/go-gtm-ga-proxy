package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/tdewolff/minify/v2"
)

type gtmSourceCodeCache struct {
	lastUpdate int64
	src        []byte
	headers    http.Header
	mux        sync.Mutex
}

var srcGtmCache = make(map[string]gtmSourceCodeCache)
var gtmMapSync sync.Mutex

func googleTagManagerHandle(w http.ResponseWriter, r *http.Request, path string) {
	var GtmContainerID string
	var GtmURLAddition string
	var GtmCookies []string
	var GtmCache gtmSourceCodeCache

	var sourceCodeToReturn []byte
	var statusCodeToReturn int = 200
	var headersToReturn http.Header
	var usedCache bool
	var endpointURI = settingsGGGP.EndpointURI

	if settingsGGGP.EndpointURI == "" {
		endpointURI = r.Host
	}

	if innerID, ok := r.URL.Query()[`id`]; ok {
		if len(innerID[0]) >= 4 {
			if innerID[0][:4] == `GTM-` {
				GtmContainerID = innerID[0][4:]
			} else {
				GtmContainerID = innerID[0]
			}
		} else {
			fmt.Println(`No correct get-parameter 'id' set.`)

			w.Write([]byte(`No correct get-parameter 'id' set.`))

			return
		}
	} else {
		fmt.Println(`No get-parameter 'id' set.`)

		w.Write([]byte(`No get-parameter 'id' set.`))

		return
	}

	for URLKey, URLValue := range r.URL.Query() {
		if URLKey != `id` {
			GtmURLAddition = GtmURLAddition + `&` + URLKey + `=` + URLValue[0]
		}
	}

	GtmCache, CacheExists := srcGtmCache[endpointURI+"/"+GtmContainerID+GtmURLAddition]

	if CacheExists == false {
		gtmMapSync.Lock()
		srcGtmCache[endpointURI+"/"+GtmContainerID+GtmURLAddition] = gtmSourceCodeCache{lastUpdate: 0}
		gtmMapSync.Unlock()

		GtmCache, _ = srcGtmCache[endpointURI+"/"+GtmContainerID+GtmURLAddition]
	}

	if !isInSlice(settingsGGGP.AllowedGtmIds, r.URL.Query()[`id`][0]) && !isInSlice(settingsGGGP.AllowedGtmIds, r.URL.Query()[`id`][0][4:]) && settingsGGGP.RestrictGtmIds {
		fmt.Println(`Tried to open disallowed GTM ID: ` + r.URL.Query()[`id`][0])

		w.Write([]byte(`ID (` + r.URL.Query()[`id`][0] + `) needs to be whitelisted.`))
		setResponseHeaders(w, headersToReturn)
		w.WriteHeader(404)

		return
	}

	// Picking gtm_* Cookies
	for _, cookie := range r.Cookies() {
		name := []rune(cookie.Name)

		if string(name[0:4]) == `gtm_` {
			GtmCookies = append(GtmCookies, cookie.Name+`=`+cookie.Value)
		}
	}

	if settingsGGGP.EnableDebugOutput {
		fmt.Println(`Locking Cache MUX`)
	}

	GtmCache.mux.Lock()

	if settingsGGGP.EnableDebugOutput {
		fmt.Println(`Locked Cache MUX`)
	}

	if len(GtmCookies) == 0 && GtmCache.lastUpdate > (time.Now().Unix()-settingsGGGP.GtmCacheTime) {
		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Unlocking Cache MUX (Cache)`)
		}
		GtmCache.mux.Unlock()
		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Unlocked Cache MUX (Cache)`)
		}

		sourceCodeToReturn = GtmCache.src
		headersToReturn = GtmCache.headers
		usedCache = true
	} else {
		client := &http.Client{}
		var req *http.Request
		var err error

		switch path {
		case `default`:
			if settingsGGGP.EnableDebugOutput {
				fmt.Println(`Requesting: https://www.googletagmanager.com/gtm.js?id=GTM-` + GtmContainerID + GtmURLAddition)
			}

			req, err = http.NewRequest(`GET`, `https://www.googletagmanager.com/gtm.js?id=GTM-`+GtmContainerID+GtmURLAddition, nil)
		case `default_a`:
			if settingsGGGP.EnableDebugOutput {
				fmt.Println(`Requesting: https://www.googletagmanager.com/a?id=GTM-` + GtmContainerID + GtmURLAddition)
			}

			req, err = http.NewRequest(`GET`, `https://www.googletagmanager.com/a?id=GTM-`+GtmContainerID+GtmURLAddition, nil)
		case `gtag`:
			if settingsGGGP.EnableDebugOutput {
				fmt.Println(`Requesting: https://www.googletagmanager.com/gtag/js?id=` + GtmContainerID + GtmURLAddition)
			}

			req, err = http.NewRequest(`GET`, `https://www.googletagmanager.com/gtag/js?id=`+GtmContainerID+GtmURLAddition, nil)
		}

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

		// It seems like the /a endpoint returns a pixel instead of js code. Temporarily disable proxying for it.
		// re := regexp.MustCompile(`www.googletagmanager.com\/a`)
		// body = re.ReplaceAll([]byte(body), []byte(settingsGGGP.EndpointURI+`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GtmAFilename))

		re := regexp.MustCompile(`(www\.)?googletagmanager.com`)
		body = re.ReplaceAll([]byte(body), []byte(endpointURI))

		re = regexp.MustCompile(endpointURI + `\/a`)
		body = re.ReplaceAll([]byte(body), []byte(`www.googletagmanager.com\/a`))

		re = regexp.MustCompile(`\/gtm.js`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GtmFilename))

		re = regexp.MustCompile(`www.google-analytics.com`)
		body = re.ReplaceAll([]byte(body), []byte(endpointURI))

		re = regexp.MustCompile(`(\/)?analytics.js`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GaFilename))

		re = regexp.MustCompile(`u\/analytics_debug.js`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GaDebugFilename))

		re = regexp.MustCompile(`\/gtag\/js`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GtagFilename))

		if settingsGGGP.JsEnableMinify {
			m := minify.New()
			m.AddCmd(`application/javascript`, exec.Command("uglifyjs"))

			var previousLengthOfJs int
			if settingsGGGP.EnableDebugOutput {
				previousLengthOfJs = len(body)
			}

			body, err = m.Bytes(`application/javascript`, body)
			if err != nil {
				panic(err)
			}

			if settingsGGGP.EnableDebugOutput {
				afterLengthOfJs := len(body)
				compressChange := fmt.Sprintf(`%f`, (float64(previousLengthOfJs-afterLengthOfJs)/float64(previousLengthOfJs))*float64(100))
				fmt.Println(`Compressed the Google Tag Manager JS File of ID ` + GtmContainerID + ` and reduced it by ` + compressChange + `%.`)
			}
		}

		if resp.StatusCode == 200 && len(GtmCookies) == 0 {
			GtmCache.headers = resp.Header
			GtmCache.src = body
			GtmCache.lastUpdate = time.Now().Unix()
		}

		headersToReturn = resp.Header
		statusCodeToReturn = resp.StatusCode
		sourceCodeToReturn = body
		usedCache = false

		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Unlocking Cache MUX (Cache)`)
		}
		GtmCache.mux.Unlock()
		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Unlocked Cache MUX (Cache)`)
		}

		// Reassigning the copy of the struct back to map
		gtmMapSync.Lock()
		srcGtmCache[endpointURI+"/"+GtmContainerID+GtmURLAddition] = GtmCache
		gtmMapSync.Unlock()
	}

	setResponseHeaders(w, headersToReturn)

	if usedCache {
		w.Header().Add(`X-Cache-Hit`, `true`)
	} else {
		w.Header().Add(`X-Cache-Hit`, `false`)
	}

	for _, f := range settingsGGGP.pluginEngine.dispatcher[`after_gtm_js`] {
		f(&w, r, &statusCodeToReturn, &sourceCodeToReturn)
	}

	w.WriteHeader(statusCodeToReturn)

	w.Write([]byte(sourceCodeToReturn))
}
