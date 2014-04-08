
function makeRequest() {
  var xhr = new XMLHttpRequest();
  var payload = "url=gwd.lphbs.com";
  xhr.open("POST", "http://localhost:8000/api/v1/url", true);
  xhr.setRequestHeader("Access-Control-Allow-Origin", "*");
  xhr.onload = function() {

  };
  xhr.onerror = function() {

  };
  xhr.send(payload);
}

makeRequest();
