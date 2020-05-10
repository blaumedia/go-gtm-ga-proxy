function addToStorage() {
  var re = /^https?:\/\/([A-Za-z0-9\-_\.]+)?([A-Za-z0-9\-_]+)\.([A-Za-z]{2,})$/;

  if (re.test(document.getElementById('new-entry').value)) {
    let pushElement = {
      url: document.getElementById('new-entry').value,
    };

    if (
      document.getElementById('advancedSetupSettings').getAttribute('style') ===
      'display:block'
    ) {
      if (document.getElementById('new-entry-gtm-ids').value !== '') {
        let newGtmIds = [];
        let GtmRegex = /([A-Z0-9]+)/g;

        document
          .getElementById('new-entry-gtm-ids')
          .value.match(GtmRegex)
          .forEach(function (el) {
            newGtmIds.push(el.trim());
          });

        pushElement['gtmIds'] = newGtmIds;
      } else {
        document.getElementById('errorMessage').innerHTML =
          'Error: GTM-ID field seems to be empty. Either hide the advanced setup settings or enter atleast one GTM-ID!';
        return false;
      }
    }

    chrome.storage.sync.get('GTM_SYNC_URIS', function (storageItem) {
      let isDuplicate = false;

      storageItem.GTM_SYNC_URIS.forEach(function (element) {
        if (element.url === pushElement['url']) {
          isDuplicate = true;
        }
      });

      if (!isDuplicate) {
        storageItem.GTM_SYNC_URIS.push(pushElement);

        chrome.storage.sync.set(
          {
            GTM_SYNC_URIS: storageItem.GTM_SYNC_URIS,
          },
          function () {
            document.getElementById('new-entry').value = '';
            document.getElementById('new-entry-gtm-ids').value = '';
            document.getElementById('errorMessage').innerHTML = '';
          }
        );
      } else {
        document.getElementById('errorMessage').innerHTML =
          'Error: URL already exists. Please remove and try again!';
      }
    });
  } else {
    if (
      document
        .getElementById('new-entry')
        .value.substring(
          document.getElementById('new-entry').value.length - 1
        ) === '/'
    ) {
      document.getElementById('errorMessage').innerHTML =
        'Error: URL contains ending slash. Please remove and try again!';
    } else if (
      document.getElementById('new-entry').value.substring(0, 7) !==
        'http://' &&
      document.getElementById('new-entry').value.substring(0, 8) !== 'https://'
    ) {
      document.getElementById('errorMessage').innerHTML =
        'Error: URL has to start with http/https. Please add and try again!';
    } else {
      document.getElementById('errorMessage').innerHTML =
        "Error: Seems like it isn't a correct URL. It has to be in the scheme protocol://domain.tld, so f.e.: https://google.com.";
    }
  }
}

document
  .getElementById('openAdvancedSetup')
  .addEventListener('click', function () {
    if (
      document.getElementById('advancedSetupSettings').getAttribute('style') ===
      'display:none'
    ) {
      document
        .getElementById('advancedSetupSettings')
        .setAttribute('style', 'display:block');

      document.getElementById('openAdvancedSetup').innerHTML =
        'Advanced setup <i data-feather="chevrons-up"></i>';
      feather.replace();
    } else {
      document
        .getElementById('advancedSetupSettings')
        .setAttribute('style', 'display:none');

      document.getElementById('openAdvancedSetup').innerHTML =
        'Advanced setup <i data-feather="chevrons-down"></i>';
      feather.replace();
    }
  });

document.getElementById('add-button').addEventListener('click', function () {
  addToStorage();
});

document
  .getElementById('new-entry')
  .addEventListener('keypress', function (ev) {
    if (ev.key === 'Enter') {
      addToStorage();
    }
  });

document
  .getElementById('new-entry-gtm-ids')
  .addEventListener('keypress', function (ev) {
    if (ev.key === 'Enter') {
      addToStorage();
    }
  });

chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
  let activeURL = tabs[0].url.match(
    /^(https?:\/\/[A-Za-z0-9\-_\.]+)\/?.*?$/
  )[1];

  document.getElementById('new-entry').value = activeURL;
});

document
  .getElementById('open-options')
  .addEventListener('click', function (event) {
    chrome.runtime.openOptionsPage();
  });

feather.replace();
