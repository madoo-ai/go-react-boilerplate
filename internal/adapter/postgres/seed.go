package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-react-boilerplate/internal/entity"
)

type SeedResult struct {
	Permissions int
	Roles       int
	NewUsers    int
}

type demoUser struct {
	Name       string
	Email      string
	Role       string
	Status     entity.UserStatus
	JobTitle   string
	Department string
}

var demoUsers = []demoUser{
	{"Super Admin", "superadmin@example.com", entity.RoleSuperadmin, entity.StatusActive, "Platform Owner", "Operations"},
	{"HR Staff", "hr@example.com", entity.RoleHR, entity.StatusActive, "People Operations", "People"},
	{"Regular User", "user@example.com", entity.RoleUser, entity.StatusActive, "Software Engineer", "Engineering"},
	{"Andi Pratama", "andi@example.com", entity.RoleUser, entity.StatusActive, "Backend Engineer", "Engineering"},
	{"Siti Rahma", "siti@example.com", entity.RoleHR, entity.StatusActive, "Recruiter", "People"},
	{"Budi Santoso", "budi@example.com", entity.RoleUser, entity.StatusActive, "Finance Analyst", "Finance"},
	{"Dewi Lestari", "dewi@example.com", entity.RoleUser, entity.StatusActive, "Product Designer", "Design"},
	{"Rizky Hidayat", "rizky@example.com", entity.RoleUser, entity.StatusActive, "Account Executive", "Sales"},
	{"Maya Putri", "maya@example.com", entity.RoleUser, entity.StatusActive, "Support Lead", "Support"},
	{"Fajar Nugroho", "fajar@example.com", entity.RoleUser, entity.StatusActive, "Data Engineer", "Engineering"},
	{"Nabila Syah", "nabila@example.com", entity.RoleUser, entity.StatusPending, "Marketing Associate", "Marketing"},
	{"Agus Wijaya", "agus@example.com", entity.RoleUser, entity.StatusPending, "Operations Associate", "Operations"},
	{"Putri Anggraini", "putri@example.com", entity.RoleUser, entity.StatusPending, "QA Engineer", "Engineering"},
	{"Hendra Kusuma", "hendra@example.com", entity.RoleUser, entity.StatusSuspended, "Sales Representative", "Sales"},
}

// Seed is idempotent. It creates the permission registry, the three system
// roles, the registration-approval setting and enough demo accounts for the
// dashboard to look real. Existing demo users are left alone.
func Seed(ctx context.Context, db *gorm.DB, passwordHash string) (SeedResult, error) {
	result := SeedResult{}
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()
		permissionIDs := map[entity.Permission]string{}
		for index, permission := range entity.AllPermissions {
			var row permissionRow
			err := tx.Where("key = ?", string(permission)).First(&row).Error
			switch {
			case err == nil:
				if row.Label != entity.PermissionLabels[permission] {
					if err := tx.Model(&permissionRow{}).Where("id = ?", row.ID).Update("label", entity.PermissionLabels[permission]).Error; err != nil {
						return err
					}
				}
			case isNotFound(err):
				row = permissionRow{
					ID: entity.NewUUID(), Key: string(permission), Label: entity.PermissionLabels[permission],
					CreatedAt: now.Add(time.Duration(index) * time.Millisecond), UpdatedAt: now,
				}
				if err := tx.Create(&row).Error; err != nil {
					return err
				}
			default:
				return err
			}
			permissionIDs[permission] = row.ID
			result.Permissions++
		}

		roleIDs := map[string]string{}
		for index, definition := range entity.RoleDefinitions {
			var row roleRow
			err := tx.Where("name = ?", definition.Name).First(&row).Error
			switch {
			case err == nil:
				if err := tx.Model(&roleRow{}).Where("id = ?", row.ID).Updates(map[string]any{
					"label": definition.Label, "description": ptr(definition.Description),
					"is_system": definition.IsSystem, "updated_at": now,
				}).Error; err != nil {
					return err
				}
			case isNotFound(err):
				row = roleRow{
					ID: entity.NewUUID(), Name: definition.Name, Label: definition.Label,
					Description: ptr(definition.Description), IsSystem: definition.IsSystem,
					CreatedAt: now.Add(time.Duration(index) * time.Millisecond), UpdatedAt: now,
				}
				if err := tx.Create(&row).Error; err != nil {
					return err
				}
			default:
				return err
			}
			roleIDs[definition.Name] = row.ID
			if err := tx.Where("role_id = ?", row.ID).Delete(&rolePermissionRow{}).Error; err != nil {
				return err
			}
			for _, permission := range definition.Permissions {
				if id, ok := permissionIDs[permission]; ok {
					if err := tx.Create(&rolePermissionRow{RoleID: row.ID, PermissionID: id}).Error; err != nil {
						return err
					}
				}
			}
			result.Roles++
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&settingRow{
			Key: entity.SettingRegistrationApproval, Value: ptr("false"), CreatedAt: now, UpdatedAt: now,
		}).Error; err != nil {
			return err
		}

		day := 24 * time.Hour
		for index, demo := range demoUsers {
			createdAt := now.Add(-time.Duration(index*2+1) * day)
			user := userRow{
				ID: entity.NewUUID(), Name: demo.Name, Email: demo.Email,
				EmailVerified: demo.Status == entity.StatusActive, RoleID: ptr(roleIDs[demo.Role]),
				Status: string(demo.Status), JobTitle: ptr(demo.JobTitle), Department: ptr(demo.Department),
				CreatedAt: createdAt, UpdatedAt: createdAt,
			}
			inserted := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
			if inserted.Error != nil {
				return inserted.Error
			}
			if inserted.RowsAffected == 0 {
				continue
			}
			if err := tx.Create(&localCredentialRow{
				UserID: user.ID, PasswordHash: passwordHash, CreatedAt: createdAt, UpdatedAt: createdAt,
			}).Error; err != nil {
				return err
			}
			result.NewUsers++
		}
		return nil
	})
	return result, err
}
