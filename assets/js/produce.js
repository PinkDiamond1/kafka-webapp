window.addEventListener("load", function(_) {
    var ws = new WebSocket("ws://" + window.location.host + "/produce.ws");
    var autoOn = document.getElementById("auto-on");
    var autoOff = document.getElementById("auto-off");
    var auto = false;
    var autoMessage = "new";

    document.getElementById("new").onclick = function(_) {
        console.log("sending 'new'");
        ws.send("new");
        return false;
    };
    document.getElementById("old").onclick = function(_) {
        console.log("sending 'old'");
        ws.send("old");
        return false;
    };
    document.onkeypress = function(event) {
        switch (event.key) {
            case "n":
            case "N":
                console.log("sending 'new'");
                ws.send("new");
                break;
            case "o":
            case "O":
                console.log("sending 'old'");
                ws.send("old");
                break;
        }
        return false;
    }
    document.getElementById("auto-on").onclick = function() {
        auto = true;
        autoOn.style.display = "none";
        autoOff.style.display = "inline";
        return false
    }
    document.getElementById("auto-off").onclick = function() {
        auto = false;
        autoOff.style.display = "none";
        autoOn.style.display = "inline";
        return false
    }

    setInterval(function() {
        if (auto) {
            console.log("sending '" + autoMessage + "'");
            ws.send(autoMessage);
            if (autoMessage == "new") {
                autoMessage = "old";
            } else {
                autoMessage = "new";
            }
        }
    }, 200);
});