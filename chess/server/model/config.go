package model

type Cfg struct {
	Mysql     MysqlCfg `json:"mysql"`
	Redis     Redis    `json:"redis"`
	Signature string   `json:"signature"`
}

type MysqlCfg struct {
	Dsn string `json:"dsn"`
}

type Redis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
}
