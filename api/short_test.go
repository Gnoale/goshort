package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sample = map[string]int64{
	"https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74": 139138397222,
	"https://github.com/go-chi/chi/tree/master":                                                         11157,
}

type mockedRepo struct {
	sample map[string]int64
}

var mr = mockedRepo{sample}

func (m *mockedRepo) CreateURL(ctx context.Context, url string) (int64, error) {
	i := m.sample[url]
	if i == 0 {
		return 0, errors.New("not found")
	}
	return i, nil
}

func (m *mockedRepo) GetURL(ctx context.Context, id int64) (string, error) {
	return "http://foo.com/SluG", nil
}

func TestSlug(t *testing.T) {

	h := handler{&mr}

	for url := range sample {
		w := httptest.NewRecorder()
		b := strings.NewReader(fmt.Sprintf("{\"long_url\":\"%s\"}", url))
		req := httptest.NewRequest("POST", "http://example.com/foo", b)
		h.Shorten(w, req)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(resp.StatusCode)
		fmt.Println(string(body))
	}
}

//func TestRedirect(t *testing.T) {
//
//	h := handler{&mr}
//
//	for url := range sample {
//		w := httptest.NewRecorder()
//		b := strings.NewReader(fmt.Sprintf("{\"long_url\":\"%s\"}", url))
//		req := httptest.NewRequest("POST", "http://example.com/foo", b)
//		h.Shorten(w, req)
//		resp := w.Result()
//		body, _ := io.ReadAll(resp.Body)
//		fmt.Println(resp.StatusCode)
//		fmt.Println(string(body))
//	}
//}

func TestEncodeDecode(t *testing.T) {

	assert.Equal(t, "9", encode(9))
	assert.Equal(t, "W3", encode(1987))
	assert.Equal(t, "2tx", encode(11157))

	assert.Equal(t, int64(11157), decode("2tx"))
	assert.Equal(t, int64(8), decode("8"))
	assert.Equal(t, int64(1987), decode("W3"))

}
