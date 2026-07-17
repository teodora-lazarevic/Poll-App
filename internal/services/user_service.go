package services

import (
	"context"
	"errors"
	"regexp"

	"github.com/teodora-lazarevic/Poll-App/ent"
	"github.com/teodora-lazarevic/Poll-App/ent/user"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Custom errors so the handler knows what went wrong
var (
	ErrInvalidInput = errors.New("invalid username, email, or password")
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	DB *ent.Client
}

func NewUserService(db *ent.Client) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) userDataIsValid(username, email, password string) bool {
	return username != "" && email != "" && password != "" && emailRegex.MatchString(email)
}

func (s *UserService) Register(ctx context.Context, username, email, password string) error {
	if !s.userDataIsValid(username, email, password) {
		return ErrInvalidInput
	}

	exists, err := s.DB.User.Query().Where(user.Or(user.UsernameEQ(username), user.EmailEQ(email))).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.DB.User.Create().
		SetUsername(username).
		SetEmail(email).
		SetPasswordHash(string(hash)).
		Save(ctx)

	return err
}

func (s *UserService) Authenticate(ctx context.Context, identifier, password string) (*ent.User, error) {
	if identifier == "" || password == "" {
		return nil, ErrInvalidInput
	}

	u, err := s.DB.User.Query().
		Where(user.Or(user.UsernameEQ(identifier), user.EmailEQ(identifier))).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, ErrInvalidCreds
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCreds
	}

	return u, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*ent.User, error) {
	return s.DB.User.Query().All(ctx)
}

func (s *UserService) GetUserById(ctx context.Context, id int) (*ent.User, error) {
	u, err := s.DB.User.Query().Where(user.ID(id)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	return u, err
}
