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

// Config is the configuration for the reader.
type Config struct {
	Addresses       []string `envconfig:"NSQ_ADDRESSES"                   default:"127.0.0.1:4150"` // a list of nsqd addresses to connect to
	LookupAddresses []string `envconfig:"NSQ_LOOKUP_ADDRESSES"            default:"127.0.0.1:4161"` // a list of nsqlookupd addresses to connect to
	Topic           string   `envconfig:"NSQ_TOPIC"                       default:"chatc"`          // topic to consume from
	Channel         string   `envconfig:"NSQ_CHANNEL"                     default:"chatc"`          // channel to consume from
	UserAgent       string   `envconfig:"NSQ_USER_AGENT"`                                           // a user agent to assume when connecting.
	MaxInFlight     int      `envconfig:"NSQ_MAX_IN_FLIGHT"               default:"100"`            // maximum number of pending messages to consume at any given time.
	MaxAttempts     uint16   `envconfig:"NSQ_MAX_ATTEMPTS"                default:"5"`              // maximum number of attempts to successfully consume a messages
	TLS             btls.Config
}
