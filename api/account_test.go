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

	"github.com/joelpatel/go-bank/currency"
	"github.com/joelpatel/go-bank/db"
	"github.com/joelpatel/go-bank/db/mockdb"
	"github.com/joelpatel/go-bank/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func beforeEach(t *testing.T) (*db.Account, *mockdb.MockStore, *Server, *httptest.ResponseRecorder) {
	account := &db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: currency.USD,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	server := NewServer(store)

	recorder := httptest.NewRecorder()

	return account, store, server, recorder
}

func requireBodyMatchAccount(t *testing.T, expectedAccount *db.Account, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var actualAccount db.Account
	err = json.Unmarshal(data, &actualAccount)
	require.NoError(t, err)

	require.Equal(t, *expectedAccount, actualAccount)
}

// API should return status OK and with the account associated with the account ID.
func TestGetAccountByIDOK(t *testing.T) {
	account, store, server, recorder := beforeEach(t)

	// build stubs
	store.EXPECT().
		GetAccountByID(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// send request
	url := fmt.Sprintf("/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, account, recorder.Body)
}

// When invalid uri param is sent in the request, the server should respond with bad request status code.
func TestGetAccountByIDInvalidURI(t *testing.T) {
	_, _, server, recorder := beforeEach(t)

	// send request
	url := fmt.Sprintf("/account/%d", 0)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

// When the request account is not present in the databse, the server should return a status not found with apt message.
func TestGetAccocuntByIDNotFound(t *testing.T) {
	account, store, server, recorder := beforeEach(t)

	store.EXPECT().
		GetAccountByID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(nil, sql.ErrNoRows)

	// send request
	url := fmt.Sprintf("/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	require.Equal(t, http.StatusNotFound, recorder.Code)
	bodyBytes, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)

	var body struct {
		Message string `json:"message"`
	}
	err = json.Unmarshal(bodyBytes, &body)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("Account with id %d not found.", account.ID), body.Message)
}

// When a random error like connection lose with the database occurs, it should return internal server error status code with apt message.
func TestGetAccountByIDInternalServerError(t *testing.T) {
	account, store, server, recorder := beforeEach(t)

	store.EXPECT().
		GetAccountByID(gomock.Any(), gomock.Any()).
		Times(1).
		Return(nil, sql.ErrConnDone)

	// send request
	url := fmt.Sprintf("/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	bodyBytes, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)

	var body struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(bodyBytes, &body)
	require.NoError(t, err)
	require.Equal(t, sql.ErrConnDone.Error(), body.Error)
}
