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

package nsqcc

import (
	"context"

	"github.com/nsqio/go-nsq"
)

type Service interface {
	// Connect attempts to establish a connection to the source, if
	// unsuccessful returns an error. If the attempt is successful (or not
	// necessary) returns nil.
	Connect(ctx context.Context) error

	// Close triggers the shut-down of this component and blocks until
	// completion or context cancellation.
	Close(ctx context.Context) error
}

// Async is a type that reads Benthos messages from an external source and
// allows acknowledgements for a message batch to be propagated asynchronously.
type Async interface {
	Service
	// ReadBatch attempts to read a new message from the source. If
	// successful a message is returned along with a function used to
	// acknowledge receipt of the returned message. It's safe to process the
	// returned message and read the next message asynchronously.
	ReadBatch(ctx context.Context) (*nsq.Message, AsyncAckFn, error)
}

// AsyncAckFn is a function used to acknowledge receipt of a message batch. The
// provided response indicates whether the message batch was successfully
// delivered. Returns an error if the acknowledgment was not propagated.
type AsyncAckFn func(context.Context, error) error

// noopAsyncAckFn is a no-op acknowledgment function.
var noopAsyncAckFn AsyncAckFn = func(context.Context, error) error {
	return nil
}

type AsyncSink interface {
	Service
	// WriteWithContext should block until either the message is sent (and
	// acknowledged) to a sink, or a transport specific error has occurred, or
	// the Type is closed.
	WriteWithContext(ctx context.Context, topic string, msg []byte) error
}
