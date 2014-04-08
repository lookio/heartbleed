chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
  if ( request.type === "pageActionEvent" ) {
    switch (request.action) {
      case "addPropToList":
        if ( detectedProps.indexOf(request.data) !== -1 ) {
          addPropToList(request.data);
        }
        break;
    }
  }
});
