package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zhashkevych/creatly-backend/internal/config"
	handler "github.com/zhashkevych/creatly-backend/internal/delivery/http"
	"github.com/zhashkevych/creatly-backend/internal/service"
	"github.com/zhashkevych/creatly-backend/pkg/auth"
)

func TestNewHandler(t *testing.T) {
	h := handler.NewHandler(&service.Services{}, &auth.Manager{})

	require.IsType(t, &handler.Handler{}, h)
}

func TestNewHandler_Init(t *testing.T) {
	h := handler.NewHandler(&service.Services{}, &auth.Manager{})

	router := h.Init(&config.Config{
		Limiter: config.LimiterConfig{
			RPS:   2,
			Burst: 4,
			TTL:   10 * time.Minute,
		},
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
