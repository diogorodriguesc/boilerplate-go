//go:build functional

package users_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	chiserver "github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server"
	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/handlers/users/responses"
	"github.com/diogorodriguesc/boilerplate-go/internal/application/api"
	"github.com/diogorodriguesc/boilerplate-go/tests"
)

const (
	postgresFixtureReset    = "postgres/000_reset.sql"
	postgresFixtureAddUsers = "postgres/100_add_users.sql"
)

var env *tests.PostgresEnvironment

func TestMain(m *testing.M) {
	os.Exit(tests.RunFunctionalMain(m, func(e *tests.PostgresEnvironment) {
		env = e
	}))
}

func seedUsers(t *testing.T) {
	t.Helper()
	require.NoError(t, env.ApplyFixture(fixtureAbsPath(postgresFixtureReset)))
	require.NoError(t, env.ApplyFixture(fixtureAbsPath(postgresFixtureAddUsers)))
}

func fixtureAbsPath(relative string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	testdataDir := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "testdata"))
	return filepath.Join(testdataDir, relative)
}

func performRequest(t *testing.T, email string) *httptest.ResponseRecorder {
	t.Helper()

	body := fmt.Sprintf(`{"email": %q}`, email)
	req := httptest.NewRequest(http.MethodPost, "/v1/users/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	application, _, err := api.NewApplication(t.Context(), env.Storage)
	if err != nil {
		t.Fatal(err)
	}

	httpServer := chiserver.NewHttpServer(t.Context(), application)

	recorder := httptest.NewRecorder()
	httpServer.SetRouter().ServeHTTP(recorder, req)
	return recorder
}

func performCreateUserRequest(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	application, _, err := api.NewApplication(t.Context(), env.Storage)
	if err != nil {
		t.Fatal(err)
	}

	httpServer := chiserver.NewHttpServer(t.Context(), application)

	recorder := httptest.NewRecorder()
	httpServer.SetRouter().ServeHTTP(recorder, req)
	return recorder
}

func performGetUserRequest(t *testing.T, id string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+id, nil)

	application, _, err := api.NewApplication(t.Context(), env.Storage)
	if err != nil {
		t.Fatal(err)
	}

	httpServer := chiserver.NewHttpServer(t.Context(), application)

	recorder := httptest.NewRecorder()
	httpServer.SetRouter().ServeHTTP(recorder, req)
	return recorder
}

func performDeleteUserRequest(t *testing.T, id string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodDelete, "/v1/users/"+id, nil)

	application, _, err := api.NewApplication(t.Context(), env.Storage)
	if err != nil {
		t.Fatal(err)
	}

	httpServer := chiserver.NewHttpServer(t.Context(), application)

	recorder := httptest.NewRecorder()
	httpServer.SetRouter().ServeHTTP(recorder, req)
	return recorder
}

func performListUsersRequest(t *testing.T, query string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/v1/users"+query, nil)

	application, _, err := api.NewApplication(t.Context(), env.Storage)
	if err != nil {
		t.Fatal(err)
	}

	httpServer := chiserver.NewHttpServer(t.Context(), application)

	recorder := httptest.NewRecorder()
	httpServer.SetRouter().ServeHTTP(recorder, req)
	return recorder
}

func TestHandler_Functional_ListUsers_Defaults(t *testing.T) {
	seedUsers(t)
	recorder := performListUsersRequest(t, "")
	require.Equal(t, http.StatusOK, recorder.Code)

	var response responses.ListUsersResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, 1, response.Page)
	require.Equal(t, 20, response.PageSize)
	require.EqualValues(t, 2, response.Total)
	require.Len(t, response.Data, 2)
}

func TestHandler_Functional_ListUsers_Paginated(t *testing.T) {
	seedUsers(t)

	recorder := performListUsersRequest(t, "?page=1&pageSize=1")
	require.Equal(t, http.StatusOK, recorder.Code)

	var page1 responses.ListUsersResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &page1))
	require.Equal(t, 1, page1.Page)
	require.Equal(t, 1, page1.PageSize)
	require.EqualValues(t, 2, page1.Total)
	require.Len(t, page1.Data, 1)

	recorder = performListUsersRequest(t, "?page=2&pageSize=1")
	require.Equal(t, http.StatusOK, recorder.Code)

	var page2 responses.ListUsersResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &page2))
	require.Len(t, page2.Data, 1)
	require.NotEqual(t, page1.Data[0].ID, page2.Data[0].ID)

	recorder = performListUsersRequest(t, "?page=3&pageSize=1")
	require.Equal(t, http.StatusOK, recorder.Code)

	var page3 responses.ListUsersResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &page3))
	require.Empty(t, page3.Data)
}

func TestHandler_Functional_ListUsers_InvalidPage(t *testing.T) {
	seedUsers(t)
	recorder := performListUsersRequest(t, "?page=0")
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Functional_ListUsers_InvalidPageSize(t *testing.T) {
	seedUsers(t)
	recorder := performListUsersRequest(t, "?pageSize=101")
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Functional_ListUsers_NonIntegerPage(t *testing.T) {
	seedUsers(t)
	recorder := performListUsersRequest(t, "?page=abc")
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Functional_GetUser(t *testing.T) {
	seedUsers(t)
	recorder := performGetUserRequest(t, "1")
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{
	"id": 1,
	"username": "foo",
	"email": "foo@gmail.com"
	}`, recorder.Body.String())
}

func TestHandler_Functional_GetUser_NotFound(t *testing.T) {
	seedUsers(t)
	recorder := performGetUserRequest(t, "999")
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestHandler_Functional_GetUser_InvalidID(t *testing.T) {
	seedUsers(t)
	recorder := performGetUserRequest(t, "not-a-number")
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestHandler_Functional_DeleteUser(t *testing.T) {
	seedUsers(t)
	recorder := performDeleteUserRequest(t, "1")
	require.Equal(t, http.StatusNoContent, recorder.Code)

	getRecorder := performGetUserRequest(t, "1")
	require.Equal(t, http.StatusNotFound, getRecorder.Code)
}

func TestHandler_Functional_DeleteUser_NotFound(t *testing.T) {
	seedUsers(t)
	recorder := performDeleteUserRequest(t, "999")
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestHandler_Functional_DeleteUser_InvalidID(t *testing.T) {
	seedUsers(t)
	recorder := performDeleteUserRequest(t, "not-a-number")
	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestHandler_Functional_SearchUser(t *testing.T) {
	seedUsers(t)
	recorder := performRequest(t, "foo@gmail.com")
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `[{
	"id": 1,
	"username": "foo",
	"email": "foo@gmail.com"
	}]`, recorder.Body.String())
}

func TestHandler_Functional_SearchUser_NoMatch(t *testing.T) {
	seedUsers(t)
	recorder := performRequest(t, "missing@gmail.com")
	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `[]`, recorder.Body.String())
}

func TestHandler_Functional_CreateUser(t *testing.T) {
	seedUsers(t)
	recorder := performCreateUserRequest(t, `{"username": "baz", "email": "baz@gmail.com"}`)
	require.Equal(t, http.StatusCreated, recorder.Code)

	var response map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, "baz", response["username"])
	require.Equal(t, "baz@gmail.com", response["email"])
}

func TestHandler_Functional_CreateUser_DuplicateEmail(t *testing.T) {
	seedUsers(t)
	recorder := performCreateUserRequest(t, `{"username": "foo2", "email": "foo@gmail.com"}`)
	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestHandler_Functional_CreateUser_InvalidPayload(t *testing.T) {
	seedUsers(t)
	recorder := performCreateUserRequest(t, `{"username": "", "email": "not-an-email"}`)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
