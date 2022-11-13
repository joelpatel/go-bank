package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/joelpatel/go-bank/db/mock"
	db "github.com/joelpatel/go-bank/db/sqlc"
	"github.com/joelpatel/go-bank/util"
	"github.com/stretchr/testify/require"
)

func RandomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func TestGetAccountAPI(t *testing.T) {
	account := RandomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(mockStore *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(mockStore *mockdb.MockStore) {
				// building stubs
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					// meaning: expect GetAccount function to be called with any context and this specific account id argument ✅
					Times(1).
					Return(account, nil) // A-Okay
				// now the mock store is built
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, &account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(mockStore *mockdb.MockStore) {
				// building stubs
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					// meaning: expect GetAccount function to be called with any context and this specific account id argument ✅
					Times(1).
					Return(db.Account{}, sql.ErrNoRows) // no rows found
				// now the mock store is built
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(mockStore *mockdb.MockStore) {
				// building stubs
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					// meaning: expect GetAccount function to be called with any context and this specific account id argument ✅
					Times(1).
					Return(db.Account{}, sql.ErrConnDone) // say connection terminated
				// now the mock store is built
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(mockStore *mockdb.MockStore) {
				// building stubs
				mockStore.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()). // user sent invalid request
					Times(0)
				// now the mock store is built
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStubs(mockStore)

			// start test http server and send requests
			server := NewServer(mockStore)
			recorder := httptest.NewRecorder() // create a new response recorder instead of starting a server

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request) // send http request through server router and copy its response in recorder
			// check response
			tc.checkResponse(t, recorder)
		})
	}

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account *db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, gotAccount, *account)
}
