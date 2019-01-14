package service

import "os"

var (
	Version string = getEnv("VERSION", "v0.1")
        KakaPeers string = os.Getenv("KAFKA_PEERS")
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
