package repository

import (
	"context"
	"fmt"
	"log"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	log.Printf("Input ID: %s", input.Id)
	if err != nil {
		fmt.Println("ERRRR", err)
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
