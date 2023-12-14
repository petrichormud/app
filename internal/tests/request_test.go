package tests

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/shared"
)

func TestCreateCommentUnauthorized(t *testing.T) {
	// TODO: Re-implement this when the routes are plugged in
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCreateCommentSuccess(t *testing.T) {
	// TODO: Re-implement this when the routes are plugged in
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestCreateCommentNotFound(t *testing.T) {
	// TODO: Re-implement this when the routes are plugged in
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid+1)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCreateCommentFatal(t *testing.T) {
	// TODO: Re-implement this when the routes are plugged in
	t.Skip()

	i := shared.SetupInterfaces()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)

	i.Close()

	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid+1)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCreateFieldCommentUnauthorized(t *testing.T) {
	// TODO: Re-implement this when the routes are plugged in
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("field", "name")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCreateFieldCommentSuccess(t *testing.T) {
	// TODO: Re-implement this when the routes are plugged in
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("field", "name")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestCreateFieldCommentNotFound(t *testing.T) {
	// TODO: Re-implement this when the routes are plugged in
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid+1)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("field", "name")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCreateFieldCommentFatal(t *testing.T) {
	// TODO: Re-implement this when the routes are plugged in
	t.Skip()

	i := shared.SetupInterfaces()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, sessionCookie := CreateTestPlayerAndCharacterApplication(t, &i, a)

	i.Close()

	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("field", "name")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func TestCreateReplyCommentUnauthorized(t *testing.T) {
	// TODO: Figure out why this test isn't working
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	cid, _ := CreateTestComment(t, &i, a, rid)
	strcid := strconv.FormatInt(cid, 10)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("cid", strcid)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
}

func TestCreateReplyCommentSuccess(t *testing.T) {
	// TODO: Figure out why this test isn't working
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	cid, sessionCookie := CreateTestComment(t, &i, a, rid)
	strcid := strconv.FormatInt(cid, 10)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("cid", strcid)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusCreated, res.StatusCode)
}

func TestCreateReplyCommentNotFound(t *testing.T) {
	// TODO: Figure out why this test isn't working
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	cid, sessionCookie := CreateTestComment(t, &i, a, rid)
	strcid := strconv.FormatInt(cid, 10)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid+1)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("cid", strcid)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCreateReplyCommentCIDNotFound(t *testing.T) {
	// TODO: Figure out why this test isn't working
	t.Skip()

	i := shared.SetupInterfaces()
	defer i.Close()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	cid, sessionCookie := CreateTestComment(t, &i, a, rid)
	strcid := strconv.FormatInt(cid+1, 10)
	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("cid", strcid)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusNotFound, res.StatusCode)
}

func TestCreateReplyCommentFatal(t *testing.T) {
	// TODO: Figure out why this test isn't working
	t.Skip()

	i := shared.SetupInterfaces()

	a := fiber.New(configs.Fiber())
	app.Middleware(a, &i)
	app.Handlers(a, &i)

	rid, _ := CreateTestPlayerAndCharacterApplication(t, &i, a)
	cid, sessionCookie := CreateTestComment(t, &i, a, rid)
	strcid := strconv.FormatInt(cid, 10)

	i.Close()

	// TODO: Get this in a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid+1)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.WriteField("cid", strcid)
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	res, err := a.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
}

func SetupTestRequests(t *testing.T, i *shared.Interfaces, u string) {
	p, err := i.Queries.GetPlayerByUsername(context.Background(), TestUsername)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}

	pid := p.ID
	reqs, err := i.Queries.ListCharacterApplicationsForPlayer(context.Background(), pid)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Delete requests of all types from this player
	for _, req := range reqs {
		query := fmt.Sprintf("DELETE FROM character_application_content WHERE rid = %d;", req.Request.ID)
		_, err := i.Database.Exec(query)
		if err != nil {
			t.Fatal(err)
		}

		query = fmt.Sprintf("DELETE FROM character_application_content_history WHERE rid = %d;", req.Request.ID)
		_, err = i.Database.Exec(query)
		if err != nil {
			t.Fatal(err)
		}

		comments, err := i.Queries.ListCommentsForRequest(context.Background(), req.Request.ID)
		if err != nil {
			t.Fatal(err)
		}

		for _, comment := range comments {
			replies, err := i.Queries.ListRepliesToComment(context.Background(), comment.ID)
			if err != nil {
				t.Fatal(err)
			}

			for _, reply := range replies {
				query = fmt.Sprintf("DELETE FROM request_comment_content_history WHERE cid = %d;", reply.ID)
				_, err = i.Database.Exec(query)
				if err != nil {
					t.Fatal(err)
				}
			}

			query = fmt.Sprintf("DELETE FROM request_comments WHERE cid = %d;", comment.ID)
			_, err = i.Database.Exec(query)
			if err != nil {
				t.Fatal(err)
			}

			query = fmt.Sprintf("DELETE FROM request_comment_content_history WHERE cid = %d;", comment.ID)
			_, err = i.Database.Exec(query)
			if err != nil {
				t.Fatal(err)
			}
		}

		query = fmt.Sprintf("DELETE FROM request_comments WHERE rid = %d;", req.Request.ID)
		_, err = i.Database.Exec(query)
		if err != nil {
			t.Fatal(err)
		}
	}

	query := fmt.Sprintf("DELETE FROM requests WHERE pid = %d;", pid)
	_, err = i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}

	query = fmt.Sprintf("DELETE FROM players WHERE username = '%s';", u)
	_, err = i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func NewCommentRequest(rid int64) *http.Request {
	// TODO: Get this into a generator
	url := fmt.Sprintf("%s/request/%d/comments/new", TestURL, rid)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("text", "Test comment.")
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func CreateTestComment(t *testing.T, i *shared.Interfaces, app *fiber.App, rid int64) (int64, *http.Cookie) {
	res := CallLogin(t, app, TestUsername, TestPassword)
	cookies := res.Cookies()
	sessionCookie := cookies[0]
	req := NewCommentRequest(rid)
	req.AddCookie(sessionCookie)
	_, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	comments, err := i.Queries.ListCommentsForRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}

	c := comments[0]
	return c.ID, sessionCookie
}
