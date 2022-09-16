package orm

import (
	"fmt"

	"github.com/kos-v/dsnparser"
)

//统一dsn解析

const (
	mysql_tpl    = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	postgres_tpl = "host=%s port=%s dbname=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai"
)

func ParseDsn(sqlType, originDsn string) string {
	before := dsnparser.Parse(originDsn)

	if sqlType == "postgres" {
		return fmt.Sprintf(postgres_tpl, before.GetHost(), before.GetPort(), before.GetPath(), before.GetUser(), before.GetPassword())
	} else {
		return fmt.Sprintf(mysql_tpl, before.GetUser(), before.GetPassword(), before.GetHost(), before.GetPort(), before.GetPath())
	}
}
