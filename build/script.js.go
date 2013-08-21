package build

const script = `
var gui = require('nw.gui');
var win = gui.Window.get();

win.on('loaded', function() {
    // Restore window location on startup.
    if (localStorage.width && localStorage.height) {
        win.resizeTo(parseInt(localStorage.width), parseInt(localStorage.height));
        win.moveTo(parseInt(localStorage.x), parseInt(localStorage.y));
    }

	// Add a basic menu
	var menu = new gui.Menu({
        type: 'menubar'
    });
	
	var helpMenu = new gui.Menu();
    helpMenu.append(new gui.MenuItem({
        label: 'Show Dev Tools',
        click: function() {
            win.showDevTools();
        }
    }));
    
    menu.append(new gui.MenuItem({
        label: 'Help',
        submenu: helpMenu
    }));

    win.menu = menu;
	
    // Ensure we are visible
    win.show();

    // Start client
    var client = require('./client.js');

    var msg = function(s) {
        var state = document.getElementById('state');
        state.appendChild(document.createTextNode(s + '\n'));
    };

    var clientProcess = client.createClient();
    clientProcess.
    on('error', function(err) {
        msg('Failed to start client: ' + err);
    }).
    once('starting', function(bin) {
        msg('Starting client: ' + bin);
    }).
    once('listening', function(url) {
        msg('Started');
        window.location.href = url;
    });
    
    // And kill client when we close the window
    win.on('close', function() {
        clientProcess.kill();
    });
});

// Save size on close.
win.on('close', function() {
    localStorage.x = win.x;
    localStorage.y = win.y;
    localStorage.width = win.width;
    localStorage.height = win.height;
    this.close(true);
});
`
