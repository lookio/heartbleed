
function showNotifcation( host ) {
  var notification = webkitNotifications.createNotification(
    'heartbleed-48.png',
    'Code Blue!',
    [
      host,
      ' is not a secure site.',
      '\n\n',
      'More info at www.heartbleed.com \n'
    ].join('')
  );
  notification.show();
}


chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
  if ( request && request.showAlert  ) {
    showNotifcation(request.hostName);
    chrome.pageAction.show(sender.tab.id);
  }
});


chrome.pageAction.onClicked.addListener(function() {
  chrome.tabs.query({active: true, currentWindow: true}, function(tabs){
    chrome.tabs.sendMessage(tabs[0].id, { triggerAlert : true });
  });

});
