package postg_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	culc "github.com/ragnack97/protoculc/gen/go"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type expression struct {
	ex       string
	exeptErr string
}
type User struct {
	name            string
	email           string
	password        string
	expression      []expression
	exeptEx         []string
	exeptErrinlogin string
}

func TestRegister(t *testing.T) {
	t.Helper()
	t.Parallel()
	Users := []User{
		{name: "Pass Register", email: "NewPerson2", password: "123445", exeptErrinlogin: "no"},
		{name: "Error in Register 1", email: "NewPerson2", password: "", exeptErrinlogin: "yes"},
		{name: "Error in Register 2", email: "", password: "12345", exeptErrinlogin: "yes"},
		{name: "No unique email", email: "NewPerson2", password: "123445", exeptErrinlogin: "yes"},
	}
	ctx, st := New(t)
	for _, test := range Users {
		t.Run(test.name, func(t *testing.T) {
			_, err := st.AuthUser.Register(ctx, &culc.RegisterReq{
				Email:    test.email,
				Password: test.password,
			})
			if test.exeptErrinlogin == "yes" {
				assert.Error(t, err)
				errorFlags := map[string]bool{
					"empty email or password": false,
					" error in register":      false,
				}
				for errorStr := range errorFlags {
					if strings.Contains(err.Error(), errorStr) {

						errorFlags[errorStr] = true
					}
				}
				atLeastOneErrorFound := false
				for _, found := range errorFlags {
					if found {
						atLeastOneErrorFound = true
						break
					}
				}
				assert.True(t, atLeastOneErrorFound)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	t.Helper()
	t.Parallel()
	Users := []User{
		{name: "Pass Login", email: "NewPerson2", password: "123445", exeptErrinlogin: "no"},
		{name: "Error in Login 1", email: "NewPerson2", password: "", exeptErrinlogin: "yes"},
		{name: "Error in Login 2", email: "", password: "123445", exeptErrinlogin: "yes"},
		{name: "Wrong password", email: "NewPerson2", password: "1235", exeptErrinlogin: "yes"},
	}
	ctx, st := New(t)
	_, err := st.AuthUser.Register(ctx, &culc.RegisterReq{
		Email:    "NewPerson2",
		Password: "123445",
	})
	require.NoError(t, err)
	for _, test := range Users {
		t.Run(test.name, func(t *testing.T) {
			LoginReq, err := st.AuthUser.Login(ctx, &culc.LoginReq{
				Email:    test.email,
				Password: test.password,
			})
			if test.exeptErrinlogin == "yes" {
				assert.Error(t, err)
				errorFlags := map[string]bool{
					"empty email or password": false,
					"error in login":          false,
					"error in token":          false,
				}
				for errorStr := range errorFlags {
					if strings.Contains(err.Error(), errorStr) {

						errorFlags[errorStr] = true
					}
				}
				atLeastOneErrorFound := false
				for _, found := range errorFlags {
					if found {
						atLeastOneErrorFound = true
						break
					}
				}
				assert.True(t, atLeastOneErrorFound)
			} else {
				assert.NoError(t, err)
				//chek valid token
				token := LoginReq.GetToken()
				require.NotEmpty(t, token)
				tokenPars, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte("culc"), nil
				})
				logintime := time.Now()
				require.NoError(t, err)
				cl, ok := tokenPars.Claims.(jwt.MapClaims)
				require.True(t, ok)
				assert.NotEmpty(t, cl["uid"].(float64))
				assert.Equal(t, test.email, cl["email"].(string))
				assert.InDelta(t, logintime.Add(st.cfg.Token_time).Unix(), cl["exp"].(float64), 1)

			}
		})

	}

}

func TestCulc(t *testing.T) {
	t.Helper()

	Users := []User{
		{name: "true Login and true Colculate", email: "NewPerson2", password: "123445", expression: []expression{{ex: "5*5/2", exeptErr: "no"}, {ex: "1-3+4", exeptErr: "no"}, {ex: "6*7-2/4", exeptErr: "no"}}, exeptEx: []string{"12.5", "2", "41.5"}, exeptErrinlogin: "no"},
		{name: "true Login and false Colculate", email: "NewPerson2", password: "123445", expression: []expression{{ex: "5*5/)2", exeptErr: "yes"}, {ex: "f1-3+4", exeptErr: "yes"}, {ex: "6*7-2/0", exeptErr: "yes"}, {ex: "5-1", exeptErr: "no"}}, exeptEx: []string{"", "", "", "4"}, exeptErrinlogin: "no"},
	}
	ctx, st := New(t)
	email := "NewPerson2"
	password := "123445"
	newuser := &culc.RegisterReq{
		Email:    email,
		Password: password,
	}
	RegResp, err := st.AuthUser.Register(ctx, newuser)
	assert.NoError(t, err)

	assert.NotEmpty(t, RegResp.GetUserId())
	for _, test := range Users {
		t.Run(test.name, func(t *testing.T) {

			ResLog, err := st.AuthUser.Login(ctx, &culc.LoginReq{Email: test.email,
				Password: test.password})

			require.NoError(t, err)

			token := ResLog.GetToken()
			assert.NotEmpty(t, token)
			tokenPars, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte("culc"), nil
			})
			logintime := time.Now()
			require.NoError(t, err)
			cl, ok := tokenPars.Claims.(jwt.MapClaims)
			require.True(t, ok)
			assert.NotEmpty(t, cl["uid"].(float64))
			assert.Equal(t, test.email, cl["email"].(string))
			assert.InDelta(t, logintime.Add(st.cfg.Token_time).Unix(), cl["exp"].(float64), 1)
			for i, testex := range test.expression {
				t.Run(fmt.Sprintf("%s выражение номер %d", test.name, i), func(t *testing.T) {
					ctxforCulc := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
						"authorization", "Bearer "+token,
					))

					RespEx, err := st.AuthUser.Calculate(ctxforCulc, &culc.CalculateReq{
						Expression: testex.ex,
					})

					if testex.exeptErr == "yes" {

						assert.Error(t, err, "error add ex to db")
						errorFlags := map[string]bool{
							"error add ex to db": false,
							"calculation error":  false,
						}
						for errorStr := range errorFlags {
							if strings.Contains(err.Error(), errorStr) {

								errorFlags[errorStr] = true
							}
						}
						atLeastOneErrorFound := false
						for _, found := range errorFlags {
							if found {
								atLeastOneErrorFound = true
								break
							}
						}
						assert.True(t, atLeastOneErrorFound)
					} else {
						require.NoError(t, err)
						assert.Equal(t, RespEx.Result, test.exeptEx[i])
					}

				})
			}
		})

	}

}
func TesConfig(t *testing.T) {
	t.Helper()
	t.Parallel()
	ctx, st := New(t)
	st.AuthUser.UpdateConfig(ctx, &culc.ConfigReq{
		Plus:     10,
		Minus:    20,
		Div:      10,
		Exponent: 10,
		MultP:    100,
	})
}
