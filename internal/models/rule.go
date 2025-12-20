package models

import (
	"regexp"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Rule struct {
	gorm.Model
	Name        string      `json:"name" gorm:"column:name"`
	Uuid        string      `json:"uuid" gorm:"column:uuid;unique;not null;type:varchar(100);default:null"`
	Category    string      `json:"-" gorm:"column:category"`
	Destination string      `json:"destination" gorm:"column:destination"`
	Criteria    []Criterion `json:"criteria" gorm:"foreignKey:RuleID;reference=ID"`
}

type Rules []Rule

type Criterion struct {
	gorm.Model
	RuleID            uint             `json:"-" gorm:"column:rule_id"`
	ParentCriterionID uint             `json:"-" gorm:"column:parent_criterion_id"`
	Type              string           `json:"type" gorm:"column:type"`
	Input             string           `json:"input" gorm:"column:input"`
	Values            pq.StringArray   `json:"values" gorm:"type:text[];column:values"`
	BuiltValues       []*regexp.Regexp `json:"-" gorm:"-:all"`
	SubCriteria       []Criterion      `json:"subcriteria" gorm:"foreignKey:ID;reference:ParentCriterionID"`
	Count             *int             `json:"count" gorm:"column:count"`
	ShouldExactCount  bool             `json:"shouldExactCount" gorm:"column:should_exact_count"`
}

func SaveRule(tx *gorm.DB, rule *Rule) error {
	return tx.Create(rule).Error
}

func SaveAllRules(tx *gorm.DB, rules []Rule) error {
	for _, rule := range rules {
		if err := tx.Create(&rule).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetAllRules(tx *gorm.DB) ([]Rule, error) {
	var out []Rule
	err := tx.Preload("Criteria").Find(&out).Error
	return out, err
}

func GetAllRulesUID(tx *gorm.DB) ([]string, error) {
	var out []string
	err := tx.Model(&Rule{}).Select("uuid").Find(&out).Error
	return out, err
}

func GetRuleByID(tx *gorm.DB, ID uint) (Rule, error) {
	var out Rule
	err := tx.Preload("Criteria").Where("id = ?", ID).Find(&out).Error
	return out, err
}

func DeleteRuleByID(tx *gorm.DB, ID uint) error {
	return tx.Select(clause.Associations).Delete(&Rule{Model: gorm.Model{ID: ID}}).Error
}

func DeleteAllRules(tx *gorm.DB) error {
	var rules []Rule
	err := tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where("1 = 1").Delete(&rules).Error
	if err != nil {
		return err
	}

	for _, rule := range rules {
		err = tx.Where(&Criterion{RuleID: rule.ID}).Delete(&Criterion{}).Error
		if err != nil {
			return err
		}
	}
	return nil

}
