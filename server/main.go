package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
)

// DebugOutput enables debug output to stdout
var DebugOutput bool

// EndpointURI represents the public URL
var EndpointURI = os.Getenv(`ENDPOINT_URI`)

// JsSubdirectory contains the analytics.js and gtm.js
var JsSubdirectory = os.Getenv(`JS_SUBDIRECTORY`)

// GaCacheTime contains the time in seconds to cache the analytics.js and analytics_debug.js file
var GaCacheTime, _ = strconv.ParseInt(os.Getenv(`GA_CACHE_TIME`), 10, 64)

// GtmCacheTime contains the time in seconds to cache the gtm.js file
var GtmCacheTime, _ = strconv.ParseInt(os.Getenv(`GTM_CACHE_TIME`), 10, 64)

// JsEnableMinify = true, to use tdewolff/minify and mishoo/UglifyJS to optimize the js files from google
var JsEnableMinify bool

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

// EnableServerSideGaCookies = true(string) if this proxy shall set the _ga cookie as HttpOnly cookie
var EnableServerSideGaCookies bool

// ServerSideGaCookieName is the name where the cookie contents of _ga are being saved
var ServerSideGaCookieName = os.Getenv(`GA_SERVER_SIDE_COOKIE`)

// CookieDomain is the domain where the cookie is set for
var CookieDomain = os.Getenv(`COOKIE_DOMAIN`)

// CookieSecure = true/false to set cookie for https only connections
var CookieSecure bool

// ClientSideGaCookieName defaults to _ga, optional can be changed through Environment variable 'GA_CLIENT_SIDE_COOKIE', see main() func
var ClientSideGaCookieName = `_ga`

// PluginsEnabled is true, if you want to enable the "plugin engine". Therefore the server searchs for plugins in ./plugins dir.
var PluginsEnabled bool

var pluginEngine = pluginSystem{
	dispatcher: make(map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte)),
}

// GaCookieVersion is the cookie version in the _ga cookie
const GaCookieVersion = "1"

type pluginSystem struct {
	plugins    []*plugin.Plugin
	dispatcher map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte)
}

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

	w.Header().Set(`X-Powered-By`, `GoGtmGaProxy `+os.Getenv(`APP_VERSION`))
}

func main() {
	if EndpointURI == `` || JsSubdirectory == `` || GtmFilename == `` || GaFilename == `` || GaDebugFilename == `` || GaCollectEndpoint == `` || GaCollectEndpointRedirect == `` || GaCollectEndpointJ == `` || len(AllowedGtmIds) < 1 {
		fmt.Println(`ERROR: Seems the environment variables aren't set. Exiting.`)
		os.Exit(1)
	}

	// Replace ClientSideGaCookieName if environment variable is set
	if os.Getenv(`GA_CLIENT_SIDE_COOKIE`) != `` {
		ClientSideGaCookieName = os.Getenv(`GA_CLIENT_SIDE_COOKIE`)
	}

	if strings.ToLower(os.Getenv(`ENABLE_DEBUG_OUTPUT`)) == `true` || strings.ToLower(os.Getenv(`ENABLE_DEBUG_OUTPUT`)) == `1` {
		DebugOutput = true
	}

	if strings.ToLower(os.Getenv(`JS_MINIFY`)) == `true` || strings.ToLower(os.Getenv(`JS_MINIFY`)) == `1` {
		JsEnableMinify = true
	}

	if strings.ToLower(os.Getenv(`ENABLE_SERVER_SIDE_GA_COOKIES`)) == `true` || strings.ToLower(os.Getenv(`ENABLE_SERVER_SIDE_GA_COOKIES`)) == `1` {
		EnableServerSideGaCookies = true
	}

	if strings.ToLower(os.Getenv(`COOKIE_SECURE`)) == `true` || strings.ToLower(os.Getenv(`COOKIE_SECURE`)) == `1` {
		CookieSecure = true
	}

	if strings.ToLower(os.Getenv(`ENABLE_PLUGINS`)) == `true` || strings.ToLower(os.Getenv(`ENABLE_PLUGINS`)) == `1` {
		PluginsEnabled = true
	}

	if PluginsEnabled {
		_, err := os.Stat(`/app/plugins`)

		if os.IsNotExist(err) == false {
			err := filepath.Walk(`./plugins/`, func(path string, info os.FileInfo, _ error) error {
				if info.IsDir() == false && info.Name()[len(info.Name())-2:] == "so" {
					fmt.Println(`/app/plugins/` + info.Name())
					p, err := plugin.Open(`/app/plugins/` + info.Name())

					if err != nil {
						fmt.Println(`ERROR: Failure on opening plugin!`)
						panic(err)
					}

					mainFunc, err := p.Lookup(`Main`)

					if err != nil {
						fmt.Println(`ERROR: Failure on starting main routine of plugin!`)
						panic(err)
					}

					mainFunc.(func())()

					pluginDispatcher, err := p.Lookup(`Dispatcher`)

					if err != nil {
						fmt.Println(`ERROR: Failure on reading Dispatcher of plugin!`)
						panic(err)
					}

					for k, v := range *pluginDispatcher.(*map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte)) {
						for _, f := range v {
							pluginEngine.dispatcher[k] = append(pluginEngine.dispatcher[k], f)
						}
					}
				}
				return nil
			})

			if err != nil {
				fmt.Println(`ERROR: Failure on opening plugin!`)
				panic(err)
			}
		} else {
			fmt.Println(`ERROR: plugins-directory doesn't exist! Did you pick the right docker image that is made for plugin usage or did you forget to disable ENABLE_PLUGINS environment variable while switching images?`)
			os.Exit(1)
		}
	}

	http.HandleFunc(JsSubdirectory, javascriptFilesHandle)

	http.HandleFunc(GaCollectEndpoint, collectHandle)
	http.HandleFunc(GaCollectEndpointRedirect, collectHandle)
	http.HandleFunc(GaCollectEndpointJ, collectHandle)

	if err := http.ListenAndServe(`:8080`, nil); err != nil {
		panic(err)
	}
}
