package main

import (
	"fmt"
	"heshi/errors"
	"math"
	"strconv"
	"strings"
	"util"
)

func (p *PriceSetting) composeInsertQuery() string {
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
		}
	}
	return fmt.Sprintf("%s) %s)", q, va)
}

func (p *PriceSetting) composeUpdateQuery() string {
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
		}
	}

	return fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), p.ID)
}

func (p *PriceSetting) paramsKV() map[string]interface{} {
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
	if p.Cuts != "" {
		params["cut"] = p.Cuts
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
	if p.Certificates != "" {
		params["certificate"] = p.Certificates
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

func (p *PriceSetting) validatePriceSetting() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage

	sValue, err := strconv.ParseFloat(p.CaratFromStr, 64)
	if err != nil {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("carat from input %s is not valid", p.CaratFromStr)
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}
	p.CaratFrom = math.Abs(sValue)

	sValue, err = strconv.ParseFloat(p.CaratToStr, 64)
	if err != nil {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("carat to input %s is not valid", p.CaratToStr)
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}
	p.CaratTo = math.Abs(sValue)

	sValue, err = strconv.ParseFloat(p.TheParaValueStr, 64)
	if err != nil {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("para value input %s is not valid", p.TheParaValueStr)
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}
	p.TheParaValue = math.Abs(sValue)

	pValue, err := strconv.Atoi(p.PriorityStr)
	if err != nil {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("priority input %s is not valid", p.PriorityStr)
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}
	p.Priority = pValue

	var invalid []string
	for _, v := range strings.Split(p.Certificates, ",") {
		if !util.IsInArrayString(v, VALID_GRADING_LAB) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("certificates input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}

	invalid = []string{}
	for _, v := range strings.Split(p.Clarities, ",") {
		if !util.IsInArrayString(v, VALID_CLARITY) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("clarity input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}

	invalid = []string{}
	for _, v := range strings.Split(p.Colors, ",") {
		if !util.IsInArrayString(v, VALID_COLOR) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("color input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}

	// invalid = []string{}
	// for _, v := range strings.Split(p.Cuts, ",") {
	// 	if !util.IsInArrayString(v, VALID_CUT_NUMBER) {
	// 		invalid = append(invalid, v)
	// 	}
	// }
	// if len(invalid) != 0 {
	// 	VEMSG_NOT_VALID.Message = fmt.Sprintf("cut input has invalid value: %s", strings.Join(invalid, ","))
	// 	vemsg = append(vemsg, VEMSG_NOT_VALID)
	// }
	invalid = []string{}
	for _, v := range strings.Split(p.Fluos, ",") {
		if !util.IsInArrayString(v, VALID_FLUORESCENCE_INTENSITY) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("fluo input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}

	invalid = []string{}
	for _, v := range strings.Split(p.Polishs, ",") {
		if !util.IsInArrayString(v, VALID_POLISH) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("polish input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}

	invalid = []string{}
	for _, v := range strings.Split(p.Symmetries, ",") {
		if !util.IsInArrayString(v, VALID_SYMMETRY) {
			invalid = append(invalid, v)
		}
	}
	if len(invalid) != 0 {
		VEMSG_NOT_VALID.Message = fmt.Sprintf("symmetry input has invalid value: %s", strings.Join(invalid, ","))
		vemsg = append(vemsg, VEMSG_NOT_VALID)
	}
	return vemsg, nil
}
