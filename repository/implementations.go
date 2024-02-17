package repository

import (
	"context"
	"log"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (r *Repository) SaveRegister(ctx context.Context, input SaveRegisterInput) (output SaveRegisterOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"INSERT INTO users (fullname, password, phone_number) VALUES ($1, $2, $3) RETURNING id",
		input.Fullname, input.HashPassword, input.PhoneNumber,
	).Scan(&output.Id)

	if err != nil {
		log.Printf("Error inserting data into the database: %v", err)
		return output, err
	}

	return output, nil
}

func (r *Repository) Login(ctx context.Context, input GetUsersByPhoneInput) (output GetUsersByPhoneOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"SELECT id, password FROM users WHERE phone_number = $1",
		input.PhoneNumber,
	).Scan(&output.Id, &output.HashPassword)

	if err != nil {
		log.Printf("Error querying user data from the database: %v", err)
		return output, err
	}

	return output, nil
}

func (r *Repository) GetProfileById(ctx context.Context, input GetProfiletByIdInput) (output GetProfileByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT fullname, phone_number FROM users WHERE id = $1", input.Id).Scan(&output.Fullname, &output.PhoneNumber)
	if err != nil {
		log.Printf("Error querying user data from the database: %v", err)
		return output, err
	}
	return output, nil
}
