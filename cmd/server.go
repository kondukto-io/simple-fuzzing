package cmd

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kondukto-io/simple-fuzzing/handlers"
)

const (
	port = ":8888"
)

func Execute() error {
	//dbPath := filepath.Join(os.TempDir(), "db.db")
	//defer os.RemoveAll(dbPath)

	//db, err := sql.Open("sqlite3", dbPath)
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()
	//e.HideBanner = true
	//e.HidePort = true

	// middlewares
	e.Use(middleware.Logger())

	err = handlers.MigrateDB(db)
	if err != nil {
		panic(err)
	}

	// Initialize handler
	h := handlers.NewHandler(db)

	// Routes
	e.POST("/create", h.CreateUser)
	e.GET("/user/:id", h.GetUserByID)

	return e.Start(port)
}
