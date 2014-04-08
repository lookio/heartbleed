chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
  if ( request && request.propExists ) {
    //console.log("Prop exists on site: ", request.propName);
    chrome.browserAction.setIcon({path: 'cog.png'});
  }

});
