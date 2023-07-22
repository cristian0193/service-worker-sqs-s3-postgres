package utils

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
	"strconv"
)

func GetString(name string) (string, error) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("env var %s not found", name)
	}
	return v, nil
}

func GetInt(name string) (int, error) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return 0, fmt.Errorf("env var %s not found", name)
	}
	intV, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("en var %s must be a number", name)
	}
	return intV, nil
}

func GetParam(c echo.Context, name string) (string, error) {
	strParam := c.Param(name)
	if strParam == "" {
		return strParam, errors.New(fmt.Sprintf("param '%s' is required", name))
	}
	return strParam, nil
}
