package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"redirect/model"
	"redirect/utils"
	"strings"
)

func FindAllUsers(search string, page, size int) ([]model.User, error) {
	var (
		found  []model.User
		offset = (page - 1) * size
		query  = `SELECT  id,account,name FROM public.users l WHERE 1=1 `
	)

	if !utils.EmptyString(search) {
		likeSearch := "%" + search + "%"
		query += fmt.Sprintf(` AND l.account LIKE '%s' OR l.name LIKE '%s'`, likeSearch, likeSearch)
	}

	query += ` ORDER BY l.id DESC LIMIT $1 OFFSET $2`
	return found, DbSelect(query, &found, size, offset)
}

func FindAllUsersTotals(search string) (int, error) {
	var query = `SELECT  count(l.id) as total_count FROM public.users l WHERE 1=1 `
	if !utils.EmptyString(search) {
		likeSearch := "'%" + search + "%'"
		query += fmt.Sprintf(` AND l.account LIKE '%s' OR l.name LIKE '%s'`, likeSearch, likeSearch)
	}
	var count int
	return count, DbGet(query, &count)
}

// NewUser 新建用户
func NewUser(user *model.User) error {
	query := `INSERT INTO public.users (account, "password") VALUES(:account,:password)`
	return DbNamedExec(query, user)
}

// DelUsers 批量删除用户
func DelUsers(ids []int) error {
	sql := `DELETE FROM public.users  WHERE  id in (?)`
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return err
	}
	formatQuery := RebindQuery(query)
	return DbExec(formatQuery, args...)
}

func UpdateUser(user model.User) error {
	query := `UPDATE public.users SET name = :name , "password" = :password WHERE id = :id`
	return DbNamedExec(query, user)
}

// FindUserByAccount 根据账号查找用户
func FindUserByAccount(account string) (model.User, error) {
	var user model.User
	query := `SELECT id,account,name,password FROM public.users u WHERE lower(u.account) = $1`
	return user, DbGet(query, &user, strings.ToLower(account))
}

// FindUserByUserId 根据用户ID查找用户
func FindUserByUserId(userId int) (model.User, error) {
	var user model.User
	query := `SELECT id,account,name,password FROM public.users u WHERE id = $1`
	return user, DbGet(query, &user, userId)
}
