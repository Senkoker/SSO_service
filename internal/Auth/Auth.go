package Auth

import (
	"GRPC_Service_sso/config"
	"GRPC_Service_sso/internal/DB_err"
	"GRPC_Service_sso/internal/jwt_token"
	"GRPC_Service_sso/internal/mail_sender"
	"GRPC_Service_sso/internal/module"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Auth struct {
	log    *slog.Logger
	cfg    config.Cfg
	St_log Storage_login
	St_reg Storage_register
	St_acc Storage_accept
	St_adm Storage_admin
}

func NewAuth(log *slog.Logger, cfg config.Cfg, st_l Storage_login, st_r Storage_register, st_ac Storage_accept, st_adm Storage_admin) *Auth {
	return &Auth{log: log, cfg: cfg, St_log: st_l, St_acc: st_ac, St_adm: st_adm, St_reg: st_r}
}

func Generatecode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func User_id_generator() int64 {
	time_now := time.Now()
	given_time := time.Date(2025, 03, 21, 20, 0, 0, 0, time.UTC)
	duration := time_now.Sub(given_time)
	f_part := int(duration.Seconds())
	rand.Seed(time.Now().UnixNano())
	s_part := rand.Intn(100000)
	id, _ := strconv.ParseInt(fmt.Sprintf("%d%d", f_part, s_part), 10, 64)
	return id
}

type Storage_login interface {
	St_login(ctx context.Context, email string) (module.User, error)
	St_app(ctx context.Context, appid string) (module.AppID, error)
}
type Storage_register interface {
	St_reg(ctx context.Context, email, code string, pass_hash []byte) (string, error)
	St_check_user(ctx context.Context, email string) (bool, error)
	St_retry(ctx context.Context, code string, email string) (string, error)
}
type Storage_accept interface {
	St_accept_copy(ctx context.Context, email string) (string, error)
	St_relocate_user(ctx context.Context, email string) error
	St_update_change_pass(ctx context.Context, email string) error
}
type Storage_admin interface {
	St_adm(ctx context.Context, email string) (bool, error)
}

// DONE
func (s *Auth) Auth_login(ctx context.Context, email, pass, appid string) (string, error) {
	const op = "login_user"
	logger := s.log.With("op", op)
	user, err := s.St_log.St_login(ctx, email)
	if err != nil {
		logger.Error("Failed to login", "err", err)
		return "", fmt.Errorf("filed to login err: %w", err)
	}
	pass_hash := user.Pass_hash
	err = bcrypt.CompareHashAndPassword(pass_hash, []byte(pass))
	if err != nil {
		logger.Error("Wrong password", "err", err)
		return "", DB_err.Invalid_password
	}

	app_id, err := s.St_log.St_app(ctx, appid)
	if err != nil {
		logger.Error("Internal app id error", "err", err)
		return "", fmt.Errorf("Internal error: %w", err)
	}

	token, err := jwt_token.NewJWT(user, app_id, s.cfg.Server.TokenTTL)
	if err != nil {
		logger.Error("Failed to create token", "err", err)
		return "", fmt.Errorf("Failed to create token,%w", err)
	}
	return token, nil
}
func (s *Auth) Auth_change_pass(ctx context.Context, email, pass string) (bool, error) {
	op := "change_pass"
	logger := s.log.With("op", op)
	user, err := s.St_log.St_login(ctx, email)
	if err != nil {
		logger.Error("Fialed to chech user", "err", err)
		return false, err
	}
	code := Generatecode(10)
	pass_code := email + "/" + "change_pass" + "/" + code
	pass_hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", "err", err)
		return false, fmt.Errorf("Failed to hash password,%w", err)
	}
	_, err = s.St_reg.St_reg(ctx, user.Email, code, pass_hash)
	if err != nil {
		logger.Error("Failed to change password", "err", err)
		return false, fmt.Errorf("Failed to change password,%w", err)
	}
	err = mail_sender.Mail_sender(user.Email, s.cfg.Server.Url_accepter, pass_code, s.cfg.Server.Mail_sender, s.cfg.Server.Mail_password)
	if err != nil {
		logger.Error("Failed to send mail", "err", err)
		return false, fmt.Errorf("Failed to send mail,%w", err)
	}
	return true, nil
}

