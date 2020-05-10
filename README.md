# go-gtm-ga-proxy
### Attention: ALPHA! Under heavy development!
Bypass any tracking-blocking browser plugins with this first-party-tracking-proxy for Google Tag Manager and Google Analytics.

## Expected behaviour
- ~Proxying gtm.js and analytics.js to original google servers~
- - ~replace any google-tagmanager.com or google-analytics.com URLs~
- - ~send redirection to user if google analytics /collect emdpoint answers with 302 redirection for google ads~
- ~set \_ga cookie so it's server side and ITP safe~
- ~deploy-ready docker container~
- ~chrome plugin for sync of gtm cookies, so preview mode works?~
- maybe (!) geoip2 integration to send more detailed and more accurated geo information to google (https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#geoid)

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