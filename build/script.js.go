package build

const script = `
"use strict";

var gui = require('nw.gui');
var win = gui.Window.get();

win.on('loaded', function() {
    // Restore window location on startup.
    if (window.localStorage.width && window.localStorage.height) {
        win.resizeTo(parseInt(window.localStorage.width), parseInt(window.localStorage.height));
        win.moveTo(parseInt(window.localStorage.x), parseInt(window.localStorage.y));
    }

    // Ensure we are visible
    win.show();

    // Start client
    var client = require('./client.js');

    var msg = function(s) {
        var state = document.getElementById('state');
        state.appendChild(document.createTextNode(s + '\n'));
    };

    var clientProcess = client.createClient(gui.App.argv);
    clientProcess.
    on('error', function(err) {
        msg('Error: ' + err);
    }).
    on('redirect', function(url) {
        window.location.href = url;
    });

    // And kill client when we close the window
    win.on('close', function() {
        clientProcess.kill();
    });
});

// Save size on close.
win.on('close', function() {
    window.localStorage.x = win.x;
    window.localStorage.y = win.y;
    window.localStorage.width = win.width;
    window.localStorage.height = win.height;
    this.close(true);
});
`
