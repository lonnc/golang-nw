package build

const client = `
exports.createClient = function() {
    var events = require('events');
    var channel = new events.EventEmitter();

    var path = require('path');
    var exe = path.join(path.dirname(process.cwd), '{{ .Bin }}');
    console.log('Using client: ' + exe);

    // Now start the client process
    var childProcess = require('child_process');

    var p = childProcess.spawn(exe, ['-http=localhost:0']);

    p.stdout.once('data', function(data) {
        channel.emit('starting', exe);
    });

    p.stdout.on('data', function(data) {
        var s = data.toString();
        console.log(s);

        // Check to see if we are listening now
        var listening = /HTTP listening on (\S+)/.exec(s);
        if (listening) {
            var url = 'http://' + listening[1] + '/';
            console.log('Detected that client has started on: ' + url);
            channel.emit('listening', url);
        }
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

    return channel;
};
`
