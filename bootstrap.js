var child_process = require('child_process');

exports.handler = function(event, context) {
    var proc = child_process.spawn(
                   "./dynamic-dynamodb-lambda",
                   [ JSON.stringify(event) ],
                   { stdio: "inherit" }
               );

    proc.on('exit', function(ret) {
        if ( ret !== 0 ) {
            return context.fail(
                       new Error("process exited with non-zero status")
                   );
        }

        context.succeed();
    });
}
