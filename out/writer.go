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
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/deepauto-io/nsqcc"

	"github.com/deepauto-io/nsqcc/filepath/ifs"
	"github.com/nsqio/go-nsq"
)

type nsqWriter struct {
	conf     Config
	tlsConf  *tls.Config
	connMut  sync.RWMutex
	producer *nsq.Producer
}

func NewNSQWriter(conf Config, mgr ifs.FS) (nsqcc.AsyncSink, error) {
	n := nsqWriter{
		conf: conf,
	}

	if conf.TLS.Enabled {
		var err error
		if n.tlsConf, err = conf.TLS.Get(mgr); err != nil {
			return nil, err
		}
	}
	return &n, nil
}

func (n *nsqWriter) Connect(ctx context.Context) error {
	n.connMut.Lock()
	defer n.connMut.Unlock()

	cfg := nsq.NewConfig()
	cfg.UserAgent = n.conf.UserAgent
	cfg.MaxInFlight = n.conf.MaxInFlight
	if n.tlsConf != nil {
		cfg.TlsV1 = true
		cfg.TlsConfig = n.tlsConf
	}

	producer, err := nsq.NewProducer(n.conf.Address, cfg)
	if err != nil {
		return err
	}

	producer.SetLogger(log.New(io.Discard, "", log.Flags()), nsq.LogLevelError)

	if err := producer.Ping(); err != nil {
		return err
	}
	n.producer = producer
	return nil
}

func (n *nsqWriter) WriteWithContext(ctx context.Context, topic string, msg []byte) error {
	n.connMut.RLock()
	prod := n.producer
	n.connMut.RUnlock()

	if prod == nil {
		return nsqcc.ErrNotConnected
	}

	if govalidator.IsNull(topic) {
		return fmt.Errorf("topic is required")
	}

	if len(msg) == 0 {
		return nil
	}
	return prod.Publish(topic, msg)
}

func (n *nsqWriter) Close(ctx context.Context) error {
	go func() {
		n.connMut.Lock()
		if n.producer != nil {
			n.producer.Stop()
			n.producer = nil
		}
		n.connMut.Unlock()
	}()
	return nil
}
