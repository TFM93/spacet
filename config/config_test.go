package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// create temp file
	invalidConfig := []byte("invalid_yaml: -")
	validConfig := []byte(`
app:
  name: 'spacet'
  version: '1.0.0'
  log_level: 'debug'
  `)
	invalidTmpFile, err := os.CreateTemp("", "invalid_config.yaml")
	assert.NoError(t, err)
	defer os.Remove(invalidTmpFile.Name())

	validTmpFile, err := os.CreateTemp("", "config*.yaml")
	assert.NoError(t, err)
	defer os.Remove(validTmpFile.Name())

	_, err = invalidTmpFile.Write(invalidConfig)
	assert.NoError(t, err)
	_, err = validTmpFile.Write(validConfig)
	assert.NoError(t, err)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr error
	}{
		{
			name: "no path provided",
			args: args{
				path: "",
			},
			want:    nil,
			wantErr: fmt.Errorf("config path not provided"),
		}, {
			name: "config error",
			args: args{
				path: invalidTmpFile.Name(),
			},
			want:    nil,
			wantErr: fmt.Errorf("config error: file format '%s' doesn't supported by the parser", filepath.Ext(invalidTmpFile.Name())),
		},
		{
			name: "config success",
			args: args{
				path: validTmpFile.Name(),
			},
			want: &Config{
				App: App{Name: "spacet", Version: "1.0.0", LogLevel: "debug"},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.path)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
