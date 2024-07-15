package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kondukto-io/simple-fuzzing/util"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Blog  string `json:"blog"`
}

func (h *Handler) CreateUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		// in the production you should not dump the error message directly
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	stmt, err := h.db.Prepare("INSERT INTO users(id, name, email) values (?, ?, ?, ?)")
	if err != nil {
		// in the production you should not dump the error message directly
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.ID, u.Name, u.Email, u.Blog)
	if err != nil {
		// in the production you should not dump the error message directly
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	return c.JSON(http.StatusOK, u)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	cid := c.Param("id")

	if !util.VaildID(cid) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Not a valid ID"}
	}

	stmt, err := h.db.Prepare("SELECT * FROM users WHERE id=?")
	if err != nil {
		// in the production you should not dump the error message directly
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	defer stmt.Close()

	var id, name, email, blog string

	err = stmt.QueryRow(cid).Scan(&id, &name, &email, &blog)
	if err != nil {
		// in the production you should not dump the error message directly
		return &echo.HTTPError{Code: http.StatusNotFound, Message: err.Error()}
	}

	return c.JSON(http.StatusOK, User{ID: id, Name: name, Email: email, Blog: blog})
}
