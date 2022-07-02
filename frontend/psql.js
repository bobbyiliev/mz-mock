log("Welcome to Materialize!");
log("Type 'show help;' for a list of commands.");

// Help command
register_cmd("show help;", function(cmd) {
    block_log("<u>Available commands:</u>");
    log("  &nbsp;show help; - Show this help message");
    log("  &nbsp;show docs; - Open the Materialize documentation");
    log("  &nbsp;show github; - Open the Materialize GitHub repository");
    log("  &nbsp;show demos; - Open the Materialize demos");
});

// Docs command
register_cmd("show docs;", function(cmd) {
    log("Opening Materialize documentation...");
    log("https://materialize.com/docs");
    window.open("https://materialize.com/docs", "_blank");
});

// GitHub command
register_cmd("show github;", function(cmd) {
    log("Opening Materialize GitHub repository...");
    log("https://github.com/MaterializeInc/materialize");
    window.open("https://github.com/MaterializeInc/materialize", "_blank");
});

// Demos command

register_cmd("show demos;", function(cmd) {
    log("Opening Materialize demos...");
    log("https://materialize.com/docs/quickstarts");
    window.open("https://materialize.com/docs/quickstarts", "_blank");
});

// Example Command - Hey

register_cmd("hey", function(cmd) {
    var parameters = cmd.split(" ").slice(1);
    for (var i = 0; i < parameters.length; i++) {
        block_log("Hello " + parameters[i]);
    }
});

// Example Command - Sum

register_cmd("sum", function(cmd) {
    var parameters = cmd.split(" ").slice(1);
    var sum = 0;
    for (var i = 0; i < parameters.length; i++) {
        sum += parseInt(parameters[i]);
    }
    block_log("Sum: " + sum);
});
