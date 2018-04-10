package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"strings"
	"util"
)

func validateGemHeaders(headers []string) []string {
	var missingHeaders []string
	for k, header := range jewelryHeaders {
		if !util.IsInArrayString(header, headers) && k < 8 {
			missingHeaders = append(missingHeaders, header)
		}
	}
	return missingHeaders
}

func importGemProducts(uid, file string) ([]util.Row, error) {
	oldStockIDList, err := getAllGemStockID()
	if err != nil {
		return nil, err
	}
	rows, err := util.ParseCSVToStruct(file)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, errors.New("uploaded file has no rows")
	}

	unimportRows := []util.Row{}
	//get headers
	originalHeaders := rows[0]

	//process rows
	util.Println("start process gem")
	for index := 1; index < len(rows); index++ {
		g := gem{}
		row := rows[index]
		record := row.Value
		util.Printf("processsing row: %d, %s", index, record)
		for i, header := range originalHeaders.Value {
			switch header {
			case "stock_id":
				if record[i] == "" {
					row.Ignored = true
					row.Message = append(row.Message, "jewelry stock id cannot be empty")
					break
				}
				g.StockID = strings.ToUpper(record[i])
			case "name":
				g.Name = strings.ToUpper(record[i])
			case "size":
				g.SizeStr = strings.ToUpper(record[i])
			case "material":
				g.Material = strings.ToUpper(record[i])
			case "shape":
				g.Shape = strings.ToUpper(record[i])
			case "certificate":
				g.Certificate = strings.ToUpper(record[i])
			case "price":
				g.PriceStr = record[i]
			case "text":
				g.Text = record[i]
			case "status":
				g.Status = strings.ToUpper(record[i])
			case "verified":
				g.Verified = strings.ToUpper(record[i])
			case "featured":
				g.Featured = strings.ToUpper(record[i])
			case "stock_quantity":
				g.StockQuantityStr = record[i]
			case "profitable":
				g.Profitable = strings.ToUpper(record[i])
			case "free_acc":
				g.FreeAcc = strings.ToUpper(record[i])
			case "image", "image1", "image2", "image3", "image4", "image5":
				g.Images = append(g.Images, record[i])
			}
		}

		if row.Ignored {
			unimportRows = append(unimportRows, row)
			//move on to next row
			continue
		}

		//validate row
		var id string
		q := fmt.Sprintf("SELECT id FROM gems WHERE stock_id='%s'", g.StockID)
		if err := dbQueryRow(q).Scan(&id); err != nil {
			//new record
			if err == sql.ErrNoRows {
				//validate as new request
				if vemsg, err := g.validateGemReq(false); err != nil {
					return nil, err
				} else if len(vemsg) != 0 {
					row.Ignored = true
					for _, v := range vemsg {
						row.Message = append(row.Message, v.Message)
					}
					unimportRows = append(unimportRows, row)
					//move on to next row
					continue
				}
				//pass validation, insert into db
				g.ID = newV4()
				g.gemImages()
				q := g.composeInsertQuery()
				if _, err := dbExec(q); err != nil {
					util.Printf("fail to add gem item. stock id: %s; err: %s", g.StockID, errors.GetMessage(err))
					return nil, err
				}
				util.Printf("new gem item added. stock id: %s", g.StockID)
			} else {
				return nil, err
			}
		}

		//already exist, validate as update request
		if vemsg, err := g.validateGemReq(true); err != nil {
			return nil, err
		} else if len(vemsg) != 0 {
			row.Ignored = true
			for _, v := range vemsg {
				row.Message = append(row.Message, v.Message)
			}
			unimportRows = append(unimportRows, row)
			//move on to next row
			continue
		}
		//pass validation, update db
		g.ID = id
		g.gemImages()
		q = g.composeUpdateQuery()
		if _, err := dbExec(q); err != nil {
			util.Printf("fail to update gem item. stock id: %s; err; %s", g.StockID, errors.GetMessage(err))
			return nil, err
		}
		util.Printf("gem item updated. stock id: %s", g.StockID)
		go newHistoryRecords(uid, "gems", g.ID, g.parmsKV())
		//remove updated stockID from old one as this has been scanned and processed
		delete(oldStockIDList, g.StockID)
	}
	util.Println("finish process gem")
	if err := offlineGemsNoLongerExist(oldStockIDList); err != nil {
		return unimportRows, err
	}
	return unimportRows, nil
}

func getAllGemStockID() (map[string]struct{}, error) {
	rows, err := dbQuery("SELECT stock_id FROM gems WHERE status IN ('AVAILABLE','OFFLINE') ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stockIDs := make(map[string]struct{})
	for rows.Next() {
		var stockID string
		if err := rows.Scan(&stockID); err != nil {
			return nil, err
		}
		//empty struct comsumes 0 bytes
		var s struct{}
		stockIDs[stockID] = s
	}
	return stockIDs, nil
}

func (g *gem) gemImages() {
	var imageNames []string
	for _, imageName := range g.Images {
		name := fmt.Sprintf("beyoudiamond-image-%s-%s", g.StockID, imageName)
		imageNames = append(imageNames, name)
	}
	g.Images = imageNames
}

//下线不存在的钻石 //TODO return or just trace err ???
func offlineGemsNoLongerExist(stockIDList map[string]struct{}) error {
	util.Traceln("Start to offline all gems no longer exists.")
	for k := range stockIDList {
		q := fmt.Sprintf("UPDATE gems SET offline='YES',updated_at=(CURRENT_TIMESTAMP) WHERE stock_id ='%s'", k)
		util.Tracef("Offline gem stock_id: %s.\n", k)
		if _, err := dbExec(q); err != nil {
			util.Tracef("error when offline gem. stock_id: %s. err: \n", k, errors.GetMessage(err))
			return err
		}
	}
	util.Traceln("Finished offline all gems no longer exists")
	return nil
}
