package app

import (
	"net/http"
	"time"

	"github.com/arseniizyk/CoCMetaAnalyser/internal/cli"
	"github.com/arseniizyk/CoCMetaAnalyser/internal/config"
	"github.com/arseniizyk/CoCMetaAnalyser/internal/pkg/coc"
	"github.com/arseniizyk/CoCMetaAnalyser/internal/service"
)

type App struct {
	c   *config.Config
	svc service.Service
}

func New() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{Timeout: 15 * time.Second}
	cocClient := coc.New(&httpClient, cfg.CocToken)
	svc := service.New(cocClient)

	return &App{
		c:   cfg,
		svc: svc,
	}, nil
}

func (a *App) Run() error {
	cli := cli.New(a.svc)
	if err := cli.Run(); err != nil {
		return err
	}

	return nil
}
