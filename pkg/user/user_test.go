package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type StubRepository struct {
	users     map[uuid.UUID]User
	addedUser User
	testId    uuid.UUID
}

func (u *StubRepository) GetUserById(ctx context.Context, id uuid.UUID) *User {

	user := u.users[id]
	return &user
}

func (u *StubRepository) AddUser(ctx context.Context, user *User) {
	user.Id = u.testId
	u.addedUser = *user
}

func createTestServer(t *testing.T, configRoute func(router *gin.Engine), request func() (*http.Request, error)) *httptest.ResponseRecorder {
	t.Helper()
	router := gin.Default()

	configRoute(router)

	w := httptest.NewRecorder()
	req, err := request()

	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	return w
}

func TestGetUser(t *testing.T) {
	testId := uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")

	routeConfig := func(router *gin.Engine) {
		repo := &StubRepository{
			users: map[uuid.UUID]User{testId: {Id: testId, Name: "testuser"}},
		}

		userRoute := UserRoute{repository: repo}
		userRoute.AddRouter(router)

	}

	t.Run("GET User", func(t *testing.T) {

		w := createTestServer(
			t,
			routeConfig,
			func() (*http.Request, error) {
				request, err := http.NewRequest("GET", "/user?id="+testId.String(), nil)
				return request, err
			})
		expect := &User{Id: testId, Name: "testuser"}

		actual := &User{}
		err := json.Unmarshal(w.Body.Bytes(), actual)
		assert.NoError(t, err)

		assert.Equal(t, expect, actual)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get user without id", func(t *testing.T) {
		w := createTestServer(t,
			routeConfig,
			func() (*http.Request, error) {
				return http.NewRequest("GET", "/user", nil)
			})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAddUser(t *testing.T) {
	testId := uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
	routeConfig := func(router *gin.Engine) {
		repo := &StubRepository{
			testId: testId,
		}

		userRoute := UserRoute{repository: repo}
		userRoute.AddRouter(router)

	}
	t.Run("Add user", func(t *testing.T) {
		user := &User{
			Name:            "Berra",
			Role:            "QA",
			YearsExperience: 42,
			Keywords: []string{
				"Pro",
				"Golang"},
			Certificates: []Certificate{
				{
					Name:   "Golang",
					Issuer: "Google",
				},
			},
		}

		w := createTestServer(
			t,
			routeConfig,
			func() (*http.Request, error) {
				data, err := json.Marshal(user)
				assert.NoError(t, err)
				return http.NewRequest("POST", "/user", bytes.NewReader(data))
			})

		actual := &User{}

		expected := user
		expected.Id = testId
		err := json.Unmarshal(w.Body.Bytes(), actual)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("Add nil user", func(t *testing.T) {
		w := createTestServer(
			t,
			routeConfig,
			func() (*http.Request, error) {
				return http.NewRequest("POST", "/user", nil)
			})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	testId := uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")

	user := User{
		Id:              testId,
		Name:            "PUT",
		Role:            "Gurksvarvare",
		YearsExperience: 42,
		Keywords: []string{
			"Svarvare",
			"CNC",
			"Botanik",
		},
		Certificates: []Certificate{
			{
				Name: "B-körkort",
			},
		},
	}

	t.Run("Update existing user", func(t *testing.T) {

		routeConfig := func(router *gin.Engine) {
			repo := &StubRepository{
				testId: testId,
				users:  map[uuid.UUID]User{testId: user},
			}

			userRoute := UserRoute{repository: repo}
			userRoute.AddRouter(router)

		}
		expected := user
		expected.Role = "Tester"

		w := createTestServer(t,
			routeConfig,
			func() (*http.Request, error) {

				data, err := json.Marshal(&expected)
				assert.NoError(t, err)
				req, err := http.NewRequest("PUT", "/user", bytes.NewReader(data))
				return req, err
			})

		actual := &User{}

		json.Unmarshal(w.Body.Bytes(), actual)
		assert.Equal(t, expected, actual)
	})
}
