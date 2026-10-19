package postgresrepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/postgresql/sqlc_gen"
	"github.com/google/uuid"
)

type UserPersistanceAdapter struct {
	conn *sql.DB
	q    *sqlc_gen.Queries
}

func NewUserPersistanceAdapter(conn *sql.DB) *UserPersistanceAdapter {
	return &UserPersistanceAdapter{
		conn: conn,
		q:    sqlc_gen.New(conn),
	}
}

func (adp *UserPersistanceAdapter) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	userData, err := adp.q.GetUserDataByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerrors.NotFound
		}
		return nil, err
	}

	return &models.User{
		ID:           userData.ID,
		Email:        userData.Email,
		PasswordHash: userData.PasswordHash,
		GroupID:      userData.GroupID,
		CreatedAt:    userData.CreatedAt,
	}, nil
}

func (adp *UserPersistanceAdapter) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	userData, err := adp.q.GetUserDataByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           userData.ID,
		Email:        userData.Email,
		PasswordHash: userData.PasswordHash,
		GroupID:      userData.GroupID,
		CreatedAt:    userData.CreatedAt,
	}, nil
}

func (adp *UserPersistanceAdapter) GetUsers(ctx context.Context, createdAt time.Time, count int) ([]*models.User, error) {
	userData, err := adp.q.GetUserList(ctx, sqlc_gen.GetUserListParams{
		CreatedAt: createdAt,
		Limit:     int32(count),
	})
	if err != nil {
		return nil, err
	}

	userList := make([]*models.User, 0, count)

	for _, userRecord := range userData {
		userList = append(userList, &models.User{
			ID:           userRecord.ID,
			Email:        userRecord.Email,
			PasswordHash: userRecord.PasswordHash,
			GroupID:      userRecord.GroupID,
			CreatedAt:    userRecord.CreatedAt,
		})
	}

	return userList, nil
}

func (adp *UserPersistanceAdapter) CreateUser(ctx context.Context, userData *models.UserData) (*models.User, error) {
	tx, err := adp.conn.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:   userData.Email,
		GroupID: userData.GroupID,
	}

	err = user.SetPassword(userData.Password)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback())
	}

	dbData, err := adp.q.WithTx(tx).CreateUser(ctx, sqlc_gen.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		GroupID:      user.GroupID,
	})
	if err != nil {
		return nil, errors.Join(err, tx.Rollback())
	}

	user.CreatedAt = dbData.CreatedAt
	user.ID = dbData.ID

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (adp *UserPersistanceAdapter) UpdateUser(ctx context.Context, userData *models.UserUpdateRequestData) error {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	updatedUser := models.User{
		ID:      userData.UserToChangeID,
		Email:   userData.NewEmail,
		GroupID: userData.NewGroupID,
	}

	err = updatedUser.SetPassword(userData.NewPassword)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	updatedUser.CreatedAt, err = adp.q.WithTx(tx).UpdateUser(ctx, sqlc_gen.UpdateUserParams{
		Email:        updatedUser.Email,
		PasswordHash: updatedUser.PasswordHash,
		GroupID:      updatedUser.GroupID,
		ID:           updatedUser.ID,
	})

	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

func (adp *UserPersistanceAdapter) CheckIdempotencyKey(ctx context.Context, key uuid.UUID) error {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	wasCreated, err := adp.q.WithTx(tx).CheckUserWasCreated(ctx, key)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if wasCreated {
		return errors.New("user was created")
	}
	return nil
}

func (adp *UserPersistanceAdapter) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = adp.q.WithTx(tx).DeleteUser(ctx, userID)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return nil
}
