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
func handleAgentLevel() error {
	q := fmt.Sprintf(`SELECT id, level FROM users WHERE userType = '%s' AND status='active'`, AGENT)
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
			if err := automaticAgentLevel(id, level); err != nil {
				util.Printf("Error: automaticAgentLevel - %s", errors.GetMessage(err))
			}
			<-conChan
		}(ID, level)
	}
	wg.Wait()
	return nil
}

//daily or weekly check
func handleCustomerLevel() error {
	q := fmt.Sprintf("SELECT id, level FROM users WHERE user_type='%s' AND status='active'", CUSTOMER)
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
			if err := automaticCustomerLevel(id, level); err != nil {
				util.Printf("Error: automaticCustomerLevel - %s", errors.GetMessage(err))
			}
			<-conChan
		}(ID, level)
	}
	wg.Wait()
	return nil
}

func automaticAgentLevel(agentID, level string) error {
	q := fmt.Sprintf(`SELECT id FROM users WHERE recommended_by='%s' AND user_type='%s'`, agentID, AGENT)
	rows, err := dbQuery(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	ids = append(ids, agentID)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			if err != sql.ErrNoRows {
				return err
			}
		}
		ids = append(ids, id)
	}

	where := fmt.Sprintf("buyer_id IN ('%s')", strings.Join(ids, "','"))
	q = fmt.Sprintf(`SELECT sold_price FROM orders 
	WHERE %s 
	 AND status = 'SOLD' 
	 AND created_at >SUBDATE(NOW(), INTERVAL 1 YEAR)`, where)

	rows, err = dbQuery(q)
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
	return updateUserLevel(agentID, level, newLevel)
}

func automaticCustomerLevel(customerID, level string) error {
	q := fmt.Sprintf(`SELECT id FROM users WHERE recommended_by='%s' AND user_type='%s'`, customerID, CUSTOMER)
	rows, err := dbQuery(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	ids = append(ids, customerID)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			if err != sql.ErrNoRows {
				return err
			}
		}
		ids = append(ids, id)
	}

	where := fmt.Sprintf("buyer_id IN ('%s')", strings.Join(ids, "','"))
	q = fmt.Sprintf(`SELECT sold_price FROM orders 
	WHERE %s 
	 AND status = 'SOLD' 
	 AND created_at >SUBDATE(NOW(), INTERVAL 1 YEAR)`, where)

	rows, err = dbQuery(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var totalAmount float64
	for rows.Next() {
		var soldPrice float64
		if err := rows.Scan(&soldPrice); err != nil {
			return err
		}
		totalAmount = totalAmount + soldPrice
	}

	newLevel, err := customerLevelByAmount(totalAmount)
	if err != nil {
		return err
	}
	return updateUserLevel(customerID, level, newLevel)
}

func updateUserLevel(userID, oldLevel, newLevel string) error {
	if oldLevel == newLevel {
		return nil
	}
	q := fmt.Sprintf("UPDATE users set level='%s' WHERE id='%s'", newLevel, userID)
	if _, err := dbExec(q); err != nil {
		return nil
	}
	info := fmt.Sprintf("users: %s upgrade level from %s to %s", userID, oldLevel, newLevel)
	return addActionLog(userID, "upgrade", info)
}

func addActionLog(userID, action, info string) error {
	q := fmt.Sprintf("INSERT INTO action_logs (id, action, user_id, info) VALUES('%s','%s','%s','%s')",
		newV4(), action, userID, info)
	if _, err := dbExec(q); err != nil {
		return err
	}
	return nil
}
