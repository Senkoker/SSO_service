package storage

import (
	"GRPC_Service_sso/internal/DB_err"
	"GRPC_Service_sso/internal/module"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)
g
type Storage struct {
	Db *sql.DB
}

func NewSt(database_url string) (Storage, error) {
	db, err := sql.Open("pgx", database_url)
	if err != nil {
		return Storage{}, fmt.Errorf("Problemt to connect to database %v", err)
	}
	err = db.Ping()
	if err != nil {
		return Storage{}, fmt.Errorf("Problemt to ping database %v", err)
	}
	return Storage{Db: db}, nil
}

// ////////////////////////////////////////////REGISTER/////////////////////////////////////////////////////////////////
func (s *Storage) St_check_user(ctx context.Context, email string) (bool, error) {
	const op = "cheking_user_reg"
	stmt, err := s.Db.Prepare(checkUserRegister)
	defer stmt.Close()
	if err != nil {
		return false, fmt.Errorf("Problem to prepare statement %v,%s", err, op)
	}
	row := stmt.QueryRowContext(ctx, email)
	var id string
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("Internal err %v %s", err, op)
	}
	return false, DB_err.Already_exists
}
func (s *Storage) St_reg(ctx context.Context, email, code string, pass_hash []byte) (string, error) {
	stmt, err := s.Db.Prepare(registrationInsert)
	defer stmt.Close()
	if err != nil {
		return "", fmt.Errorf("Problemt to prepare statement %v", err)
	}
	result, err := stmt.QueryContext(ctx, email, pass_hash, code)
	pgerr, ok := err.(*pgconn.PgError)
	if err != nil {
		if pgerr.Code == "23505" && ok {
			return "", DB_err.Dublicate_name
		}
		return "", fmt.Errorf("Problemt to execute statement %w", err)
	}
	var id string
	for result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return "", fmt.Errorf("Problem to return id %w", err)
		}
	}
	return id, nil
}

// ///////////////////////////////////////////////////////////// LOGIN/////////////////////////////////////
func (s *Storage) St_login(ctx context.Context, email string) (module.User, error) {
	stmt, err := s.Db.Prepare(loginSelect)
	defer stmt.Close()
	if err != nil {
		return module.User{}, fmt.Errorf("Problemt to prepare statement %v", err)
	}
	var user module.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.Pass_hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return module.User{}, DB_err.Does_not_exist
		}
		return module.User{}, fmt.Errorf("Problem to return row:%w", err)
	}
	return user, nil
}

func (s *Storage) St_app(ctx context.Context, appid string) (module.AppID, error) {
	stmt, err := s.Db.Prepare(appSelect)
	defer stmt.Close()
	if err != nil {
		return module.AppID{}, fmt.Errorf("Problemt to prepare statement %v", err)
	}
	var app module.AppID
	err = stmt.QueryRowContext(ctx, appid).Scan(&app.Id, &app.Secret)
	if err != nil {
		if err == sql.ErrNoRows {
			return module.AppID{}, DB_err.Does_not_exist
		}
		return module.AppID{}, fmt.Errorf("Problemt to exec statement to db %v", err)
	}
	return app, nil
}

////////////////////////////////////Cheking user code //////////////////////////////////////

func (s *Storage) St_update_change_pass(ctx context.Context, email string) error {
	tx, err := s.Db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("Problemt to start transaction %v", err)
	}
	stmt, err := tx.Prepare(updateChangePass)
	if err != nil {
		return fmt.Errorf("Problemt to prepare statement %v", err)
	}
	result, err := stmt.ExecContext(ctx, email)
	if err != nil {
		return fmt.Errorf("Problemt to execute statement %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Problemt to return id %v", err)
	}
	if id == 0 {
		return DB_err.Does_not_exist
	}
	stmt, err = tx.Prepare(updateChangePassDeleteHash)
	if err != nil {
		return fmt.Errorf("Problemt to prepare statement %v", err)
	}
	result, err = stmt.ExecContext(ctx, email)
	if err != nil {
		return fmt.Errorf("Problemt to execute statement %v", err)
	}
	id, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Problemt to return id %v", err)
	}
	if id == 0 {
		return DB_err.Does_not_exist
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Problemt to commit transaction %v", err)
	}
	return nil

}

func (s *Storage) St_accept_copy(ctx context.Context, email string) (string, error) {
	stmt, err := s.Db.Prepare(acceptCodeSelect)
	if err != nil {
		return "", fmt.Errorf("Problemt to prepare statement %v", err)
	}
	var code string
	err = stmt.QueryRowContext(ctx, email).Scan(&code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", DB_err.Does_not_exist
		}
		return "", fmt.Errorf("Problemt to get code  %v", err)
	}
	return code, nil
}
func (s *Storage) St_relocate_user(ctx context.Context, email string) error {
	tx, err := s.Db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("Problem to start transaction %v", err)
	}
	result, err := tx.ExecContext(ctx, relocateUser, email)
	if err != nil {
		return fmt.Errorf("Problem to execute anonymous block %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Problem to execute anonymous block %v", err)
	}
	if id == 0 {
		return DB_err.Does_not_exist
	}
	_, err = tx.ExecContext(ctx, deleteUserHashRelocate, email)
	if err != nil {
		return fmt.Errorf("Problem to execute anonymous block %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Problem to commit transaction %v", err)
	}
	return nil
}

// //////////////////////////////////////RETRY  /////////////////////////////////////////////////////////
func (s *Storage) St_retry(ctx context.Context, code string, email string) (string, error) {
	stmt, err := s.Db.Prepare(updateUsercode)
	if err != nil {
		return "", fmt.Errorf("Problemt to prepare statement %v", err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, code, email)
	if err != nil {
		return "", fmt.Errorf("Problem to return data %v", err)
	}
	nums, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("Problem to count row %v", err)
	}
	if nums == 0 {
		return "", DB_err.Does_not_exist
	}
	return email, nil
}

// ////////////////////////////////////////////ADMIN/////////////////////////
func (s *Storage) St_adm(ctx context.Context, email string) (bool, error) {
	stmt, err := s.Db.Prepare(selectIdAdmin)
	if err != nil {
		return false, err
	}
	var id int64
	err = stmt.QueryRowContext(ctx, email).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, DB_err.Does_not_exist
		}
		return false, fmt.Errorf("Problem to return data %v", err)
	}
	return true, nil
}
