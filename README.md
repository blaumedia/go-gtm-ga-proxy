# Google Tag Manager & Google Analytics Proxy
![Header](http://i.blaumedia.com/gogtmgaproxy-github-header.png)
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Chrome Extension](#chrome-extension)
- [Performance](#performance)
- [Plugins](#plugins)
- [To do](#to-do)
- [Support](#support)
- [Donate](#donate)


## Introduction
The fight against tracking blockers goes into the next round! Companies depend on having a solid data base to make the right decisions and to understand their customers accurately. This proxy server makes it possible to achieve this with minimal effort.

But how exactly does this "proxy" for the Google Tag Manager and Google Analytics work? And what is a proxy actually?

* IMAGE *

To prevent tracking, blockers try to block the URLs of www.googletagmanager.com and www.google-analytics.com among several other sites. So this proxy server integrates itself into your website and makes https://www.googletagmanager.com/gtm.js?id=GTM-ABC123 (f.e.) to https://www.yourdomain.tld/K3ZQ8HD.js?id=ABC123. Thus a Tracking Blocker does not realize that the Google Tag Manager is behind the script. The same happens with the script of Google Analytics (https://www.google-analytics.com/analytics.js). The proxy will change the URLs in the scripts to route the traffic through the proxy instead of requesting Google servers. Additionally the proxy of the Google Analytics Javascript has a *built-in solution* for the client-side cookie problematic that came up with Apples Intelligent Tracking Prevention (ITP).

## Features
- Bypass tracking blockers to load your Google Tag Manager
- Bypass tracking blockers to load the Google Analytics script and send pageviews/events/... to Analytics
- Adjustable cache for the above scripts to increase performance (see [Performance](#performance) section to see the benefits and comparison to [ZitRos proxy solution](https://github.com/ZitRos/save-analytics-from-content-blockers))
- Use of server-side cookie to keep the google client id in users browsers to up to 2 years (ITP 2.3 bypass)
- Chrome Extension for the trouble-free use of the Google Tag Manager preview mode (see [Chrome Extension](#chrome-extension) for the why)
- High performant and full RFC compliant thanks to Go and [net/http](https://golang.org/pkg/net/http/)
- JS minifier for optimizations of up to 30% (thanks to [tdewolff/minify](https://github.com/tdewolff/minify) and [mishoo/UglifyJS](https://github.com/mishoo/UglifyJS))
- Small 50MB (nice ;)) [ready-2-use](https://hub.docker.com/repository/docker/blaumedia/go-gtm-ga-proxy) docker image for crazy fast deployments, updates and close-to-zero harddrive occupation
- Many environment variables to adjust the proxy exactly to your needs
- Expandable through [Plugins](#plugins)

## Installation
First of all, you'll need an installed docker daemon on your server. See the [docker docs](https://docs.docker.com/install/) for more information.

There are many possible solutions to integrate the proxy into your system. Please read the documentation of *each environment variable* carefully.
Starting the proxy is possible with following command and environment variables:
```shell
docker run \
    -e ENDPOINT_URI=www.yourdomain.tld \
    -e JS_SUBDIRECTORY=/js/ \
    -e GA_CACHE_TIME=3600 \
    -e GTM_CACHE_TIME=3600 \
    -e GTM_FILENAME=gtm.js \
    -e GTM_A_FILENAME=gtm_a.js \
    -e GA_FILENAME=ga.js \
    -e GADEBUG_FILENAME=ga_debug.js \
    -e GA_PLUGINS_DIRECTORYNAME=/links/ \
    -e GTAG_FILENAME=tag.js \
    -e RESTRICT_GTM_IDS=false \
    -e GTM_IDS=YOUR_GTM_ID_WITHOUT_GTM- \
    -e GA_COLLECT_ENDPOINT=/fetch \
    -e GA_COLLECT_REDIRECT_ENDPOINT=/fetch_r \
    -e GA_COLLECT_J_ENDPOINT=/fetch_j \
    -e PROXY_IP_HEADER=X-Forwarded-For \
    -e PROXY_IP_HEADER_INDEX=0 \
    -e ENABLE_SERVER_SIDE_GA_COOKIES=true \
    -e GA_SERVER_SIDE_COOKIE=_gggp \
    -e GA_CLIENT_SIDE_COOKIE=_ga \
    -e COOKIE_DOMAIN=yourdomain.tld \
    -e COOKIE_SECURE=true \
    -p "8080:8080" \
    blaumedia/go-gtm-ga-proxy:1.0.0
```


|Variable|Description|Example|
|-----------------|-----|-----|
|```ENABLE_DEBUG_OUTPUT```|Set this to true to let the proxy print debug details for setting up or debugging.|false|
|```ENDPOINT_URI```|URI where the proxy will be reachable.|www.google.com|
|```JS_SUBDIRECTORY```|It is intended to serve the javascript files within a subdirectory of the URI. Here you can define a name for it.|/js/|
|```GA_CACHE_TIME```|Time in seconds the proxy caches the Google Analytics client javascript.|3600|
|```GTM_CACHE_TIME```|Time in seconds the proxy caches the Google Tag Manager client javascript.|3600|
|```GTM_FILENAME```|The filename the GTM javascript file is reachable.|gtm_inject.js|
|```GTM_A_FILENAME```|The alternative filename the GTM javascript file is reachable. (www.googletagmanager.com/a)|container_a.js|
|```GA_FILENAME```|The filename the Google Analytics javascript file is reachable.|ga_inject.js|
|```GADEBUG_FILENAME```|The filename the [Google Analytics debug](https://developers.google.com/analytics/devguides/collection/analyticsjs/debugging) javascript file is reachable.|gadebug_inject.js|
|```GA_PLUGINS_DIRECTORYNAME```|The directory name where google analytics plugins will be accessable at.|/links/|
|```GTAG_FILENAME```|The directory name where google analytics plugins will be accessable at.|tag.js|
|```RESTRICT_GTM_IDS```|Set to true if you want to enable a whitelist for the GTM-IDs.|false|
|```GTM_IDS```|Here you can setup the whitelist for GTM ids. Comma-separate if you want to add multiple ids. Just put the ids without the leading 'GTM-'.|NNQJ5LT,N5ZZT3|
|```GA_COLLECT_ENDPOINT```|Set the new name for the /collect endpoint.|/fetch|
|```GA_COLLECT_REDIRECT_ENDPOINT```|Set the new name for the /r/collect endpoint.|/fetch_r|
|```GA_COLLECT_J_ENDPOINT```|Set the new name for the /j/collect endpoint. Couldn't find out yet for what this endpoint is. If you know more, please share! :)|/fetch_j|
|```PROXY_IP_HEADER```|The header variable where the proxy will find the IP address of the user.|X-Forwarded-For|
|```PROXY_IP_HEADER_INDEX```|The header variable value will be split by comma. With this variable you can set the index of the users IP.|0|
|```ENABLE_SERVER_SIDE_GA_COOKIES```|If the proxy should transfer the client id to a serverside cookie, set this to true.|true|
|```GA_SERVER_SIDE_COOKIE```|Set the cookie name for the serverside cookie.|_gggp|
|```GA_CLIENT_SIDE_COOKIE```|Set the cookie name for the clientside cookie, where google analytics will take the client id from.|_ga|
|```COOKIE_DOMAIN```|The [domain](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie#Attributes) where the cookies are setup to.|google.com|
|```COOKIE_SECURE```|If your website is accessable through https://, you should set this variable to true. [More info](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie#Secure)|true|


There are two ways to integrate the proxy into your existing system. You have the choice between integration using a subdomain or a subdirectory on your website. When integrating through a subdomain, make sure that you place the cookies on the root domain and not under the subdomain; this can be configured with the environment variable ```COOKIE_DOMAIN```. Another important setting that must be configured in the frontend proxy is the deactivation of rewrite of forwardings. The Google Measurement Protocol will return a 302 redirect on /r/collect requests, leading to a Google Ads endpoint. In order for the end user to be redirected there, the redirect rewrite must be disabled. Example configs for nginx and apache2 are below. If you successfully integrated the proxy into another solution, feel free to fork, edit the docs & submit a pull request! :)

### Example integration using subdirectory
Please note that 127.0.0.1 is an example IP in the following configs. Obviously you have to change it to the ip address where the proxy is reachable.

#### nginx
```nginx
[..]

location /goproxy/ {
    proxy_set_header  X-Forwarded-For $remote_addr;
    proxy_set_header  X-Forwarded-Host $remote_addr;

    proxy_pass http://127.0.0.1:8080/;
    proxy_redirect off;
}

[..]
```

#### apache2
Put the following into the httpd.conf and not into the virtualhosts-files.
```apache2
[..]

ProxyPass "/goproxy" "http://127.0.0.1:8080"

[..]
```

### Example integration using subdomain
Please note that 127.0.0.1 is an example IP in the following configs. Obviously you have to change it to the ip address where the proxy is reachable.

#### nginx
This vhost configuration for nginx does require the following:
```bash
apt-get update && apt-get install openssl certbot

openssl dhparam -out /etc/nginx/dhparam.pem 4096
certbot certonly --standalone -d gggp.example.com
```

```nginx
server {
        listen 443 ssl http2;
        listen [::]:443 ssl http2;

        ssl_certificate /etc/letsencrypt/live/gggp.example.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/gggp.example.com/privkey.pem;

        ssl_protocols TLSv1.2;
        ssl_dhparam /etc/nginx/dhparam.pem;
        ssl_ciphers EECDH+AESGCM:EDH+AESGCM;
        ssl_prefer_server_ciphers on;

        ssl_ecdh_curve secp384r1;

        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload";

        ssl_stapling on;
        ssl_trusted_certificate /etc/letsencrypt/live/gggp.example.com/chain.pem;
        ssl_stapling_verify on;

        resolver 8.8.8.8;

        server_name gggp.example.com;

        location / {
                proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header  X-Forwarded-Host $proxy_add_x_forwarded_for;

                proxy_pass http://127.0.0.1:8080/;
        }
}

server {
        listen 80;
        listen [::]:80;

        server_name gggp.example.com;

        rewrite     ^   https://$server_name$request_uri? permanent;
}
```

# go-gtm-ga-proxy
### Attention: ALPHA! Under heavy development!
Bypass any tracking-blocking browser plugins with this first-party-tracking-proxy for Google Tag Manager and Google Analytics.

## Expected behaviour
- ~Proxying gtm.js and analytics.js to original google servers~
- - ~replace any google-tagmanager.com or google-analytics.com URLs~
- - ~send redirection to user if google analytics /collect endpoint answers with 302 redirection for google ads~
- ~set \_ga cookie so it's server side and ITP safe~
- ~deploy-ready docker container~
- ~chrome plugin for sync of gtm cookies, so preview mode works?~
- maybe (!) geoip2 integration to send more detailed and more accurated geo information to google (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#geoid)
- uBlock dataLayer Blocker https://i.blaumedia.com/h6hyn_Fernstudium_f%C3%BCr_Bachelor_und_Master___IUBH_F.png
- ~Fix Chrome Plugin~
- Add Debug-Output
- ~Implement JS-minifier~

## Client Id Cookie Integration
### Direct Google Analytics Integration
```javascript
ga('create', 'UA-XXXXX-Y', {
    'cookieUpdate': false
});
```

### GTM Integration
Create a centralized "GA Settings" variable that is used by all your tags that get in contact with your google analytics instance. Then set another field according to:
![Google Tag Manager Integration](https://i.blaumedia.com/6800y_Tag_Manager_360_-_Google_Chrome_2020-04-20_0.png)

This will prevent that analytics.js will touch the _ga cookie in any form. The proxy manages the _ga cookie.