/*******************************************************************************
* Copyright (c) Siemens AG 2023 ALL RIGHTS RESERVED.
*******************************************************************************/

package server

import (
	"time"

	"github.com/rs/zerolog/log"
)

func checkExpired(service serviceEntry) (expired bool) {
	timeExpired := service.timeAdded.Add(time.Duration(serviceExpireTime) * time.Second)
	log.Debug().
		Time("Expiring", timeExpired).
		Time("Service added", service.timeAdded).
		Str("appInstanceIds", service.driver.AppInstanceId).
		Msg("Checking if service entry is expired.")

	return time.Now().After(timeExpired)
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, element := range elements {
		if !encountered[element] {
			encountered[element] = true
			result = append(result, element)
		}
	}

	return result
}
