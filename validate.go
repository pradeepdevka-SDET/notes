package main

import (
	"errors"
	"strings"
)

func validateTitile(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title is required")
	}
	return nil
}
