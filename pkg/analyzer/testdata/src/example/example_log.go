package example

import (
	"log"
)

var (
	password = "secret123"
	apiKey   = "key-abc-xyz"
	token    = "tok-987-654"
)

func testLog() {
	log.Print("Starting server")   // want "log message should start with lowercase letter"
	log.Print("Failed to connect") // want "log message should start with lowercase letter"
	log.Print("starting server")   // OK
	log.Print("failed to connect") // OK

	log.Print("запуск сервера")     // want "log message should contain only English letters"
	log.Print("ошибка подключения") // want "log message should contain only English letters"
	log.Print("server started")     // OK

	log.Print("server started! 🚀")                // want "log message should not contain special symbols or emojis"
	log.Print("connection failed!!!")             // want "log message should not contain special symbols or emojis"
	log.Print("warning: something went wrong...") // want "log message should not contain special symbols or emojis"
	log.Print("server started")                   // OK
	log.Print("connection failed")                // OK

	log.Print("user password: " + password) // want "log message should not contain sensitive information" "log message should not contain sensitive information"
	log.Print("api_key=" + apiKey)          // want "log message should not contain sensitive information" "log message should not contain sensitive information"
	log.Print("token: " + token)            // want "log message should not contain sensitive information" "log message should not contain sensitive information"
	log.Print("user authenticated")         // OK
	log.Print("api request completed")      // OK
	log.Print("token validated")            // OK

	log.Print("Запуск сервера! 🚀") // want "log message should start with lowercase letter" "log message should contain only English letters" "log message should not contain special symbols or emojis"

	log.Print("Password: " + password) // want "log message should not contain sensitive information" "log message should not contain sensitive information" "log message should start with lowercase letter"
}
