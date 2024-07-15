package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"unicode/utf8"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"

	"github.com/kondukto-io/simple-fuzzing/util"
)

var (
	tests = []struct {
		name    string
		args    User
		wantErr bool
	}{
		{
			name: "success",
			args: User{
				ID:    "1111",
				Name:  "kondukto",
				Email: "helo@kondukto.io",
				Blog:  "http://www.myblog.com",
			},
			wantErr: false,
		},
		{
			name: "success",
			args: User{
				ID:    "1112",
				Name:  "kondukto",
				Email: "helo@kondukto.io",
				Blog:  "https://myblog.com",
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: User{
				ID:    "1212121212121212121212121111",
				Name:  "kondukto",
				Email: "helo@kondukto.io",
				Blog:  "www.myblog.com",
			},
			wantErr: true,
		},
		{
			name: "fail",
			args: User{
				ID:    "s1111", // not a valid ID
				Name:  "kondukto",
				Email: "helo@kondukto.io",
				Blog:  "www.myblog.com",
			},
			wantErr: true,
		},
	}
)

func TestCreateUser(t *testing.T) {
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Run test in parallel
			t.Helper()

			// setup
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a mock db conn", err)
			}
			defer db.Close()

			mock.ExpectPrepare(regexp.QuoteMeta("INTO users(id, name, email) values (?, ?, ?, ?)"))

			h := NewHandler(db)

			body, err := json.Marshal(tt.args)
			if err != nil {
				t.Fatalf("error %v", err)
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/create")

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(id, name, email) values (?, ?, ?, ?)")).
				WithArgs(tt.args.ID, tt.args.Name, tt.args.Email, tt.args.Blog).WillReturnResult(sqlmock.NewResult(1, 1))

			// testing the function
			if err := h.CreateUser(c); err != nil {
				t.Errorf("CreateUser() err = %v, wantErr %v", err, tt.wantErr)
			}

			// ensure all expectations have been met
			if err = mock.ExpectationsWereMet(); err != nil {
				fmt.Printf("unmet expectation error: %s", err)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Run test in parallel
			t.Parallel()
			t.Helper()

			// setup
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a mock db conn", err)
			}
			defer db.Close()

			mock.ExpectPrepare("SELECT (.+) FROM users WHERE id=?")

			h := NewHandler(db)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/user/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.args.ID)

			rows := sqlmock.NewRows([]string{"id", "name", "email", "blog"}).
				AddRow(tt.args.ID, tt.args.Name, tt.args.Email, tt.args.Blog)

			mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
				WithArgs(tt.args.ID).
				WillReturnRows(rows)

			// testing the function
			if err := h.GetUserByID(c); (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() err = %v, wantErr %v", err, tt.wantErr)
			}

			//if assert.NoError(t, h.GetUserByID(c)) {
			//	assert.Equal(t, http.StatusOK, rec.Code)

			//	user := new(User)
			//	json.Unmarshal([]byte(rec.Body.String()), user)

			//	if user.ID != tt.args.Name {
			//		t.Fatal("hop")
			//	}
			//}
		})
	}
}

func FuzzGetUserByID(f *testing.F) {
	// setup
	db, mock, err := sqlmock.New()
	if err != nil {
		f.Fatalf("an error '%s' was not expected when opening a mock db conn", err)
	}
	defer db.Close()

	for _, tt := range tests {
		f.Add(tt.args.ID)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		if !util.VaildID(orig) {
			return
		}

		mock.ExpectPrepare("SELECT (.+) FROM users WHERE id=?")
		h := NewHandler(db)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		t.Log("\t ==== value is:", orig)

		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues(orig)

		rows := sqlmock.NewRows([]string{"id", "name", "email", "blog"}).
			AddRow(orig, "kondukto", "test@kondukto.io", "myblog.com")

		mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").
			WithArgs(orig).
			WillReturnRows(rows)

		if err := h.GetUserByID(c); err != nil {
			t.Fatalf("error occured '%s' and Value: %s", err, orig)
		}

		user := new(User)
		json.Unmarshal([]byte(rec.Body.String()), user)

		expected := &User{
			ID:    orig,
			Name:  "kondukto",
			Email: "test@kondukto.io",
			Blog:  "myblog.com",
		}
		if user.ID != orig {
			t.Fatalf("test failed expected %v -- got %v", expected, user)
		}
	})
}

func FuzzCreateUser(f *testing.F) {
	// setup
	db, mock, err := sqlmock.New()
	if err != nil {
		f.Fatalf("an error '%s' was not expected when opening a mock db conn", err)
	}
	defer db.Close()

	for _, tt := range tests {
		f.Add(tt.args.ID, tt.args.Name, tt.args.Email, tt.args.Blog)
	}

	f.Fuzz(func(t *testing.T, id, name, email, blog string) {
		if !util.VaildID(id) || !utf8.ValidString(name) || !utf8.ValidString(email) || !util.ValidURL(blog) {
			return
		}

		mock.ExpectPrepare(regexp.QuoteMeta("INTO users(id, name, email) values (?, ?, ?, ?)"))

		h := NewHandler(db)
		input := User{
			ID:    id,
			Name:  name,
			Email: email,
			Blog:  blog,
			//Blog: "https://myblog.com",
		}

		t.Log(input)

		body, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("error %v", err)
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/create")

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(id, name, email) values (?, ?, ?, ?)")).
			WithArgs(input.ID, input.Name, input.Email, input.Blog).WillReturnResult(sqlmock.NewResult(1, 1))

		// testing the function
		if err := h.CreateUser(c); err != nil {
			t.Errorf("CreateUser() err = %v", err)
		}

		// ensure all expectations have been met
		if err = mock.ExpectationsWereMet(); err != nil {
			fmt.Printf("unmet expectation error: %s", err)
		}
	})
}
