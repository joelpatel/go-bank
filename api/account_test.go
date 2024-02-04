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

	"github.com/gin-gonic/gin"
	"github.com/joelpatel/go-bank/currency"
	"github.com/joelpatel/go-bank/db"
	"github.com/joelpatel/go-bank/db/mockdb"
	"github.com/joelpatel/go-bank/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type messageStruct struct {
	Message string `json:"message"`
}

type errorStruct struct {
	Error string `json:"error"`
}

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

func requireBodyMatchAccounts[A db.Account | []db.Account](t *testing.T, expected *A, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	var actual A
	err = json.Unmarshal(data, &actual)
	assert.NoError(t, err)

	assert.Equal(t, *expected, actual)
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
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccounts(t, account, recorder.Body)
}

// When invalid uri param is sent in the request, the server should respond with bad request status code.
func TestGetAccountByIDInvalidURI(t *testing.T) {
	_, _, server, recorder := beforeEach(t)

	// send request
	url := fmt.Sprintf("/account/%d", 0)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
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
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	responseBody, err := io.ReadAll(recorder.Body)
	assert.NoError(t, err)

	var response messageStruct
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Account with id %d not found.", account.ID), response.Message)
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
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	responseBody, err := io.ReadAll(recorder.Body)
	assert.NoError(t, err)

	var response errorStruct
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, sql.ErrConnDone.Error(), response.Error)
}

// When required owner and the currency is supported, the server should create a new account and return status ok with the created account.
func TestCreateAccountOK(t *testing.T) {
	account, store, server, recorder := beforeEach(t)
	account.Balance = int64(0)

	store.EXPECT().
		CreateAccount(gomock.Any(), gomock.Eq(account.Owner), gomock.Eq(int64(0)), gomock.Eq(account.Currency)).
		Times(1).
		Return(account, nil)

	body := gin.H{"owner": account.Owner, "currency": account.Currency}
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	// send request
	url := "/account/create"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccounts(t, account, recorder.Body)
}

// When the requested currency is not supported, it should respond with status bad request with apt error message.
func TestCreateAccountUnsupportedCurrency(t *testing.T) {
	account, _, server, recorder := beforeEach(t)

	body := gin.H{"owner": account.Owner, "currency": "XYZ"}
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	// send request
	url := "/account/create"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var response errorStruct
	responseBody, err := io.ReadAll(recorder.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, "XYZ is an unsupported currency.", response.Error)
}

// When the required JSON object is not present in the request, it should respond with status bad request with apt message.
func TestCreateAccountBadRequestBody(t *testing.T) {
	_, _, server, recorder := beforeEach(t)

	body := gin.H{}
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	// send request
	url := "/account/create"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

// If an error was occured during creation of new account, like DB connection terminated, it should respond with internal server error status and apt message.
func TestCreateAccountInternalServerError(t *testing.T) {
	account, store, server, recorder := beforeEach(t)

	store.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(nil, sql.ErrConnDone)

	body := gin.H{"owner": account.Owner, "currency": account.Currency}
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	// send request
	url := "/account/create"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	var response errorStruct
	responseBody, err := io.ReadAll(recorder.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, sql.ErrConnDone.Error(), response.Error)
}
