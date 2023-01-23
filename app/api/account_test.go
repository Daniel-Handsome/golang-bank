package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/daniel/master-golang/db/mock"
	db "github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/token"
	"github.com/daniel/master-golang/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestServer_getAccount(t *testing.T) {
	user, _ := randomUser(t)
	account := randAccount(user.Name)

	testCases := []struct{
		name string
		accountID int64
		setUpAuth func (t *testing.T, request *http.Request,tokenMaker token.Marker)
		buildStub func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, record *httptest.ResponseRecorder)
	}{
		{
			name : "OK",
			accountID : account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			buildStub : func(store *mockdb.MockStore) {
				store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(account, nil)
			},
			checkResponse: func(t *testing.T, record *httptest.ResponseRecorder) {
				require.Equal(t, record.Code, http.StatusOK)
				responseBodyMatchAccount(t, record.Body, account)
			},
		},
		{
			name : "InvaildID",
			accountID : 0,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			buildStub : func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse: func(t *testing.T, record *httptest.ResponseRecorder) {
				require.Equal(t, record.Code, http.StatusBadRequest)
			},
		},
		{
			name : "internal error",
			accountID : account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			buildStub : func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, record *httptest.ResponseRecorder) {
				require.Equal(t, record.Code, http.StatusInternalServerError)
			},
		},
		{
			name : "Not found",
			accountID : account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			buildStub : func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, record *httptest.ResponseRecorder) {
				require.Equal(t, record.Code, http.StatusNotFound)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// new store
			controller := gomock.NewController(t)
		
			store := mockdb.NewMockStore(controller)
			testCase.buildStub(store)
			
			// start server
			server := newTestServer(t, store)
			record := httptest.NewRecorder()
		
			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			//httptest.NewRequest的第三个参数可以用来传递body数据，必须实现io.Reader接口。
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// set header
			testCase.setUpAuth(t, request, server.tokenMaker)

		
			server.router.ServeHTTP(record, request)
			testCase.checkResponse(t, record)
		})
	}
}

func responseBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var responseAccount db.Account
	err = json.Unmarshal(data, &responseAccount)
	require.NoError(t, err)
	require.Equal(t, responseAccount.ID, account.ID)
}

func randAccount(name string) db.Account {
	return db.Account{
        ID:       int64(utils.RandInt(1, 1000)),
		Owner :  name, 
		Balance:  utils.RandBalance(),
		Currency:    utils.RandCurrency(),
    }
}