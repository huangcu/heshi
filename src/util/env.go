package util

import "os"

func IsProductionStage() bool {
	return os.Getenv("STAGE") == "production"
}

func IsDevelopmentStage() bool {
	return os.Getenv("STAGE") == "development"
}

func ShouldTrace() bool {
	return os.Getenv("TRACE") == "true"
}
