package main

import (
	"fmt"
	"heshi/errors"
	"strconv"
	"strings"
	"time"
	"util"
)

func (p *priceSetting) composeInsertQuery() string {
	params := p.paramsKV()
	q := `INSERT INTO price_settings_universal (id `
	va := fmt.Sprintf(`VALUES ('%s'`, p.ID)
	for k, v := range params {
		q = fmt.Sprintf("%s, %s", q, k)
		switch v.(type) {
		case string:
			va = fmt.Sprintf("%s, '%s'", va, v.(string))
		case float64:
			va = fmt.Sprintf("%s, '%f'", va, v.(float64))
		case int:
			va = fmt.Sprintf("%s, '%d'", va, v.(int))
		case int64:
			va = fmt.Sprintf("%s, '%d'", va, v.(int64))
		case time.Time:
			va = fmt.Sprintf("%s, '%s'", va, v.(time.Time).Format(timeFormat))
		}
	}
	return fmt.Sprintf("%s) %s)", q, va)
}

func (p *priceSetting) composeUpdateQuery() string {
	params := p.paramsKV()
	q := `UPDATE price_settings_universal SET`
	for k, v := range params {
		switch v.(type) {
		case string:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(string))
		case float64:
			q = fmt.Sprintf("%s %s='%f',", q, k, v.(float64))
		case int:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int))
		case int64:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int64))
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}

	return fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), p.ID)
}
func (p *priceSetting) composeUpdateQueryTrack(updatedBy string) string {
	params := p.paramsKV()
	q := `UPDATE price_settings_universal SET`
	for k, v := range params {
		switch v.(type) {
		case string:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(string))
		case float64:
			q = fmt.Sprintf("%s %s='%f',", q, k, v.(float64))
		case int:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int))
		case int64:
			q = fmt.Sprintf("%s %s='%d',", q, k, v.(int64))
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}
	newHistoryRecords(updatedBy, "price_settings_universal", p.ID, params)
	return fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), p.ID)
}

func (p *priceSetting) paramsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if p.SupplierID != "" {
		params["supplier_id"] = p.SupplierID
	}
	if p.CaratFrom != 0 {
		params["carat_from"] = p.CaratFrom
	}
	if p.CaratTo != 0 {
		params["carat_to"] = p.CaratTo
	}
	if p.Colors != "" {
		params["color"] = p.Colors
	}
	if p.Clarities != "" {
		params["clarity"] = p.Clarities
	}
	if p.CutGrades != "" {
		params["cut_grade"] = p.CutGrades
	}
	if p.Symmetries != "" {
		params["symmetry"] = p.Symmetries
	}
	if p.Polishs != "" {
		params["polish"] = p.Polishs
	}
	if p.Fluos != "" {
		params["fluo"] = p.Fluos
	}
	if p.GradingLabs != "" {
		params["grading_lab"] = p.GradingLabs
	}
	if p.TheParaValue != 0 {
		params["the_para_value"] = p.TheParaValue
	}
	if p.Priority != 0 {
		params["priority"] = p.Priority
	}
	// if p.Status != "" {
	// 	params["status"] = p.Status
	// }
	return params
}

func (p *priceSetting) validatePriceSetting() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage

	sValue, err := util.StringToFloat(p.CaratFromStr)
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("carat from input %s is not valid", p.CaratFromStr)
		vemsg = append(vemsg, vemsgNotValid)
	}
	p.CaratFrom = sValue

	sValue, err = util.StringToFloat(p.CaratToStr)
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("carat to input %s is not valid", p.CaratToStr)
		vemsg = append(vemsg, vemsgNotValid)
	}
	p.CaratTo = sValue

	sValue, err = util.StringToFloat(p.TheParaValueStr)
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("para value input %s is not valid", p.TheParaValueStr)
		vemsg = append(vemsg, vemsgNotValid)
	}
	p.TheParaValue = sValue

	pValue, err := strconv.Atoi(p.PriorityStr)
	if err != nil {
		vemsgNotValid.Message = fmt.Sprintf("priority input %s is not valid", p.PriorityStr)
		vemsg = append(vemsg, vemsgNotValid)
	}
	p.Priority = pValue

	invalid := itemsNotInArray(p.GradingLabs, VALID_GRADING_LAB)
	if len(invalid) != 0 {
		vemsgNotValid.Message = fmt.Sprintf("grading lab input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, vemsgNotValid)
	}

	invalid = itemsNotInArray(p.Clarities, VALID_CLARITY)
	if len(invalid) != 0 {
		vemsgNotValid.Message = fmt.Sprintf("clarity input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, vemsgNotValid)
	}

	invalid = itemsNotInArray(p.Colors, VALID_COLOR)
	if len(invalid) != 0 {
		vemsgNotValid.Message = fmt.Sprintf("color input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, vemsgNotValid)
	}

	invalid = itemsNotInArray(p.CutGrades, VALID_CUT_GRADE)
	// p.CutGrades = strings.Join(cutGrades, ",")
	if len(invalid) != 0 {
		vemsgNotValid.Message = fmt.Sprintf("cut grade input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, vemsgNotValid)
	}

	invalid = []string{}
	for _, v := range strings.Split(p.Fluos, ",") {
		if !util.IsInArrayString(v, VALID_FLUORESCENCE_INTENSITY) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		vemsgNotValid.Message = fmt.Sprintf("fluo input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, vemsgNotValid)
	}

	invalid = []string{}
	for _, v := range strings.Split(p.Polishs, ",") {
		if !util.IsInArrayString(v, VALID_POLISH) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		vemsgNotValid.Message = fmt.Sprintf("polish input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, vemsgNotValid)
	}

	invalid = []string{}
	for _, v := range strings.Split(p.Symmetries, ",") {
		if !util.IsInArrayString(v, VALID_SYMMETRY) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		vemsgNotValid.Message = fmt.Sprintf("symmetry input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, vemsgNotValid)
	}
	return vemsg, nil
}

func isSupplierPriceRuleExistByID(id string) (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT COUNT(*) FROM price_settings_universal WHERE id='%s'", id)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}
