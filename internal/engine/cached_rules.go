package engine

import (
	"time"

	"github.com/sirupsen/logrus"

	"maildefender/engine/internal/configuration"
	db "maildefender/engine/internal/database"
	"maildefender/engine/internal/models"
)

var cachedRules = struct {
	rules       models.Rules
	lastRefresh time.Time
}{
	rules: nil,
}

func refreshCachedRules() error {
	if cachedRules.lastRefresh.IsZero() || cachedRules.lastRefresh.Add(time.Second*time.Duration(configuration.RulesRefreshDelay())).Before(time.Now()) {
		logrus.Info("rules has expired, refreshing...")
		r, err := models.GetAllRules(db.Instance().Gorm)
		if err != nil {
			return err
		}
		cachedRules.rules = r
		cachedRules.lastRefresh = time.Now()
	}
	return nil
}
