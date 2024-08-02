package User

import (
	"database/sql"
	"fmt"
	bookPackage "library/Book"
)

type UserRepository struct {
	db *sql.DB
}
type CityRequest struct {
	City string `json:"userCity"`
}
type DistrictRequest struct {
	District string `json:"userDistrict"`
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

// TODO diğer kolonlar eklenmeli
func (us *UserRepository) InsertUser(user User) error {

	stmt, err := us.db.Prepare("INSERT INTO  users ( name, last_name, identity_number, city_id, district_id, adress_description) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("error while insert user %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UserName, user.UserLastName, user.UserIdentityNumber, CityId, DistrictId, user.UserAdress.UserAdressDescription)
	if err != nil {

		return fmt.Errorf("Error while execute query: %v", err)

	}

	return nil
}

// GET users
func (us *UserRepository) GetUsers() ([]User, error) {
	rows, err := us.db.Query("SELECT * FROM users ORDER BY created_at DESC")
	if err != nil {
		fmt.Errorf("Bir hata oluştu.", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserId, &user.UserName, &user.UserLastName, &user.UserIdentityNumber, &user.UserAdress.UserCity, &user.UserAdress.UserDistrict, &user.UserAdress.UserAdressDescription, &user.UserCreatedAt)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}

	return users, nil
}

func (us *UserRepository) GetUserIdById(req bookPackage.AssignRequest) (int, error) {
	userId := req.UserId
	rows, err := us.db.Query("SELECT id From users where id = ?", userId)
	if err != nil {
		fmt.Errorf("Error while querying user: %v", err)
	}

	var id int
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("Error while querying user: %v", err)
		}

	}

	return id, nil

}
func (us *UserRepository) GetAssignedUser() ([]Assignments, error) {
	rows, err := us.db.Query("SELECT user_id,book_id FROM assignments ")
	if err != nil {
		return nil, fmt.Errorf("Error while querying book: %v", err)
	}
	defer rows.Close()
	var users []Assignments

	for rows.Next() {
		var id Assignments
		err := rows.Scan(&id.UserId, &id.BookId)
		if err != nil {
			return nil, fmt.Errorf("Error while querying book: %v", err)
		}
		users = append(users, id)
	}
	return users, nil
}
func (us *UserRepository) GetUserById(req IdRequest) ([]User, error) {
	id := req.UserId
	rows, err := us.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("Error while querying user: %v", err)
	}
	defer rows.Close()

	var userInfo []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserId, &user.UserName, &user.UserLastName, &user.UserIdentityNumber, &user.UserAdress.UserCity, &user.UserAdress.UserDistrict, &user.UserAdress.UserAdressDescription, &user.UserCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Error while scaning user: %v", err)
		}
		userInfo = append(userInfo, user)
	}
	return userInfo, nil

}
