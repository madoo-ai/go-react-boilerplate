package postgres

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-react-boilerplate/internal/entity"
)

// SessionStore is the gorm-backed session adapter.
type SessionStore struct {
	db *gorm.DB
}

func NewSessionStore(db *gorm.DB) *SessionStore { return &SessionStore{db: db} }

func (s *SessionStore) Create(ctx context.Context, session *entity.Session) error {
	row := sessionRow{
		ID: session.ID, UserID: session.UserID, TokenHash: session.TokenHash,
		ExpiresAt: session.ExpiresAt, CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt, IPAddress: ptr(session.IPAddress),
		UserAgent: ptr(session.UserAgent),
	}
	return s.db.WithContext(ctx).Create(&row).Error
}

func (s *SessionStore) GetContext(ctx context.Context, tokenHash string) (*entity.SessionContext, error) {
	var row sessionRow
	if err := s.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&row).Error; err != nil {
		if isNotFound(err) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	users := NewUserStore(s.db)
	user, err := users.GetByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}
	session := entity.Session{
		ID: row.ID, UserID: row.UserID, TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		IPAddress: deref(row.IPAddress), UserAgent: deref(row.UserAgent),
	}
	return &entity.SessionContext{
		Session: session, User: *user, Permissions: user.Permissions,
	}, nil
}

func (s *SessionStore) UpdateExpiry(ctx context.Context, id string, expiresAt time.Time) error {
	res := s.db.WithContext(ctx).Model(&sessionRow{}).Where("id = ?", id).Updates(map[string]any{
		"expires_at": expiresAt, "updated_at": time.Now().UTC(),
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (s *SessionStore) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	res := s.db.WithContext(ctx).Where("token_hash = ?", tokenHash).Delete(&sessionRow{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (s *SessionStore) DeleteByUserID(ctx context.Context, userID string) error {
	return s.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&sessionRow{}).Error
}

func (s *SessionStore) DeleteOthers(ctx context.Context, userID, keepSessionID string) error {
	return s.db.WithContext(ctx).
		Where("user_id = ? AND id <> ?", userID, keepSessionID).
		Delete(&sessionRow{}).Error
}

func (s *SessionStore) DeleteExpired(ctx context.Context, now time.Time) error {
	return s.db.WithContext(ctx).Where("expires_at < ?", now).Delete(&sessionRow{}).Error
}

// SettingStore is the gorm-backed key/value settings adapter.
type SettingStore struct {
	db *gorm.DB
}

func NewSettingStore(db *gorm.DB) *SettingStore { return &SettingStore{db: db} }

func (s *SettingStore) Get(ctx context.Context, key string) (string, error) {
	var row settingRow
	if err := s.db.WithContext(ctx).Where("key = ?", key).First(&row).Error; err != nil {
		if isNotFound(err) {
			return "", entity.ErrNotFound
		}
		return "", err
	}
	return deref(row.Value), nil
}

func (s *SettingStore) Set(ctx context.Context, key, value string) error {
	row := settingRow{Key: key, Value: ptr(value), UpdatedAt: time.Now().UTC()}
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&row).Error
}

// VerificationStore is the gorm-backed one-shot token adapter.
type VerificationStore struct {
	db *gorm.DB
}

func NewVerificationStore(db *gorm.DB) *VerificationStore { return &VerificationStore{db: db} }

func (s *VerificationStore) Create(ctx context.Context, verification *entity.Verification) error {
	row := verificationRow{
		ID: verification.ID, Identifier: verification.Identifier,
		Purpose: string(verification.Purpose), TokenHash: verification.TokenHash,
		ExpiresAt: verification.ExpiresAt, CreatedAt: verification.CreatedAt,
	}
	return s.db.WithContext(ctx).Create(&row).Error
}

func (s *VerificationStore) Consume(ctx context.Context, tokenHash string, purpose entity.VerificationPurpose, now time.Time) (*entity.Verification, error) {
	var verification *entity.Verification
	expired := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row verificationRow
		if err := tx.Where("token_hash = ? AND purpose = ?", tokenHash, string(purpose)).First(&row).Error; err != nil {
			if isNotFound(err) {
				return entity.ErrNotFound
			}
			return err
		}
		if err := tx.Where("id = ?", row.ID).Delete(&verificationRow{}).Error; err != nil {
			return err
		}
		if !now.Before(row.ExpiresAt) {
			expired = true
			return nil
		}
		verification = &entity.Verification{
			ID: row.ID, Identifier: row.Identifier, Purpose: entity.VerificationPurpose(row.Purpose),
			TokenHash: row.TokenHash, ExpiresAt: row.ExpiresAt, CreatedAt: row.CreatedAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if expired || verification == nil {
		return nil, entity.ErrNotFound
	}
	return verification, nil
}

func (s *VerificationStore) DeleteForIdentifier(ctx context.Context, identifier string, purpose entity.VerificationPurpose) error {
	return s.db.WithContext(ctx).
		Where("identifier = ? AND purpose = ?", identifier, string(purpose)).
		Delete(&verificationRow{}).Error
}

// TwoFactorStore is the gorm-backed TOTP adapter.
type TwoFactorStore struct {
	db *gorm.DB
}

func NewTwoFactorStore(db *gorm.DB) *TwoFactorStore { return &TwoFactorStore{db: db} }

func (s *TwoFactorStore) GetByUserID(ctx context.Context, userID string) (*entity.TwoFactor, error) {
	var row twoFactorRow
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&row).Error; err != nil {
		if isNotFound(err) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	codes := []string{}
	if err := json.Unmarshal([]byte(row.BackupCodes), &codes); err != nil {
		return nil, err
	}
	return &entity.TwoFactor{
		ID: row.ID, UserID: row.UserID, Secret: row.Secret, BackupCodeHashes: codes,
		Verified: row.Verified, FailedVerificationCount: row.FailedVerificationCount,
		LockedUntil: row.LockedUntil, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, nil
}

func (s *TwoFactorStore) Save(ctx context.Context, twoFactor *entity.TwoFactor) error {
	encoded, err := json.Marshal(twoFactor.BackupCodeHashes)
	if err != nil {
		return err
	}
	row := twoFactorRow{
		ID: twoFactor.ID, UserID: twoFactor.UserID, Secret: twoFactor.Secret,
		BackupCodes: string(encoded), Verified: twoFactor.Verified,
		FailedVerificationCount: twoFactor.FailedVerificationCount,
		LockedUntil:             twoFactor.LockedUntil, CreatedAt: twoFactor.CreatedAt,
		UpdatedAt: twoFactor.UpdatedAt,
	}
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"secret", "backup_codes", "verified", "failed_verification_count", "locked_until", "updated_at",
		}),
	}).Create(&row).Error
}

func (s *TwoFactorStore) DeleteByUserID(ctx context.Context, userID string) error {
	return s.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&twoFactorRow{}).Error
}

