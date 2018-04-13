package main

import (
	"fmt"
	"heshi/errors"
	"strconv"
	"strings"
	"time"
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
		case time.Time:
			va = fmt.Sprintf("%s, '%s'", va, v.(time.Time).Format(timeFormat))
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
		case time.Time:
			q = fmt.Sprintf("%s %s='%s',", q, k, v.(time.Time).Format(timeFormat))
		}
	}
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, g.ID)
	return q
}

//only track price/promotion_id(track in promotion section) change
func (g *gem) composeUpdateQueryTrack(updatedBy string) string {
	trackMap := make(map[string]interface{})
	params := g.parmsKV()
	q := `UPDATE gems SET`
	for k, v := range params {
		if k == "price" {
			trackMap["price"] = v
		}
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
	if len(trackMap) != 0 {
		newHistoryRecords(updatedBy, "gems", g.ID, trackMap)
	}
	q = fmt.Sprintf("%s updated_at=(CURRENT_TIMESTAMP) WHERE id='%s'", q, g.ID)
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
	if len(g.Images) != 0 {
		params["images"] = strings.Join(g.Images, ";")
	}
	if g.Certificate != "" {
		params["certificate"] = g.Certificate
	}
	if g.Status != "" {
		params["status"] = g.Status
	}
	if g.Verified != "" {
		params["verified"] = g.Verified
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

func (g *gem) validateGemReq(update bool) ([]errors.HSMessage, error) {
	var vemsg []errors.HSMessage
	if !update && g.SizeStr == "" {
		vemsg = append(vemsg, vemsgGemSizeEmpty)
	} else if g.SizeStr != "" {
		cValue, err := util.StringToFloat(g.SizeStr)
		if err != nil {
			vemsg = append(vemsg, vemsgGemSizeNotValid)
		} else if cValue == 0 {
			vemsg = append(vemsg, vemsgGemSizeNotValid)
		} else {
			g.Size = cValue
		}
	}
	if !update && g.PriceStr == "" {
		vemsg = append(vemsg, vemsgGemPriceEmpty)
	} else if g.PriceStr != "" {
		pValue, err := util.StringToFloat(g.PriceStr)
		if err != nil {
			vemsg = append(vemsg, vemsgGemPriceNotValid)
		} else if pValue == 0 {
			vemsg = append(vemsg, vemsgGemPriceNotValid)
		} else {
			g.Price = pValue
		}
	}
	if !update && g.StockQuantityStr == "" {
		vemsg = append(vemsg, vemsgStockQuantityEmpty)
	} else if g.StockQuantityStr != "" {
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
	if !update && g.Shape == "" {
		vemsgNotValid.Message = "gem shape cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	} else if g.Shape != "" {
		s, err := jewelryShape(g.Shape)
		if err != nil {
			return nil, err
		}
		g.Shape = s
	}

	if !update && g.Text == "" {
		vemsgNotValid.Message = "gem text cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if !update && g.Certificate == "" {
		vemsgNotValid.Message = "gem certificate cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if g.Certificate != "" {
		if update {
			if e, err := isItemExistInDbByPropertyWithDifferentID("gems", "certificate", g.Certificate, g.ID); err != nil {
				return nil, err
			} else if e {
				vemsgAlreadyExist.Message = "gem certificate already exist"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		} else {
			if e, err := isItemExistInDbByProperty("gems", "certificate", g.Certificate); err != nil {
				return nil, err
			} else if e {
				vemsgAlreadyExist.Message = "gem certificate already exist"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		}
	}

	if !update && g.StockID == "" {
		vemsgNotValid.Message = "gem stock id cannot be empty"
		vemsg = append(vemsg, vemsgNotValid)
	}
	if g.StockID != "" {
		if update {
			if e, err := isItemExistInDbByPropertyWithDifferentID("gems", "stock_id", g.StockID, g.ID); err != nil {
				return nil, err
			} else if e {
				vemsgAlreadyExist.Message = "gem stock id already exist"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}

		} else {
			if e, err := isItemExistInDbByProperty("gems", "stock_id", g.StockID); err != nil {
				return nil, err
			} else if e {
				vemsgAlreadyExist.Message = "gem stock id already exist"
				vemsg = append(vemsg, vemsgAlreadyExist)
			}
		}
	}

	if !update && g.Name == "" {
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
		if e, err := isItemExistInDbByPropertyWithDifferentID("gems", "stock_id", g.StockID, g.ID); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem certificate already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	if g.StockID != "" {
		if e, err := isItemExistInDbByPropertyWithDifferentID("gems", "stock_id", g.StockID, g.ID); err != nil {
			return nil, err
		} else if e {
			vemsgAlreadyExist.Message = "gem stock id already exist"
			vemsg = append(vemsg, vemsgAlreadyExist)
		}
	}

	return vemsg, nil
}

func (g *gem) isGemExistByStockID() error {
	var id string
	if err := dbQueryRow(fmt.Sprintf("SELECT id FROM gems WHERE diamond_id='%s'", g.StockID)).Scan(&id); err != nil {
		return err
	}
	g.ID = id
	return nil
}

func isGemExistByID(id string) (bool, error) {
	var count int
	if err := dbQueryRow(fmt.Sprintf("SELECT COUNT(*) FROM gems WHERE id='%s'", id)).Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}
