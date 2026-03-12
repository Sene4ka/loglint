package basic

import (
	"go.uber.org/zap"
)

func testZapBasic(logger *zap.Logger) {
	logger.Info("Starting server")    // want "log message should start with lowercase letter"
	logger.Error("Failed to connect") // want "log message should start with lowercase letter"
	logger.Info("starting server")    // OK
	logger.Error("failed to connect") // OK

	logger.Info("Starting " + "server")        // want "log message should start with lowercase letter"
	logger.Error("Failed " + "to" + "connect") // want "log message should start with lowercase letter"
	logger.Info("starting " + "server")        // OK
	logger.Error("failed" + "to" + "connect")  // OK

	logger.Info("запуск сервера")      // want "log message should contain only English letters"
	logger.Error("ошибка подключения") // want "log message should contain only English letters"
	logger.Info("server started")      // OK

	logger.Info("запуск " + "сервера")      // want "log message should contain only English letters"
	logger.Error("ошибка " + "подключения") // want "log message should contain only English letters"
	logger.Info("server " + "started")      // OK

	logger.Info("server started! 🚀")                // want "log message should not contain special symbols or emojis"
	logger.Error("connection failed!!!")            // want "log message should not contain special symbols or emojis"
	logger.Warn("warning: something went wrong...") // want "log message should not contain special symbols or emojis"
	logger.Info("server started")                   // OK
	logger.Error("connection failed")               // OK

	logger.Info("server " + "started!" + "🚀")                 // want "log message should not contain special symbols or emojis"
	logger.Error("connection " + " failed!!!")                // want "log message should not contain special symbols or emojis"
	logger.Warn("warning: " + "something went wrong." + "..") // want "log message should not contain special symbols or emojis"
	logger.Info("server " + "started")                        // OK
	logger.Error("connection" + " failed")                    // OK

	logger.Info("user password: " + password) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Debug("api_key=" + apiKey)         // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("token: " + token)            // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("user authenticated")         // OK
	logger.Debug("api request completed")     // OK
	logger.Info("token validated")            // OK

	logger.Info("us" + "er pass" + "word: " + password) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Debug("ap" + "i_" + "key=" + apiKey)         // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("to" + "ken: " + token)                 // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	logger.Info("us" + "er au" + "thenti" + "cated")    // OK
	logger.Debug("ap" + "i reques" + "t completed")     // OK
	logger.Info("tok" + "en valid" + " ated")           // OK

	logger.Info("Запуск сервера! 🚀") // want "log message should start with lowercase letter" "log message should contain only English letters" "log message should not contain special symbols or emojis"

	logger.Info("Password: " + password) // want "log message should start with lowercase letter" "log message should not contain sensitive information" "log variable name suggests sensitive data"
}