// /DONE
func (s *Auth) Auth_register(ctx context.Context, email, pass string) (string, error) {
	const op = "auth_register"
	logger := s.log.With("op", op)
	resp, err := s.St_reg.St_check_user(ctx, email)
	if err != nil && !resp {
		logger.Error("Failed to check user", "err", err)
		return "", err
	}
	pass_hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", "err", err)
		fmt.Errorf("Internal_eror:%w", err)
		return "", err
	}
	code := Generatecode(10)
	rand.Seed(time.Now().UnixNano())
	var id string
	id, err = s.St_reg.St_reg(ctx, email, code, pass_hash)
	if err != nil {
		logger.Error("Register error", "err", err)
		fmt.Errorf("Register_eror:%w", err)
		return "", err
	}
	reuslt_id := email + "/" + code
	err = mail_sender.Mail_sender(email, s.cfg.Server.Url_accepter, reuslt_id, s.cfg.Server.Mail_sender, s.cfg.Server.Mail_password)
	if err != nil {
		logger.Error("Problem to send code to email recepient", "err", err)
		return id, fmt.Errorf("Problem to send code to emailrecepient,err:%w", err)
	}
	logger.Info("mail sender succsesfully")
	return id, nil
}

// DONE
func (s *Auth) Auth_Accept(ctx context.Context, usercode string) (bool, error) {
	if strings.Contains(usercode, "change_pass") {
		const op = "Change_pass_accept"
		logger := s.log.With("op", op)
		pattern := "/"
		result := strings.SplitN(usercode, pattern, 3)
		email := result[0]
		code := result[2]
		db_code, err := s.St_acc.St_accept_copy(ctx, email)
		if err != nil {
			logger.Error("Failed to accept user", "err", err)
			return false, err
		}
		if db_code != code {
			return false, DB_err.Invalid_usercode
		}
		err = s.St_acc.St_update_change_pass(ctx, email)
		if err != nil {
			logger.Error("Failed to relocate user", "err", err)
			return false, err
		}
		return true, nil
	}
	const op = "auth_Accept"
	logger := s.log.With("op", op)
	pattern := "/"
	result := strings.SplitN(usercode, pattern, 2)
	email := result[0]
	code := result[1]
	if strings.Contains(usercode, "change_pass") {

	}
	db_code, err := s.St_acc.St_accept_copy(ctx, email)
	if err != nil {
		logger.Error("Failed to accept user", "err", err)
		return false, err
	}
	if db_code != code {
		return false, DB_err.Invalid_usercode
	}
	err = s.St_acc.St_relocate_user(ctx, email)
	if err != nil {
		logger.Error("Failed to relocate user", "err", err)
		return false, err
	}
	return true, nil
}

// DONE
func (s *Auth) Auth_Retry(ctx context.Context, email string) (bool, error) {
	if strings.Contains(email, "Re_ch_pass") {
		const op = "auth_Retry"
		code := Generatecode(10)
		data := strings.SplitN(email, "/", 2)
		email = data[1]
		resp, err := s.St_reg.St_retry(ctx, code, email)
		logger := s.log.With("operation", op)
		if err != nil {
			logger.Error("Failed to retry", "err", err)
			return false, fmt.Errorf("Failed to retry:%w", err)
		}
		usercode := email + "/" + "change_pass" + "/" + code
		err = mail_sender.Mail_sender(resp, s.cfg.Server.Url_accepter, usercode, s.cfg.Server.Mail_sender, s.cfg.Server.Mail_password)
		if err != nil {
			logger.Error("Problem to send code to email recepient", "err", err)
			return false, fmt.Errorf("Problem to send code to emailrecepient,err:%w", err)
		}
		return true, nil
	}
	const op = "auth_Retry"
	code := Generatecode(10)
	resp, err := s.St_reg.St_retry(ctx, code, email)
	logger := s.log.With("operation", op)
	if err != nil {
		logger.Error("Failed to retry", "err", err)
		return false, fmt.Errorf("Failed to retry:%w", err)
	}
	usercode := email + "/" + code
	err = mail_sender.Mail_sender(resp, s.cfg.Server.Url_accepter, usercode, s.cfg.Server.Mail_sender, s.cfg.Server.Mail_password)
	if err != nil {
		logger.Error("Problem to send code to email recepient", "err", err)
		return false, fmt.Errorf("Problem to send code to emailrecepient,err:%w", err)
	}
	return true, nil

}
func (s *Auth) Auth_IsAdmin(ctx context.Context, email string) (bool, error) {
	const op = "auth_IsAdmin"
	resp, err := s.St_adm.St_adm(ctx, email)
	logger := s.log.With("op", op)
	if err != nil {
		logger.Error("Failed to adm", "err", err)
		return false, err
	}
	return resp, nil
}
