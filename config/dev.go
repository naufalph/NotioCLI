package config

/*
Package config provides application configuration utilities.

The `dev.go` file contains functions to retrieve the location of the `.env` file for the development environment.
The `.env` file is used to store environment variables necessary for database and application configuration.

### Example `.env` File

```env
# Database configuration
MYSQL_USER=[your_mysql_user]
MYSQL_PASSWORD=[your_mysql_password]
MYSQL_DB_NAME=[your_database_name]
MYSQL_DB_PORT=[your_database_port]
MYSQL_DB_HOST=[your_database_host]
MYSQL_PROTOCOL=tcp
*/

func DevEnv() string {
	return "../.env" //Enter your dev environment
}
