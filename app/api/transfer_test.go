package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/daniel/master-golang/db/mock"
	"github.com/daniel/master-golang/db/sqlc"
	"github.com/daniel/master-golang/token"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestServer_Transfer(t *testing.T) {
	amount := int(10)
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	// user3, _ := randomUser(t)

	account1 := randAccount(user1.Name)
	account2:= randAccount(user2.Name)
	// account3 := randAccount(user3.Name)

	// test currency USD
	account1.Currency = "USD"
	account2.Currency = "USD"


	testCases := []struct {
		Name string
		setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Marker)
		request map[string]interface{}
		buildStubs func (*mockdb.MockStore)
		checkResponse func(record *httptest.ResponseRecorder)
	}{
		{
			Name: "ok",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account1.Owner, time.Minute)
			},
			request: map[string]interface{}{
				"from_account_id" : account1.ID,
				"to_account_id" : account2.ID,
				"amount" : amount,
				"currency" : "USD",
			},
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				
				ms.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: int64(amount),
				}

				ms.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse : func(record *httptest.ResponseRecorder) {
				// fmt.Println(record.Code)
				require.Equal(t, record.Code, http.StatusOK)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.Name, func(t *testing.T) {

			controller := gomock.NewController(t)
			store := mockdb.NewMockStore(controller)
			testcase.buildStubs(store)

			server := newTestServer(t, store)

			data, err := json.Marshal(testcase.request)
			require.NoError(t, err)

			
			request, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(data))
			require.NoError(t, err)
			
			record := httptest.NewRecorder()
			
			testcase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(record, request)
		})
	}
}
