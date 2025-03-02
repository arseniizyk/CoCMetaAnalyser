package app

import (
	"net/http"
	"time"

	"github.com/arseniizyk/CoCMetaAnalyser/internal/config"
	"github.com/arseniizyk/CoCMetaAnalyser/internal/pkg/coc"
	"github.com/arseniizyk/CoCMetaAnalyser/internal/service"
)

type App struct {
	c    *config.Config
	http http.Client
	svc  service.Service
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
		c:    cfg,
		http: httpClient,
		svc:  svc,
	}, nil
}

func (a *App) Run() error {
	// _, err := a.svc.GetLeagueSeasonRanking(10000)
	// if err != nil {
	// 	return err
	// }
	// _, err := a.svc.GetLeagueSeasons()
	// if err != nil {
	// 	return err
	// }
	// _, err := a.svc.GetLegendaryLeague()
	// if err != nil {
	// 	return err
	// }
	// league, err := a.service.GetLegendaryLeague()
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("Legendary League:", league)

	return nil
}
