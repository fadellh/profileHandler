// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type SaveRegisterInput struct {
	Fullname     string
	HashPassword string
	PhoneNumber  string
}

type SaveRegisterOutput struct {
	Id int
}

type GetUsersByPhoneInput struct {
	PhoneNumber string
}

type GetUsersByPhoneOutput struct {
	Id           int
	HashPassword string
}
