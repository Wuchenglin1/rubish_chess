package dao

import "testing"

func TestInitMysql(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"successful"},
	}
	for range tests {
		InitMysql()
	}
}

func TestInitRedis(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"successful"},
	}
	for range tests {
		InitRedis()
	}
}
