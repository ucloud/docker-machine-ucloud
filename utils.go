package ucloud

import (
	"errors"
	"math/rand"
	"time"
)

var (
	errInvalidRegion = errors.New("invalid region specified")
)

var regions = []string{
	"cn-north-01",
	"cn-north-02",
	"cn-north-03",
	"cn-east-01",
	"cn-south-01",
	"cn-south-02",
	"hk-01",
	"us-west-01",
}

func validateUCloudRegion(region string) (string, error) {
	for _, v := range regions {
		if v == region {
			return region, nil
		}
	}

	return "", errInvalidRegion
}

func generateRandomPassword(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+}{:?><")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
