package main

import (
	"fmt"
	"heshi/errors"
	"strconv"
	"strings"
	"util"
)

func (g *gem) composeInsertQuery() string {
	params := g.parmsKV()
	q := `INSERT INTO gems (id`
	va := fmt.Sprintf(`VALUES ('%s'`, g.ID)
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
	q = fmt.Sprintf("%s) %s)", q, va)
	return q
}

func (g *gem) composeUpdateQuery() string {
	params := g.parmsKV()
	q := `UPDATE gems SET`
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
	q = fmt.Sprintf("%s WHERE id='%s'", strings.TrimSuffix(q, ","), g.ID)
	return q
}

// 	params := make(map[string]interface{})
//TODO validate input
func (g *gem) parmsKV() map[string]interface{} {
	params := make(map[string]interface{})
	if g.StockID != "" {
		params["stock_id"] = g.StockID
	}
	if g.Material != "" {
		params["material"] = g.Material
	}
	if g.Name != "" {
		params["name"] = g.Name
	}
	if g.Shape != "" {
		params["shape"] = g.Shape
	}
	if g.Size != 0 {
		params["size"] = g.Size
	}
	if g.Price != 0 {
		params["price"] = g.Price
	}
	if g.Featured != "" {
		params["featured"] = g.Featured
	}
	if g.Text != "" {
		params["text"] = g.Text
	}
	if g.Certificate != "" {
		params["certificate"] = g.Certificate
	}
	if g.Online != "" {
		params["online"] = g.Online
	}
	if g.Verified != "" {
		params["verified"] = g.Verified
	}
	if g.InStock != "" {
		params["in_stock"] = g.InStock
	}
	if g.Featured != "" {
		params["featured"] = g.Featured
	}
	if g.StockQuantity != 0 {
		params["stock_quantity"] = g.StockQuantity
	}
	if g.Profitable != "" {
		params["profitable"] = g.Profitable
	}
	if g.TotallyScanned != 0 {
		params["totally_scanned"] = g.TotallyScanned
	}
	if g.FreeAcc != "" {
		params["free_acc"] = g.FreeAcc
	}
	return params
}

func (g *gem) validateGemReq() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	cValue, err := util.StringToFloat(g.SizeStr)
	if err != nil {
		vemsgNotValid.Message = "gem size input is not valid"
		vemsg = append(vemsg, vemsgNotValid)
	} else if cValue == 0 {
		vemsgNotValid.Message = "gem size input is not valid"
		vemsg = append(vemsg, vemsgNotValid)
	}
	g.Size = cValue
	pValue, err := util.StringToFloat(g.PriceStr)
	if err != nil {
		vemsgNotValid.Message = "gem price input is not valid"
		vemsg = append(vemsg, vemsgNotValid)
	} else if pValue == 0 {
		vemsgNotValid.Message = "gem price input is not valid"
		vemsg = append(vemsg, vemsgNotValid)
	}
	g.Price = pValue
	prValue, err := strconv.Atoi(g.StockQuantityStr)
	if err != nil {
		vemsgNotValid.Message = "gem stock quantity input is not valid"
		vemsg = append(vemsg, vemsgNotValid)
	} else if prValue == 0 {
		vemsgNotValid.Message = "gem stock quantity input is not valid"
		vemsg = append(vemsg, vemsgNotValid)
	}
	g.StockQuantity = prValue
	// return vemsg, nil

	//be an array
	if g.Shape == "" {
		vemsgNotValid.Message = "must input gem shape"
		vemsg = append(vemsg, vemsgNotValid)
	} else {
		shape := strings.Replace(g.Shape, " ", "", -1)
		shapes := strings.Split(strings.ToUpper(shape), ",")
		if !util.IsIn(shapes, VALID_DIAMOND_SHAPE) {
			vemsgNotValid.Message = "gem shape input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		}
		g.Shape = strings.Join(shapes, ",")
	}

	if g.Text == "" {
		vemsgNotValid.Message = "must input gem text"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if g.Certificate == "" {
		vemsgNotValid.Message = "must input gem certificate"
		vemsg = append(vemsg, vemsgNotValid)
	} else {
		if e, err := isGemExistByCertificate(g.Certificate); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem certificate already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	if g.StockID == "" {
		vemsgNotValid.Message = "must input gem stock id"
		vemsg = append(vemsg, vemsgNotValid)
	} else {
		if e, err := isGemExistByStockID(g.StockID); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem stock id already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	if g.Name == "" {
		vemsgNotValid.Message = "must input gem name"
		vemsg = append(vemsg, vemsgNotValid)
	}
	return vemsg, nil
}

func (g *gem) validateGemUpdateReq() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if g.SizeStr != "" {
		cValue, err := util.StringToFloat(g.SizeStr)
		if err != nil {
			vemsgNotValid.Message = "gem size input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else if cValue == 0 {
			vemsgNotValid.Message = "gem size input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		}
		g.Size = cValue
	}
	if g.PriceStr != "" {
		pValue, err := util.StringToFloat(g.PriceStr)
		if err != nil {
			vemsgNotValid.Message = "gem price input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else if pValue == 0 {
			vemsgNotValid.Message = "gem price input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		}
		g.Price = pValue
	}
	if g.StockQuantityStr != "" {
		prValue, err := strconv.Atoi(g.StockQuantityStr)
		if err != nil {
			vemsgNotValid.Message = "gem stock quantity input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		} else if prValue == 0 {
			vemsgNotValid.Message = "gem stock quantity input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		}
		g.StockQuantity = prValue
	}
	//be an array
	if g.Shape != "" {
		shape := strings.Replace(g.Shape, " ", "", -1)
		shapes := strings.Split(strings.ToUpper(shape), ",")
		if !util.IsIn(shapes, VALID_DIAMOND_SHAPE) {
			vemsgNotValid.Message = "gem shape input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		}
		g.Shape = strings.Join(shapes, ",")
	}

	if g.Certificate != "" {
		if e, err := isGemExistByCertificate(g.Certificate); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem certificate already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	if g.StockID != "" {
		if e, err := isGemExistByStockID(g.StockID); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem stock id already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	return vemsg, nil
}

func isGemExistByStockID(stockID string) (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM gems WHERE stock_id='%s'", stockID)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func isGemExistByCertificate(certificate string) (bool, error) {
	var count int
	q := fmt.Sprintf("SELECT count(*) FROM gems WHERE certificate='%s'", certificate)
	if err := dbQueryRow(q).Scan(&count); err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
