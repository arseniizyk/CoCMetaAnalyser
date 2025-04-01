package app

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

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
	// minimum 100 players and maximum 25k
	// meta, err := a.svc.GetMetaItems(10000, "2025-03")
	// if err != nil {
	// 	return err
	// }

	// file, _ := os.Create("meta/2025-03/meta10k.json")
	// defer file.Close()
	// formatted, _ := json.MarshalIndent(meta, "", "  ")
	// file.Write(formatted)

	// meta, err = a.svc.GetMetaItemPairs(10000, "2025-03")
	// if err != nil {
	// 	return err
	// }

	// file, _ = os.Create("metapairs/2025-03/meta10k.json")
	// defer file.Close()
	// formatted, _ = json.MarshalIndent(meta, "", "  ")
	// file.Write(formatted)
	// return nil
}
