package gowindams

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type EnvironmentConfig struct {
	Name string                `json:"name"                yaml:"name"`
	ServiceURI string          `json:"serviceURI"          yaml:"serviceURI"`
	ClientId string            `json:"clientId"            yaml:"clientId"`
	ClientSecret string        `json:"clientSecret"        yaml:"clientSecret"`
	TenantId string            `json:"tenantId"            yaml:"tenantId"`
	ServiceAppId string        `json:"serviceAppId"        yaml:"serviceAppId"`
	AccessTokenProvider string `json:"accessTokenProvider" yaml:"accessTokenProvider"`
}

type EnvironmentConfigs []EnvironmentConfig

type Environment struct {
	Name string
	ClientId string
	ServiceAppId string
	ServiceURI string
	TenantId string
	accessTokenProvider accessTokenProvider
	assetServiceClient *AssetServiceClient
	assetInspectionServiceClient *AssetInspectionServiceClient
	componentServiceClient *ComponentServiceClient
	componentInspectionServiceClient *ComponentInspectionServiceClient
	inspectionEventResourceServiceClient *InspectionEventResourceServiceClient
	processQueueServiceClient *ProcessQueueServiceClient
	resourceServiceClient *ResourceServiceClient
	siteServiceClient *SiteServiceClient
	workOrderServiceClient *WorkOrderServiceClient
}

func (env Environment) AssetInspectionServiceClient() *AssetInspectionServiceClient {
	return env.assetInspectionServiceClient
}

func (env Environment) AssetServiceServiceClient() *AssetServiceClient {
	return env.assetServiceClient
}

func (env Environment) ComponentInspectionServiceClient() *ComponentInspectionServiceClient {
	return env.componentInspectionServiceClient
}

func (env Environment) ComponentServiceClient() *ComponentServiceClient {
	return env.componentServiceClient
}

func (env Environment) GetAuthenticationProviderType() AuthenticationProviderType {
	if env.accessTokenProvider == nil {
		return AP_Other
	}
	return env.accessTokenProvider.getAuthenticationProviderType()
}

func (env Environment) IsServerToServer() bool {
	return env.accessTokenProvider != nil && env.accessTokenProvider.isServerToServer()
}

func (env Environment) IsUserAuthenticated() bool {
	return env.accessTokenProvider != nil && env.accessTokenProvider.isUserAuthenticated()
}

func (env Environment) ObtainAccessToken() (string, error) {
	if env.accessTokenProvider == nil {
		// No provider
		return "", fmt.Errorf("No access token provider available for the environment %s", env.Name)
	} else {
		token, err := obtainAccessToken(env.accessTokenProvider, env.ServiceAppId)
		return token, err
	}
}

func (env Environment) ObtainSigningKeys() (map[string]interface{}, error) {
	if env.accessTokenProvider == nil {
		keys := make(map[string]interface{})
		return keys, nil
	} else {
		keys, err := obtainSigningKeys(env.accessTokenProvider)
		return keys, err
	}
}

func (env Environment) InspectionEventResourceServiceClient() *InspectionEventResourceServiceClient {
	return env.inspectionEventResourceServiceClient
}

func (env Environment) ProcessQueueServiceClient() *ProcessQueueServiceClient {
	return env.processQueueServiceClient
}

func (env Environment) ResourceServiceClient() *ResourceServiceClient {
	return env.resourceServiceClient
}

func (env Environment) SiteServiceClient() *SiteServiceClient {
	return env.siteServiceClient
}

func (env Environment) WorkOrderServiceClient() *WorkOrderServiceClient {
	return env.workOrderServiceClient
}

type Environments []Environment

const DEFAULT_CONFIG_PATH = "/etc/windams/environments.yaml"

func (envs *Environments) Find(name string) *Environment {
	for _, env := range *envs {
		if env.Name == name {
			return &env
		}
	}
	return nil
}

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
	count := len(*configs)
	log.Printf("Loaded configurations for %d environments", count)
	environments := make(Environments, count)
	i := 0
	for _, cfg := range *configs {
		env := Environment{
			Name:                cfg.Name,
			ClientId:            cfg.ClientId,
			ServiceAppId:        cfg.ServiceAppId,
			ServiceURI:          strings.TrimRight(cfg.ServiceURI, "/"),
			TenantId:            cfg.TenantId,
			accessTokenProvider: NewProvider(&cfg),
		}
		env.assetInspectionServiceClient = &AssetInspectionServiceClient{
			env: &env,
		}
		env.assetServiceClient = &AssetServiceClient{
			env: &env,
		}
		env.componentInspectionServiceClient = &ComponentInspectionServiceClient{
			env: &env,
		}
		env.componentServiceClient = &ComponentServiceClient{
			env: &env,
		}
		env.inspectionEventResourceServiceClient = &InspectionEventResourceServiceClient{
			env: &env,
		}
		env.processQueueServiceClient = &ProcessQueueServiceClient{
			env: &env,
		}
		env.resourceServiceClient = &ResourceServiceClient{
			env: &env,
		}
		env.siteServiceClient = &SiteServiceClient{
			env: &env,
		}
		env.workOrderServiceClient = &WorkOrderServiceClient{
			env: &env,
		}
		environments[i] = env
		log.Printf("Configured environment %d: %s\t%s", i, env.Name, env.ServiceURI)
		i++
	}
	return &environments, nil
}
