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
		"SELECT id, password, number_login FROM users WHERE phone_number = $1",
		input.PhoneNumber,
	).Scan(&output.Id, &output.HashPassword, &output.NumberLogin)

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

func (r *Repository) UpdateNumberLogin(ctx context.Context, input UpdateNumberLoginInput) (err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"UPDATE users SET number_login = number_login + 1 WHERE id = $1",
		input.Id,
	).Err()

	if err != nil {
		log.Printf("Error updating number_login in the database: %v", err)
		return err
	}

	return nil
}
