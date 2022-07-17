package service

import (
	"chess/server/model"
	"reflect"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	u := model.User{
		UserName: "clinyu",
		Password: "zxc123",
		RoomNum:  0,
	}
	u1 := model.User{}

	type args struct {
		u         model.User
		tokenType string
		duration  time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "right case", args: args{u: u, tokenType: "access_token", duration: time.Hour}, want: u.UserName},
	}

	for _, tt := range tests {
		got, err := CreateToken(tt.args.u, tt.args.tokenType, tt.args.duration)
		u1, err = ParseToken(got)
		got = u1.UserName
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CreateToken() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. CreateToken() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestParseToken(t *testing.T) {
	var str string
	u := model.User{
		UserName: "cLinYu",
		Password: "zxc123",
		RoomNum:  0,
	}
	str, _ = CreateToken(u, "refresh_token", time.Hour)
	type args struct {
		tokenStr string
	}
	tests := []struct {
		name    string
		args    args
		want    model.User
		wantErr bool
	}{
		{"right case", args{tokenStr: str}, u, false},
	}
	for _, tt := range tests {
		got, err := ParseToken(tt.args.tokenStr)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ParseToken() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ParseToken() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
