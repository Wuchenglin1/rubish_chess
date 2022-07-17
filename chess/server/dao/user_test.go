package dao

import (
	"chess/server/model"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

var mock sqlmock.Sqlmock
var gormDB *gorm.DB

func TestSearchUserByName(t *testing.T) {
	//mysql模拟
	var err error
	var db *sql.DB
	db, mock, err = sqlmock.New()

	if err != nil {
		log.Fatalln("into sqlmock(mysql) db err:", err)
	}
	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("init DB with sqlmock(gorm) fail err:", err)
	}
	Db = gormDB
	//redis 模拟
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.Set("testUser1", "0;zxc123")

	Rdb = redis.NewClient(&redis.Options{
		Addr: s.Addr(), // mock redis server的地址
	})
	//上面一堆配置
	//--------------------------------------------

	u := model.User{
		Model:    gorm.Model{},
		UserName: "testUser1",
		Password: "zxc123",
		RoomNum:  0,
	}
	u1 := model.User{
		Model:    gorm.Model{},
		UserName: "testUser2",
		Password: "zxc123",
		RoomNum:  0,
	}

	row := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_name", "password", "room_num"}).AddRow(0, time.Now(), time.Now().Add(time.Minute*5), nil, "testUser2", "zxc123", 0)
	fmt.Println(row)
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE user_name \\= \\? AND `users`\\.`deleted_at` IS NULL").WithArgs(u1.UserName).WillReturnRows(row)
	type args struct {
		u *model.User
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"redis success case", args{u: &u}, true, false},
		{"mysql success case", args{u: &u1}, true, false},
	}
	for _, tt := range tests {
		got, err := SearchUserByName(tt.args.u)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. SearchUserByName() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. SearchUserByName() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCreateUser(t *testing.T) {
	var err error
	var db *sql.DB
	db, mock, err = sqlmock.New()

	if err != nil {
		log.Fatalln("into sqlmock(mysql) db err:", err)
	}
	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("init DB with sqlmock(gorm) fail err:", err)
	}
	Db = gormDB
	//redis 模拟
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	s.Set("testUser1", "0;zxc123")

	Rdb = redis.NewClient(&redis.Options{
		Addr: s.Addr(), // mock redis server的地址
	})
	//上面都是配置
	//---------------------------------------------------
	u1 := model.User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName: "testUser1",
		Password: "zxc123456",
		RoomNum:  0,
	}
	u2 := model.User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserName: "testUser2",
		Password: "zxc123456",
		RoomNum:  0,
	}

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		u model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success case", args: args{u: u1}, wantErr: false},
		{name: "bad case", args: args{u: u2}, wantErr: false},
	}
	for _, tt := range tests {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users` \\(`created_at`,`updated_at`,`deleted_at`,`user_name`,`password`,`room_num`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").WithArgs(tt.args.u.CreatedAt, tt.args.u.UpdatedAt, tt.args.u.DeletedAt, tt.args.u.UserName, tt.args.u.Password, tt.args.u.RoomNum)
		mock.ExpectCommit()

		if err := CreateUser(tt.args.u); (err != nil) != tt.wantErr {
			t.Errorf("%q. CreateUser() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
