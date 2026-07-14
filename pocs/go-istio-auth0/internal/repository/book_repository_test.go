package repository

import (
	"testing"

	"book-service/internal/db"
	"book-service/internal/models"
)

func newRepo(t *testing.T) *BookRepository {
	t.Helper()
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return NewBookRepository(conn)
}

func TestCreateAndGet(t *testing.T) {
	repo := newRepo(t)
	id, err := repo.Create(&models.Book{Title: "Go", Author: "Alan", ISBN: "123", Year: 2015})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "Go" || got.Author != "Alan" || got.Year != 2015 {
		t.Fatalf("unexpected book: %+v", got)
	}
}

func TestListPaginationBoundaries(t *testing.T) {
	repo := newRepo(t)
	for i := 0; i < 25; i++ {
		if _, err := repo.Create(&models.Book{Title: "T", Author: "A", ISBN: "I", Year: 2000 + i}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	total, err := repo.Count()
	if err != nil || total != 25 {
		t.Fatalf("count got %d err %v", total, err)
	}

	page2, err := repo.List(10, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(page2) != 10 {
		t.Fatalf("expected 10 rows on page 2, got %d", len(page2))
	}
	if page2[0].Year != 2010 {
		t.Fatalf("pagination offset wrong, first year %d want 2010", page2[0].Year)
	}

	lastPage, err := repo.List(10, 20)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(lastPage) != 5 {
		t.Fatalf("expected 5 rows on last page, got %d", len(lastPage))
	}
}

func TestDelete(t *testing.T) {
	repo := newRepo(t)
	id, _ := repo.Create(&models.Book{Title: "T", Author: "A", ISBN: "I", Year: 2000})
	affected, err := repo.Delete(id)
	if err != nil || affected != 1 {
		t.Fatalf("delete got %d err %v", affected, err)
	}
	missing, err := repo.Delete(id)
	if err != nil || missing != 0 {
		t.Fatalf("delete missing got %d err %v", missing, err)
	}
}
