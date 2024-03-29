// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	SaveRegister(ctx context.Context, input SaveRegisterInput) (output SaveRegisterOutput, err error)
	Login(ctx context.Context, input GetUsersByPhoneInput) (output GetUsersByPhoneOutput, err error)
	GetProfileById(ctx context.Context, input GetProfiletByIdInput) (output GetProfileByIdOutput, err error)
	UpdateNumberLogin(ctx context.Context, input UpdateNumberLoginInput) (err error)
	GetProfileByPhone(ctx context.Context, input GetProfileByPhoneInput) (output GetProfileByPhoneOutput, err error)
	UpdateProfile(ctx context.Context, input UpdateProfileInput) error
}
