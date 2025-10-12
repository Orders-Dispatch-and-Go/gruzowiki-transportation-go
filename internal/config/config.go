package config

import (
	"github.com/go-ozzo/ozzo-validation"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const (
	defaultServerPort = 8080
)

type Config struct {
	ServerPort int `yaml:"server_port" env:"TRANSPORTATION_SERVICE_PORT"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(
		&c,
		//validation.Field(&c.DSN, validation.Required),
		//validation.Field(&c.JWTSigningKey, validation.Required),
	)
}

func Load(file string, logger *log.Logger) (*Config, error) {
	c := Config{
		ServerPort: defaultServerPort,
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	//if err = env.New("APP_", logger.Infof).Load(&c); err != nil {
	//	return nil, err
	//}

	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}
