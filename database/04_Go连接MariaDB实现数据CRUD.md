在 Macbook 安装 MariaDB，使用 Go 连接数据库，并实现 CRUD。大致**实现框架和逻辑**：

1. 准备工作：在新安装的 MariaDB 中创建用户 iam（与之对应的密码是：`iam59!z$`），并使用账户创建名为 iam 的数据库；
2. 导入数据库操作的第三方框架，使用 GORM 框架操作数据库；
3. 写 DSN，连接数据库；
4. 创建 User 表，并插入数据，实现 CRUD 操作。

**实现记录**如下：

# 1 MariaDB 准备

~~~bash
$ mysql -h127.0.0.1 -P3306 -uroot -p'iam59!z$'
MariaDB [(none)]> grant all on iam.* TO iam@127.0.0.1 identified by 'iam59!z$';
Query OK, 0 rows affected (0.000 sec)
MariaDB [(none)]> flush privileges;
Query OK, 0 rows affected (0.000 sec)
~~~

连接 MariaDB，-h 指定主机，-P 指定监听端口，-u 指定登录用户，-p 指定登录密码。使用 ant 账户登录有，同时创建了名为 iam 的账户，并为其赋予所有的权限，指定其登录密码是：`iam59!z$`。**反向验证**：

~~~bash
MariaDB [(none)]> select User from  mysql.user;
+-------------+
| User        |
+-------------+
| iam         |
| ant         |
| mariadb.sys |
| root        |
+-------------+
4 rows in set (0.002 sec)
~~~

使用 iam 账户登录，并执行 source.sql 创建数据库，并创建对应的表：

~~~sql
-- Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file.

CREATE DATABASE  IF NOT EXISTS `iam` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `iam`;
-- MySQL dump 10.13  Distrib 5.7.29, for Win64 (x86_64)
--
-- Host: 106.52.30.200    Database: iam
-- ------------------------------------------------------
-- Server version	5.5.5-10.5.9-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `policy`
--

DROP TABLE IF EXISTS `policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `instanceID` varchar(20) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `username` varchar(255) NOT NULL,
  `policyShadow` longtext DEFAULT NULL,
  `extendShadow` longtext DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name_username` (`name`,`username`),
  UNIQUE KEY `instanceID_UNIQUE` (`instanceID`),
  KEY `fk_policy_user_idx` (`username`),
  CONSTRAINT `fk_policy_user` FOREIGN KEY (`username`) REFERENCES `user` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `policy`
--

