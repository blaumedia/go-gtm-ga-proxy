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

function removeFromStorage(event) {
  let parentEl = event.target.parentElement.getAttribute('data-id')
    ? event.target.parentElement.getAttribute('data-id')
    : event.target.parentElement.parentElement.getAttribute('data-id');

  chrome.storage.sync.get('GTM_SYNC_URIS', function (storageItem) {
    let newStorageItem = [];
    storageItem.GTM_SYNC_URIS.forEach(function (value, index) {
      if (index != parentEl) {
        newStorageItem.push(value);
      }
    });

    chrome.storage.sync.set({
      GTM_SYNC_URIS: newStorageItem,
    });
  });
}

function renderItems() {
  chrome.storage.sync.get('GTM_SYNC_URIS', function (storageItem) {
    let allItems = '';
    storageItem.GTM_SYNC_URIS.forEach(function (el, index) {
      if (el.gtmIds) {
        allItems +=
          '<li data-id="' +
          index +
          '">' +
          el.url +
          ' (' +
          el.gtmIds.join(', ') +
          ') <i data-feather="x"></i></li>';
      } else {
        allItems +=
          '<li data-id="' +
          index +
          '">' +
          el.url +
          ' <i data-feather="x"></i></li>';
      }
    });

    document
      .getElementsByClassName('lds-ellipsis')[0]
      .setAttribute('style', 'display:none');

    document.getElementById('itemList').innerHTML = allItems;

    feather.replace();

    let allDeleteIcons = document.getElementsByClassName('feather');
    for (var i = 0; i < allDeleteIcons.length; i++) {
      allDeleteIcons[i].addEventListener('click', removeFromStorage, false);
    }
  });
}

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

document.getElementById('addButton').addEventListener('click', function () {
  addToStorage();
});

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

chrome.storage.onChanged.addListener(function (changes) {
  renderItems();
});

renderItems();
