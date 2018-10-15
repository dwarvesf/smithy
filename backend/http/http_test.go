package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-kit/kit/log"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	auth "github.com/dwarvesf/smithy/backend/auth"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	utilPG "github.com/dwarvesf/smithy/backend/config/database/pg/util"
	"github.com/dwarvesf/smithy/backend/endpoints"
	"github.com/dwarvesf/smithy/backend/service"
	utilTest "github.com/dwarvesf/smithy/common/utils/database/pg/test/set1"
)

const (
	secretKey string = "lalala"
	Admin     string = "admin"
	User      string = "user"
)

func TestNewHTTPHandler(t *testing.T) {
	//make up-dashboard
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	tsDashboard := httptest.NewServer(initDashboardServer(t, cfg))
	defer tsDashboard.Close()

	tests := []struct {
		name       string
		header     http.Header
		wantErr    string
		wantStatus int
	}{
		{
			name:       "Success",
			header:     newAuthHeader(auth.New(secretKey, "aaa", Admin, false).Encode()),
			wantStatus: http.StatusOK,
		},
		{
			name:       "Header is nil",
			header:     nil,
			wantErr:    "jwtauth: no token found",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Header is empty",
			header:     http.Header{},
			wantErr:    "jwtauth: no token found",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong secret key",
			header:     newAuthHeader(auth.New("wrong", "aaa", Admin, false).Encode()),
			wantErr:    "signature is invalid",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong token string",
			header:     newAuthHeader("blabla"),
			wantErr:    "token contains an invalid number of segments",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong algorithm",
			header: newAuthHeader(newJwt512Token([]byte(secretKey), jwtauth.Claims{
				"username": "aaa",
				"role":     Admin,
			})),
			wantErr:    "jwtauth: token is unauthorized",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong secret key and algorithm",
			header: newAuthHeader(newJwt512Token([]byte("wrong"), jwtauth.Claims{
				"username": "aaa",
				"role":     Admin,
			})),
			wantErr:    "signature is invalid",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Client can't agent-sync",
			header:     newAuthHeader(auth.New(secretKey, "bbb", User, false).Encode()),
			wantErr:    "Unauthorized",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if status, resp := testRequest(t, tsDashboard, "GET", "/models", tt.header, nil); status != tt.wantStatus || resp != tt.wantErr {
				t.Errorf("NewHTTPHandler() = (%v, %d), want (%v, %d)", resp, status, tt.wantErr, tt.wantStatus)
			}
		})
	}

}