// DashboardStore is the gorm-backed overview adapter.
type DashboardStore struct {
	db *gorm.DB
}

func NewDashboardStore(db *gorm.DB) *DashboardStore { return &DashboardStore{db: db} }

func (s *DashboardStore) Stats(ctx context.Context) (int, int, int, error) {
	var row struct {
		Total   int64 `gorm:"column:total"`
		Active  int64 `gorm:"column:active"`
		Pending int64 `gorm:"column:pending"`
	}
	err := s.db.WithContext(ctx).Model(&userRow{}).
		Select("count(*) AS total, count(*) FILTER (WHERE status = 'active') AS active, count(*) FILTER (WHERE status = 'pending') AS pending").
		Scan(&row).Error
	return int(row.Total), int(row.Active), int(row.Pending), err
}

func (s *DashboardStore) RoleCounts(ctx context.Context) ([]entity.RoleDatum, error) {
	var rows []struct {
		Name  string `gorm:"column:name"`
		Label string `gorm:"column:label"`
		Users int64  `gorm:"column:users"`
	}
	err := s.db.WithContext(ctx).
		Table("roles").
		Select("roles.name, roles.label, count(users.id) AS users").
		Joins("LEFT JOIN users ON users.role_id = roles.id").
		Group("roles.id, roles.name, roles.label").
		Order("roles.name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.RoleDatum, 0, len(rows))
	for _, row := range rows {
		out = append(out, entity.RoleDatum{Name: row.Name, Label: row.Label, Users: int(row.Users)})
	}
	return out, nil
}

func (s *DashboardStore) SignupTimes(ctx context.Context, since time.Time) ([]time.Time, error) {
	var times []time.Time
	err := s.db.WithContext(ctx).Model(&userRow{}).
		Where("created_at >= ?", since).
		Order("created_at ASC").
		Pluck("created_at", &times).Error
	return times, err
}

func (s *DashboardStore) RecentUsers(ctx context.Context, limit int) ([]entity.User, error) {
	var rows []userRow
	if err := s.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	return loadUsers(ctx, s.db, rows)
}
