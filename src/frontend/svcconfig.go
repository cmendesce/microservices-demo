package main

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func GetEnvInt(key string) (int, bool) {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value, true
	}
	return 0, false
}

func BuildServiceConfig() string {
	maxAttempts, maxAttemptsOk := GetEnvInt("GRPC_MAX_ATTEMPTS")
	initialBackoff := GetEnv("GRPC_INITIAL_BACKOFF", "")
	maxBackoff := GetEnv("GRPC_MAX_BACKOFF", "")
	backoffMultiplier, backoffMultiplierOk := GetEnvInt("GRPC_BACKOFF_MULTIPLIER")

	hedgingMaxAttempts := GetEnv("GRPC_HEDGING_MAX_ATTEMPTS", "")
	hedgingDelay := GetEnv("GRPC_HEDGING_DELAY", "")

	if maxAttemptsOk && backoffMultiplierOk && initialBackoff != "" && maxBackoff != "" {
		return fmt.Sprintf(`{
			"methodConfig": [
				{
					"name": [
						{ "service": "", "method": "" }
					],
					"retryPolicy": {
						"maxAttempts": %d,
						"initialBackoff": "%s",
						"maxBackoff": "%s",
						"backoffMultiplier": %d,
						"retryableStatusCodes": [
							"UNAVAILABLE"
						]
					},
					"waitForReady": true
				}
			]
		}`, maxAttempts, initialBackoff, maxBackoff, backoffMultiplier)
	} else if hedgingMaxAttempts != "" && hedgingDelay != "" {
		return fmt.Sprintf(`{
			"methodConfig": [
				{
					"name": [
						{ "service": "", "method": "" }
					],
					"hedgingPolicy": {
						"maxAttempts": %d,
						"hedgingDelay": "%s",
					},
					"waitForReady": true
				}
			]
		}`, hedgingMaxAttempts, hedgingDelay)
	} else {
		return ""
	}
}
