package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"my_law_server/internal/database"
	"my_law_server/internal/handler"
	"my_law_server/internal/repository"
	"my_law_server/internal/server"
	"my_law_server/internal/service"
)

type envelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func TestListTypePreviewsReturnsConcreteTypesWithAtMostTwentyItems(t *testing.T) {
	router := newTestRouter(t, t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/types/previews", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", recorder.Code, http.StatusOK)
	}

	var resp envelope[[]service.TypePreview]
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if resp.Code != 0 {
		t.Fatalf("unexpected business code: got %d want 0", resp.Code)
	}

	if got := len(resp.Data); got != 3 {
		t.Fatalf("unexpected preview type count: got %d want 3", got)
	}

	typeIDs := []int{resp.Data[0].TypeID, resp.Data[1].TypeID, resp.Data[2].TypeID}
	wantTypeIDs := []int{110, 120, 230}
	assertIntSlice(t, typeIDs, wantTypeIDs)

	if got := len(resp.Data[2].Items); got != 20 {
		t.Fatalf("unexpected preview item count for type 230: got %d want 20", got)
	}

	if got := resp.Data[2].Total; got != 22 {
		t.Fatalf("unexpected total for type 230: got %d want 22", got)
	}

	if got := resp.Data[2].Items[0].VersionID; got != "law-230-22" {
		t.Fatalf("unexpected first preview item: got %s want %s", got, "law-230-22")
	}
}

func TestListLawsByTypeSupportsPagination(t *testing.T) {
	router := newTestRouter(t, t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/types/230/laws?page=2&pageSize=5", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", recorder.Code, http.StatusOK)
	}

	var resp envelope[service.PaginatedLawList]
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if resp.Data.Page != 2 || resp.Data.PageSize != 5 {
		t.Fatalf("unexpected pagination metadata: got page=%d pageSize=%d", resp.Data.Page, resp.Data.PageSize)
	}

	if resp.Data.Total != 22 || resp.Data.TotalPages != 5 {
		t.Fatalf("unexpected total metadata: got total=%d totalPages=%d", resp.Data.Total, resp.Data.TotalPages)
	}

	if got := len(resp.Data.Items); got != 5 {
		t.Fatalf("unexpected item count: got %d want 5", got)
	}

	if got := resp.Data.Items[0].VersionID; got != "law-230-17" {
		t.Fatalf("unexpected first item on page 2: got %s want %s", got, "law-230-17")
	}
}

func TestListLawsByTypeSortsEmptyEffectDateLast(t *testing.T) {
	router := newTestRouter(t, t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/types/120/laws", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", recorder.Code, http.StatusOK)
	}

	var resp envelope[service.PaginatedLawList]
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	gotOrder := []string{
		resp.Data.Items[0].VersionID,
		resp.Data.Items[1].VersionID,
		resp.Data.Items[2].VersionID,
	}
	wantOrder := []string{"law-120-b", "law-120-c", "law-120-a"}
	assertStringSlice(t, gotOrder, wantOrder)
}

func TestGetParsedLawReturnsJSONWhenFileExists(t *testing.T) {
	detailDir := t.TempDir()
	writeTestJSONFile(t, detailDir, "law-parsed-1", []byte(`{"chapters":[{"id":1}],"title":"解析版法律"}`))
	router := newTestRouter(t, detailDir)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laws/law-parsed-1/parsed", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", recorder.Code, http.StatusOK)
	}

	var resp envelope[service.ParsedLawDetail]
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if !resp.Data.Available || resp.Data.Content == nil {
		t.Fatalf("expected parsed content to be available")
	}

	var content map[string]any
	decodeJSON(t, *resp.Data.Content, &content)

	if got := content["title"]; got != "解析版法律" {
		t.Fatalf("unexpected parsed title: got %v want %v", got, "解析版法律")
	}
}

func TestGetParsedLawReturnsPlaceholderWhenFileMissing(t *testing.T) {
	router := newTestRouter(t, t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laws/law-parsed-missing/parsed", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", recorder.Code, http.StatusOK)
	}

	var resp envelope[service.ParsedLawDetail]
	decodeJSON(t, recorder.Body.Bytes(), &resp)

	if resp.Message != "暂无解析数据" {
		t.Fatalf("unexpected message: got %s want %s", resp.Message, "暂无解析数据")
	}

	if resp.Data.Available {
		t.Fatalf("expected parsed content to be unavailable")
	}

	if resp.Data.Content != nil {
		t.Fatalf("expected parsed content to be nil when file is missing")
	}
}

func TestGetParsedLawReturnsNotFoundWhenVersionIDDoesNotExist(t *testing.T) {
	router := newTestRouter(t, t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/laws/not-exists/parsed", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: got %d want %d", recorder.Code, http.StatusNotFound)
	}
}

