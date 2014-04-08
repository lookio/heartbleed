var LP_PROPS = [
  'lpTag',
  '_LP_CFG_',
  'LPMobile'
];

function main() {
  //Append script to host page to open communcation channel
  var script = document.createElement('script');
  script.src = chrome.extension.getURL("lp-discover.js");
  document.body.appendChild(script);

  var dataElem = document.createElement('div');
  dataElem.id = 'lpd-plugin-data';
  dataElem.style.display = 'none';

  dataElem.addEventListener('data', function() {
    var data = JSON.parse(dataElem.innerHTML);
    checkForProps(data);
    dataElem.innerHTML = '';
  });

  document.body.appendChild(dataElem);
}

function checkForProps(windowProps) {
  windowProps.forEach(function(prop) {
    if ( LP_PROPS.indexOf(prop) !== -1 ) {
      //console.log('Prop exists', prop);
      notifyExtension(prop);
    }
  });
}

function notifyExtension(prop) {
  chrome.runtime.sendMessage({
    propExists : true, 
    propName: prop
  });
}

main();
