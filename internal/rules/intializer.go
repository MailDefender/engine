package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"maildefender/engine/internal/models"
	"maildefender/engine/internal/utils"
)

func PopulateFromDir(tx *gorm.DB, dir string, rulesUidToIgnore []string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if err = PopulateFromDir(tx, entry.Name(), rulesUidToIgnore); err != nil {
				return err
			}
		} else {
			if err = populateFromFile(tx, fmt.Sprintf("%s/%s", dir, entry.Name()), rulesUidToIgnore); err != nil {
				return err
			}
		}
	}
	return nil
}

func populateFromFile(tx *gorm.DB, filepath string, rulesUidToIgnore []string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	fmt.Print(string(content))

	var rules RulesPerCatergory
	if err = json.Unmarshal(content, &rules); err != nil {
		return err
	}

	for _, r := range rules.Align() {
		if slices.Contains(rulesUidToIgnore, r.Uuid) {
			logrus.WithFields(logrus.Fields{"rule_uuid": r.Uuid}).Info("this rule already exist, skipping...")
			continue
		}
		logrus.WithField("rule_uuid", r.Uuid).Info("populating this rule")
		if err = models.SaveRule(tx, &r); err != nil && !utils.IsUniqueViolationErr(err) {
			logrus.WithError(err).WithField("rule_uuid", r.Uuid).Error("cannot insert this rule")
			return err
		}
	}
	return nil
}
