package main

import (
	"database/sql"
	"fmt"
	"heshi/errors"
	"strings"
	"sync"
	"util"
)

//daily or weekly check
func agentDailyCheck() error {
	q := fmt.Sprintf(`SELECT id, level FROM users WHERE user_type = '%s' AND status='ACTIVE'`, AGENT)
	rows, err := dbQuery(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	agents := make(map[string]string)
	for rows.Next() {
		var ID, level string
		if err := rows.Scan(&ID, &level); err != nil {
			return err
		}
		agents[ID] = level
	}
	var wg sync.WaitGroup
	conChan := make(chan bool, 10)
	for ID, level := range agents {
		wg.Add(1)
		conChan <- true
		go func(id, level string) {
			defer wg.Done()
			if err := automaticAgentLevelAndPurchaseAmount(id, level); err != nil {
				util.Printf("Error: automaticAgentLevelAndPurchaseAmount - %s", errors.GetMessage(err))
			}
			if err := returnPointForCustomer(id); err != nil {
				util.Printf("Error: returnPointForCustomer - %s", errors.GetMessage(err))
			}
			<-conChan
		}(ID, level)
	}
	wg.Wait()
	return nil
}

//daily or weekly check
func customerDailyCheck() error {
	q := fmt.Sprintf("SELECT id, level FROM users WHERE user_type='%s' AND status='ACTIVE'", CUSTOMER)
	rows, err := dbQuery(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	customers := make(map[string]string)
	for rows.Next() {
		var ID, level string
		if err := rows.Scan(&ID, &level); err != nil {
			return err
		}
		customers[ID] = level
	}

	var wg sync.WaitGroup
	conChan := make(chan bool, 10)
	for ID, level := range customers {
		wg.Add(1)
		conChan <- true
		go func(id, level string) {
			defer wg.Done()
			if err := automaticCustomerLevelAndPurchaseAmount(id, level); err != nil {
				util.Printf("Error: automaticCustomerLevelAndPurchaseAmount - %s", errors.GetMessage(err))
			}
			if err := returnPointForAgent(id, level); err != nil {
				util.Printf("Error: returnPointForAgent - %s", errors.GetMessage(err))
			}
			<-conChan
		}(ID, level)
	}
	wg.Wait()
	return nil
}

// =====================================
func automaticAgentLevelAndPurchaseAmount(agentID, level string) error {
	q := fmt.Sprintf(`SELECT count(*) AS count, sum(sold_price_usd) AS total_amount 
		FROM orders 
		WHERE buyer_id='%s'  
		AND status IN ('%s','%s','%s','%s','%s','%s')  
		AND created_at >SUBDATE(NOW(), INTERVAL 1 YEAR)
		GROUP BY item_category`,
		agentID, MPAID, PAID, MDELIVERED, DELIVERED, MRECEIVED, RECEIVED)
	rows, err := dbQuery(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var totalAmount float64
	count := 0
	for rows.Next() {
		var soldPrice float64
		if err := rows.Scan(&soldPrice); err != nil {
			return err
		}
		count++
		totalAmount = totalAmount + soldPrice
	}
	newLevel, err := agentLevelByAmountAndPieces(totalAmount, count)
	if err != nil {
		return err
	}
	return updateUserLevelAndPurchaseAmount(agentID, level, newLevel, totalAmount)
}

func agentLevelByAmountAndPieces(amount float64, pieces int) (string, error) {
	q := fmt.Sprintf(`SELECT level FROM level_rate_rules 
		WHERE rule_type='%s' 
		AND (amount < '%f' OR pieces < '%d') order by level DESC`,
		AGENT, amount, pieces)
	rows, err := dbQuery(q)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var level string
		if err := rows.Scan(&level); err != nil {
			return level, nil
		}
	}
	return LEVEL1, nil
}

// =========================================
func automaticCustomerLevelAndPurchaseAmount(customerID, level string) error {
	q := fmt.Sprintf(`SELECT sum(sold_price_usd) AS total_amount 
		FROM orders 
		WHERE buyer_id='%s'  
		AND status IN ('%s','%s','%s','%s','%s','%s')  
		AND created_at >SUBDATE(NOW(), INTERVAL 1 YEAR)`,
		customerID, MPAID, PAID, MDELIVERED, DELIVERED, MRECEIVED, RECEIVED)

	var totalAmount float64
	if err := dbQueryRow(q).Scan(&totalAmount); err != nil {
		return err
	}

	newLevel, err := customerLevelByAmount(totalAmount)
	if err != nil {
		return err
	}
	return updateUserLevelAndPurchaseAmount(customerID, level, newLevel, totalAmount)
}

func customerLevelByAmount(amount float64) (string, error) {
	q := fmt.Sprintf("SELECT level FROM level_rate_rules WHERE rule_type='%s' AND amount < '%f' order by level DESC",
		CUSTOMER, amount)
	rows, err := dbQuery(q)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var level string
		if err := rows.Scan(&level); err != nil {
			return level, nil
		}
	}
	// default return the lowest level
	return LEVEL1, nil
	// if amount < 5000 {
	// 	return LEVEL0
	// } else if amount >= 5000 && amount < 10000 {
	// 	return LEVEL1
	// } else if amount >= 10000 && amount < 20000 {
	// 	return LEVEL2
	// }
	// return LEVEL3
}

func updateUserLevelAndPurchaseAmount(userID, oldLevel, newLevel string, totalAmount float64) error {
	if oldLevel == newLevel {
		return nil
	}
	q := fmt.Sprintf("UPDATE users set level='%s', total_purchase_amount=%f WHERE id='%s'", newLevel, totalAmount, userID)
	if _, err := dbExec(q); err != nil {
		return nil
	}
	action := "UPGRADE"
	if strings.Compare(newLevel, oldLevel) < 0 {
		action = "DOWNGRADE"
	}
	info := fmt.Sprintf("user: %s %s level from %s to %s", userID, strings.ToLower(action), oldLevel, newLevel)
	return addActionLog(userID, action, info)
}

func addActionLog(userID, action, info string) error {
	q := fmt.Sprintf("INSERT INTO action_logs (id, action, user_id, info) VALUES('%s','%s','%s','%s')",
		newV4(), action, userID, info)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}

// ======================================
//return point caculation
// TODO start of year
// TODO price_usd, cny, eur
// TODO sold price caculate -- return point caculation
func returnPointForCustomer(customerID string) error {
	q := fmt.Sprintf(`SELECT SUM(sold_price_usd) AS total_amount 
	FROM orders
	WHERE status IN ('%s','%s','%s','%s','%s','%s') 
	AND buyer_id IN (SELECT id FROM users WHERE recommended_by='%s' AND user_type='%s') 
	AND created_at > timestampadd(year, -1, now())`,
		MPAID, PAID, MDELIVERED, DELIVERED, MRECEIVED, RECEIVED, customerID, CUSTOMER)

	var totalAmount sql.NullFloat64
	if err := dbQueryRow(q).Scan(&totalAmount); err != nil {
		return err
	}
	if totalAmount.Float64 != 0 {
		// TODO rule here;
		returnPoint := totalAmount.Float64
		// update users return point,返点： 整型）
		q = fmt.Sprintf(`UPDATE users set point=point+%f WHERE id='%s'`, returnPoint, customerID)
		if _, err := dbExec(q); err != nil {
			return err
		}
	}
	return nil
}

func returnPointForAgent(agentID, level string) error {
	// return point from recommended customer
	q := fmt.Sprintf(`SELECT SUM(sold_price_usd) AS total_amount 
	FROM orders
	WHERE status IN ('%s','%s','%s','%s','%s','%s') 
	AND buyer_id IN (SELECT id FROM users WHERE recommended_by='%s' AND user_type='%s') 
	AND created_at > timestampadd(year, -1, now())`,
		MPAID, PAID, MDELIVERED, DELIVERED, MRECEIVED, RECEIVED, agentID, CUSTOMER)

	var totalAmount sql.NullFloat64
	if err := dbQueryRow(q).Scan(&totalAmount); err != nil {
		return err
	}
	// select item_quantity, (select id from users order by created_at desc limit 1)
	// return point based on agent's own level accordingt to system rule
	// subquery must limit to 1, possible will be null
	q = fmt.Sprintf(`SELECT SUM(orders.sold_price_usd) AS total_amount, 
	(SELECT return_point_percent FROM level_rate_rules WHERE level='%s' 
		AND rule_type='%s' ORDER BY created_at DESC LIMIT 1) AS return_point_percent 
	FROM orders JOIN level_rate_rules ON orders.buyer_id == 
	WHERE status IN ('%s','%s','%s','%s','%s','%s') 
	AND buyer_id='%s' 
	AND created_at > timestampadd(year, -1, now())`,
		level, AGENT, MPAID, PAID, MDELIVERED, DELIVERED, MRECEIVED, RECEIVED, agentID)

	var totalAmount2 sql.NullFloat64
	var returnPointPercent sql.NullInt64
	if err := dbQueryRow(q).Scan(&totalAmount2, &returnPointPercent); err != nil {
		return err
	}

	var returnPoint float64
	if totalAmount2.Float64 != 0 && returnPointPercent.Int64 != 0 {
		returnPoint = totalAmount2.Float64 * float64(returnPointPercent.Int64/100)
	}

	if totalAmount.Float64 != 0 {
		// TODO rule here;
		returnPoint = returnPoint + totalAmount.Float64
		// update users return point,返点： 整型）
	}
	if returnPoint != 0 {
		q = fmt.Sprintf(`UPDATE users set point=point+%f WHERE id='%s'`, returnPoint, agentID)
		if _, err := dbExec(q); err != nil {
			return err
		}
	}
	return nil
}
