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

package in

import (
	"fmt"

	ntls "github.com/deepauto-io/nsqcc/tls"
)

// Config is the configuration for the reader.
type Config struct {
	Addresses       []string `envconfig:"NSQ_ADDRESSES"                   default:"127.0.0.1:4150"`   // Nsqd 地址列表
	LookupAddresses []string `envconfig:"NSQ_LOOKUP_ADDRESSES"            default:"127.0.0.1:4161"`   // NSQLookupd 地址列表
	Topic           string   `envconfig:"NSQ_TOPIC"`                                                  // 消费的主题名
	Channel         string   `envconfig:"NSQ_CHANNEL"`                                                // 消费的频道名
	UserAgent       string   `envconfig:"NSQ_USER_AGENT"                  default:"DeepAuto NSQ/1.0"` // 连接时使用的用户UA
	MaxInFlight     int      `envconfig:"NSQ_MAX_IN_FLIGHT"               default:"64"`               // 同时处理的最大消息数量.
	MaxAttempts     uint16   `envconfig:"NSQ_MAX_ATTEMPTS"                default:"3"`                // 消息最大重试次数
	TLS             ntls.Config
}

// NewConfig creates a new Config with default values.
func NewConfig() Config {
	return Config{
		Addresses:       []string{"127.0.0.1:4150"},
		LookupAddresses: []string{"127.0.0.1:4161"},
		UserAgent:       "DeepAuto NSQ/1.0",
		MaxInFlight:     64,
		MaxAttempts:     5,
	}
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if len(c.Addresses) == 0 {
		return fmt.Errorf("nsq address is required")
	}

	if len(c.LookupAddresses) == 0 {
		return fmt.Errorf("nsq lookupd addresses is required")
	}

	return nil
}
