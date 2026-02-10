package database

import (
	"context"
	"spahtmx/internal/domain"

	"github.com/uptrace/bun"
)

type UserBunRepository struct {
	DB *bun.DB
}

type UserBun struct {

	bun.BaseModel `bun:"table:users"`

	ID       int64		 `bun:"id,pk,autoincrement"`
	Username string       
	Password string       
	Email    string       
	Status   bool         
}

func ToUserDomain(u UserBun) domain.User {

	return domain.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Status:   u.Status,
	}

}

func FromUserDomain(user domain.User) (*UserBun, error) {

	return &UserBun{
		ID:       user.ID, // This is a placeholder. You should implement a proper conversion from string to int64.
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Status:   user.Status,
	}, nil
}

func (r UserBunRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	var users []UserBun
	err := r.DB.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}

	var domainUsers []domain.User
	for _, u := range users {
		domainUsers = append(domainUsers, ToUserDomain(u))
	}
	return domainUsers, nil
}

func (r UserBunRepository) GetUser(ctx context.Context, id string) (domain.User, error) {

	var user UserBun
	err := r.DB.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return domain.User{}, err
	}

	return ToUserDomain(user), nil
}



func (r UserBunRepository) GetByUsername(ctx context.Context, username string) (domain.User, error)	 {

	var user UserBun
	err := r.DB.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return domain.User{}, err
	}

	return ToUserDomain(user), nil
}


func (r UserBunRepository) CreateUser(ctx context.Context, user domain.User) error {

	userBun, err := FromUserDomain(user)
	if err != nil {
		return err
	}

	_, err = r.DB.NewInsert().Model(userBun).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}	

func (r UserBunRepository) UpdateUser(ctx context.Context, user domain.User) error {

	userBun, err := FromUserDomain(user)
	if err != nil {
		return err
	}

	_, err = r.DB.NewUpdate().Model(userBun).Where("id = ?", userBun.ID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

	