LOCK TABLES `policy` WRITE;
/*!40000 ALTER TABLE `policy` DISABLE KEYS */;
/*!40000 ALTER TABLE `policy` ENABLE KEYS */;
UNLOCK TABLES;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`iam`@`127.0.0.1`*/ /*!50003 TRIGGER `iam`.`policy_BEFORE_DELETE` BEFORE DELETE ON `policy` FOR EACH ROW
BEGIN
	insert into policy_audit values(old.id, old.instanceID, old.name, old.username, old.policyShadow, old.extendShadow, old.createdAt, old.updatedAt, curtime());
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `policy_audit`
--

DROP TABLE IF EXISTS `policy_audit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy_audit` (
  `id` bigint(20) unsigned NOT NULL,
  `instanceID` varchar(20) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `username` varchar(255) NOT NULL,
  `policyShadow` longtext DEFAULT NULL,
  `extendShadow` longtext DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deletedAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  KEY `fk_policy_user_idx` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `policy_audit`
--

LOCK TABLES `policy_audit` WRITE;
/*!40000 ALTER TABLE `policy_audit` DISABLE KEYS */;
/*!40000 ALTER TABLE `policy_audit` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `secret`
--

DROP TABLE IF EXISTS `secret`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `secret` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `instanceID` varchar(20) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `username` varchar(255) NOT NULL,
  `secretID` varchar(36) NOT NULL,
  `secretKey` varchar(255) NOT NULL,
  `expires` int(64) unsigned NOT NULL DEFAULT 1534308590,
  `description` varchar(255) NOT NULL,
  `extendShadow` longtext DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name_username` (`name`,`username`),
  UNIQUE KEY `instanceID_UNIQUE` (`instanceID`),
  KEY `fk_secret_user_idx` (`username`),
  CONSTRAINT `fk_secret_user` FOREIGN KEY (`username`) REFERENCES `user` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `secret`
--

LOCK TABLES `secret` WRITE;
/*!40000 ALTER TABLE `secret` DISABLE KEYS */;
/*!40000 ALTER TABLE `secret` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `instanceID` varchar(20) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `nickname` varchar(30) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(256) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `isAdmin` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '1: administrator\\\\n0: non-administrator',
  `extendShadow` longtext DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  UNIQUE KEY `instanceID_UNIQUE` (`instanceID`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (0,'user-lingfei','admin','admin','$2a$10$WnQD2DCfWVhlGmkQ8pdLkesIGPf9KJB7N1mhSOqulbgN7ZMo44Mv2','admin@foxmail.com','1812884xxxx',1,'{}','2021-05-27 10:01:40','2021-05-05 21:13:14');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'PIPES_AS_CONCAT,ANSI_QUOTES,ONLY_FULL_GROUP_BY,NO_AUTO_VALUE_ON_ZERO,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER="iam"@"127.0.0.1"*/ /*!50003 TRIGGER `iam`.`user_BEFORE_DELETE` BEFORE DELETE ON `user` FOR EACH ROW
BEGIN
	delete from secret where username = old.name;
    delete from policy where username = old.name;
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Dumping events for database 'iam'
--

--
-- Dumping routines for database 'iam'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-05-27 18:07:55
~~~

> 从上述 .sql 脚本文件可看出：支持 # 到该行结束、-- 到该行结束 以及 /* 行中间或多个行 */ 的**注释方式**。

对应的**执行**：

~~~go
ant@MacBook-Pro tmp % mysql -h127.0.0.1 -P3306 -uiam -p'iam59!z$'
Welcome to the MariaDB monitor.  Commands end with ; or \g.
Your MariaDB connection id is 37
Server version: 10.6.4-MariaDB Homebrew

Copyright (c) 2000, 2018, Oracle, MariaDB Corporation Ab and others.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

MariaDB [(none)]> source ./iam.sql;
Query OK, 1 row affected (0.002 sec)
...
~~~

至此，整个项目的数据库准备工作就完成了。

此前在部署其他项目时，遇到了数据库操作失败的问题，而且也考虑到了需要前期需要创建数据库、创建表等。在 iam-apiserver 这个项目中，从部署流程来看，确实是这样的。

# 2 数据库连接配置

此部分相当于是解析 yaml 配置文件，并 Unmarshal 到结构对象，此处省略直接做配置：

定义数据通用配置结构体：

~~~go
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
~~~

在 main 中设置：

~~~go
package main

import (
	"context"
	"log"

	"github.com/go-examples-with-tests/database/v4/pkg"
	"github.com/go-examples-with-tests/database/v4/store"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	log.Println("hello, world!")

	options := pkg.NewMySQLOptions()
	initOptions(options)
	log.Printf("MariaDB options:%s", options)

	...
}

func initOptions(options *pkg.MySQLOptions) {
	if options == nil {
		return
	}
	options.Username = "iam"
	options.Password = "iam59!z$"
	options.Database = "iam"
	options.LogLevel = 4
}
~~~

# 3 连接数据库

连接数据的配置项设置好后，就可以创建数据库实例：

~~~go
package main

