chrome.cookies.onChanged.addListener(function(changeInfo) {
    console.log(changeInfo)
    if (changeInfo.cookie.domain === 'www.googletagmanager.com' && changeInfo.cookie.name.substring(0, 4) === 'gtm_') {
        if (changeInfo.cause === 'overwrite' || changeInfo.cookie.value === '') {
            chrome.storage.sync.get('GTM_SYNC_URIS', function(items) {
                items.GTM_SYNC_URIS.forEach(el => {
                    chrome.cookies.remove(
                        {
                            url: el.url,
                            name: changeInfo.cookie.name
                        },
                        function(cookieRemoveReturn) {
                            if (cookieRemoveReturn === null) {
                                alert("Error in syncing cookies from GTM!\nThere was an error while removing the cookie '" + changeInfo.cookie.name + "' for domain '" + domain[1] + "'.\n\nChrome: " + chrome.runtime.lastError)
                            }
                        }
                    );
                });
            });
        } else {
            chrome.storage.sync.get('GTM_SYNC_URIS', function(items) {
                var re = /^(?:(https?):\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)/;

                items.GTM_SYNC_URIS.forEach(el => {
                    let domain = re.exec(el.url);
                    if (domain[1] && domain[2]) {
                        chrome.cookies.set(
                            {
                                url: el.url,
                                name: changeInfo.cookie.name,
                                value: changeInfo.cookie.value,
                                domain: domain[2],
                                path: '/',
                                secure: domain[1] === 'https' ? true : false,
                                httpOnly: false,
                            },
                            function(cookieSetReturn) {
                                if (cookieSetReturn === null) {
                                    alert("Error in syncing cookies from GTM!\nThere was an error while adding the cookie '" + changeInfo.cookie.name + "' for domain '" + domain[1] + "'.\n\nChrome: " + chrome.runtime.lastError)
                                }
                            }
                        );
                    } else {
                        alert("Error in syncing cookies from GTM!\nCouldn't extract protocol or domain from given URL: " + el.url + "\nPlease make sure that the URL is in the format: protocol://domain.tld/\nf.e.: https://blaumedia.com/")
                    }
                });
            });
        }
    }
})