func newTestRouter(t *testing.T, detailDir string) *gin.Engine {
	t.Helper()

	dbPath := createTestSQLiteDatabase(t)
	db, err := database.OpenReadOnlySQLite(dbPath)
	if err != nil {
		t.Fatalf("open read-only sqlite: %v", err)
	}

	typeRepo := repository.NewTypeRepository(db)
	lawRepo := repository.NewLawRepository(db)
	parsedLawRepo := repository.NewParsedLawRepository(detailDir)
	lawService := service.NewLawService(typeRepo, lawRepo, parsedLawRepo)
	lawHandler := handler.NewLawHandler(lawService)

	return server.NewRouter(lawHandler)
}

func createTestSQLiteDatabase(t *testing.T) string {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "laws.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test sqlite: %v", err)
	}

	for _, statement := range []string{
		`CREATE TABLE laws_list (
			versionId TEXT PRIMARY KEY,
			title TEXT,
			lawTypeId INTEGER,
			lawType TEXT,
			publishDate TEXT,
			effectDate TEXT,
			detailJson TEXT,
			effectiveStatus INTEGER,
			authorityId INTEGER,
			authorityName TEXT,
			parse_state INTEGER DEFAULT 0
		);`,
		`CREATE TABLE types (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			parent_id INTEGER
		);`,
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("execute schema statement: %v", err)
		}
	}

	insertType(t, db, 101, "法律(全选)", nil)
	insertType(t, db, 102, "法律", intPtr(101))
	insertType(t, db, 110, "宪法相关法", intPtr(102))
	insertType(t, db, 120, "民法商法", intPtr(102))
	insertType(t, db, 222, "地方法规", nil)
	insertType(t, db, 230, "地方性法规", intPtr(222))

	insertLaw(t, db, "law-parsed-1", "解析文件法律", 110, "法律", "2024-01-02", "2024-01-03", 1, "全国人大")
	insertLaw(t, db, "law-parsed-missing", "缺失解析文件法律", 110, "法律", "2024-01-01", "2024-01-02", 1, "全国人大")
	insertLaw(t, db, "law-110-1", "宪法相关法一", 110, "法律", "2024-01-05", "2024-01-06", 1, "全国人大")
	insertLaw(t, db, "law-110-2", "宪法相关法二", 110, "法律", "2024-01-04", "2024-01-05", 1, "全国人大")

	insertLaw(t, db, "law-120-a", "民法商法 A", 120, "法律", "2024-03-01", "", 1, "国务院")
	insertLaw(t, db, "law-120-b", "民法商法 B", 120, "法律", "2024-01-01", "2024-02-01", 1, "国务院")
	insertLaw(t, db, "law-120-c", "民法商法 C", 120, "法律", "2024-04-01", "", 1, "国务院")

	for i := 1; i <= 22; i++ {
		versionID := fmt.Sprintf("law-230-%02d", i)
		title := fmt.Sprintf("地方性法规 %02d", i)
		date := fmt.Sprintf("2024-03-%02d", i)
		insertLaw(t, db, versionID, title, 230, "地方性法规", date, date, 1, "地方人大")
	}

	return dbPath
}

func insertType(t *testing.T, db *gorm.DB, id int, name string, parentID *int) {
	t.Helper()

	if err := db.Exec(
		`INSERT INTO types(id, name, parent_id) VALUES (?, ?, ?)`,
		id, name, parentID,
	).Error; err != nil {
		t.Fatalf("insert type %d: %v", id, err)
	}
}

func insertLaw(t *testing.T, db *gorm.DB, versionID, title string, lawTypeID int, lawType, publishDate, effectDate string, effectiveStatus int, authorityName string) {
	t.Helper()

	if err := db.Exec(
		`INSERT INTO laws_list(versionId, title, lawTypeId, lawType, publishDate, effectDate, detailJson, effectiveStatus, authorityId, authorityName, parse_state)
		 VALUES (?, ?, ?, ?, ?, ?, '', ?, 0, ?, 1)`,
		versionID, title, lawTypeID, lawType, publishDate, effectDate, effectiveStatus, authorityName,
	).Error; err != nil {
		t.Fatalf("insert law %s: %v", versionID, err)
	}
}

func writeTestJSONFile(t *testing.T, dir, versionID string, data []byte) {
	t.Helper()

	filePath := filepath.Join(dir, versionID+".json")
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		t.Fatalf("write parsed law json: %v", err)
	}
}

func decodeJSON(t *testing.T, data []byte, target any) {
	t.Helper()

	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(target); err != nil {
		t.Fatalf("decode json: %v", err)
	}
}

func assertIntSlice(t *testing.T, got, want []int) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("unexpected slice length: got %d want %d", len(got), len(want))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("unexpected slice value at %d: got %d want %d", i, got[i], want[i])
		}
	}
}

func assertStringSlice(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("unexpected slice length: got %d want %d", len(got), len(want))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("unexpected slice value at %d: got %s want %s", i, got[i], want[i])
		}
	}
}

func intPtr(v int) *int {
	return &v
}
