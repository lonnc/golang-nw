package build

const client = `
"use strict";

exports.createClient = function(args) {
    var events = require('events');
    var channel = new events.EventEmitter();
	var http = require('http');
    var server = http.createServer(function (request, response) {
	  if (request.method!='POST') {
	    response.writeHead(404);
        response.end('');
		return;
	  }
		
	  var body = '';
	  request.on('data', function(chunk) { body += chunk.toString(); });
	  request.on('end', function() {
		switch (request.url) {
		case '/redirect': 
		  channel.emit('redirect', body);
		  response.writeHead(204);
		  response.end('');
		  break;
		case '/error': 
		  channel.emit('error', body);
		  response.writeHead(204);
		  response.end('');
		  break;
		default:
		  response.writeHead(404);
		  response.end('');
		 };
	  });
    });

    server.listen(0, '127.0.0.1', 1, function() {
        var nodeWebkitAddr = 'http://127.0.0.1:'+server.address().port;
        console.log('Listening for golang-nw on '+nodeWebkitAddr);
        startClient(channel, nodeWebkitAddr, args);
    });
	
	return channel;
};
	
function startClient(channel, nodeWebkitAddr, args) {
    var path = require('path');
    var exe = '.'+path.sep+'{{ .Bin }}';
    console.log('Using client: ' + exe);

    // Make the exe executable
    var fs = require('fs');
    fs.chmodSync(exe, '755');

    // Now start the client process
    var childProcess = require('child_process');

	var env = process.env;
	env['{{ .EnvVar }}'] = nodeWebkitAddr;
    var p = childProcess.spawn(exe, args, {env: env});

    p.stdout.on('data', function(data) {
        console.log(data.toString());
    });
	
    p.stderr.on('data', function(data) {
        console.error(data.toString());
    });

    p.on('error', function(err) {
        console.error('child error: ' + err);
        channel.emit('error', err);
    });

    p.on('close', function(code) {
        console.log('child process closed with code ' + code);
        channel.emit('close', code);
    });

    p.on('exit', function(code) {
        console.log('child process exited with code ' + code);
        channel.emit('exit', code);
    });

    channel.kill = function() {
        p.kill();
    }
};
`
