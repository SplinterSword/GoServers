package database

import "errors"

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	Is_Chirpy_Red  bool   `json:"is_chirpy_red"`
}

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:             id,
		Email:          email,
		HashedPassword: hashedPassword,
		Is_Chirpy_Red:  false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUserStatus(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("couldn't Retrive the Data")
	}

	user := dbStructure.Users[id]

	if user.Is_Chirpy_Red {
		return user, errors.New("user already using chirpy red")
	}

	user.Is_Chirpy_Red = true
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, errors.New("something went wrong while updating the data")
	}
	return user, nil
}
