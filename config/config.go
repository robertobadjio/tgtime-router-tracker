package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Config struct {
	RouterHost           string
	RouterPort           string
	RouterUserName       string
	RouterPassword       string
	TgTimeAggregatorHost string
	TgTimeAggregatorPort string
	DelaySeconds         time.Duration
}

const projectDirName = "tgtime-router-tracker"

func init() {
	loadEnv()
}

func New() *Config {
	tgTimeAggregatorHost := getEnv("TGTIME_AGGREGATOR_SERVICE_HOST", "")
	if tgTimeAggregatorHost == "" {
		panic("Param TGTIME_AGGREGATOR_SERVICE_HOST not set")
	}

	tgTimeAggregatorPort := getEnv("TGTIME_AGGREGATOR_SERVICE_PORT", "")
	if tgTimeAggregatorPort == "" {
		panic("Param TGTIME_AGGREGATOR_SERVICE_PORT not set")
	}

	return &Config{
		RouterHost:           getEnv("ROUTER_HOST", ""),
		RouterPort:           getEnv("ROUTER_PORT", ""),
		RouterUserName:       getEnv("ROUTER_USER_NAME", ""),
		RouterPassword:       getEnv("ROUTER_PASSWORD", ""),
		TgTimeAggregatorHost: tgTimeAggregatorHost,
		TgTimeAggregatorPort: tgTimeAggregatorPort,
		DelaySeconds:         getTimeDurationSecondsEnv("DELAY_SECONDS", 30),
	}
}

func loadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Problem loading .env file")
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getTimeDurationSecondsEnv(key string, defaultVal int) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return time.Duration(defaultVal) * time.Second
	}

	seconds, err := strconv.Atoi(value)
	if err != nil {
		return time.Duration(defaultVal) * time.Second
	}

	return time.Duration(seconds) * time.Second
}
