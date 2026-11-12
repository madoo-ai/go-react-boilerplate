// Package health reports liveness and database readiness.
package health

import (
	"context"
)

type Checker interface {
	Ping(ctx context.Context) error
}

type Status struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

type UseCase interface {
	Check(ctx context.Context) Status
}

type useCase struct {
	checker Checker
	version string
}

func New(checker Checker, version string) UseCase {
	return &useCase{checker: checker, version: version}
}

func (u *useCase) Check(ctx context.Context) Status {
	status := Status{Status: "ok", Database: "up", Version: u.version}
	if err := u.checker.Ping(ctx); err != nil {
		status.Status = "degraded"
		status.Database = "down"
	}
	return status
}
