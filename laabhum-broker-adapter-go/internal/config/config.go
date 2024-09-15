package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	BrokerConfig         BrokerConfig
	KafkaConfig          KafkaConfig
	RedisConfig          RedisConfig
	CircuitBreakerConfig CircuitBreakerConfig
	ServerConfig         ServerConfig // Added ServerConfig
	Kafka                struct {
		Brokers []string `json:"brokers"`
	} `json:"kafka"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

type BrokerConfig struct {
	WebSocketURL string
	APIBaseURL   string
	APIKey       string
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type CircuitBreakerConfig struct {
	MaxFailures int
	Timeout     int
}

type ServerConfig struct {
	Port string
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg *RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

// LoadConfig loads the configuration from a JSON file.
func LoadConfig() (*Config, error) {
	file, err := os.Open("internal/config/config.json") // Adjust to your config file path
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
