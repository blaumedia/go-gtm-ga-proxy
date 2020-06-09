chrome.runtime.onInstalled.addListener(function (details) {
  chrome.storage.sync.get('GTM_SYNC_URIS', function (items) {
    if (!('GTM_SYNC_URIS' in items)) {
      chrome.storage.sync.set({
        GTM_SYNC_URIS: [],
      });
    }
  });
});

chrome.cookies.onChanged.addListener(function (changeInfo) {
  if (
    changeInfo.cookie.domain === 'www.googletagmanager.com' &&
    changeInfo.cookie.name.substring(0, 4) === 'gtm_'
  ) {
    if (changeInfo.cause === 'overwrite' || changeInfo.cookie.value === '') {
      chrome.storage.sync.get('GTM_SYNC_URIS', function (items) {
        items.GTM_SYNC_URIS.forEach((URLElement) => {
          chrome.cookies.remove(
            {
              url: URLElement.url,
              name: changeInfo.cookie.name,
            },
            function (cookieRemoveReturn) {
              if (cookieRemoveReturn === null) {
                alert(
                  "Error in syncing cookies from GTM!\nThere was an error while removing the cookie '" +
                    changeInfo.cookie.name +
                    "' for domain '" +
                    domain[1] +
                    "'.\n\nChrome: " +
                    chrome.runtime.lastError
                );
              }
            }
          );
        });
      });
    } else {
      chrome.storage.sync.get('GTM_SYNC_URIS', function (items) {
        var re = /^(?:(https?):\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)/;

        items.GTM_SYNC_URIS.forEach((URLElement) => {
          let domain = re.exec(URLElement.url);
          if (domain[1] && domain[2]) {
            cookieValue = [];

            if (URLElement.gtmIds) {
              changeInfo.cookie.value
                .split(':')
                .forEach(function (cookieElement) {
                  if (
                    URLElement.gtmIds.includes(
                      cookieElement.match(/GTM\-(.*)=.*/)[1]
                    )
                  ) {
                    cookieValue.push(cookieElement);
                  }
                });
            } else {
              cookieValue.push(changeInfo.cookie.value);
            }

            chrome.cookies.set(
              {
                url: URLElement.url,
                name: changeInfo.cookie.name,
                value: cookieValue.join(':'),
                domain: domain[2],
                path: '/',
                secure: domain[1] === 'https' ? true : false,
                httpOnly: false,
              },
              function (cookieSetReturn) {
                if (cookieSetReturn === null) {
                  alert(
                    "Error in syncing cookies from GTM!\nThere was an error while adding the cookie '" +
                      changeInfo.cookie.name +
                      "' for domain '" +
                      domain[1] +
                      "'.\n\nChrome: " +
                      chrome.runtime.lastError
                  );
                }
              }
            );
          } else {
            alert(
              "Error in syncing cookies from GTM!\nCouldn't extract protocol or domain from given URL: " +
                URLElement.url +
                '\nPlease make sure that the URL is in the format: protocol://domain.tld/\nf.e.: https://blaumedia.com/'
            );
          }
        });
      });
    }
  }
});
