package handlers

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v2"
)

type testHTTPHandler struct {
	w http.ResponseWriter
	r *http.Request
}

func (h testHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	w.Write(body)
	return
}
func TestVerboseFunc(t *testing.T) {
	loggingFunc := verboseMiddleware(testHTTPHandler{})
	switch loggingFunc.(type) {
	case http.HandlerFunc:
		break
	default:
		t.Errorf("verboseMiddlevare should return http.HandlerFunc not %T", loggingFunc)
	}
	body := `{
		"some_data"
	}`
	r := httptest.NewRequest(http.MethodGet, "/someurl/", bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	loggingFunc.ServeHTTP(w, r)
	if w.Body.String() != body {
		t.Errorf("Can't read request body after logging func execution")
	}
}

func TestFetchUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db = sqlx.NewDb(mockDB, "postgres")
	defer db.Close()
	id := uuid.New()
	user := User{id, RequestUser{"L", "O", "L", ""}}
	secondId := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "first_name", "second_name", "phone", "email"}).
		AddRow(user.ID, user.FirstName, user.SecondName, user.Phone, user.Email)
	mock.ExpectQuery(`SELECT \* FROM users WHERE id=`).WithArgs(id).WillReturnRows(rows)
	mock.ExpectQuery(`SELECT \* FROM users WHERE id=`).WithArgs(secondId).WillReturnError(sql.ErrNoRows)
	fetched, err := fetchUser(id.String())
	if err != nil {
		t.Errorf("Can't get user data with correct id (%v). ERROR: %s", id, err)
	}
	if fetched != user {
		t.Errorf("User data was fetched incorrectly. Original user:\n %v \n Fetched user:\n %v", user, fetched)
	}
	_, err = fetchUser(secondId.String())
	if err == nil {
		t.Errorf("Incorrect user absence handling")
	}
}
func TestFetchAccount(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db = sqlx.NewDb(mockDB, "postgres")
	defer db.Close()
	id := uuid.New()
	account := Account{id, uuid.New(), RequestAccount{"", 50}}
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "balance"}).
		AddRow(account.ID, account.UserID, account.Name, account.Balance).
		AddRow(account.ID, account.UserID, account.Name, account.Balance)
	secondId := uuid.New()
	mock.ExpectQuery(`SELECT [[:ascii:]]* FROM personal_accounts WHERE id=`).WithArgs(id).WillReturnRows(rows)
	fetched, err := fetchAccount(id.String())
	if err != nil {
		t.Errorf("Can't get account data with correct id (%v). ERROR: %s", id, err)
	}
	if fetched != account {
		t.Errorf("Account data was fetched incorrectly. Original account:\n %v \n Fetched account:\n %v", account, fetched)
	}
	mock.ExpectQuery(`SELECT [[:ascii:]]* FROM personal_accounts WHERE id=`).WithArgs(secondId).WillReturnError(sql.ErrNoRows)
	_, err = fetchAccount(secondId.String())
	if err == nil {
		t.Errorf("Incorrect account absence handling")
	}
}
