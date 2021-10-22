/**
 * @Author: Lee
 * @Description:
 * @File:  mysql
 * @Version: 1.0.0
 * @Date: 2021/10/19 10:51 下午
 */

package databases

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type Mysql struct {
	host      string
	port      int
	user      string
	password  string
	dbName    string
	maxIdle   int
	maxOpen   int
	sqlClient *gorm.DB
}

func NewMysql(host string, port int, user string, password string, dbName string, maxIdle int, maxOpen int, isLogger bool) (
	sql *Mysql, err error) {
	sql = &Mysql{
		host: host, port: port, user: user, password: password, dbName: dbName, maxIdle: maxIdle, maxOpen: maxOpen,
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		user, password, host, port, dbName)
	config := &gorm.Config{}
	if isLogger {
		config.Logger = logger.Default.LogMode(logger.Info)
	}
	sql.sqlClient, err = gorm.Open(mysql.Open(url), config)
	if err != nil {
		panic(err)
	}
	sqlDB, err := sql.sqlClient.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(maxIdle)              //最大空闲连接数
	sqlDB.SetMaxOpenConns(maxOpen)              //最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * 200) //设置连接空闲超时
	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		panic(err)
	}
	return
}

func (sql *Mysql) AutoMigrate(dst ...interface{}) error {
	return sql.sqlClient.AutoMigrate(dst...)
}

// InitRecord 初始化数据库数据
func (sql *Mysql) InitRecord(record map[schema.Tabler]string) (err error) {
	var (
		count int64
	)
	// 遍历数据库数据初始化字典，将所有需要初始化的表数据初始化
	for model, sqlFile := range record {
		err = sql.sqlClient.Model(model).Limit(1).Count(&count).Error
		if err != nil {
			log.Fatal(err)
			return
		}
		if count == 0 {
			dbExecSQLFile(sql.sqlClient, sqlFile)
		}
	}
	return
}

// dbExecSQLFile 执行sql文件，文件内容只允许写SQL语句与换行字符，不允许注释
func dbExecSQLFile(db *gorm.DB, filePath string) {
	path, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatal(err)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(string(file), "\n") {
		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, "\n", "")
		if len(line) > 0 {
			err = db.Exec(line).Error
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