import (
	"context"
	"log"

	"github.com/go-examples-with-tests/database/v4/pkg"
	"github.com/go-examples-with-tests/database/v4/store"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

...

func main() {
	...
	// connection to MariaDB
	dbFactory, err := store.GetMySQLFactoryOr(options)
	if err != nil {
		log.Fatal(err)
	}
	...
}
~~~

整个创建 DB 的动作都在此处封装：

~~~go
package store

import (
	"fmt"
	"sync"

	"github.com/go-examples-with-tests/database/v4/pkg"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Users() UserStore {
	return newUsers(ds) // 用于和 MariaDB 交互
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}
	return db.Close()
}

var (
	mysqlFactory Factory
	once         sync.Once
)

func GetMySQLFactoryOr(opts *pkg.MySQLOptions) (Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB

	// 此处是同步还是异步的？会阻塞吗？
	once.Do(func() {
		options := &pkg.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
		}
		// pkg.New --> *gorm.DB
		dbIns, err = pkg.New(options)
		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}
~~~

设计目的：

1. 构造出了一个 Factory 接口，这个接口是直接面向底层数据库的，可在其中设计出不同业务流对象，比如 `Users() UserStore`，又或者是其他的数据流；
2. `*datastore` 实现了 Factory 接口，应用层在使用时，获取到的  `*datastore` 实例。

连接数据库的核心代码是：

~~~go
package pkg

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Options defines optsions for mysql database.
type Options struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
}

// New create a new gorm db instance with the given options.
func New(opts *Options) (*gorm.DB, error) {
	// DSN Data Source Name: username:password@protocol(address)/dbname?param=value
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		true,
		"Local")
	// gorm ORM 在内部仍然使用的驱动器是：github.com/go-sql-driver/mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: opts.Logger, //FIXME 有什么作用？
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		//FIXME 此处返回 error，为什么没有打印日志？
		return nil, err
	}
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)
    
	return db, nil
}
~~~

其核心就是使用 `gorm.io/gorm` 框架连接底层的 MariaDB。

与之相关的代码是：

~~~go
package store

var client Factory

type Factory interface {
	Users() UserStore
	Close() error
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
~~~

# 4 获取数据源

如下是给应用层的访问接口：

~~~~go
type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Users() UserStore {
	return newUsers(ds) // 用于和 MariaDB 交互
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}
	return db.Close()
}
~~~~

`*datastore` 实例调用 `Users` 函数，获取这条业务线的数据操作通道。

~~~go
package store

import (
	"context"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{db: ds.db}
}

func (u *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return nil
}
func (u *users) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	return nil
}
func (u *users) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	return nil
}

func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	return nil
}

func (u *users) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	user := &v1.User{}
	err := u.db.Where("name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(110001, err.Error())
		}
		return nil, errors.WithCode(100101, err.Error())
	}
	return user, nil
}

func (u *users) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	return nil, nil
}
~~~

其中 UserStore 实际上是一个接口，用于封装业务层需要的业务内容：

~~~go
package store

import (
	"context"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

type UserStore interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
	Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
}
~~~

最后，在 main 中获取到数据：

~~~go
package main

import (
	"context"
	"log"

	"github.com/go-examples-with-tests/database/v4/pkg"
	"github.com/go-examples-with-tests/database/v4/store"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

func main() {
	...

	// CRUD，拿着 userStore 就可以和 MariaDB 交互
	userStore := dbFactory.Users()
	user, err := userStore.Get(context.TODO(), "admin", metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(user)
	}

	// close DB connection
	dbFactory.Close()
}
ant@MacBook-Pro v4 % go run main.go
2021/10/22 14:38:43 main.go:16: hello, world!
2021/10/22 14:38:43 main.go:20: MariaDB options:Host:127.0.0.1:3306, Username:iam, Password:iam59!z$, Database:iam, MaxIdleConnections:100, MaxOpenConnections:100, MaxConnectionLifeTime:10s, LogLevel:4
2021/10/22 14:38:43 main.go:34: &{{0 user-lingfei admin {} {} 2021-05-27 18:01:40 +0800 CST 2021-05-06 05:13:14 +0800 CST} admin $2a$10$WnQD2DCfWVhlGmkQ8pdLkesIGPf9KJB7N1mhSOqulbgN7ZMo44Mv2 admin@foxmail.com 1812884xxxx 1 0}
~~~

