////////////
//Tools
////////////
function intArrayToString(intArray) {
  var binary = '';
  var length = intArray.length;
  for (var t = 0; t < length; t++) {
    i = intArray[t];
    binary += String.fromCharCode((i & 0xFF000000) >>> 24);
    binary += String.fromCharCode((i & 0x00FF0000) >>> 16);
    binary += String.fromCharCode((i & 0x0000FF00) >>> 8);
    binary += String.fromCharCode(i & 0x000000FF);
  }
  return binary;
}

function intArrayToBase64(intArray) {
  var binary = intArrayToString(intArray);
  return window.btoa( binary );
}

//////////////
//encrypt page
//////////////
function encryptAndSend() {
  var current_date = (new Date()).valueOf().toString();
  var random = Math.random().toString();
  var key = intArrayToString(sjcl.hash.sha256.hash(current_date + random));
  var text = document.getElementById('word').value;
  var val = sjcl.encrypt(key, text);

  current_date = (new Date()).valueOf().toString();
  random = Math.random().toString();
  var pkey = sjcl.hash.sha256.hash(current_date + random);

  var link = "localhost:8080/paste/" + 
  encodeURIComponent(intArrayToBase64(pkey)) + "#" + 
  encodeURIComponent(btoa(key));

  //request server for stuff
  var obj = {"key" : intArrayToBase64(pkey), "value" : val};

  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 201) {
      document.getElementById("submit").style.visibility = "hidden";
      document.getElementById("link").style.visibility = "visible";
      document.getElementById("link_text").value = link;
    }
  };
  xhttp.open("POST", "api/post/", true);
  xhttp.setRequestHeader("Content-type", "application/json");
  xhttp.send(JSON.stringify(obj));
}

//////////////
//decrypt page
//////////////
function decryptAndDisplay() {
  var path = window.location.pathname;
  var parts = path.split('/');
  //console.log(parts);
  if (parts.length < 3) {
    document.getElementById("paste").value = "you seem to have an invalid link";
    return;
  }
  var key = parts[2];
  var obj = {"key" : decodeURIComponent(key)}
  
  var xhttp = new XMLHttpRequest();
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      s = JSON.parse(this.responseText).value
      decryptAndDisplayHelper(JSON.stringify(s));
    }
  };
  xhttp.open("POST", "../api/retrieve/", true);
  xhttp.setRequestHeader("Content-type", "application/json");
  xhttp.send(JSON.stringify(obj));
}

function decryptAndDisplayHelper(message) {
  try {
    if (window.location.hash) {
      var key = atob(decodeURIComponent(window.location.hash).substring(1));
      var text = sjcl.decrypt(key, message);
      document.getElementById("paste").value = text;
    } else {
      document.getElementById("paste").value = "You don't seem to have the key or have an incorrect key";
    } 
  }
  catch (err) {
    document.getElementById("paste").value = "You don't seem to have the key or have an incorrect key";
  }
}
