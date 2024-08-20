package psql

import (
	"database/sql"

	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"github.com/jmoiron/sqlx"
)

type UserReads struct {
	db *sqlx.DB
}

func NewUsersReads(db *sqlx.DB) *UserReads {
	return &UserReads{
		db: db,
	}
}

func (u *UserReads) GetUserByUsername(username string) (entity.User, error) {
	var user = entity.User{}
	query := `
	SELECT
		u.id AS user_id,
		u.email,
		u.username,
		u.user_type,
		u.password,
		u.pin,
		m.merchant_id,
		m.merchant_name,
		m.merchant_secret,
		m.currency,
		m.status AS merchant_status,
		r.id AS role_id,
		r.role_name,
		u.status AS user_status,
		u.created_at AS user_created_at,
		u.updated_at AS user_updated_at
	FROM
		users u
	LEFT JOIN
		merchants m ON u.merchantid = m.merchant_id
	LEFT JOIN
		roles r ON u.role_id = r.id
	WHERE
		u.username = $1;
	`

	err := u.db.Get(&user, query, username)
	if err != nil && err != sql.ErrNoRows {
		slog.Errorw("unexpected error", "stack_trace", err.Error())
		return user, err
	}

	return user, nil
}

func (u *UserReads) GetPermissionByRoleId(id int) ([]entity.Permission, error) {
	var permissions []entity.Permission

	query := `
	SELECT
		p.ID AS permission_id,
		p.permission_desc
	FROM
		roles r
		JOIN role_permissions rp ON r.ID = rp.role_id
		JOIN permissions p ON rp.permission_id = p.ID
	WHERE
		r.ID = $1;
	`

	err := u.db.Select(&permissions, query, id)
	if err != nil && err != sql.ErrNoRows {
		slog.Errorw("unexpected error", "stack_trace", err.Error())
		return permissions, err
	}

	return permissions, nil
}

func (u *UserReads) GetRolesRepo() ([]entity.RolesEntity, error) {
	var listRoles []entity.RolesEntity

	query := `
	SELECT
		rl.ID,
		rl.role_name,
		p.permission_desc
	FROM
		roles rl
		JOIN role_permissions rp ON rp.role_id = rl.ID
		JOIN permissions p ON p.ID = rp.permission_id
	ORDER BY rl.created_at DESC
	`

	err := u.db.Select(&listRoles, query)
	if err != nil {
		return listRoles, err
	}

	return listRoles, nil
}

func (u *UserReads) GetListUserByMerchantIdRepo(merchantId string) ([]entity.ListUsersEntity, error) {
	var listUsers []entity.ListUsersEntity

	query := `
	SELECT
		u.id,
		u.username,
		u.email,
		u.created_at,
		u.status,
		r.role_name
	FROM
		users u
		JOIN roles r ON r.ID = u.role_id
		JOIN merchants m ON u.merchantid = m.merchant_id
	WHERE
		m.merchant_id = $1;
	`

	err := u.db.Select(&listUsers, query, merchantId)
	if err != nil && err != sql.ErrNoRows {
		slog.Errorw("unexpected error", "stack_trace", err.Error())
		return listUsers, err
	}

	return listUsers, nil
}
