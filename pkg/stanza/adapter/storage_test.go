// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/storagetest"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	r := createReceiver(t)
	host := storagetest.NewStorageHost().
		WithFileBackedStorageExtension("test", t.TempDir())
	require.NoError(t, r.Start(ctx, host))

	myBytes := []byte("my_value")

	require.NoError(t, r.storageClient.Set(ctx, "key", myBytes))
	val, err := r.storageClient.Get(ctx, "key")
	require.NoError(t, err)
	require.Equal(t, myBytes, val)

	// Cycle the receiver
	require.NoError(t, r.Shutdown(ctx))
	for _, e := range host.GetExtensions() {
		require.NoError(t, e.Shutdown(ctx))
	}

	r = createReceiver(t)
	err = r.Start(ctx, host)
	require.NoError(t, err)

	// Value has persisted
	val, err = r.storageClient.Get(ctx, "key")
	require.NoError(t, err)
	require.Equal(t, myBytes, val)

	err = r.storageClient.Delete(ctx, "key")
	require.NoError(t, err)

	// Value is gone
	val, err = r.storageClient.Get(ctx, "key")
	require.NoError(t, err)
	require.Nil(t, val)

	require.NoError(t, r.Shutdown(ctx))

	_, err = r.storageClient.Get(ctx, "key")
	require.Error(t, err)
	require.Equal(t, "client closed", err.Error())
}

func TestFailOnMultipleStorageExtensions(t *testing.T) {
	r := createReceiver(t)
	host := storagetest.NewStorageHost().
		WithInMemoryStorageExtension("one").
		WithInMemoryStorageExtension("two")
	err := r.Start(context.Background(), host)
	require.Error(t, err)
	require.Equal(t, "storage client: multiple storage extensions found", err.Error())
}

func createReceiver(t *testing.T) *receiver {
	params := component.ReceiverCreateSettings{
		TelemetrySettings: componenttest.NewNopTelemetrySettings(),
	}

	factory := NewFactory(TestReceiverType{}, component.StabilityLevelInDevelopment)

	logsReceiver, err := factory.CreateLogsReceiver(
		context.Background(),
		params,
		factory.CreateDefaultConfig(),
		consumertest.NewNop(),
	)
	require.NoError(t, err, "receiver should successfully build")

	r, ok := logsReceiver.(*receiver)
	require.True(t, ok)
	return r
}
