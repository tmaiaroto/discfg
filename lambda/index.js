
const MAX_FAILS = 4;

var AWS = require('aws-sdk');
var child_process = require('child_process'),
	go_proc = null,
	done = console.log.bind(console),
	fails = 0;

(function new_go_proc() {

	// pipe stdin/out, blind passthru stderr
	go_proc = child_process.spawn('./main', { stdio: ['pipe', 'pipe', process.stderr] });

	go_proc.on('error', function(err) {
		process.stderr.write("go_proc errored: "+JSON.stringify(err)+"\n");
		if (++fails > MAX_FAILS) {
			process.exit(1); // force container restart after too many fails
		}
		new_go_proc();
		done(err);
	});

	go_proc.on('exit', function(code) {
		process.stderr.write("go_proc exited prematurely with code: "+code+"\n");
		if (++fails > MAX_FAILS) {
			process.exit(1); // force container restart after too many fails
		}
		new_go_proc();
		done(new Error("Exited with code "+code));
	});

	go_proc.stdin.on('error', function(err) {
		process.stderr.write("go_proc stdin write error: "+JSON.stringify(err)+"\n");
		if (++fails > MAX_FAILS) {
			process.exit(1); // force container restart after too many fails
		}
		new_go_proc();
		done(err);
	});

	var data = null;
	go_proc.stdout.on('data', function(chunk) {
		fails = 0; // reset fails
		if (data === null) {
			data = new Buffer(chunk);
		} else {
			data.write(chunk);
		}
		// check for newline ascii char 10
		if (data.length && data[data.length-1] == 10) {
			var output = JSON.parse(data); // already a string from discfg //JSON.parse(data.toString('UTF-8'));
			data = null;
			done(null, output);
		};
	});
})();

exports.handler = function(event, context) {

	// always output to current context's done
	done = context.done.bind(context);

	// NOTE: AWS Credentials are present in environment variables. I'm unsure if they are also available to the Go process or not...
	// I assume so. If not, they can always be sent in the JSON data to the Go process.
	// var AWS = require('aws-sdk');
	// var creds = new AWS.EnvironmentCredentials('AWS');
	event.creds = new AWS.EnvironmentCredentials('AWS');

	go_proc.stdin.write(JSON.stringify({
		"event": event,
		"context": context
	})+"\n");

}
