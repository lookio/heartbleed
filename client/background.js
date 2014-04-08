
function makeRequest() {
  var xhr = new XMLHttpRequest();
  xhr.open("POST", "http://localhost:8000/api/v1/url", false);
  xhr.addEventListener('load', function(e) {
    console.log(e);
  }, false);
  xhr.send("url=gwd.lphbs.com");
}

makeRequest();
