package pkg

import (
	"fmt"
	"time"
)

// MySQLOptions defines options for mysql database.
type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"-"                                  mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
}

func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "", // 置为空的部分，是留给用户自定义的
		Password:              "",
		Database:              "", // 指定要操作的数据库名称
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

func (options *MySQLOptions) String() string {
	return fmt.Sprintf("Host:%s, Username:%s, Password:%s, Database:%s, MaxIdleConnections:%d, MaxOpenConnections:%d, MaxConnectionLifeTime:%v, LogLevel:%d", options.Host, options.Username, options.Password, options.Database, options.MaxIdleConnections, options.MaxOpenConnections, options.MaxConnectionLifeTime, options.LogLevel)
}
