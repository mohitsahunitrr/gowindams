package gowindams

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type EnvironmentConfig struct {
	Name string            `json:"name" yaml:"name"`
	ServiceURI string      `json:"serviceURI" yaml:"serviceURI"`
	ClientId string        `json:"clientId" yaml:"clientId"`
	ClientSecret string    `json:"clientSecret" yaml:"clientSecret"`
	TenantId string        `json:"tenantId" yaml:"tenantId"`
	ServiceAppId string    `json:"serviceAppId" yaml:"serviceAppId"`
}

type EnvironmentConfigs []EnvironmentConfig

type Environment struct {
	Name string
	ServiceAppId string
	ServiceURI string
	accessTokenProvider *accessTokenProvider
	resourceServiceClient *ResourceServiceClient
	processQueueServiceClient *ProcessQueueServiceClient
}

func (env Environment) obtainAccessToken() (string, error) {
	return env.accessTokenProvider.obtainAccessToken(env.ServiceAppId)
}

func (env Environment) ResourceServiceClient() *ResourceServiceClient {
	return env.resourceServiceClient
}

func (env Environment) ProcessQueueServiceClient() *ProcessQueueServiceClient {
	return env.processQueueServiceClient
}

type Environments []Environment

const DEFAULT_CONFIG_PATH = "/etc/windams/environments.yaml"

func LoadEnvironments(configFilePath string) (*Environments, error) {
	if "" == configFilePath {
		configFilePath = DEFAULT_CONFIG_PATH
	}

	// Read configuration from config path.
	body, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	configs := new(EnvironmentConfigs)
	err = yaml.Unmarshal(body, &configs)
	if err != nil {
		return nil, err
	}
	environments := make(Environments, len(*configs))
	i := 0
	for _, cfg := range *configs {
		env := Environment{
			Name:                cfg.Name,
			ServiceAppId:        cfg.ServiceAppId,
			ServiceURI:          cfg.ServiceURI,
			accessTokenProvider: NewProvider(&cfg),
		}
		env.processQueueServiceClient = &ProcessQueueServiceClient{
			env: &env,
		}
		env.resourceServiceClient = &ResourceServiceClient{
			env: &env,
		}
		environments[i] = env
		i++
	}
	return &environments, nil
}
