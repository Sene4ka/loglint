package example

import (
	"log/slog"
)

func testSlog() {
	slog.Info("Starting server")    // want "log message should start with lowercase letter"
	slog.Error("Failed to connect") // want "log message should start with lowercase letter"
	slog.Info("starting server")    // OK
	slog.Error("failed to connect") // OK

	slog.Info("запуск сервера")      // want "log message should contain only English letters"
	slog.Error("ошибка подключения") // want "log message should contain only English letters"
	slog.Info("server started")      // OK

	slog.Info("server started! 🚀")                // want "log message should not contain special symbols or emojis"
	slog.Error("connection failed!!!")            // want "log message should not contain special symbols or emojis"
	slog.Warn("warning: something went wrong...") // want "log message should not contain special symbols or emojis"
	slog.Info("server started")                   // OK
	slog.Error("connection failed")               // OK

	slog.Info("user password: " + password) // want "log message should not contain sensitive information" "log message should not contain sensitive information"
	slog.Debug("api_key=" + apiKey)         // want "log message should not contain sensitive information" "log message should not contain sensitive information"
	slog.Info("token: " + token)            // want "log message should not contain sensitive information" "log message should not contain sensitive information"
	slog.Info("user authenticated")         // OK
	slog.Debug("api request completed")     // OK
	slog.Info("token validated")            // OK

	slog.Info("Запуск сервера! 🚀") // want "log message should start with lowercase letter" "log message should contain only English letters" "log message should not contain special symbols or emojis"

	slog.Info("Password: " + password) // want "log message should start with lowercase letter" "log message should not contain sensitive information" "log message should not contain sensitive information"
}
