package build

const client = `
"use strict";

exports.createClient = function() {
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
	    console.log('Listening for golang-nw on http://127.0.0.1:'+server.address().port);
		startClient(channel, 'http://127.0.0.1:'+server.address().port);
	});
	
	return channel;
};
	
function startClient(channel, nodeWebkitAddr) {
    var path = require('path');
    var exe = path.join(path.dirname(process.cwd), '{{ .Bin }}');
    console.log('Using client: ' + exe);

    // Now start the client process
    var childProcess = require('child_process');

	var env = process.env;
	env['GOLANG-NW'] = nodeWebkitAddr;
    var p = childProcess.spawn(exe, [], {env: env});

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
