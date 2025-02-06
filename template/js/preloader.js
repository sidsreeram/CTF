var strings = [
  "Initializing request",
  "Resolving internet address 127.0.0.1",
  "Requesting access to server",
  "Entering credentials",
  "Login successful",
  "Initializing..."
];

var preloader = document.getElementById('preloader');
var delay = 100; // Adjusted delay for better visibility
var count = 0;

function addLog() {
  var row = createLog('ok', count);
  preloader.appendChild(row);
  
  goScrollToBottom();
  
  count++;
  
  if (count < strings.length) {
    setTimeout(addLog, delay);
  } else {
    setTimeout(function() {
      preloader.style.display = 'none'; // Hide preloader
      document.getElementById('main').style.display = 'block'; // Show main content
    }, 1000);
  }
}

function createLog(type, index) {
  var row = document.createElement('div');
  
  var spanStatus = document.createElement('span');
  spanStatus.className = type;
  spanStatus.innerHTML = type.toUpperCase();
  
  var message = (index != null) 
              ? strings[index] 
              : 'kernel: Initializing...';

  var spanMessage = document.createElement('span');
  spanMessage.innerHTML = message;
  
  row.appendChild(spanStatus);
  row.appendChild(spanMessage);
  
  return row;
}

function goScrollToBottom() {
  window.scrollTo(0, document.body.scrollHeight);
}

function setCookie(cname, cvalue, exdays) {
  var d = new Date();
  d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
  var expires = "expires=" + d.toGMTString();
  document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function getCookie(cname) {
  var name = cname + "=";
  var decodedCookie = decodeURIComponent(document.cookie);
  var ca = decodedCookie.split(';');
  for (var i = 0; i < ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

function checkCookie() {
  var user = getCookie("visited");
  if (user == 1) {
    setCookie("visited", 1, 30); // Update the cookie
    document.getElementById('main').style.display = 'block'; // Show main content
  } else {
    addLog();
    setCookie("visited", 1, 30);
  }
}

checkCookie();