package api

import (
	"github.com/Uspacy/uspacy-go-sdk/auth"
	"github.com/Uspacy/uspacy-go-sdk/crm"
	"github.com/Uspacy/uspacy-go-sdk/group"
	"github.com/Uspacy/uspacy-go-sdk/newsfeed"
	"github.com/Uspacy/uspacy-go-sdk/task"
)

var emptyHeaders map[string]string

type Service string

func (s Service) getService() string {
	service := map[Service]string{
		"crm":      crm.VersionUrl,
		"auth":     auth.VersionUrl,
		"tasks":    task.VersionUrl,
		"newsfeed": newsfeed.VersionUrl,
		"groups":   group.VersionUrl,
	}

	if entity, ok := service[s]; !ok {
		return "unknown"
	} else {
		return entity
	}
}
