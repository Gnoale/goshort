package api

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

var urlToID = map[string]int64{
	"https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74": 139138397222,
	"https://github.com/go-chi/chi/tree/master":                                                         11157,
}

var idToURL = map[int64]string{
	139138397222: "https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74",
	11157:        "https://github.com/go-chi/chi/tree/master",
}

type mockedRepo struct{}

var mr = mockedRepo{}

func (m *mockedRepo) CreateURL(ctx context.Context, url string) (int64, error) {
	id, ok := urlToID[url]
	if !ok {
		return 0, sql.ErrNoRows
	}
	return id, nil
}

func (m *mockedRepo) GetURL(ctx context.Context, id int64) (string, error) {
	url, ok := idToURL[id]
	if !ok {
		return "", sql.ErrNoRows
	}
	return url, nil
}

func TestSlug(t *testing.T) {

	h := handler{&mr}

	for url := range urlToID {
		w := httptest.NewRecorder()
		b := strings.NewReader(fmt.Sprintf("{\"long_url\":\"%s\"}", url))
		req := httptest.NewRequest("POST", "http://127.0.0.1:8000/foo", b)
		h.Shorten(w, req)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 201, resp.StatusCode)
		fmt.Println(string(body))
	}
}

var testRedirect = []struct {
	slug string
	url  string
}{
	{
		"2tx",
		"https://github.com/go-chi/chi/tree/master",
	},
	{
		"2RsIXB8",
		"https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74",
	},
}

func TestRedirect(t *testing.T) {
	h := handler{&mr}
	for _, test := range testRedirect {
		rc := chi.NewRouteContext()
		rc.URLParams.Add("slug", test.slug)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://example/foo", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))
		h.Slug(w, req)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 301, resp.StatusCode)
		assert.Equal(t, test.url, resp.Header.Get("Location"))
		fmt.Println(string(body))
	}
}

var encodingTest = []struct {
	id       int64
	expected string
}{
	{
		9,
		"9",
	},
	{
		1987,
		"W3",
	},
	{
		11157,
		"2tx",
	},
	{
		139138397222,
		"2RsIXB8",
	},
}

var decodingTest = []struct {
	slug     string
	expected int64
}{
	{
		"2tx",
		11157,
	},
	{
		"W3",
		1987,
	},
	{
		"8",
		8,
	},
}

func TestEncodeDecode(t *testing.T) {
	for _, test := range encodingTest {
		v, err := encode(test.id)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, v)
	}
	for _, test := range decodingTest {
		v, err := decode(test.slug)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, v)
	}
}
