package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NailUsmanov/practicum-shortener-url/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	// Создаем хранилище с моком, который возвращает фиксированный ключ
	mockStore := storage.NewMemoryStorage()

	app := NewApp(mockStore, "http://test")

	t.Run("Test routes", func(t *testing.T) {
		// 1. Тестируем создание короткой ссылки
		reqPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
		wPost := httptest.NewRecorder()
		app.router.ServeHTTP(wPost, reqPost)

		resPost := wPost.Result()
		defer resPost.Body.Close()

		assert.Equal(t, http.StatusCreated, resPost.StatusCode)

		// Получаем shortURL из ответа
		body, err := io.ReadAll(resPost.Body)
		require.NoError(t, err)
		shortURL := string(body) // "http://test/mock123"

		// Извлекаем ключ из shortURL
		parts := strings.Split(shortURL, "/")
		key := parts[len(parts)-1]

		// 2. Тестируем редирект по полученному ключу
		reqGet := httptest.NewRequest(http.MethodGet, "/"+key, nil)
		wGet := httptest.NewRecorder()
		app.router.ServeHTTP(wGet, reqGet)

		resGet := wGet.Result()
		defer resGet.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, resGet.StatusCode)
		assert.Equal(t, "https://example.com", resGet.Header.Get("Location"))
	})
}
