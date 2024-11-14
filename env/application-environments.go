package env

import "os"

func SetEnv(file string) {

	// By default, environment variables will be set for LocalStack.
	os.Setenv("LOCALSTACK_ENDPOINT", "http://118.69.182.236:4566/")
	os.Setenv("LOCALSTACK_DEFAULT_REGION", "us-east-1")
	os.Setenv("LOCALSTACK_ACCESS_KEY_ID", "AKID")
	os.Setenv("LOCALSTACK_SECRET_ACCESS_KEY", "SECRET")

	if len(file) == 0 {
		return
	}

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return // File does not exist (case default)
	}
	if err != nil {
		panic(err) // Error reading the file
	}

	// Read the file and set the environment variables
	// ...
}
