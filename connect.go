package gomysql

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/go-sql-driver/mysql"
)

// ConnectDB initializes the MySQL connection and returns a *sql.DB instance.
func ConnectDB(net string, addr string, db_name string) (*sql.DB, error) {
	if os.Getenv("DB_USER") == "" {
		var user string;
		fmt.Println("DB_USER environment variable not set.")
		fmt.Println("What is the username for the database? (default: root): ")
		fmt.Scanln(&user);
		if user == "" {
			user = "root";
		} else {
			os.Setenv("DB_USER", user);
		}
	} else {
		fmt.Println("DB_USER environment variable set to:", os.Getenv("DB_USER"))
	}

	if os.Getenv("DB_PASS") == "" {
		var pass string;
		fmt.Println("DB_PASS environment variable not set.")
		fmt.Println("What is the password for the database? (default: password): ")
		fmt.Scanln(&pass);
		if pass == "" {
			pass = "password";
		} else {
			os.Setenv("DB_PASS", pass);
		}
	} else {
		fmt.Println("DB_PASS environment variable set to:", os.Getenv("DB_PASS"))
	}

	// Configure database connection using environment variables
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    net,
		Addr:   addr,
		DBName: db_name,
	}

	// Open a connection to the database
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close() // Close DB before returning error
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to the database")
	return db, nil
}