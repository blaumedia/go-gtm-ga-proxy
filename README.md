# go-gtm-ga-proxy
Bypass any tracking-blocking browser plugins with this first-party-tracking-proxy for Google Tag Manager and Google Analytics.

## Expected behaviour
- Proxying gtm.js and analytics.js to original google servers
- - replace any google-tagmanager.com or google-analytics.com URLs
- - send redirection to user if google analytics /collect emdpoint answers with 302 redirection for google ads
- set \_ga cookie so it's server side and ITP safe
- deploy-ready docker container
- chrome plugin for sync of gtm cookies, so preview mode works?
- maybe (!) geoip2 integration to send more detailed and more accurated geo information to google (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#geoid)
