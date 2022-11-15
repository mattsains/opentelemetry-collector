// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlphttpexporter

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

func TestUnmarshalDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NoError(t, component.UnmarshalExporterConfig(confmap.New(), cfg))
	assert.Equal(t, factory.CreateDefaultConfig(), cfg)
	// Default/Empty config is invalid.
	assert.Error(t, cfg.Validate())
}

func TestUnmarshalConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	expectedHTTPConfig := confighttp.NewDefaultHTTPClientSettings()
	expectedHTTPConfig.Headers = map[string]string{
		"can you have a . here?": "F0000000-0000-0000-0000-000000000000",
		"header1":                "234",
		"another":                "somevalue",
	}
	expectedHTTPConfig.Endpoint = "https://1.2.3.4:1234"
	expectedHTTPConfig.TLSSetting = configtls.TLSClientSetting{
		TLSSetting: configtls.TLSSetting{
			CAFile:   "/var/lib/mycert.pem",
			CertFile: "certfile",
			KeyFile:  "keyfile",
		},
		Insecure: true,
	}
	expectedHTTPConfig.ReadBufferSize = 123
	expectedHTTPConfig.WriteBufferSize = 345
	expectedHTTPConfig.Timeout = time.Second * 10
	expectedHTTPConfig.Compression = "gzip"

	assert.NoError(t, component.UnmarshalExporterConfig(cm, cfg))
	assert.Equal(t,
		&Config{
			ExporterSettings: config.NewExporterSettings(component.NewID(typeStr)),
			RetrySettings: exporterhelper.RetrySettings{
				Enabled:         true,
				InitialInterval: 10 * time.Second,
				MaxInterval:     1 * time.Minute,
				MaxElapsedTime:  10 * time.Minute,
			},
			QueueSettings: exporterhelper.QueueSettings{
				Enabled:      true,
				NumConsumers: 2,
				QueueSize:    10,
			},
			HTTPClientSettings: expectedHTTPConfig,
		}, cfg)
}