func TestAuthorized(t *testing.T) {
	//make up-dashboard
	cfg, _ := utilTest.CreateConfig(t)

	dbTest := []string{"test1", "test2"}
	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}
	}

	cfg.DBUsername = "agent_db_manager"
	cfg.DBPassword = "1"

	clearACL := utilTest.ACLUsersTable(t, cfg)
	defer clearACL()

	// re-connect with DBusername = agent_db_manager
	err := cfg.UpdateDB()
	if err != nil {
		t.Fatalf("Fail to update connection. %s", err.Error())
	}

	for _, dbase := range cfg.Databases {
		// set schema for current db connection
		err = cfg.DB(dbase.DBName).Exec("SET search_path TO " + dbase.SchemaName).Error
		if err != nil {
			t.Fatalf("Fail to set search_path to created schema. %s", err.Error())
		}
	}

	tsDashboard := httptest.NewServer(initDashboardServer(t, cfg))
	defer tsDashboard.Close()

	headerAdmin := newAuthHeader(auth.New(secretKey, "aaa", Admin).Encode())
	headerUser := newAuthHeader(auth.New(secretKey, "bbb", User).Encode())
	header := newAuthHeader(auth.New(secretKey, "118168790790272259317", User, true).Encode())

	//user ACL : cru
	//table ACL : cr

	type args struct {
		HTTPMethod, url string
		data            []byte
		header          http.Header
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    error
	}{
		{
			name: "create record with permission",
			args: args{
				HTTPMethod: "POST",
				url:        fmt.Sprintf("/databases/%s/table/users/create", dbTest[0]),
				header:     headerAdmin,
				data: []byte(`{
					"fields": 	[ "name" ],
					"data":     [ "lorem ipsum" ]
				}`),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail to update DB. User has cru permisstion but table permission just has cr permisstion",
			args: args{
				HTTPMethod: "PUT",
				url:        fmt.Sprintf("/databases/%s/table/users/update", dbTest[0]),
				header:     headerAdmin,
				data: []byte(`{
						"fields": ["id", "name" ],
						"data":     [1, "aaaaaa" ]
				}`),
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    auth.ErrUnauthorized,
		},
		{
			name: "Fail to delete DB. User has cru permisstion. cant delete",
			args: args{
				HTTPMethod: "DELETE",
				url:        fmt.Sprintf("/databases/%s/table/users/delete", dbTest[0]),
				header:     headerAdmin,
				data: []byte(`{
					"filter": {
						"fields": [ "id" ],
						"data":     [ "1" ]
				   }
				}`),
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    auth.ErrUnauthorized,
		},
		{
			name: "query with permission. group has permission but user not",
			args: args{
				HTTPMethod: "POST",
				url:        fmt.Sprintf("/databases/%s/table/users/query", dbTest[0]),
				header:     headerAdmin,
				data: []byte(`{
					"fields": 	[ "name" ]
				}`),
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    auth.ErrUnauthorized,
		},
		{
			name: "create group with admin permission",
			args: args{
				HTTPMethod: "POST",
				url:        fmt.Sprintf("/groups"),
				header:     headerAdmin,
				data: []byte(`{
					"group": {
						"name": "admin1",
						"description": "lora"
					}
				}`),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "create group with user permission",
			args: args{
				HTTPMethod: "POST",
				url:        fmt.Sprintf("/groups"),
				header:     headerUser,
				data: []byte(`{
					"group": {
						"name": "admin1",
						"description": "lora"
					}
				}`),
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    auth.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, resp := testRequest(t, tsDashboard, tt.args.HTTPMethod, tt.args.url, tt.args.header, tt.args.data)
			if status != tt.wantStatus || (tt.wantErr != nil && resp != tt.wantErr.Error()) {
				t.Errorf("Authorized() = (%v, %d), want (%v, %d)", resp, status, tt.wantErr, tt.wantStatus)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	//make up-dashboard
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	tsDashboard := httptest.NewServer(initDashboardServer(t, cfg))
	defer tsDashboard.Close()

	// test login api
	loginTests := []struct {
		name       string
		jsonString []byte
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "Login success",
			jsonString: []byte(`{"username":"aaa", "password": "abc"}`),
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name:       "Login wrong username",
			jsonString: []byte(`{"username":"adfs", "password": "abc"}`),
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name:       "Login wrong password",
			jsonString: []byte(`{"username":"aaa", "password": "conmeocon"}`),
			wantStatus: 401,
			wantErr:    true,
		},
	}

	for _, tt := range loginTests {
		t.Run(tt.name, func(t *testing.T) {
			if status := loginTestRequest(t, tsDashboard, "/auth/login", tt.jsonString); status != tt.wantStatus && tt.wantErr {
				t.Errorf("Login() = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

//
// Test helper functions
//

func testRequest(t *testing.T, ts *httptest.Server, method, path string, header http.Header, body []byte) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v[0])
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}
	defer resp.Body.Close()

	data := &auth.ErrAuthentication{}

	if err = json.Unmarshal(respBody, data); err != nil {
		return resp.StatusCode, string(respBody)
	}

	return resp.StatusCode, data.Error
}

func newAuthHeader(tokenStr string) http.Header {
	h := http.Header{}
	h.Set("Authorization", "BEARER "+tokenStr)
	return h
}

func initDashboardServer(t *testing.T, cfg *backendConfig.Config) http.Handler {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	pg, _, _ := utilPG.CreateTestDatabase(t)
	// defer closeDB()

	if err := utilPG.SeedCreateTable(pg); err != nil {
		panic(fmt.Sprintf("fail to migrate table by error %v", err))
	}

	s, err := service.NewService(cfg, pg)
	if err != nil {
		t.Fatal(err)
	}

	return NewHTTPHandler(
		endpoints.MakeServerEndpoints(s),
		logger,
		os.Getenv("ENV") == "local",
		cfg,
		s,
	)
}

func newJwt512Token(secret []byte, claims ...jwtauth.Claims) string {
	// use-case: when token is signed with a different alg than expected
	token := jwt.New(jwt.GetSigningMethod("HS512"))
	if len(claims) > 0 {
		token.Claims = claims[0]
	}
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("error at newJwt512Token")
	}
	return tokenStr
}

func loginTestRequest(t *testing.T, ts *httptest.Server, path string, body []byte) int {
	req, err := http.NewRequest("POST", ts.URL+path, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return 0
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, errRespone := client.Do(req)
	if errRespone != nil {
		t.Fatal(err)
		return 0
	}

	return response.StatusCode
}
