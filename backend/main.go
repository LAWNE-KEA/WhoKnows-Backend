package main

import (
	"fmt"
	"net/http"

	"whoKnows/api"
	"whoKnows/api/configs"
	"whoKnows/database"
)

// to fix the css issue see: https://stackoverflow.com/questions/13302020/rendering-css-in-a-go-web-application or https://stackoverflow.com/questions/43601359/how-do-i-serve-css-and-js-in-go for a newer solution

func main() {
	fmt.Println("Initializing server...")

	// Load environment variables
	if err := configs.LoadEnv(); err != nil {
		fmt.Printf("Error loading environment variables: %s", err)
		return
	}

	database.InitDatabase()

	handler := api.CreateRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	fmt.Println("Server running on port 8080")
	server.ListenAndServe()
}

// func parseSQLCommands(sqlCommands string) []string {
// 	var commands []string
// 	var currentCommand strings.Builder
// 	inSingleQuote := false
// 	inDoubleQuote := false

// 	for _, char := range sqlCommands {
// 		switch char {
// 		case '\'':
// 			if !inDoubleQuote {
// 				inSingleQuote = !inSingleQuote
// 			}
// 		case '"':
// 			if !inSingleQuote {
// 				inDoubleQuote = !inDoubleQuote
// 			}
// 		case ';':
// 			if !inSingleQuote && !inDoubleQuote {
// 				commands = append(commands, currentCommand.String())
// 				currentCommand.Reset()
// 				continue
// 			}
// 		}
// 		currentCommand.WriteRune(char)
// 	}

// 	// Add the last command if any
// 	if currentCommand.Len() > 0 {
// 		commands = append(commands, currentCommand.String())
// 	}

// 	return commands
// }
