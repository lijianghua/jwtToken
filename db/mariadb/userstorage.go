package mariadb

import (
	"context"
	"database/sql"
	"github.com/twinj/uuid"
	"jwtToken/service/userService"
	"jwtToken/util"
	"strconv"
	"strings"
)

//package mariadb: DB层 user 处理
type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}
func (s UserStorage) Create(ctx context.Context, u *userService.User) error {
	query := `INSERT INTO user
			   (id,user_name,user_pwd,phone,email,status) 
			  VALUES (?,?,?,?,?,?)
	`
	u.UserID = uuid.NewV4().String()

	//Notice:User.Password始终存放的是hash之前的明文密码
	hashedPass, err := util.HashPass(u.Password)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, query, u.UserID, u.Username, hashedPass, u.Phone, u.Email, 1)

	return err
}

func (s UserStorage) Update(ctx context.Context, u *userService.User) error {
	query := `UPDATE user
              SET user_name = ?,
				  phone = ?,
                  status = ?
                  
	`
	args := []interface{}{
		u.Username,
	}
	if u.Password != "" {
		hashedPass, err := util.HashPass(u.Password)
		if err != nil {
			return &userService.FieldError{
				UserError: userService.UserError{
					Type:    userService.InvalidArgument,
					Code:    "invalid password",
					Message: err.Error(),
				},
				Field: "password",
			}
		}
		query += ", user_pwd = ?"
		args = append(args, hashedPass)
	}
	query += " WHERE id = ?"
	args = append(args, u.UserID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s UserStorage) Delete(ctx context.Context, id string) error {

	query := `DELETE FROM user
              WHERE id=?
	`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return userService.ErrNotFound
	}
	return nil
}

func (s UserStorage) Verify(ctx context.Context, name string, pass string) error {
	var hashedPwd string

	err := s.db.QueryRow("select user_pwd from USER where user_name = ? limit 1", name).Scan(&hashedPwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return userService.ErrNotFound
		}
		return err
	}
	//verify password
	suc := util.VerifyPass(hashedPwd, pass)
	if !suc {
		return &userService.UserError{
			Type:    userService.InvalidArgument,
			Code:    "invalid password",
			Message: "invalid password",
		}
	}
	return nil
}

func (s UserStorage) GetByName(ctx context.Context, userName string) (*userService.User, error) {
	query := `SELECT id, user_name, created_at,updated_at,last_active,email,phone,status
		FROM user
		WHERE user_name = ?
	`
	row := s.db.QueryRowContext(ctx, query, userName)
	u, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, userService.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s UserStorage) Get(ctx context.Context, id string) (*userService.User, error) {
	query := `SELECT id, user_name, created_at,updated_at,last_active,email,phone,status
		FROM user
		WHERE id = ?
	`
	row := s.db.QueryRowContext(ctx, query, id)
	u, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, userService.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s UserStorage) List(ctx context.Context, opts *userService.ListOptions) (*userService.ListResponse, error) {
	query := `
		SELECT id, user_name, created_at,updated_at,last_active,email,phone,status
		FROM user
	`
	where := []string{}
	args := []interface{}{}

	if opts.Status != "" {
		where = append(where, "status = ?")
		args = append(args, opts.Status)
	}
	if opts.Cursor != "" {
		// TODO: implement cursor based
	}
	query += " WHERE " + strings.Join(where, " AND ")
	if opts.Sort != "" {
		var mode, field string
		if strings.HasPrefix(opts.Sort, "-") {
			mode = "DESC"
			field = strings.TrimPrefix(opts.Sort, "-")
		} else {
			mode = "ASC"
			field = opts.Sort
		}
		query += " ORDER BY " + field + " " + mode
	}
	query += " LIMIT " + strconv.FormatInt(opts.PerPage, 10)
	if opts.Cursor == "" {
		query += " OFFSET " + strconv.FormatInt(int64(opts.Page)*opts.PerPage, 10)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*userService.User, 0)

	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	// TODO: get total of users using COUNT(*) from the SELECT

	return &userService.ListResponse{
		Total:   int64(len(users)),
		PerPage: opts.PerPage,
		Users:   users,
	}, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(row rowScanner) (*userService.User, error) {
	u := new(userService.User)
	err := row.Scan(
		&u.UserID,
		&u.Username,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastActive,
		&u.Email,
		&u.Phone,
		&u.Status,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}
