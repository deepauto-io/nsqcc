/*
Copyright 2024 The nsqcc Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package out

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	ntls "github.com/deepauto-io/nsqcc/tls"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Config represents the configuration for the nsqcc command.
type Config struct {
	Address     string `json:"address" envconfig:"NSQ_WRITER_ADDRESS"                     default:"127.0.0.1:4150"`           // NSQ 地址
	UserAgent   string `json:"user_agent" envconfig:"NSQ_WRITER_USER_AGENT"                  default:"DeepAuto Producer/1.0"` // 连接时使用的用户UA
	MaxInFlight int    `json:"max_in_flight" envconfig:"NSQ_WRITER_MAX_IN_FLIGHT"               default:"64"`                 // 同时处理的最大消息数量
	TLS         ntls.Config
}

// NewConfig creates a new Config with default values.
func NewConfig() Config {
	return Config{
		Address:     "127.0.0.1:4150",
		UserAgent:   "DeepAuto Producer/1.0",
		MaxInFlight: 64,
		TLS:         ntls.NewConfig(),
	}
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if govalidator.IsNull(c.Address) {
		return fmt.Errorf("nsq address is required")
	}
	return nil
}

func MustLoadConfig(cfgPath string) Config {
	cfg := Config{}
	if cfgPath == "" {
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatal(err)
		}
		return cfg
	}
	return cfg
}
