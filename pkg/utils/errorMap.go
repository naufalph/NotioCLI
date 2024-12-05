package utils

/**
This file provides utilities for managing and standardizing error messages across the application.

The `errorMap.go` file contains constants for common error messages used throughout the codebase. These constants
help ensure consistency in error logging and messaging while reducing duplication and maintenance effort.
*/

const (
	GeneralError      = "General Error"
	DBConnectionError = "Fail to connect to DB"
	DBCloseError      = "Fail to close DB"
)
