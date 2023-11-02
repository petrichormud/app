package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	html "github.com/gofiber/template/html/v2"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/sessiondata"
	"petrichormud.com/app/internal/shared"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("SET GLOBAL local_infile=true;")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
		Protocol: 3,
	})

	// TODO: Update this config to be more secure. Will depend on environment.
	s := session.New()

	i := shared.InterfacesBuilder().Database(db).Redis(r).Sessions(s).Build()

	views := html.NewFileSystem(http.FS(viewsfs), ".html")
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	readTimeout := time.Second * time.Duration(readTimeoutSecondsCount)
	config := fiber.Config{
		Views:       views,
		ViewsLayout: "web/views/layouts/main",
		ReadTimeout: readTimeout,
	}
	app := fiber.New(config)

	app.Use(logger.New())
	app.Use(csrf.New(
		csrf.Config{
			KeyLookup:         "header:X-CSRF-Token",
			CookieName:        "csrf_",
			CookieSameSite:    "Lax",
			CookieSessionOnly: true,
			CookieHTTPOnly:    true,
			Expiration:        1 * time.Hour,
			KeyGenerator:      utils.UUIDv4,
			Session:           s,
			SessionKey:        "fiber.csrf.token",
			ContextKey:        "csrf",
			HandlerContextKey: "fiber.csrf.handler",
		},
	))
	app.Use(sessiondata.New(&i))
	app.Use(bind.New())

	app.Static("/", "./web/static")
	app.Static("/loaders", "./web/svg/loaders")

	app.Get("/", handlers.Home())

	app.Post("/login", handlers.Login(&i))
	app.Post("/logout", handlers.Logout(&i))
	app.Get("/logout", handlers.LogoutPage())

	player := app.Group("/player")
	player.Post("/new", handlers.CreatePlayer(&i))
	player.Post("/reserved", handlers.UsernameReserved(&i))

	email := player.Group("/email")
	email.Post("/new", handlers.AddEmail(&i))
	email.Delete("/:id", handlers.DeleteEmail(&i))
	email.Put("/:id", handlers.EditEmail(&i))
	email.Post("/:id/resend", handlers.ResendEmailVerification(&i))

	// TODO: Move this behind the email group
	app.Get("/verify", handlers.Verify(&i))
	app.Post("/verify", handlers.VerifyEmail(&i))

	app.Get("/profile", handlers.Profile(&i))

	log.Fatal(app.Listen(":8008"))
}
