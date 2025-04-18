package test

import (
	"GRPC_Service_sso/internal/test/kit"
	sso_v1_ssov1 "github.com/Senkoker/sso_proto/proto/proto_go/protobufcontract/protobufcontract"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

type User struct {
	id       string
	email    string
	password string
}

const (
	passlen        = 10
	app_id         = "1"
	secret         = "my_secret"
	delta          = 1
	app_id_int int = 1
)

func Test_register_login_happy(t *testing.T) {
	kit_st := kit.Kit_new(t)
	for i := 0; i < 20; i++ {
		user := User{app_id, gofakeit.Email(), gofakeit.Password(true, true, true, true, false, passlen)}
		resp, err := kit_st.Auth_client.Register(kit_st.Ctx, &sso_v1_ssov1.Registrequest{Email: user.email, Password: user.password})
		id := resp.GetUserid()
		require.NoError(t, err)
		assert.NotEmpty(t, resp.GetUserid())
		var usercode string
		query := `SELECT code FROM users_hash WHERE email = $1`
		stmt, err := kit_st.Db.Prepare(query)
		if err != nil {
			log.Fatalln("get data", err)
		}
		err = stmt.QueryRow(user.email).Scan(&usercode)
		if err != nil {
			log.Fatalln("scan_data", err, user.email)
		}
		defer stmt.Close()
		usercode = user.email + "/" + usercode
		acceptrep, err := kit_st.Auth_client.Accept(kit_st.Ctx, &sso_v1_ssov1.Acceptrequest{Usercode: usercode})
		require.NoError(t, err)
		assert.Equal(t, acceptrep.GetAccresp(), true)
		logintime := time.Now()
		login_resp, err := kit_st.Auth_client.Login(kit_st.Ctx, &sso_v1_ssov1.Loginrequest{Email: user.email, Password: user.password, Appid: user.id})
		require.NoError(t, err)
		token := login_resp.GetToken()
		parsedtoken, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		require.NoError(t, err)
		claims, _ := parsedtoken.Claims.(jwt.MapClaims)
		assert.Equal(t, claims["email"], user.email)
		assert.Equal(t, int64(claims["user.id"].(float64)), id)
		assert.Equal(t, int(claims["app.id"].(float64)), app_id_int)
		assert.InDelta(t, claims["exp"], logintime.Add(kit_st.Cfg.TokenTtl).Unix(), delta)
	}
}
