package users

import (
	"karmaclips/database"
)

func CreateUser(user *database.Users) (*database.Users, error) {
	db, err := database.DBConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = database.InsertStruct(db, "user", user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserById(uid string) (*database.Users, error) {
	db, err := database.DBConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM users WHERE id = $1"
	rows, err := db.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user []*database.Users

	if err := database.ParseRows(rows, &user); err != nil {
		return nil, err
	}

	return user[0], nil
}

func GetUserByEmail(email string) (*database.Users, error) {
	db, err := database.DBConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM users WHERE email = $1"
	rows, err := db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user []*database.Users

	if err := database.ParseRows(rows, &user); err != nil {
		return nil, err
	}

	return user[0], nil
}
