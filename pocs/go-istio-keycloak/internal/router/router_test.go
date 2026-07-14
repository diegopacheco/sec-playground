package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"book-service/internal/config"
	"book-service/internal/db"
	"book-service/internal/models"
	"book-service/internal/repository"
)

func newServer(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	repo := repository.NewBookRepository(conn)
	return New(config.Config{AuthHeader: "X-Auth-User"}, repo)
}

func TestListRejectsMissingIdentity(t *testing.T) {
	srv := newServer(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/books", nil)
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("no identity must be rejected by the mesh contract, got %d", w.Code)
	}
}

func TestCreateAndPaginateWithIdentity(t *testing.T) {
	srv := newServer(t)

	for i := 0; i < 12; i++ {
		body := `{"title":"T","author":"A","isbn":"I","year":2001}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/books", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-User", "auth0|abc")
		srv.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("create got %d body %s", w.Code, w.Body.String())
		}
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/books?page=2&page_size=5", nil)
	req.Header.Set("X-Auth-User", "auth0|abc")
	srv.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list got %d", w.Code)
	}

	var page models.Page
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if page.Total != 12 || page.TotalPages != 3 || page.Page != 2 {
		t.Fatalf("pagination metadata wrong: %+v", page)
	}
	rows, ok := page.Data.([]any)
	if !ok || len(rows) != 5 {
		t.Fatalf("expected 5 rows on page 2, got %v", page.Data)
	}
}
