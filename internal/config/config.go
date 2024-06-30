package config

import (
	authzuc "github.com/azarc-io/verathread-next-common/usecase/authz"
	dapruc "github.com/azarc-io/verathread-next-common/usecase/dapr"
	devuc "github.com/azarc-io/verathread-next-common/usecase/dev"
	luc "github.com/azarc-io/verathread-next-common/usecase/logging"
	mongouc "github.com/azarc-io/verathread-next-common/usecase/mongo"
	natsuc "github.com/azarc-io/verathread-next-common/usecase/nats"
)

type (
	Config struct {
		Name        string               `yaml:"name"`
		ID          string               `yaml:"id"`
		DataDir     string               `yaml:"data_dir"`
		Gateway     *APIGatewayConfig    `yaml:"gateway"`
		Logger      *Logging             `yaml:"logger"`
		Nats        *natsuc.NatsConfig   `yaml:"nats"`
		Development *devuc.Config        `yaml:"development"`
		Auth        *authzuc.Config      `yaml:"auth"`
		Dapr        *dapruc.Config       `yaml:"dapr"`
		Database    *mongouc.MongoConfig `yaml:"database"`
	}

	Logging struct {
		Level  string          `yaml:"level"`
		Mode   luc.LoggingMode `yaml:"mode"`
		Caller bool            `yaml:"caller"`
	}
)
