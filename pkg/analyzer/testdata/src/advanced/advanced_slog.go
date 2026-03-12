package advanced

import (
	"fmt"
	"log/slog"
)

func testSlogAdvanced() {
	slog.Info("Starting server")    // want "log message should start with lowercase letter"
	slog.Error("Failed to connect") // want "log message should start with lowercase letter"
	slog.Info("starting server")    // OK
	slog.Error("failed to connect") // OK

	slog.Info("Starting " + "server")        // want "log message should start with lowercase letter"
	slog.Error("Failed " + "to" + "connect") // want "log message should start with lowercase letter"
	slog.Info("starting " + "server")        // OK
	slog.Error("failed" + "to" + "connect")  // OK

	slog.Info("user logged in", "user", userName) // OK
	slog.Info("api call", "key", apiKey)          // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	slog.Info("token check", "token", token)      // want "log message should not contain sensitive information" "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("user action", "user", userName+"123") // OK
	slog.Info("auth", "password", password)          // want "log message should not contain sensitive information" "log message should not contain sensitive information" "log variable name suggests sensitive data"
	slog.Info("api", "api_key", apiKey)              // want "log variable name suggests sensitive data" "log message should not contain sensitive information"
	slog.Info("auth", "token", token)                // want "log message should not contain sensitive information" "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("us"+"er", "user", userName)      // OK
	slog.Info("pass"+"word", "value", password) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	slog.Info("ap"+"i_key", "value", apiKey)    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprintf("server %s", serverName)) // want "log message should start with lowercase letter"
	slog.Error("Failed " + fmt.Sprintf("to %s", "connect"))       // want "log message should start with lowercase letter"
	slog.Info("user " + fmt.Sprintf("%s logged in", userName))    // OK
	slog.Info("password " + fmt.Sprintf("value: %s", password))   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	slog.Info("api " + fmt.Sprintf("key=%s", apiKey))             // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info(fmt.Sprintf("Starting %s", serverName)) // want "log message should start with lowercase letter"
	slog.Info(fmt.Sprintf("password: %s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	slog.Info(fmt.Sprintf("api_key=%s", apiKey))      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("msg: " + fmt.Sprintf("Starting %s", serverName))   // OK
	slog.Info("msg: " + fmt.Sprintf("password: %s", password))    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	slog.Info("msg: prefix - " + fmt.Sprintf("token: %s", token)) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info(fmt.Sprintf("Starting %s", serverName), "env", "prod") // want "log message should start with lowercase letter"
	slog.Info(fmt.Sprintf("password: %s", password), "env", "prod")  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	slog.Info(fmt.Sprintf("api_key=%s", apiKey), "env", "prod")      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("user "+userName+" has pass "+password, "action", "login") // want "log variable name suggests sensitive data"
	slog.Info("api "+"service with key "+apiKey, "action", "call")       // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprintf("server %s", serverName) + " ok") // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprintf("value %s", password) + " end")   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprintf("server %s", fmt.Sprintf("%s", serverName))) // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprintf("value %s", fmt.Sprintf("%s", password)))    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("user " + userName + fmt.Sprintf(" %s", "logged"))     // OK
	slog.Info("password " + password + fmt.Sprintf("=%s", "secret")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprint("server ", serverName))  // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprint("value ", password))     // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	slog.Info("Starting " + fmt.Sprintln("server", serverName)) // want "log message should start with lowercase letter"

	slog.Info("Starting "+fmt.Sprintf("server %s", serverName), "status", "ok")  // want "log message should start with lowercase letter"
	slog.Info("password "+fmt.Sprintf("value %s", password), "status", "hidden") // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprintf("srv_%s", serverName)) // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprintf("val_%s", password))   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + serverName + " (prod)") // want "log message should start with lowercase letter" "log message should not contain special symbols or emojis"
	slog.Info("password " + password + " (hidden)") // want "log message should not contain sensitive information" "log variable name suggests sensitive data" "log message should not contain special symbols or emojis"

	slog.Info("Starting " + fmt.Sprintf("server %s", fmt.Sprintf("%s", serverName)) + " end") // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprintf("value %s", fmt.Sprintf("%s", password)) + " end")    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + "srv_" + fmt.Sprintf("%s", serverName)) // want "log message should start with lowercase letter"
	slog.Info("password " + "pass_" + fmt.Sprintf("%s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("user " + userName + " pass " + fmt.Sprintf("val_%s", password)) // want "log variable name suggests sensitive data"
	slog.Info("api " + "svc" + " key " + fmt.Sprintf("k_%s", apiKey))          // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprintf("srv:%s", serverName) + ":8080") // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprintf("val:%s", password) + ":end")    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprintf("server %s", serverName) + fmt.Sprintf(" %s", "ok")) // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprintf("pass %s", password) + fmt.Sprintf(" %s", "hidden")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprintf("server %s", fmt.Sprintf("%s", serverName)) + " end") // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprintf("pass %s", fmt.Sprintf("%s", password)) + " end")     // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + "srv_" + fmt.Sprintf("%s", serverName)) // want "log message should start with lowercase letter"
	slog.Info("password " + "pass_" + fmt.Sprintf("%s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	slog.Info("Starting " + fmt.Sprintf("srv_%s", serverName) + "_" + fmt.Sprintf("%s", "ok"))    // want "log message should start with lowercase letter"
	slog.Info("password " + fmt.Sprintf("pass_%s", password) + "_" + fmt.Sprintf("%s", "hidden")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
}
