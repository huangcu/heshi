package util

import (
	"errors"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func IsInArrayString(item string, items []string) bool {
	inFlag := false
	for _, v := range items {
		if v == item {
			inFlag = true
			break
		}
	}
	return inFlag
}

func IsInArrayInt(item int, items []int) bool {
	inFlag := false
	for _, v := range items {
		if v == item {
			inFlag = true
			break
		}
	}
	return inFlag
}

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func IsIn(items1 []string, items2 []string) bool {
	for _, v1 := range items1 {
		inFlag := false
		for _, v2 := range items2 {
			if v1 == v2 {
				inFlag = true
				break
			}
		}
		if !inFlag {
			return false
		}
	}
	return true
}

func GetRequestIP(req *http.Request) string {
	xRealIP := req.Header.Get("X-Real-Ip")
	xForwardedFor := req.Header.Get("X-Forwarded-For")

	// If both empty, return IP from remote address
	if xRealIP == "" && xForwardedFor == "" {
		var remoteIP string
		// If there are colon in remote address, remove the port number
		// otherwise, return remote address as is
		if strings.ContainsRune(req.RemoteAddr, ':') {
			remoteIP, _, _ = net.SplitHostPort(req.RemoteAddr)
		} else {
			remoteIP = req.RemoteAddr
		}

		return remoteIP
	}

	// Check list of IP in X-Forwarded-For and return the first global address
	for _, address := range strings.Split(xForwardedFor, ",") {
		address = strings.TrimSpace(address)
		isPrivate, err := isPrivateAddress(address)
		if !isPrivate && err == nil {
			return address
		}
	}

	// If nothing succeed, return X-Real-IP
	return xRealIP
}

func isPrivateAddress(address string) (bool, error) {
	maxCidrBlocks := []string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}

	cidrs := make([]*net.IPNet, len(maxCidrBlocks))
	for i, maxCidrBlock := range maxCidrBlocks {
		_, cidr, _ := net.ParseCIDR(maxCidrBlock)
		cidrs[i] = cidr
	}
	ipAddress := net.ParseIP(address)
	if ipAddress == nil {
		return false, errors.New("address is not valid")
	}

	for i := range cidrs {
		if cidrs[i].Contains(ipAddress) {
			return true, nil
		}
	}
	return false, nil
}
