package userService

import (
	"context"
	"jwtToken/service/tokenService"
)

type Storage interface {
	// Create creates a new user.
	Create(context.Context, *User) error
	// Update updates a existing user.
	Update(context.Context, *User) error
	// Delete deletes a existing user.
	Delete(_ context.Context, id string) error
	// Get retrieves user by id.
	Get(_ context.Context, id string) (*User, error)
	// GetByName retrieves user by login name.
	GetByName(_ context.Context, name string) (*User, error)
	//Verify login by login name and pass
	Verify(_ context.Context, name string, pass string) error
	// List retrieves a list of users using the criterias on ListOptions.
	List(context.Context, *ListOptions) (*ListResponse, error)
}

type Manager struct {
	storage      Storage
	TokenService *tokenService.Manager
}

var Service *Manager

func NewService(storage Storage, manager *tokenService.Manager) {
	Service = &Manager{
		storage:      storage,
		TokenService: manager,
	}
}

// Create creates a user and publish a user.created event to our message broker.
func (s Manager) Create(ctx context.Context, u *User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Create(ctx, u)
	return err
}

// Update updates a user and publish a user.updated event to our message broker.
func (s Manager) Update(ctx context.Context, u *User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Update(ctx, u)
	return err
}

// Delete deletes a user and publish a user.deleted event to our message broker.
func (s Manager) Delete(ctx context.Context, id string) error {

	err := s.storage.Delete(ctx, id)
	return err
}

// Get retrieves a user by its id.
func (s Manager) Get(ctx context.Context, id string) (*User, error) {
	return s.storage.Get(ctx, id)
}

func (s Manager) GetByName(ctx context.Context, name string) (*User, error) {
	if name == "" {
		return nil, &FieldError{
			UserError: UserError{
				Type:    InvalidArgument,
				Code:    "missing_field",
				Message: "Field was not provided",
			},
			Field: "username",
		}
	}
	return s.storage.GetByName(ctx, name)
}
func (s Manager) Verify(ctx context.Context, name string, pass string) error {
	if name == "" || pass == "" {
		return &FieldError{
			UserError: UserError{
				Type:    InvalidArgument,
				Code:    "missing_field",
				Message: "Field was not provided",
			},
			Field: "username or password",
		}
	}
	return s.storage.Verify(ctx, name, pass)
}

// List retrieves a list of user using the criteria provided on opts.
func (s Manager) List(ctx context.Context, opts *ListOptions) (*ListResponse, error) {
	if opts == nil {
		opts = NewListOptions()
	}
	return s.storage.List(ctx, opts)
}
