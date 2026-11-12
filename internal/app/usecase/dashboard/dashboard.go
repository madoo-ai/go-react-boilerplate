// Package dashboard assembles the admin overview.
package dashboard

import (
	"context"
	"time"
	_ "time/tzdata"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

type UseCase interface {
	Get(ctx context.Context, session *entity.SessionContext) (entity.DashboardData, error)
}

type useCase struct {
	store port.DashboardStore
	clock port.Clock
}

func New(store port.DashboardStore, clock port.Clock) UseCase {
	return &useCase{store: store, clock: clock}
}

func (u *useCase) Get(ctx context.Context, session *entity.SessionContext) (entity.DashboardData, error) {
	if !session.Actor().Can(entity.PermissionViewAdmin) {
		return entity.DashboardData{}, entity.ErrForbidden
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.UTC
	}
	now := u.clock.Now().In(location)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location).AddDate(0, 0, -29)

	total, active, pending, err := u.store.Stats(ctx)
	if err != nil {
		return entity.DashboardData{}, err
	}
	roles, err := u.store.RoleCounts(ctx)
	if err != nil {
		return entity.DashboardData{}, err
	}
	signups, err := u.store.SignupTimes(ctx, start)
	if err != nil {
		return entity.DashboardData{}, err
	}
	recent, err := u.store.RecentUsers(ctx, 5)
	if err != nil {
		return entity.DashboardData{}, err
	}

	buckets := make(map[string]int, 30)
	for index := 0; index < 30; index++ {
		day := start.AddDate(0, 0, index)
		buckets[day.Format("2006-01-02")] = 0
	}
	for _, created := range signups {
		key := created.In(location).Format("2006-01-02")
		if _, ok := buckets[key]; ok {
			buckets[key]++
		}
	}
	series := make([]entity.SignupDatum, 0, 30)
	cumulative := 0
	for index := 0; index < 30; index++ {
		day := start.AddDate(0, 0, index)
		key := day.Format("2006-01-02")
		cumulative += buckets[key]
		series = append(series, entity.SignupDatum{Date: key, Users: cumulative})
	}

	return entity.DashboardData{
		TotalUsers: total, ActiveUsers: active, PendingUsers: pending,
		RolesCount: len(roles), RoleData: roles, SignupData: series, RecentUsers: recent,
	}, nil
}
