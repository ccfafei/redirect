package storage

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"redirect/utils"
)

var dbService = &DatabaseService{}

// DatabaseService 数据库服务
type DatabaseService struct {
	Connection *sqlx.DB
}

type BatchQueryArgs struct {
	query string
	args  []interface{}
}

// InitDatabaseService 初始化数据库服务
func InitDatabaseService() (*DatabaseService, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		utils.DatabaseConfig.Host, utils.DatabaseConfig.Port, utils.DatabaseConfig.User,
		utils.DatabaseConfig.Password, utils.DatabaseConfig.DbName)
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return dbService, err
	}
	conn.SetMaxOpenConns(utils.DatabaseConfig.MaxOpenConns)
	conn.SetMaxIdleConns(utils.DatabaseConfig.MaxIdleConn)
	conn.SetConnMaxLifetime(0) // always REUSE
	dbService.Connection = conn
	return dbService, nil
}

// DbNamedExec 执行带有命名参数的sql语句
func DbNamedExec(query string, args interface{}) error {
	_, err := dbService.Connection.NamedExec(query, args)
	return err
}

func DbExec(query string, args ...any) error {
	_, err := dbService.Connection.Exec(query, args...)
	return err
}

// DbExecTx 执行事务
func DbExecTx(query ...string) error {
	tx := dbService.Connection.MustBegin()
	for _, s := range query {
		tx.MustExec(s)
	} // end of for
	err := tx.Commit()
	if err != nil {
		return tx.Rollback()
	}
	return nil
} // end of func

// DbBatchExecTx 批量执行带参数事务
func DbBatchExecTx(query []*BatchQueryArgs) error {
	tx := dbService.Connection.MustBegin()
	for _, item := range query {
		tx.Exec(tx.Rebind(item.query), item.args...)
	} // end of for
	err := tx.Commit()
	if err != nil {
		return tx.Rollback()
	}
	return nil
} // end of func

func GetDbTx() *sqlx.Tx {
	tx := dbService.Connection.MustBegin()
	return tx
}

// DbGet 获取单条记录
func DbGet(query string, dest interface{}, args ...interface{}) error {
	err := dbService.Connection.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// DbSelect 获取多条记录
func DbSelect(query string, dest interface{}, args ...interface{}) error {
	return dbService.Connection.Select(dest, query, args...)
}

//RebindQuery 格式化sql
func RebindQuery(query string) string {
	return dbService.Connection.Rebind(query)
}

// DbClose 关闭数据库连接
func DbClose() {
	dbService.Connection.Close()
}
