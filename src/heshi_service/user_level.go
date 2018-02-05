package main

import (
	"fmt"
	"heshi/errors"
)

func isAgent(uid string) (bool, error) {
	var userType string
	q := fmt.Sprintf("SELECT user_type FROM users WHERE id='%s'", uid)
	if err := dbQueryRow(q).Scan(&userType); err != nil {
		return false, err
	}
	return userType == AGENT, nil
}

func userLevelDiscount(uid string) (int, int, error) {
	var level, discount int
	q := fmt.Sprintf("SELECT level, discount FROM users WHERE id='%s'", uid)
	if err := dbQueryRow(q).Scan(&level, discount); err != nil {
		return 0, 0, err
	}
	return level, discount, nil
}

func agentLevelDiscount(uid string) (int, int, error) {
	if s, err := isAgent(uid); err != nil {
		return 0, 0, err
	} else if !s {
		return 0, 0, errors.New("user is not an AGENT")
	}
	var level, discount int
	q := fmt.Sprintf("SELECT level, discount FROM users WHERE id='%s'", uid)
	if err := dbQueryRow(q).Scan(&level, discount); err != nil {
		return 0, 0, err
	}
	return level, discount, nil
}

func getUserType(uid string) (string, error) {
	var userType string
	q := fmt.Sprintf("SELECT user_type FROM users WHERE id = '%s'", uid)
	if err := dbQueryRow(q).Scan(&userType); err != nil {
		return "", err
	}
	return userType, nil
}
