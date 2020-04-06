window.addEventListener("load", function(evt) {
    var oldLogo = document.getElementById("old");
    var newLogo = document.getElementById("new");

    var ws = new WebSocket("ws://" + window.location.host + "/consume.ws");
    ws.onmessage = function(evt) {
        switch (evt.data) {
            case "new":
                oldLogo.style.display = "none";
                newLogo.style.display = "block";
                break;
            case "old":
                newLogo.style.display = "none";
                oldLogo.style.display = "block";
                break;
        }
        return false
    }
});
