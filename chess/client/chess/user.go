package chess

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var Uid uint

func Login(userName, password string) error {
	_url := "http://110.42.184.72:8084/user/login"
	client := &http.Client{}
	values := url.Values{}
	values.Set("userName", userName)
	values.Add("password", password)

	req, err := http.NewRequest("POST", _url, strings.NewReader(values.Encode()))
	if err != nil {
		fmt.Printf("create request error : %v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("do request error %v", err)
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read error : %v", err)
		return err
	}

	var info ResponseLoginMessage
	err = json.Unmarshal(data, &info)
	if err != nil {
		fmt.Printf("unmarshal resp data error : %v", err)
		return err
	}
	token = info.Data.AccessToken
	Uid = uint(info.Data.Uid)
	return nil
}

func Register(userName, password string) error {
	_url := "http://110.42.184.72:8084/user/register"
	client := &http.Client{}
	values := url.Values{}
	values.Set("userName", userName)
	values.Add("password", password)
	req, err := http.NewRequest("POST", _url, strings.NewReader(values.Encode()))
	if err != nil {
		fmt.Printf("create request error : %v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("do request error %v", err)
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read error : %v", err)
		return err
	}
	var respMsg ResponseRegisterMessage
	err = json.Unmarshal(data, &respMsg)
	if err != nil {
		fmt.Printf("unmarshal error : %v", err)
	}
	if respMsg.Status == 400 {
		nErr := errors.New("account wrong")
		return nErr
	}
	return err
}
