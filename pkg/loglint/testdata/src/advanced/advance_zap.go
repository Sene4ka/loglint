package advanced

import (
	"fmt"
	"go.uber.org/zap"
)

func testZapAdvanced(logger *zap.Logger) {
	logger.Info("Starting server")    // want "log message should start with lowercase letter"
	logger.Error("Failed to connect") // want "log message should start with lowercase letter"
	logger.Info("starting server")    // OK
	logger.Error("failed to connect") // OK

	logger.Info("Starting " + "server")        // want "log message should start with lowercase letter"
	logger.Error("Failed " + "to" + "connect") // want "log message should start with lowercase letter"
	logger.Info("starting " + "server")        // OK
	logger.Error("failed" + "to" + "connect")  // OK

	logger.Info("user logged in", zap.String("user", userName)) // OK
	logger.Info("api call", zap.String("key", apiKey))          // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("token check", zap.String("token", token))      // want "log message should not contain sensitive information" "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("user action", zap.String("user", userName+"123")) // OK
	logger.Info("auth", zap.String("password", password))          // want "log message should not contain sensitive information" "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("api", zap.String("api_key", apiKey))              // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("auth", zap.String("token", token))                // want "log message should not contain sensitive information" "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("us"+"er", zap.String("user", userName))      // OK
	logger.Info("pass"+"word", zap.String("value", password)) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("ap"+"i_key", zap.String("value", apiKey))    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprintf("server %s", serverName)) // want "log message should start with lowercase letter"
	logger.Error("Failed " + fmt.Sprintf("to %s", "connect"))       // want "log message should start with lowercase letter"
	logger.Info("user " + fmt.Sprintf("%s logged in", userName))    // OK
	logger.Info("password " + fmt.Sprintf("value: %s", password))   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("api " + fmt.Sprintf("key=%s", apiKey))             // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info(fmt.Sprintf("Starting %s", serverName)) // want "log message should start with lowercase letter"
	logger.Info(fmt.Sprintf("password: %s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info(fmt.Sprintf("api_key=%s", apiKey))      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("msg: " + fmt.Sprintf("Starting %s", serverName))   // OK
	logger.Info("msg: " + fmt.Sprintf("password: %s", password))    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("msg: prefix - " + fmt.Sprintf("token: %s", token)) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info(fmt.Sprintf("Starting %s", serverName), zap.String("env", "prod")) // want "log message should start with lowercase letter"
	logger.Info(fmt.Sprintf("password: %s", password), zap.String("env", "prod"))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info(fmt.Sprintf("api_key=%s", apiKey), zap.String("env", "prod"))      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("user "+userName+" has pass "+password, zap.String("action", "login")) // want "log variable name suggests sensitive data"
	logger.Info("api "+"service with key "+apiKey, zap.String("action", "call"))       // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprintf("server %s", serverName) + " ok") // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprintf("value %s", password) + " end")   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprintf("server %s", fmt.Sprintf("%s", serverName))) // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprintf("value %s", fmt.Sprintf("%s", password)))    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("user " + userName + fmt.Sprintf(" %s", "logged"))     // OK
	logger.Info("password " + password + fmt.Sprintf("=%s", "secret")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprint("server ", serverName))  // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprint("value ", password))     // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("Starting " + fmt.Sprintln("server", serverName)) // want "log message should start with lowercase letter"

	logger.Info("Starting "+fmt.Sprintf("server %s", serverName), zap.String("status", "ok"))  // want "log message should start with lowercase letter"
	logger.Info("password "+fmt.Sprintf("value %s", password), zap.String("status", "hidden")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprintf("srv_%s", serverName)) // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprintf("val_%s", password))   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + serverName + " (prod)") // want "log message should start with lowercase letter" "log message should not contain special symbols or emojis"
	logger.Info("password " + password + " (hidden)") // want "log message should not contain sensitive information" "log variable name suggests sensitive data" "log message should not contain special symbols or emojis"

	logger.Info("Starting " + fmt.Sprintf("server %s", fmt.Sprintf("%s", serverName)) + " end") // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprintf("value %s", fmt.Sprintf("%s", password)) + " end")    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + "srv_" + fmt.Sprintf("%s", serverName)) // want "log message should start with lowercase letter"
	logger.Info("password " + "pass_" + fmt.Sprintf("%s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("user " + userName + " pass " + fmt.Sprintf("val_%s", password)) // want "log variable name suggests sensitive data"
	logger.Info("api " + "svc" + " key " + fmt.Sprintf("k_%s", apiKey))          // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprintf("srv:%s", serverName) + ":8080") // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprintf("val:%s", password) + ":end")    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprintf("server %s", serverName) + fmt.Sprintf(" %s", "ok")) // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprintf("pass %s", password) + fmt.Sprintf(" %s", "hidden")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprintf("server %s", fmt.Sprintf("%s", serverName)) + " end") // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprintf("pass %s", fmt.Sprintf("%s", password)) + " end")     // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + "srv_" + fmt.Sprintf("%s", serverName)) // want "log message should start with lowercase letter"
	logger.Info("password " + "pass_" + fmt.Sprintf("%s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	logger.Info("Starting " + fmt.Sprintf("srv_%s", serverName) + "_" + fmt.Sprintf("%s", "ok"))    // want "log message should start with lowercase letter"
	logger.Info("password " + fmt.Sprintf("pass_%s", password) + "_" + fmt.Sprintf("%s", "hidden")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
}
