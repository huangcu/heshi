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
	if g.SizeStr == "" {
		vemsg = append(vemsg, vemsgGemSizeEmpty)
	} else {
		cValue, err := util.StringToFloat(g.SizeStr)
		if err != nil {
			vemsg = append(vemsg, vemsgGemSizeNotValid)
		} else if cValue == 0 {
			vemsg = append(vemsg, vemsgGemSizeNotValid)
		} else {
			g.Size = cValue
		}
	}
	if g.PriceStr == "" {
		vemsg = append(vemsg, vemsgGemPriceEmpty)
	} else {
		pValue, err := util.StringToFloat(g.PriceStr)
		if err != nil {
			vemsg = append(vemsg, vemsgGemPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgGemPriceNotValid)
		} else {
			g.Price = pValue
		}
	}
	if g.StockQuantityStr == "" {
		vemsg = append(vemsg, vemsgStockQuantityEmpty)
	} else {
		prValue, err := strconv.Atoi(g.StockQuantityStr)
		if err != nil {
			vemsg = append(vemsg, vemsgStockQuantityNotValid)
		} else if prValue == 0 {
			vemsg = append(vemsg, vemsgStockQuantityNotValid)
		} else {
			g.StockQuantity = prValue
		}
	}
	// return vemsg, nil

	//be an array
	if g.Shape == "" {
		vemsgNotValid.Message = "gem shape cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	} else {
		shapes := strings.Split(g.Shape, ",")
		if !util.IsIn(shapes, VALID_DIAMOND_SHAPE) {
			vemsgNotValid.Message = "gem shape input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		}
		g.Shape = strings.Join(shapes, ",")
	}

	if g.Text == "" {
		vemsgNotValid.Message = "gem text cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if g.Certificate == "" {
		vemsgNotValid.Message = "gem certificate cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	} else {
		if e, err := isItemExistInDbByProperty("gems", "certificate", g.Certificate); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem certificate already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	if g.StockID == "" {
		vemsgNotValid.Message = "gem stock id cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	} else {
		if e, err := isItemExistInDbByProperty("gems", "stock_id", g.StockID); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem stock id already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	if g.Name == "" {
		vemsgNotValid.Message = "gem name cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	}
	return vemsg, nil
}

func (g *gem) validateGemUpdateReq() ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if g.SizeStr != "" {
		cValue, err := util.StringToFloat(g.SizeStr)
		if err != nil {
			vemsg = append(vemsg, vemsgGemSizeNotValid)
		} else if cValue == 0 {
			vemsg = append(vemsg, vemsgGemSizeNotValid)
		} else {
			g.Size = cValue
		}
	}
	if g.PriceStr != "" {
		pValue, err := util.StringToFloat(g.PriceStr)
		if err != nil {
			vemsg = append(vemsg, vemsgGemPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgGemPriceNotValid)
		} else {
			g.Price = pValue
		}
	}
	if g.StockQuantityStr != "" {
		prValue, err := strconv.Atoi(g.StockQuantityStr)
		if err != nil {
			vemsg = append(vemsg, vemsgStockQuantityNotValid)
		} else if prValue == 0 {
			vemsg = append(vemsg, vemsgStockQuantityNotValid)
		} else {
			g.StockQuantity = prValue
		}
	}
	//be an array
	if g.Shape != "" {
		shapes := strings.Split(g.Shape, ",")
		if !util.IsIn(shapes, VALID_DIAMOND_SHAPE) {
			vemsgNotValid.Message = "gem shape input is not valid"
			vemsg = append(vemsg, vemsgNotValid)
		}
		g.Shape = strings.Join(shapes, ",")
	}

	if g.Certificate != "" {
		if e, err := isItemExistInDbByProperty("gems", "stock_id", g.StockID); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem certificate already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	if g.StockID != "" {
		if e, err := isItemExistInDbByProperty("gems", "stock_id", g.StockID); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem stock id already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	return vemsg, nil
}
