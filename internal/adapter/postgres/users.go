package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

// UserStore is the gorm-backed account and credential adapter.
type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore { return &UserStore{db: db} }

func portUnique() error { return port.ErrUniqueViolation }

func upsertCredential() clause.OnConflict {
	return clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"password_hash", "updated_at"}),
	}
}

func loadRoles(ctx context.Context, db *gorm.DB, ids []string) (map[string]*entity.Role, error) {
	out := make(map[string]*entity.Role, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	var rows []roleRow
	if err := db.WithContext(ctx).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		role := toRole(row)
		out[row.ID] = &role
	}
	var permissions []struct {
		RoleID string `gorm:"column:role_id"`
		Key    string `gorm:"column:key"`
	}
	if err := db.WithContext(ctx).
		Table("role_permissions").
		Select("role_permissions.role_id, permissions.key").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ?", ids).
		Order("permissions.key ASC").
		Scan(&permissions).Error; err != nil {
		return nil, err
	}
	for _, row := range permissions {
		if role := out[row.RoleID]; role != nil && entity.IsPermission(row.Key) {
			role.Permissions = append(role.Permissions, entity.Permission(row.Key))
		}
	}
	var counts []struct {
		RoleID string `gorm:"column:role_id"`
		Count  int64  `gorm:"column:count"`
	}
	if err := db.WithContext(ctx).
		Model(&userRow{}).
		Select("role_id, count(*) AS count").
		Where("role_id IN ?", ids).
		Group("role_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}
	for _, row := range counts {
		if role := out[row.RoleID]; role != nil {
			role.UserCount = int(row.Count)
		}
	}
	return out, nil
}

func loadUsers(ctx context.Context, db *gorm.DB, rows []userRow) ([]entity.User, error) {
	roleIDs := make([]string, 0, len(rows))
	seen := map[string]bool{}
	for _, row := range rows {
		if row.RoleID != nil && !seen[*row.RoleID] {
			seen[*row.RoleID] = true
			roleIDs = append(roleIDs, *row.RoleID)
		}
	}
	roles, err := loadRoles(ctx, db, roleIDs)
	if err != nil {
		return nil, err
	}
	users := make([]entity.User, 0, len(rows))
	for _, row := range rows {
		var role *entity.Role
		var permissions []entity.Permission
		if row.RoleID != nil {
			role = roles[*row.RoleID]
			if role != nil {
				permissions = role.Permissions
			}
		}
		users = append(users, toUser(row, role, permissions))
	}
	return users, nil
}

func (s *UserStore) CreateWithPassword(ctx context.Context, user *entity.User, passwordHash string) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := userRow{
			ID: user.ID, Name: user.Name, Email: user.Email, EmailVerified: user.EmailVerified,
			Image: ptr(user.Image), RoleID: ptr(user.RoleID), Status: string(user.Status),
			PhoneNumber: ptr(user.PhoneNumber), JobTitle: ptr(user.JobTitle),
			Department: ptr(user.Department), TwoFactorEnabled: user.TwoFactorEnabled,
			CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		return tx.Create(&localCredentialRow{
			UserID: user.ID, PasswordHash: passwordHash,
			CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
		}).Error
	})
	if isUnique(err) {
		return errors.Join(portUnique(), err)
	}
	return err
}

func (s *UserStore) UpdateWithPassword(ctx context.Context, user *entity.User, passwordHash string) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&userRow{}).Where("id = ?", user.ID).Updates(map[string]any{
			"name": user.Name, "email": user.Email, "email_verified": user.EmailVerified,
			"image": ptr(user.Image), "role_id": ptr(user.RoleID), "status": string(user.Status),
			"phone_number": ptr(user.PhoneNumber), "job_title": ptr(user.JobTitle),
			"department": ptr(user.Department), "two_factor_enabled": user.TwoFactorEnabled,
			"updated_at": user.UpdatedAt,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return entity.ErrNotFound
		}
		if passwordHash == "" {
			return nil
		}
		credential := localCredentialRow{
			UserID: user.ID, PasswordHash: passwordHash,
			CreatedAt: user.UpdatedAt, UpdatedAt: user.UpdatedAt,
		}
		return tx.Clauses(upsertCredential()).Create(&credential).Error
	})
	if isUnique(err) {
		return errors.Join(portUnique(), err)
	}
	return err
}

func (s *UserStore) Delete(ctx context.Context, id string) error {
	res := s.db.WithContext(ctx).Delete(&userRow{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var row userRow
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if isNotFound(err) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	users, err := loadUsers(ctx, s.db, []userRow{row})
	if err != nil {
		return nil, err
	}
	return &users[0], nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var row userRow
	if err := s.db.WithContext(ctx).Where("lower(email) = ?", email).First(&row).Error; err != nil {
		if isNotFound(err) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	users, err := loadUsers(ctx, s.db, []userRow{row})
	if err != nil {
		return nil, err
	}
	return &users[0], nil
}

func (s *UserStore) List(ctx context.Context, filters entity.UserFilters) ([]entity.User, int, error) {
	base := s.db.WithContext(ctx).Model(&userRow{})
	if filters.Search != "" {
		like := "%" + filters.Search + "%"
		base = base.Where("users.name ILIKE ? OR users.email ILIKE ?", like, like)
	}
	if filters.Role != "" {
		base = base.Joins("JOIN roles ON roles.id = users.role_id").Where("roles.name = ?", filters.Role)
	}
	if filters.Status != "" {
		base = base.Where("users.status = ?", filters.Status)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageCount := (int(total) + entity.UsersPageSize - 1) / entity.UsersPageSize
	if pageCount < 1 {
		pageCount = 1
	}
	page := filters.Page
	if page < 1 {
		page = 1
	}
	if page > pageCount {
		page = pageCount
	}

	var rows []userRow
	if err := base.Select("users.*").
		Order("users.created_at DESC").
		Offset((page - 1) * entity.UsersPageSize).
		Limit(entity.UsersPageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	users, err := loadUsers(ctx, s.db, rows)
	if err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}

func (s *UserStore) CountSuperadmins(ctx context.Context) (int, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Table("users").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.name = ?", entity.RoleSuperadmin).
		Count(&count).Error
	return int(count), err
}

func (s *UserStore) GetPasswordHash(ctx context.Context, userID string) (string, error) {
	var row localCredentialRow
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&row).Error; err != nil {
		if isNotFound(err) {
			return "", entity.ErrNotFound
		}
		return "", err
	}
	return row.PasswordHash, nil
}
