package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tdewolff/minify/v2"
)

type gaSourceCodeCache struct {
	lastUpdate int64
	src        []byte
	headers    http.Header
	mux        sync.Mutex
}

var srcGaCache = make(map[string]gaSourceCodeCache)
var gaMapSync sync.Mutex

func generateGACookie() string {
	rand.Seed(time.Now().UnixNano())
	return `GA` + GaCookieVersion + `.2.` + strconv.FormatInt(int64(rand.Intn(888888888)+111111111), 10) + `.` + strconv.FormatInt(time.Now().Unix(), 10)
}

func googleAnalyticsJsHandle(w http.ResponseWriter, r *http.Request, path string) {
	var sourceCodeToReturn []byte
	var statusCodeToReturn int = 200
	var headersToReturn http.Header
	var usedCache bool
	var endpointURI = settingsGGGP.EndpointURI
	var cookieDomain = settingsGGGP.CookieDomain
	var cachePath = path

	if settingsGGGP.EndpointURI == "" {
		endpointURI = r.Host
		cachePath = endpointURI + "/" + path
	}
	if settingsGGGP.CookieDomain == "" {
		cookieDomain = r.Host
	}

	GaCache, CacheExists := srcGaCache[cachePath]

	if CacheExists == false {
		gaMapSync.Lock()
		srcGaCache[cachePath] = gaSourceCodeCache{lastUpdate: 0}
		gaMapSync.Unlock()

		GaCache, _ = srcGaCache[cachePath]
	}

	if settingsGGGP.EnableDebugOutput {
		fmt.Println(`Locking Cache MUX`)
	}

	GaCache.mux.Lock()

	if settingsGGGP.EnableDebugOutput {
		fmt.Println(`Locked Cache MUX`)
	}

	if GaCache.lastUpdate > (time.Now().Unix() - settingsGGGP.GaCacheTime) {
		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Unlocking Cache MUX (Cache)`)
		}
		GaCache.mux.Unlock()
		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Unlocked Cache MUX (Cache)`)
		}

		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Fetching GA-JS Request from Cache...`)
		}

		sourceCodeToReturn = GaCache.src
		headersToReturn = GaCache.headers
		usedCache = true
	} else {
		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Refreshing Cache for GA-JS Request...`)
		}

		client := &http.Client{}

		var req *http.Request
		var err error

		switch path {
		case `default`:
			req, err = http.NewRequest(`GET`, `https://www.google-analytics.com/analytics.js`, nil)

			if settingsGGGP.EnableDebugOutput {
				fmt.Println(`REQUESTING: https://www.google-analytics.com/analytics.js`)
			}
		case `debug`:
			req, err = http.NewRequest(`GET`, `https://www.google-analytics.com/analytics_debug.js`, nil)

			if settingsGGGP.EnableDebugOutput {
				fmt.Println(`REQUESTING: https://www.google-analytics.com/analytics_debug.js`)
			}
		default:
			re := regexp.MustCompile(settingsGGGP.GaPluginsDirectoryname)
			pathTranslated := re.ReplaceAll([]byte(path), []byte(`/plugins/`))

			re = regexp.MustCompile(`(.*)?(\/plugins\/.*\.js)(.*)?`)
			pathRequest := re.FindStringSubmatch(string(pathTranslated[:]))

			req, err = http.NewRequest(`GET`, `https://www.google-analytics.com`+pathRequest[2], nil)

			if settingsGGGP.EnableDebugOutput {
				fmt.Println(`REQUESTING: https://www.google-analytics.com` + pathRequest[2])
			}
		}

		if err != nil {
			fmt.Println(`Experienced problems on requesting analytics.js from google. Aborting.`)

			return
		}

		req.Header.Set(`User-Agent`, `GoGtmGaProxy `+os.Getenv(`APP_VERSION`)+`; github.com/blaumedia/go-gtm-ga-proxy`)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(`Experienced problems on requesting analytics.js from google. Aborting.`)

			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(`Experienced problems on requesting analytics.js from google. Aborting.`)

			return
		}

		re := regexp.MustCompile(`googletagmanager.com`)
		body = re.ReplaceAll([]byte(body), []byte(endpointURI))

		re = regexp.MustCompile(`\/gtm.js`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GtmFilename))

		re = regexp.MustCompile(`www.google-analytics.com`)
		body = re.ReplaceAll([]byte(body), []byte(endpointURI))

		re = regexp.MustCompile(`analytics.js`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GaFilename))

		re = regexp.MustCompile(`u\/analytics_debug.js`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GaDebugFilename))

		re = regexp.MustCompile(`\"/r\/collect`)
		body = re.ReplaceAll([]byte(body), []byte(`"`+settingsGGGP.GaCollectEndpointRedirect))

		re = regexp.MustCompile(`\"/j\/collect`)
		body = re.ReplaceAll([]byte(body), []byte(`"`+settingsGGGP.GaCollectEndpointJ))

		re = regexp.MustCompile(`\"/collect`)
		body = re.ReplaceAll([]byte(body), []byte(`"`+settingsGGGP.GaCollectEndpoint))

		re = regexp.MustCompile(`\/plugins\/`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GaPluginsDirectoryname+`/`))

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
				fmt.Println(`Compressed the Google Analytics JS File and reduced it by ` + compressChange + `%.`)
			}
		}

		if resp.StatusCode == 200 {
			GaCache.headers = resp.Header
			GaCache.src = body
			GaCache.lastUpdate = time.Now().Unix()
		}

		headersToReturn = resp.Header
		statusCodeToReturn = resp.StatusCode
		sourceCodeToReturn = body
		usedCache = false

		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Unlocking Cache MUX (No-Cache)`)
		}

		GaCache.mux.Unlock()

		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Unlocked Cache MUX (No-Cache)`)
		}

		// Reassigning the copy of the struct back to map
		gaMapSync.Lock()
		srcGaCache[cachePath] = GaCache
		gaMapSync.Unlock()
	}

	if settingsGGGP.EnableServerSideGaCookies {
		var newCookieContent string
		var newCookieDecodedContent string

		if cookieContent, errCookie := r.Cookie(settingsGGGP.ServerSideGaCookieName); errCookie == nil {
			if cookieDecodedContent, errCookieDecode := base64.StdEncoding.DecodeString(cookieContent.Value); errCookieDecode == nil {
				// Hardening cookie method; allow only numbers, characters and points
				if match, _ := regexp.MatchString(`[A-Za-z0-9\.]`, cookieContent.Value); match {
					newCookieContent = cookieContent.Value
					newCookieDecodedContent = string(cookieDecodedContent)
				}
			}
		} else {
			if gaCookie, gaErr := r.Cookie(`_ga`); gaErr == nil {
				// Hardening cookie method; allow only numbers, characters and points
				if match, _ := regexp.MatchString(`[A-Za-z0-9\.]`, gaCookie.Value); match {
					newCookieDecodedContent = gaCookie.Value
					newCookieContent = base64.StdEncoding.EncodeToString([]byte(newCookieDecodedContent))
				}
			}
		}

		if newCookieContent == `` && newCookieDecodedContent == `` {
			newCookieDecodedContent = generateGACookie()
			newCookieContent = base64.StdEncoding.EncodeToString([]byte(newCookieDecodedContent))
		}

		if settingsGGGP.CookieSecure {
			if settingsGGGP.EnableDebugOutput {
				fmt.Println(`Set-Cookie: ` + settingsGGGP.ServerSideGaCookieName + `=` + newCookieContent + `; Domain=` + cookieDomain + `; Secure; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			}
			w.Header().Add(`Set-Cookie`, settingsGGGP.ServerSideGaCookieName+`=`+newCookieContent+`; Domain=`+cookieDomain+`; Secure; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			w.Header().Add(`Set-Cookie`, settingsGGGP.ClientSideGaCookieName+`=`+newCookieDecodedContent+`; Domain=`+cookieDomain+`; Secure; SameSite=Lax; Path=/; Max-Age=63072000`)
		} else {
			if settingsGGGP.EnableDebugOutput {
				fmt.Println(`Set-Cookie: ` + settingsGGGP.ServerSideGaCookieName + `=` + newCookieContent + `; Domain=` + cookieDomain + `; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			}
			w.Header().Add(`Set-Cookie`, settingsGGGP.ServerSideGaCookieName+`=`+newCookieContent+`; Domain=`+cookieDomain+`; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			w.Header().Add(`Set-Cookie`, settingsGGGP.ClientSideGaCookieName+`=`+newCookieDecodedContent+`; Domain=`+cookieDomain+`; SameSite=Lax; Path=/; Max-Age=63072000`)
		}
	}

	setResponseHeaders(w, headersToReturn)

	if usedCache {
		w.Header().Add(`X-Cache-Hit`, `true`)
	} else {
		w.Header().Add(`X-Cache-Hit`, `false`)
	}

	for _, f := range settingsGGGP.pluginEngine.dispatcher[`after_ga_js`] {
		f(&w, r, &statusCodeToReturn, &sourceCodeToReturn)
	}

	w.WriteHeader(statusCodeToReturn)

	w.Write([]byte(sourceCodeToReturn))
}

func googleAnalyticsCollectHandle(w http.ResponseWriter, r *http.Request) {
	var redirectURL *url.URL

	if r.Method != `GET` && r.Method != `POST` {
		fmt.Println(`ERROR: Connection to collect endpoint through ` + r.Method + ` method. Aborting.`)
		return
	}

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}
	clientURL := ``

	switch r.URL.Path {
	case settingsGGGP.GaCollectEndpointRedirect:
		clientURL = `https://www.google-analytics.com/r/collect`
	case settingsGGGP.GaCollectEndpointJ:
		clientURL = `https://www.google-analytics.com/j/collect`
	case settingsGGGP.GaCollectEndpoint:
		fallthrough
	default:
		clientURL = `https://www.google-analytics.com/collect`
	}

	var req *http.Request
	var err error
	var bodyPayload = make(map[string]string)

	switch r.Method {
	case `GET`:
		for k, v := range r.URL.Query() {
			bodyPayload[url.QueryEscape(k)] = url.QueryEscape(v[0])
		}
	case `POST`:
		postPayloadRaw, _ := ioutil.ReadAll(r.Body)

		postPayload := strings.Split(string(postPayloadRaw), `&`)

		re := regexp.MustCompile(`(.*)=(.*)`)
		for _, item := range postPayload {
			itemParsed := re.FindStringSubmatch(item)

			switch len(itemParsed) {
			case 2:
				bodyPayload[itemParsed[1]] = ``
			case 3:
				bodyPayload[itemParsed[1]] = itemParsed[2]
			}
		}
	}

	if settingsGGGP.EnableDebugOutput {
		fmt.Println(`Collect-Redirect:`)
		fmt.Println(`Payload:`)
		fmt.Println(bodyPayload)
	}

	var formatPayLoad string

	for k, v := range bodyPayload {
		switch k {
		case `uip`:
		case `ua`:
		default:
			formatPayLoad = formatPayLoad + k + `=` + v + `&`
		}
	}

	if os.Getenv(`PROXY_IP_HEADER`) != `` {
		proxyHeaderIps := strings.Split(r.Header.Get(os.Getenv(`PROXY_IP_HEADER`)), `,`)
		var proxyHeaderIPIndex = 0

		if os.Getenv(`PROXY_IP_HEADER_INDEX`) != `` {
			n, err := strconv.Atoi(os.Getenv(`PROXY_IP_HEADER_INDEX`))
			if err == nil {
				proxyHeaderIPIndex = n
			} else {
				fmt.Println(`ERROR: Couldn't convert PROXY_IP_HEADER_INDEX environment variable to type int. Falling back to 0.`)
			}
		}

		IPToRedirect := ``
		for ipIndex, ipValue := range proxyHeaderIps {
			if ipIndex == proxyHeaderIPIndex {
				IPToRedirect = strings.TrimSpace(ipValue)
			}
		}

		if IPToRedirect == `` {
			fmt.Println(`ERROR: Given PROXY_IP_HEADER_INDEX environment variable wasn't found in header range. Fallback to index 0.`)

			IPToRedirect = strings.TrimSpace(proxyHeaderIps[0])
		}

		formatPayLoad = formatPayLoad + `ua=` + url.QueryEscape(r.Header.Get(`User-Agent`)) + `&uip=` + url.QueryEscape(IPToRedirect)
	} else {
		formatPayLoad = formatPayLoad + `ua=` + url.QueryEscape(r.Header.Get(`User-Agent`)) + `&uip=` + url.QueryEscape(strings.Split(r.RemoteAddr, `:`)[0])
	}

	switch r.Method {
	case `GET`:
		req, err = http.NewRequest(`GET`, clientURL+`?`+formatPayLoad, nil)
		if err != nil {
			fmt.Println(`Experienced problems on redirecting collect to google (GET). Aborting.`)

			return
		}
		if settingsGGGP.EnableDebugOutput {
			fmt.Println(clientURL + `?` + formatPayLoad)
		}
	case `POST`:
		req, err = http.NewRequest(`POST`, clientURL, bytes.NewBuffer([]byte(formatPayLoad)))
		if err != nil {
			fmt.Println(`Experienced problems on redirecting collect to google (POST). Aborting.`)

			return
		}
		if settingsGGGP.EnableDebugOutput {
			fmt.Println(`Format Payload:`)
			fmt.Println(formatPayLoad)
		}
	}

	req.Header.Set(`User-Agent`, `GoGtmGaProxy `+os.Getenv(`APP_VERSION`)+`; github.com/blaumedia/go-gtm-ga-proxy`)
	req.Header.Set(`Accept`, `*/*`)
	req.Header.Set(`Content-Type`, `text/plain;charset=UTF-8`)

	resp, err := client.Do(req)
	if err != nil {
		if resp.StatusCode == http.StatusFound {
			redirectURL, _ = resp.Location()
		} else {
			fmt.Println(`Experienced problems on redirecting collect to google. Aborting.`)

			return
		}
	}

	defer resp.Body.Close()

	if redirectURL != nil {
		fmt.Println(`Detected Redirect from Google Measurement Protocol (probably Google Ads). Redirecting to: ` + redirectURL.String())
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(`Experienced problems on redirecting collect to google. Aborting.`)

			return
		}

		w.WriteHeader(resp.StatusCode)
		setResponseHeaders(w, resp.Header)

		w.Write([]byte(body))
	}
}
