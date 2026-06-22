package users_postgres_repository

import "github.com/gptikhomirov/go-rest-prod/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainFromModel(model UserModel) domain.User {
	return domain.NewUser(
		model.ID,
		model.Version,
		model.FullName,
		model.PhoneNumber,
	)
}

func userDomainsFromModel(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = userDomainFromModel(user)
	}

	return userDomains
}
