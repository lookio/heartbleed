(function() {
  var dataElem = document.getElementById('lpd-plugin-data');
  function main() {
    var event = new Event('data');
    var keys = Object.keys(window);
    dataElem.innerHTML = JSON.stringify(keys);
    dataElem.dispatchEvent(event);
  }
  main();
  setTimeout(main, 5*1000);
  //setInterval(main, 5*1000);
})();
