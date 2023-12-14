package tests

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func TestReservedConflict(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	config := configs.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestReserved(&i, t)
	defer SetupTestReserved(&i, t)

	_ = CallRegister(t, a, TestUsername, TestPassword)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.Close()

	url := MakeTestURL(routes.Reserved)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusConflict, res.StatusCode)
}

func TestReservedFatal(t *testing.T) {
	i := shared.SetupInterfaces()

	config := configs.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestReserved(&i, t)

	_ = CallRegister(t, a, TestUsername, TestPassword)

	i.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.Close()

	url := MakeTestURL(routes.Reserved)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	it := shared.SetupInterfaces()
	SetupTestReserved(&it, t)
	it.Close()

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestReservedOK(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	config := configs.Fiber()
	a := fiber.New(config)
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	SetupTestReserved(&i, t)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("username", TestUsername)
	writer.Close()

	url := MakeTestURL(routes.Reserved)
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusOK, res.StatusCode)
}

func SetupTestReserved(i *shared.Interfaces, t *testing.T) {
	query := fmt.Sprintf("DELETE FROM players WHERE username = '%s';", TestUsername)
	_, err := i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
