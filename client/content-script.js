var API_ENDPOINT = "https://heartbleed.look.io/api/v1/url";
var haveShownAlert = false;

function makeAPIRequest() {

  var host = window.location.host;
  var payload = "url=" + host;

  var xhr = new XMLHttpRequest();
  xhr.open("POST", API_ENDPOINT, true);
  xhr.setRequestHeader("Access-Control-Allow-Origin", "*");
  xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

  xhr.onload = function(e) {
    if ( this.response === "1" ) {
      showAlert(host);
    }
  };

  xhr.send(payload);

}

function showAlert(host) {
  haveShownAlert = true;
  chrome.runtime.sendMessage({
    showAlert : true,
    hostName : host || window.location.host
  });
}

chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
  if ( request && request.triggerAlert && haveShownAlert ) {
    showAlert();
  }
});

window.addEventListener('load', function() {
  if ( window.location.protocol === 'https:' ) {
    makeAPIRequest();
  }
});
