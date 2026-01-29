package usecase

import (
	"context"
	appError "github.com/harryosmar/protobuf-go/error"
	userpb "github.com/harryosmar/protobuf-go/gen/user"
	"github.com/harryosmar/protobuf-go/repository"
)

// UserServiceUsecase defines the interface for user business logic
type UserServiceUsecase interface {
	// Methods matching the proto service definition
	CreateUser(ctx context.Context, req *userpb.CreateUserRequestDTO) (*userpb.CreateUserResponseDTO, error)
	GetUser(ctx context.Context, req *userpb.GetUserRequestDTO) (*userpb.GetUserResponse, error)
	UpdateUser(ctx context.Context, req *userpb.UpdateUserRequestDTO) (*userpb.UpdateUserResponseDTO, error)
	DeleteUser(ctx context.Context, req *userpb.DeleteUserRequestDTO) (*userpb.DeleteUserResponseDTO, error)
	ListUsers(ctx context.Context, req *userpb.ListUsersRequestDTO) (*userpb.ListUsersResponseDTO, error)
}

// userServiceUsecase implements UserServiceUsecase interface
type userServiceUsecase struct {
	userServiceRepo repository.ServiceRepository[userpb.UserEntityORM, uint32]
}

// NewUserServiceUsecase creates a new user usecase instance
func NewUserServiceUsecase(repo repository.ServiceRepository[userpb.UserEntityORM, uint32]) UserServiceUsecase {
	return &userServiceUsecase{
		userServiceRepo: repo,
	}
}

// CreateUser implements the CreateUser RPC method from the proto service
func (u *userServiceUsecase) CreateUser(ctx context.Context, req *userpb.CreateUserRequestDTO) (*userpb.CreateUserResponseDTO, error) {
	// Create user entity from DTO
	dto := req.User
	entity := &userpb.UserEntity{
		Id:        dto.Id,
		Name:      dto.Name,
		Email:     dto.Email,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}

	// Convert to ORM model for database operations
	userORM, err := entity.ToORM(ctx)
	if err != nil {
		return nil, err
	}

	// Save to database using repository
	newUserORM, err := u.userServiceRepo.Create(ctx, &userORM)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponseDTO{
		User: &userpb.UserDTO{
			Name:  newUserORM.Name,
			Email: newUserORM.Email,
			Id:    newUserORM.Id,
		},
	}, nil
}

// GetUser implements the GetUser RPC method from the proto service
func (u *userServiceUsecase) GetUser(ctx context.Context, req *userpb.GetUserRequestDTO) (*userpb.GetUserResponse, error) {
	// Query database for user using repository
	orm, err := u.userServiceRepo.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if orm == nil {
		return nil, appError.ErrUserNotFound
	}

	return &userpb.GetUserResponse{
		User: u.ormToDTO(orm),
	}, nil
}

func (u *userServiceUsecase) ormToDTO(orm *userpb.UserEntityORM) *userpb.UserDTO {
	return &userpb.UserDTO{
		Name:  orm.Name,
		Email: orm.Email,
		Id:    orm.Id,
	}
}

// DeleteUser implements the DeleteUser RPC method from the proto service
func (u *userServiceUsecase) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequestDTO) (*userpb.DeleteUserResponseDTO, error) {
	err := u.userServiceRepo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.DeleteUserResponseDTO{}, nil
}

// UpdateUser implements the UpdateUser RPC method from the proto service
func (u *userServiceUsecase) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequestDTO) (*userpb.UpdateUserResponseDTO, error) {
	// Get existing user first
	orm, err := u.userServiceRepo.GetById(ctx, req.User.Id)
	if err != nil {
		return nil, err
	}
	if orm == nil {
		return nil, appError.ErrUserNotFound
	}

	// Update fields
	orm.Name = req.User.Name
	orm.Email = req.User.Email

	// Update in database using repository
	if _, err := u.userServiceRepo.Update(ctx, orm); err != nil {
		return nil, err
	}

	return &userpb.UpdateUserResponseDTO{
		User: req.User,
	}, nil
}

func (u *userServiceUsecase) ormToDTOList(ormRecords []userpb.UserEntityORM) []*userpb.UserDTO {
	var dtoRecords []*userpb.UserDTO
	for _, record := range ormRecords {
		dtoRecords = append(dtoRecords, u.ormToDTO(&record))
	}
	return dtoRecords
}

// ListUsers implements the ListUsers RPC method from the proto service
func (u *userServiceUsecase) ListUsers(ctx context.Context, req *userpb.ListUsersRequestDTO) (*userpb.ListUsersResponseDTO, error) {
	ormRecords, paginator, err := u.userServiceRepo.GetPerPage(
		ctx,
		req.Pagination.Page,
		req.Pagination.Limit,
		[]repository.OrderBy{},
		[]repository.Where{},
	)
	if err != nil {
		return nil, err
	}

	return &userpb.ListUsersResponseDTO{
		Users: u.ormToDTOList(ormRecords),
		Pagination: &userpb.PaginationResponse{
			Total: paginator.Total,
			Page:  paginator.Page,
			Limit: paginator.PerPage,
		},
	}, nil
}
