package tool

import (
	"chess/server/model"
	"reflect"
	"testing"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "toolPackage"},
	}
	for range tests {
		InitConfig()
	}
}

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		want *model.Cfg
	}{
		{name: "configPackage", want: cfg},
	}
	for _, tt := range tests {
		if got := GetConfig(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetConfig() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
