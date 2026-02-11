package config

import (
	"fmt"
	"os"
)

const ENV_PREFIX = "AGENT_DEV_ENVIRONMENT_"

func GetValue(key string) string {
	value, ok := os.LookupEnv(ENV_PREFIX + key)
	if !ok {
		panic(fmt.Sprintf("environment variable '%s%s' is required", ENV_PREFIX, key))
	}

	return value
}
