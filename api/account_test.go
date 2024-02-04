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

func randomAccount() *db.Account {
	return &db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: currency.USD,
	}
}

func randomAccounts(n int) *[]db.Account {
	accounts := make([]db.Account, n)
	owner := utils.RandomOwner()
	for i := 0; i < n; i++ {
		accounts[i] = db.Account{
			ID:       utils.RandomInt(1, 1000),
			Owner:    owner,
			Balance:  utils.RandomMoney(),
			Currency: currency.USD,
		}
	}
	return &accounts
}

func beforeEach(t *testing.T) (*mockdb.MockStore, *Server, *httptest.ResponseRecorder) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	server := NewServer(store)

	recorder := httptest.NewRecorder()

	return store, server, recorder
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
	store, server, recorder := beforeEach(t)
	account := randomAccount()

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
	_, server, recorder := beforeEach(t)

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
	store, server, recorder := beforeEach(t)
	account := randomAccount()

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
	store, server, recorder := beforeEach(t)
	account := randomAccount()

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
	store, server, recorder := beforeEach(t)
	account := randomAccount()
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
	_, server, recorder := beforeEach(t)
	account := randomAccount()

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
	_, server, recorder := beforeEach(t)

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
	store, server, recorder := beforeEach(t)
	account := randomAccount()

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

// When correct query parameter and account owner is passed, it should return accounts taking into consideration pagition limit and size.
func TestListAccountsByOwnerOK(t *testing.T) {
	store, server, recorder := beforeEach(t)
	accounts := randomAccounts(5)
	owner := (*accounts)[0].Owner

	// build stubs
	store.EXPECT().
		ListAccounts(gomock.Any(), gomock.Eq(owner), gomock.Eq(int64(5)), gomock.Eq(int64(0))).
		Times(1).
		Return(accounts, nil)

	// build & send request
	body := gin.H{"owner": owner}
	data, err := json.Marshal(body)
	assert.NoError(t, err)
	url := "/accounts"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	assert.NoError(t, err)
	q := request.URL.Query()
	q.Add("page_id", "1")
	q.Add("page_size", "5")
	request.URL.RawQuery = q.Encode()
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccounts[[]db.Account](t, accounts, recorder.Body)
}

// When owner data is not provided, it should respond with status code of bad request.
func TestListAccountsByOwnerInvalidOwner(t *testing.T) {
	_, server, recorder := beforeEach(t)

	// build & send request
	url := "/accounts"
	request, err := http.NewRequest(http.MethodPost, url, nil)
	assert.NoError(t, err)
	q := request.URL.Query()
	q.Add("page_id", "1")
	q.Add("page_size", "5")
	request.URL.RawQuery = q.Encode()
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

// When page_id or page_size is not provided or invalid in the reqest query parameter. it should respond with status code of bad request.
func TestListAccountsByOwnerBadQueryParam(t *testing.T) {
	_, server, recorder := beforeEach(t)

	// build & send request
	body := gin.H{"owner": utils.RandomOwner()}
	data, err := json.Marshal(body)
	assert.NoError(t, err)
	url := "/accounts"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// with only page_id
	q := request.URL.Query()
	q.Add("page_id", "1")
	request.URL.RawQuery = q.Encode()
	server.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// with only page_size
	q.Del("page_id")
	q.Add("page_size", "5")
	request.URL.RawQuery = q.Encode()
	server.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// with invalid query parameter
	q.Add("page_id", "-1")
	q.Set("page_size", "1000")
	request.URL.RawQuery = q.Encode()
	server.router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

// When any internal server occurs like connection to DB terminated, then it should respond with status internal server error.
func TestListAccountsByOwnerInterServerError(t *testing.T) {
	store, server, recorder := beforeEach(t)

	// build stubs
	store.EXPECT().
		ListAccounts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(nil, sql.ErrConnDone)

	// build & send request
	body := gin.H{"owner": utils.RandomOwner()}
	data, err := json.Marshal(body)
	assert.NoError(t, err)
	url := "/accounts"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	assert.NoError(t, err)
	q := request.URL.Query()
	q.Add("page_id", "1")
	q.Add("page_size", "5")
	request.URL.RawQuery = q.Encode()
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

// If there are no account records for the request owner, then it should respond with status not found with apt message.
func TestListAccountsByOwnerNoRecords(t *testing.T) {
	store, server, recorder := beforeEach(t)

	// build stubs
	store.EXPECT().
		ListAccounts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(&[]db.Account{}, nil)

	// build & send request
	body := gin.H{"owner": utils.RandomOwner()}
	data, err := json.Marshal(body)
	assert.NoError(t, err)
	url := "/accounts"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	assert.NoError(t, err)
	q := request.URL.Query()
	q.Add("page_id", "1")
	q.Add("page_size", "5")
	request.URL.RawQuery = q.Encode()
	server.router.ServeHTTP(recorder, request)

	// check response
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	var response messageStruct
	responseBody, err := io.ReadAll(recorder.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, "Accounts from entry 0 not found.", response.Message)
}
