package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/daniel/master-golang/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Marker,
	authorizationTypeBearer string,
	userName string,
	duration time.Duration,
) {
	token, _ , err := tokenMaker.CreateToken(userName, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s",authorizationTypeBearer, token)
	request.Header.Add(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name string
		setUpAuth func (t *testing.T, request *http.Request,tokenMaker token.Marker)
		checkResponse func (t *testing.T, request *http.Request, recorder *httptest.ResponseRecorder)
	}{
		{
			name : "ok",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer,"user", time.Minute)
			},
			checkResponse: func (t *testing.T, request *http.Request, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name : "authorization header is not provided",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
			},
			checkResponse: func (t *testing.T, request *http.Request, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name : "unsupported authorization type",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, "unsuppered", "user", time.Minute)
			},
			checkResponse: func (t *testing.T, request *http.Request, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name : "expired token",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", -time.Minute)
			},
			checkResponse: func (t *testing.T, request *http.Request, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//建立 request 然後打 , 在打之前set headers
			//不用建立store 因為不會去查詢
			server := newTestServer(t, nil)
			
			//建立route 可以用實際的 但寫一個比較快只給這個測試用
			server.router.GET(
				"/test/auth",
				authMiddleware(server.tokenMaker),
				 func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				})

			// 然後用http
			request, err := http.NewRequest(http.MethodGet, "/test/auth", nil)
			require.NoError(t, err)

			recoreder := httptest.NewRecorder()

			// set header
			testCase.setUpAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recoreder, request)
			testCase.checkResponse(t, request, recoreder)
		})
	}
}