package users_test

import (
	"context"
	"testing"
	"time"

	"github.com/henrywhitaker3/go-template/internal/test"
	"github.com/henrywhitaker3/go-template/internal/users"
	"github.com/henrywhitaker3/go-template/internal/uuid"
	"github.com/stretchr/testify/require"
)

func TestItCreatesAUser(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	type testCase struct {
		name        string
		params      users.CreateParams
		shouldError bool
	}

	email := test.Email()

	tcs := []testCase{
		{
			name: "creates a new user",
			params: users.CreateParams{
				Name:     test.Word(),
				Email:    email,
				Password: test.Word(),
			},
			shouldError: false,
		},
		{
			name: "fails to create user with existing email",
			params: users.CreateParams{
				Name:     test.Word(),
				Email:    email,
				Password: test.Word(),
			},
			shouldError: true,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			user, err := app.Users.CreateUser(ctx, c.params)
			if c.shouldError {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, c.params.Email, user.Email)
			require.Equal(t, c.params.Name, user.Name)
		})
	}
}

func TestItGetsUsersById(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, _ := test.User(t, app)

	type testCase struct {
		name        string
		user        *users.User
		shouldError bool
	}

	tcs := []testCase{
		{
			name:        "gets a valid user",
			user:        user,
			shouldError: false,
		},
		{
			name:        "errors on invalid user",
			user:        nil,
			shouldError: true,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			var id uuid.UUID
			if c.user == nil {
				id = uuid.MustNew()
			} else {
				id = c.user.ID
			}

			user, err := app.Users.Get(ctx, id)
			if c.shouldError {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.Equal(t, c.user.ID, user.ID)
			require.Equal(t, c.user.Name, user.Name)
			require.Equal(t, c.user.Email, user.Email)
		})
	}
}

func TestItGetsUsersByEmail(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, _ := test.User(t, app)

	type testCase struct {
		name        string
		user        *users.User
		shouldError bool
	}

	tcs := []testCase{
		{
			name:        "gets a valid user",
			user:        user,
			shouldError: false,
		},
		{
			name:        "errors on invalid user",
			user:        nil,
			shouldError: true,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			var email string
			if c.user == nil {
				email = test.Email()
			} else {
				email = c.user.Email
			}

			user, err := app.Users.GetByEmail(ctx, email)
			if c.shouldError {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.Equal(t, c.user.ID, user.ID)
			require.Equal(t, c.user.Name, user.Name)
			require.Equal(t, c.user.Email, user.Email)
		})
	}
}

func TestItGetsUserByLogin(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user, password := test.User(t, app)

	delUser, delPassword := test.User(t, app)
	require.Nil(t, app.Users.DeleteUser(ctx, delUser.ID))

	type testCase struct {
		name        string
		user        *users.User
		password    string
		shouldError bool
	}

	tcs := []testCase{
		{
			name:        "gets a valid user",
			user:        user,
			password:    password,
			shouldError: false,
		},
		{
			name:        "errors on invalid password",
			user:        user,
			password:    test.Sentence(5),
			shouldError: true,
		},
		{
			name:        "errors on invalid user",
			user:        nil,
			password:    test.Sentence(5),
			shouldError: true,
		},
		{
			name:        "errors on deleted user",
			user:        delUser,
			password:    delPassword,
			shouldError: true,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			var email string
			if c.user == nil {
				email = test.Email()
			} else {
				email = c.user.Email
			}

			user, err := app.Users.Login(ctx, email, c.password)
			if c.shouldError {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.Equal(t, c.user.ID, user.ID)
			require.Equal(t, c.user.Name, user.Name)
			require.Equal(t, c.user.Email, user.Email)
		})
	}
}
