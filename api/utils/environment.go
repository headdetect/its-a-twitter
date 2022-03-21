package utils

import (
	"log"
	"strconv"
	"syscall"
)

func GetIntOrDefault(key string, defaultValue int) (int, bool) {
	val, exists := syscall.Getenv(key)

	if exists {
		num, err := strconv.Atoi(val)
		
		if err != nil {
			log.Fatal("Invalid value: " + key)
		}

		return num, true
	} else {
		return defaultValue, false
	}
}