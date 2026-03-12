package advanced

import (
	"fmt"
	"log"
)

var (
	password   = "secret123"
	apiKey     = "key-abc-xyz"
	token      = "tok-987-654"
	userName   = "john"
	serverName = "prod-server"
)

func testLogAdvanced() {
	log.Printf("Starting server %s", serverName)      // want "log message should start with lowercase letter"
	log.Printf("Failed to connect to %s", serverName) // want "log message should start with lowercase letter"
	log.Printf("starting server %s", serverName)      // OK
	log.Printf("failed to connect %s", serverName)    // OK

	log.Printf("Starting %s", "server")   // want "log message should start with lowercase letter"
	log.Printf("Failed %s", "to connect") // want "log message should start with lowercase letter"
	log.Printf("starting %s", "server")   // OK
	log.Printf("failed %s", "to connect") // OK

	log.Printf("user %s logged in", userName) // OK
	log.Printf("api key: %s", apiKey)         // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("token value: %s", token)      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("user %s", userName+"123") // OK
	log.Printf("password: %s", password)  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("api_key=%s", apiKey)      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("token: %s", token)        // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("us"+"er %s", userName)      // OK
	log.Printf("pass"+"word: %s", password) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("ap"+"i_key=%s", apiKey)     // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprintf("server %s", serverName)) // want "log message should start with lowercase letter"
	log.Printf("Failed %s", fmt.Sprintf("to %s", "connect"))        // want "log message should start with lowercase letter"
	log.Printf("user %s", fmt.Sprintf("%s logged in", userName))    // OK
	log.Printf("password %s", fmt.Sprintf("value: %s", password))   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("api %s", fmt.Sprintf("key=%s", apiKey))             // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("%s", fmt.Sprintf("Starting %s", serverName)) // OK
	log.Printf("%s", fmt.Sprintf("password: %s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("%s", fmt.Sprintf("api_key=%s", apiKey))      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("msg: %s", fmt.Sprintf("Starting %s", serverName))         // OK
	log.Printf("msg: %s", fmt.Sprintf("password: %s", password))          // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("msg: %s - %s", "prefix", fmt.Sprintf("token: %s", token)) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf(fmt.Sprintf("Starting %s", serverName)) // want "log message should start with lowercase letter"
	log.Printf(fmt.Sprintf("password: %s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf(fmt.Sprintf("api_key=%s", apiKey))      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("user %s has pass %s", userName, password) // want "log variable name suggests sensitive data"
	log.Printf("api %s with key %s", "service", apiKey)   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("token %s for user %s", token, userName)   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("user %s", userName+" has pass "+password) // want "log variable name suggests sensitive data"
	log.Printf("api %s", "service with key "+apiKey)      // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprintf("server %s", serverName)+" ok") // want "log message should start with lowercase letter"
	log.Printf("password %s", fmt.Sprintf("value %s", password)+" end")   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprintf("server %s", fmt.Sprintf("%s", serverName))) // want "log message should start with lowercase letter"
	log.Printf("password %s", fmt.Sprintf("value %s", fmt.Sprintf("%s", password)))    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("user %s", userName+fmt.Sprintf(" %s", "logged"))     // OK
	log.Printf("password %s", password+fmt.Sprintf("=%s", "secret")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprint("server ", serverName))  // want "log message should start with lowercase letter"
	log.Printf("password %s", fmt.Sprint("value ", password))     // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("Starting %s", fmt.Sprintln("server", serverName)) // want "log message should start with lowercase letter"

	log.Printf("%s %s", "Starting", fmt.Sprintf("server %s", serverName)) // OK
	log.Printf("%s %s", "password:", fmt.Sprintf("value %s", password))   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", "server"+"_"+serverName) // want "log message should start with lowercase letter"
	log.Printf("password %s", "value"+"_"+password)    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
	log.Printf("api %s", "key"+"_"+apiKey)             // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprintf("srv_%s", serverName)) // want "log message should start with lowercase letter"
	log.Printf("password %s", fmt.Sprintf("val_%s", password))   // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", serverName+" (prod)") // want "log message should start with lowercase letter" "log message should not contain special symbols or emojis"
	log.Printf("password %s", password+" (hidden)") // want "log message should not contain sensitive information" "log variable name suggests sensitive data" "log message should not contain special symbols or emojis"

	log.Printf("Starting %s", fmt.Sprintf("%s [OK]", serverName))     // want "log message should start with lowercase letter" "log message should not contain special symbols or emojis"
	log.Printf("password %s", fmt.Sprintf("%s [REDACTED]", password)) // want "log message should not contain sensitive information" "log variable name suggests sensitive data" "log message should not contain special symbols or emojis"

	log.Printf("user %s pass %s", userName, fmt.Sprintf("val_%s", password)) // want "log variable name suggests sensitive data"
	log.Printf("api %s key %s", "svc", fmt.Sprintf("k_%s", apiKey))          // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprintf("srv:%s", serverName)+":8080") // want "log message should start with lowercase letter"
	log.Printf("password %s", fmt.Sprintf("val:%s", password)+":end")    // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprintf("server %s", serverName)+fmt.Sprintf(" %s", "ok")) // want "log message should start with lowercase letter"
	log.Printf("password %s", fmt.Sprintf("pass %s", password)+fmt.Sprintf(" %s", "hidden")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprintf("server %s", fmt.Sprintf("%s", serverName))+" end") // want "log message should start with lowercase letter"
	log.Printf("password %s", fmt.Sprintf("pass %s", fmt.Sprintf("%s", password))+" end")     // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", "srv_"+fmt.Sprintf("%s", serverName)) // want "log message should start with lowercase letter"
	log.Printf("password %s", "pass_"+fmt.Sprintf("%s", password))  // want "log message should not contain sensitive information" "log variable name suggests sensitive data"

	log.Printf("Starting %s", fmt.Sprintf("srv_%s", serverName)+"_"+fmt.Sprintf("%s", "ok"))    // want "log message should start with lowercase letter"
	log.Printf("password %s", fmt.Sprintf("pass_%s", password)+"_"+fmt.Sprintf("%s", "hidden")) // want "log message should not contain sensitive information" "log variable name suggests sensitive data"
}
