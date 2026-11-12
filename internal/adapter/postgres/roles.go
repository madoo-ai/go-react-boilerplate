package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"go-react-boilerplate/internal/app/port"
	"go-react-boilerplate/internal/entity"
)

// RoleStore is the gorm-backed role and permission adapter.
type RoleStore struct {
	db *gorm.DB
}

func NewRoleStore(db *gorm.DB) *RoleStore { return &RoleStore{db: db} }

func (s *RoleStore) List(ctx context.Context) ([]entity.Role, error) {
	var rows []roleRow
	if err := s.db.WithContext(ctx).Order("name ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	roles, err := loadRoles(ctx, s.db, ids)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Role, 0, len(rows))
	for _, row := range rows {
		if role := roles[row.ID]; role != nil {
			out = append(out, *role)
		}
	}
	return out, nil
}

func (s *RoleStore) GetByID(ctx context.Context, id string) (*entity.Role, error) {
	var row roleRow
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if isNotFound(err) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	roles, err := loadRoles(ctx, s.db, []string{row.ID})
	if err != nil {
		return nil, err
	}
	role := roles[row.ID]
	if role == nil {
		return nil, entity.ErrNotFound
	}
	return role, nil
}

func (s *RoleStore) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	var row roleRow
	if err := s.db.WithContext(ctx).Where("name = ?", name).First(&row).Error; err != nil {
		if isNotFound(err) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	roles, err := loadRoles(ctx, s.db, []string{row.ID})
	if err != nil {
		return nil, err
	}
	role := roles[row.ID]
	if role == nil {
		return nil, entity.ErrNotFound
	}
	return role, nil
}

func (s *RoleStore) Create(ctx context.Context, role *entity.Role) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := roleRow{
			ID: role.ID, Name: role.Name, Label: role.Label,
			Description: ptr(role.Description), IsSystem: role.IsSystem,
			CreatedAt: role.CreatedAt, UpdatedAt: role.UpdatedAt,
		}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		permissionIDs, err := resolvePermissionIDs(ctx, tx, role.Permissions)
		if err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			if err := tx.Create(&rolePermissionRow{RoleID: role.ID, PermissionID: permissionID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if isUnique(err) {
		return errors.Join(port.ErrUniqueViolation, err)
	}
	return err
}

func (s *RoleStore) Update(ctx context.Context, role *entity.Role) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&roleRow{}).Where("id = ?", role.ID).Updates(map[string]any{
			"name": role.Name, "label": role.Label, "description": ptr(role.Description),
			"updated_at": role.UpdatedAt,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return entity.ErrNotFound
		}
		if err := tx.Where("role_id = ?", role.ID).Delete(&rolePermissionRow{}).Error; err != nil {
			return err
		}
		permissionIDs, err := resolvePermissionIDs(ctx, tx, role.Permissions)
		if err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			if err := tx.Create(&rolePermissionRow{RoleID: role.ID, PermissionID: permissionID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if isUnique(err) {
		return errors.Join(port.ErrUniqueViolation, err)
	}
	return err
}

func (s *RoleStore) Delete(ctx context.Context, id string) error {
	res := s.db.WithContext(ctx).Delete(&roleRow{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (s *RoleStore) ListPermissions(ctx context.Context) ([]entity.PermissionInfo, error) {
	var rows []permissionRow
	if err := s.db.WithContext(ctx).Order("key ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]entity.PermissionInfo, 0, len(rows))
	for _, row := range rows {
		if !entity.IsPermission(row.Key) {
			continue
		}
		label := row.Label
		if label == "" {
			label = entity.PermissionLabels[entity.Permission(row.Key)]
		}
		out = append(out, entity.PermissionInfo{Key: entity.Permission(row.Key), Label: label})
	}
	return out, nil
}

func resolvePermissionIDs(ctx context.Context, tx *gorm.DB, permissions []entity.Permission) ([]string, error) {
	if len(permissions) == 0 {
		return nil, nil
	}
	keys := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		if entity.IsPermission(string(permission)) {
			keys = append(keys, string(permission))
		}
	}
	if len(keys) == 0 {
		return nil, nil
	}
	var rows []permissionRow
	if err := tx.WithContext(ctx).Where("key IN ?", keys).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}
