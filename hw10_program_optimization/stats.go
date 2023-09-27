package hw10programoptimization

import (
	"bufio"
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
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var user UserEmail
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}

		matched := strings.Contains(user.Email, "."+domain)
		if matched {
			emailPart := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[emailPart] += 1
		}

	}

	return result, nil
}
