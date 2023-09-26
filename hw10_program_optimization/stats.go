package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type UserEmail struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]UserEmail

func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r) // порционное чтение?
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user UserEmail
		if err = json.Unmarshal([]byte(line), &user); err != nil { // try codegen?
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := strings.Contains(user.Email, "."+domain)
		if matched {
			emailPart := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])

			num := result[emailPart]
			num++
			result[emailPart] = num
		}
	}
	return result, nil
}
