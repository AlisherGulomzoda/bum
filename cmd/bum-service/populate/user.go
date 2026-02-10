//nolint:gosec //it's ok for generating test data.
package populate

import (
	"context"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"math/rand"
	"sync/atomic"

	"bum-service/internal/domain"
	"bum-service/internal/service/user"
)

const (
	randomUsersCountToGenerate = 100

	testPass = "test"
)

func generateUserCreateArgs() user.AddUserArgs {
	var (
		middleName *string
		phone      *string
	)

	if rand.Int()%5 != 0 {
		m := gofakeit.MiddleName()
		middleName = &m
	}

	prefix := prefixEmailAndPhone.Add(1)

	if rand.Int()%4 == 0 {
		p := fmt.Sprint(prefix) + gofakeit.Phone()
		phone = &p
	}

	return user.AddUserArgs{
		FirstName:  gofakeit.FirstName(),
		LastName:   gofakeit.LastName(),
		MiddleName: middleName,
		Gender:     gofakeit.Gender(),
		Phone:      phone,
		Email:      fmt.Sprint(prefix) + gofakeit.Email(),
		Password:   testPass,
	}
}

var prefixEmailAndPhone atomic.Int64

func (s *Service) genRandomUsers(ctx context.Context) error {
	for userCount := range randomUsersCountToGenerate {
		_ = userCount

		createdUser, err := s.genUser(ctx)
		if err != nil {
			return err
		}

		s.randomUsers = append(s.randomUsers, createdUser)
	}

	return nil
}

var tryCount atomic.Int64

func (s *Service) genUser(ctx context.Context) (domain.User, error) {
	var (
		createdUser domain.User
		err         error
	)

	for try := 1; try <= 10; try++ {
		createdUser, err = s.userService().AddUser(ctx, generateUserCreateArgs())
		if err != nil {
			switch {
			case
				errors.Is(err, domain.ErrUserEmailAlreadyExists),
				errors.Is(err, domain.ErrUserPhoneAlreadyExists):
				tryCount.Add(1)
				fmt.Println("try count ", tryCount.Load())
				continue
			}

			return domain.User{}, fmt.Errorf("failed to create user: %w", err)
		}

		break
	}

	return createdUser, nil
}

func (s *Service) getRandomUser() domain.User {
	return s.randomUsers[rand.Intn(len(s.randomUsers))]
}
