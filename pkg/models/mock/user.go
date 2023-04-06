package mock

import (
	"errors"
	"time"

	"github.com/JohanVong/online_bazaar/pkg/models"
)

type UserModel struct{}

var mockUser = &models.UserOutput{
	UserUID:    "uuid.v6[1]",
	Username:   "TestUser",
	Hash:       "hzNDoZShWoQPmw9HmK1RvVeE8PtMJpDHR4ru5+QVnwL0NdqVBUmb7x7rUDahYMBSfTS3zzJg7WE7DIJBexaWWQ==",
	Email:      "testuser@mail.test",
	Phone:      "87776665544",
	CountryUID: "uuid.v6[2]",
	Country:    "TestCountry",
	HistoryUID: "uuid.v6[3]",
	CreatedAt:  time.Now(),
	UpdatedAt:  *new(time.Time),
	DeletedAt:  *new(time.Time),
}

var mockUserDeleted = &models.UserOutput{
	UserUID:    "uuid.v6[4]",
	Username:   "DeletedUser",
	Hash:       "hzNDoZShWoQPmw9HmK1RvVeE8PtMJpDHR4ru5+QVnwL0NdqVBUmb7x7rUDahYMBSfTS3zzJg7WE7DIJBexaWWQ==",
	Email:      "deleteduser@mail.test",
	Phone:      "87771112233",
	CountryUID: "uuid.v6[2]",
	Country:    "TestCountry",
	HistoryUID: "uuid.v6[5]",
	CreatedAt:  time.Now(),
	UpdatedAt:  *new(time.Time),
	DeletedAt:  time.Now(),
}

func (u *UserModel) Insert(input *models.UserSignupInput) error {
	if input.Username == "Exists" {
		return errors.New("duplicate key value violates unique constraint")
	}

	if input.Country != "TestCountry" {
		return errors.New("Provided country does not exist")
	}

	return nil
}

func (u *UserModel) Get(key string, byPK bool) (*models.UserOutput, error) {
	switch {
	case key == "uuid.v6[1]" && byPK:
		time.Sleep(time.Millisecond * 100)
		return mockUser, nil

	case key == "TestUser" && !byPK:
		time.Sleep(time.Millisecond * 200)
		return mockUser, nil

	case key == "uuid.v6[4]" && byPK:
		time.Sleep(time.Millisecond * 100)
		return mockUserDeleted, nil

	case key == "DeletedUser" && !byPK:
		time.Sleep(time.Millisecond * 200)
		return mockUserDeleted, nil

	case key == "panic":
		panic("test panic!")

	default:
		return nil, errors.New("No record found")
	}
}

func (u *UserModel) Update(uid string, input *models.UserUpdateInput) error {
	counter := 0

	if input.Email != "" {
		counter++
	}

	if input.Phone != "" {
		counter++
	}

	if input.Country != "" {
		counter++
	}

	if counter <= 0 {
		return errors.New("Nothing to update")
	}

	return nil
}

func (u *UserModel) UpdatePassword(uid string, input *models.UpdateUserPasswordInput) error {
	if uid != "uuid.v6[1]" {
		return errors.New("No record found")
	}

	return nil
}

func (u *UserModel) Delete(uid string) error {
	if uid != "uuid.v6[1]" {
		return errors.New("No record found")
	}

	return nil
}
