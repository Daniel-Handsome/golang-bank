package api

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	// "reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	mockdb "github.com/daniel/master-golang/db/mock"
	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/utils"
)

// type eqCreateUserParamsMatcher struct{
// 	// arg db.CreateUserParams
// 	originPassword string
// }

// func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
// 	arg, ok := x.(db.CreateUserParams)
// 	fmt.Printf("%#v", arg)
// 	panic("ttt")
// 	if !ok {
//         return false
//     }

// 	return arg.Password == e.originPassword
// }

// func (eqCreateUserParamsMatcher) String() string {
// 	return "is anything"
// }

// func eqCreateUserParams(originPassword string) gomock.Matcher {
// 	return eqCreateUserParamsMatcher{
// 		// arg: arg,
// 		originPassword: originPassword,
// 	}
// }

type eqCreateUserParamsMatcher struct {
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.Password)

	if err != nil {
		return false
	}

	return true
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("password %v", e.password)
}

func EqCreateUserParams(password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{password}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)
	testCases := []struct {
		name string
		request map[string]interface{}
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(record *httptest.ResponseRecorder)
	}{
		{
            name: "created",
            request: map[string]interface{}{
				"name" : user.Name,
				"password" : password,
				"full_name" : user.FullName,
				"email" : user.Email,
            },
			buildStubs: func(store *mockdb.MockStore) {
				// 這邊自訂的只要丟原本的password是因為 這邊是request完成才會進來的 所以他 interface那個參數就會是request打進來的interface
				store.
                    EXPECT().
                    CreateUser(gomock.Any(), EqCreateUserParams(password)).
                    Times(1).
                    Return(user, nil)
			},
			checkResponse: func(record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, record.Code)
				recordBodyMatchUser(t, record.Body, user)
			},
		},
		{
            name: "InvalidUserName",
			// 當搜尋條件與特殊標記衝突時,如：逗號（,），或操作（|），中橫線（-）等則需要使用 UTF-8十六進位制表示形式
            request: map[string]interface{}{
				"name" : "invalid-user#1",
				"password" : password,
				"full_name" : user.FullName,
				"email" : user.Email,
            },
			buildStubs: func(store *mockdb.MockStore) {
				store.
                    EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Times(0)
			},
			checkResponse: func(record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, record.Code)
			},
		},
		{
            name: "DuplicateUsername",
            request: map[string]interface{}{
				"name" : user.Name,
				"password" : password,
				"full_name" : user.FullName,
				"email" : user.Email,
            },
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, record.Code)
			},
		},
		{
            name: "Internal error",
            request: map[string]interface{}{
				"name" : user.Name,
				"password" : password,
				"full_name" : user.FullName,
				"email" : user.Email,
            },
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, record.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.request)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}


func randomUser(t *testing.T) (user db.User, password string)  {
	password = utils.RandString(6)
	hashPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
        Name: utils.RandString(10),
		Password: hashPassword,
		FullName: utils.RandString(10),
		Email: utils.RandEmail(),
	}

	return
}

func recordBodyMatchUser(t *testing.T, record *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(record)
	require.NoError(t, err)

	var getUser db.User 
	err =json.Unmarshal(data, &getUser)
	require.NoError(t, err)

	require.Equal(t, user.Name, getUser.Name)
	require.Equal(t, user.FullName, getUser.FullName)
	require.Equal(t, user.Email, getUser.Email)

	require.Empty(t, getUser.Password)
}