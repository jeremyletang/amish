package domain

import (
	"strconv"
	"time"

	"github.com/google/go-github/github"
)

type User struct {
	Id        string
	Login     string
	Name      string
	Company   string
	AvatarUrl string
	Location  string
	Blog      string
	Email     string

	ContentUpdated int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func fromStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func UpdateUserWithGithubUser(user *User, githubUser *github.User) {
	user.Name = fromStringPtr(githubUser.Name)
	user.Company = fromStringPtr(githubUser.Company)
	user.Email = fromStringPtr(githubUser.Email)
	user.Location = fromStringPtr(githubUser.Location)
	user.Blog = fromStringPtr(githubUser.Blog)

}

func NewUserFromGithubUser(u *github.User) *User {
	return &User{
		Id:        strconv.Itoa(*u.ID),
		Login:     fromStringPtr(u.Login),
		Name:      fromStringPtr(u.Name),
		Company:   fromStringPtr(u.Company),
		AvatarUrl: fromStringPtr(u.AvatarURL),
		Location:  fromStringPtr(u.Location),
		Blog:      fromStringPtr(u.Blog),
		Email:     fromStringPtr(u.Email),
	}
}

type Repository struct {
	Id    string
	Owner string
	Name  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c Repository) TableName() string {
	return "repositories"
}

type Star struct {
	Id           string
	RepositoryId string
	UserId       string
	StarredAt    time.Time
	Valid        int

	CreatedAt time.Time
	UpdatedAt time.Time
}
