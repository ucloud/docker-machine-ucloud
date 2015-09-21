package ucloud

import (
	"errors"
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
