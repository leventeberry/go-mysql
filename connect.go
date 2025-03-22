package gomysql

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type EnvConfig struct {
	user string
	pass string
}
	

// getEnvOrPrompt checks an environment variable or prompts the user if unset.
func getEnvOrPrompt(envVar, prompt, defaultValue string) string {
	val := os.Getenv(envVar)
	if val == "" {
		fmt.Println(prompt)
		fmt.Printf("(default: %s): ", defaultValue)

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			val = defaultValue
		} else {
			val = input
		}
	}
	return val
}

// getEnv checks an environment variable or prompts the user if unset.
func checkEnv() EnvConfig {
	if os.Getenv("APP_ENV") == "dev" {
		user := getEnvOrPrompt("DB_USER", "Enter the database username", "root")
		pass := getEnvOrPrompt("DB_PASS", "Enter the database password", "password")
		return EnvConfig{user, pass}
	} else {
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		return EnvConfig{user, pass}	
	}
}

// ConnectDB initializes the MySQL connection and returns a *sql.DB instance.
func ConnectDB(net, addr, dbName string) (*sql.DB, error) {
	// Check environment variables
	env := checkEnv()
	fmt.Println("[🔐] Environment variables loaded")

	// Configure database connection
	cfg := mysql.Config{
		User:   env.user,
		Passwd: env.pass,
		Net:    net,
		Addr:   addr,
		DBName: dbName,
	}

	// Open a connection to the database
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("[✅] Connected to the database")
	return db, nil
}