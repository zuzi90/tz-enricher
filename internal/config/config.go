package config

import "github.com/caarlos0/env/v10"

type Config struct {
	ServerPORT        string   `env:"SERVER_PORT"      envDefault:":5005"`
	PgDSN             string   `env:"PG_DSN"           envDefault:"postgresql://postgres:secret@localhost:5431/tzDB?sslmode=disable"`
	RedisDSN          string   `env:"REDIS_DSN"        envDefault:"localhost:6379"`
	RedisPassword     string   `env:"REDIS_PASSWORD"   envDefault:""`
	RedisDB           int      `env:"REDIS_DB"         envDefault:"0"`
	LogLvl            string   `env:"LOG_LEVEL"        envDefault:"debug"`
	AgeURL            string   `env:"AGE_URL"          envDefault:"https://api.agify.io/?name="`
	GenderURL         string   `env:"GENDER_URL"       envDefault:"https://api.genderize.io/?name="`
	NationalityURL    string   `env:"NATIONALITY_URL"  envDefault:"https://api.nationalize.io/?name="`
	Brokers           []string `env:"BROKERS"          envDefault:"localhost:9092"`
	WorkersCount      int      `env:"WORKERS_COUNT"    envDefault:"1"`
	KafkaTopic        string   `env:"KAFKA_TOPIC"      envDefault:"FN"`
	KafkaTopicWrongFN string   `env:"KAFKA_TOPIC_WRONG_FN"      envDefault:"WRONG_FN"`
}

func NewConfig() (*Config, error) {
	cnf := Config{}

	if err := env.Parse(&cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
