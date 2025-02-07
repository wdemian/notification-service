package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/thatva/notification-service/ratelimiter"
)

func LoadRateLimiterConfig(configFile string) (ratelimiter.RateLimiterConfig, error) {
	defaults := map[string]ratelimiter.LimitConfig{
		"news":   {Limit: 1, Period: 24 * time.Hour},
		"update": {Limit: 3, Period: time.Minute},
	}
	config := ratelimiter.RateLimiterConfig{
		Limits: defaults,
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return config, fmt.Errorf("failed to load config: %w", err)
		}
		if err := viper.UnmarshalKey("limits", &config.Limits); err != nil {
			return config, fmt.Errorf("failed to unmarshal limits config: %w", err)
		}
	}

	return config, nil
}
