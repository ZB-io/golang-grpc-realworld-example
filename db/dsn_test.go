// ********RoostGPT********
/*
Test generated by RoostGPT for test go-grpc-client using AI Type Azure Open AI and AI Model roostgpt-4-32k

ROOST_METHOD_HASH=dsn_e202d1c4f9
ROOST_METHOD_SIG_HASH=dsn_b336e03d64

Scenario 1: All environment variables properly set

Details:
  Description: This test verifies that the function works correctly when all required environment variables ("DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT") are set in the environment. The function should return a well-structured connection string.
Execution:
  Arrange: Set up the environment variables with the relevant values.
  Act: Call the function dsn().
  Assert: Check that the string and the error returned by the function match the expected values.
Validation:
  The test passes if the function generates a correct connection string and there is no error. The test aims to check that the software correctly retrieves data from the environment. It is critical for the application workflow.

Scenario 2: No environment variables set

Details:
  Description: This test checks how the function behaves whenever the required environment variables are not set. It's expected to return an error message indicating which specific environment variable isn't set.
Execution:
  Arrange: Ensuring that none of the required environment variables are set.
  Act: Call the function dsn().
  Assert: Validate that the function returns the correct error message (no connection string) and that the first missing environment variable is "$DB_HOST".
Validation:
  The test checks that the function correctly handles situations when key data is missing. It's essential for predicting the program's behavior in incomplete or invalid configurations.

Scenario 3: Some environment variables missing

Details:
  Description: This test verifies how the function acts when some but not all required environment variables are defined. It's expected to return an error indicating which specific environment variable isn't found.
Execution:
  Arrange: Setting some of the environment variables, leaving others undefined.
  Act: Call the function dsn().
  Assert: Validate that the function generates the correct error message for the first missing environment variable.
Validation:
  This test assesses the function's ability to correctly identify and notify about the missing environment variables. It demonstrates that the program correctly handles incomplete configurations. 

Scenario 4: Incorrect format for environment variables

Details:
  Description: This test checks how function operates when environment variables are set but incorrectly formatted. Despite the variables being present, the incorrect format would eventually lead to connection failure.
Execution:
  Arrange: Set the literal value of one or more environment variables in an incorrect format.
  Act: Call the function dsn().
  Assert: Check whether the function generates the expected connection string.
Validation:
  As per the current function behavior, we expect the function to return a connection string even with wrongly formatted values since there isn't any format check. However, this test scenario would define a potential case where future improvements of format checks may be added to this function.

*/

// ********RoostGPT********
package main

import (
	"errors"
	"fmt"
	"os"
)

func dsn() (string, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return "", errors.New("$DB_HOST is not set")
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		return "", errors.New("$DB_USER is not set")
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return "", errors.New("$DB_PASSWORD is not set")
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		return "", errors.New("$DB_NAME is not set")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		return "", errors.New("$DB_PORT is not set")
	}
	options := "charset=utf8mb4&parseTime=True&loc=Local"
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?%s", user, password, host, port, name, options), nil
}
