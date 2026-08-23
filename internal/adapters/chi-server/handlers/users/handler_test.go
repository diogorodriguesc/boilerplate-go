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
	testdataDir := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata"))
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
