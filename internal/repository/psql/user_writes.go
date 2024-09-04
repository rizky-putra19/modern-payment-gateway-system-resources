package psql

import (
	"github.com/hypay-id/backend-dashboard-hypay/internal/constant"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/jmoiron/sqlx"
)

type UsersWrites struct {
	db *sqlx.DB
}

func NewUsersWrites(db *sqlx.DB) *UsersWrites {
	return &UsersWrites{
		db: db,
	}
}

func (uw *UsersWrites) CreateUsersMerchantRepo(payload dto.InviteMerchantUserDto, credentials dto.EmailDataHtmlDto) (int, error) {
	var userId int

	query := `
	INSERT INTO users (username, email, password, pin, merchantid, role_id, status, user_type, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')
	RETURNING id
	`

	row := uw.db.QueryRow(query, payload.Email, payload.Email, credentials.Password, credentials.Pin, payload.MerchantId, payload.RolesId, "ACTIVE", constant.UserMerchant)
	err := row.Scan(&userId)
	if err != nil || userId == 0 {
		return userId, err
	}

	return userId, nil
}

func (uw *UsersWrites) UpdatePassOrPinRepo(passHash string, pinHash string, username string) error {
	query := `
	UPDATE users
	SET
		password = $1,
		pin = $2,
		updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'
	WHERE username = $3;
	`

	_, err := uw.db.Exec(query, passHash, pinHash, username)
	if err != nil {
		return err
	}
	return nil
}
