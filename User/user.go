package Users

import (
	"errors"
	"fmt"
	"library/constants"
	"strconv"
	"strings"
	"unicode"
)

type User struct {
	Id                 int
	UserName           string
	UserLastName       string
	UserIdentityNumber string
	UserAdress         UserAdress
}

type UserAdress struct {
	UserCity     string
	UserDistrict string
}

//func NewUser(id int, username string, userLastName string, userIdentityNumber string, userCity string, userDistrict string) User {
//
//	user := User{
//		Id:                 id,
//		UserName:           username,
//		UserLastName:       userLastName,
//		UserIdentityNumber: userIdentityNumber,
//		UserAdress:         UserAdress{UserCity: userCity, UserDistrict: userDistrict},
//	}
//	return user
//}

func (user User) GetUserInformations() string {

	return fmt.Sprintf(
		"\n User ID: %d User Name: %s User Last Name: %s User Identity Number: %s User City: %s User District: %s",
		user.Id,
		user.UserName,
		user.UserLastName,
		user.UserIdentityNumber,
		user.UserAdress.UserCity,
		user.UserAdress.UserDistrict,
	)

}

func isEmpty(userName string, userLastName string, userCity string, userDistrict string) error {
	if len(userName) < 3 {
		return errors.New("name is is too short") // error mesajları düzeltielcek
	}
	if len(userLastName) < 3 {
		return errors.New("last name is too shorts")
	}
	if len(userCity) < 3 {
		return errors.New("user city is too shorts")
	}
	if len(userDistrict) < 3 {
		return errors.New("district is too shorts")
	}

	return nil
}

func (user *User) Valid() error {

	err := isEmpty(user.UserName, user.UserLastName, user.UserAdress.UserCity, user.UserAdress.UserDistrict)
	if err != nil {
		return err
	}

	cities := constants.GetCities()

	for index, _ := range user.UserName {
		if _, err := strconv.Atoi(string(user.UserName[index])); err == nil {
			return errors.New("name error,cannot contain a number")
		}
	}

	for _, digit := range user.UserLastName {
		if _, err := strconv.Atoi(string(digit)); err == nil {
			return errors.New("lastname error,cannot contain a number")
		}
	}
	if len(user.UserIdentityNumber) == 11 {

		for index, _ := range user.UserIdentityNumber {
			if _, err := strconv.Atoi(string(user.UserIdentityNumber[index])); err != nil {
				return errors.New("T.C error,cannot contain a number")
			}
		}
	} else {
		return errors.New("T.C must 11 characters")
	}

	for _, userCity := range user.UserAdress.UserCity {
		if unicode.IsDigit(userCity) {
			return errors.New("city error,cannot contain a number")
		}
	}

	cityFound := false
	var chooseCity = user.UserAdress.UserCity
	for _, city := range cities {
		if strings.ToLower(city) == strings.ToLower(user.UserAdress.UserCity) {
			cityFound = true
			chooseCity = city
			continue

		}
	}
	if !cityFound {
		return errors.New("we dont support this city")
	}

	var chooseDistrict = user.UserAdress.UserDistrict
	for _, userDistrict := range user.UserAdress.UserDistrict {
		if unicode.IsDigit(userDistrict) {
			return errors.New("district error,cannot contain a number")
		}
	}

	err = constants.DistrictEveryCity(chooseCity, chooseDistrict)

	if err != nil {
		return errors.New("we dont support this district")
	}

	return nil
}

type IdRequest struct {
	Id int `json:"id"`
}

func IdFilter(users []User, req IdRequest) (map[int][]User, error) {

	UserIdMap := make(map[int][]User)
	isFoundUserId := false

	for _, user := range users {
		if user.Id == req.Id {
			isFoundUserId = true
			UserIdMap[user.Id] = append(UserIdMap[user.Id], user)

		}
	}

	if !isFoundUserId {
		return UserIdMap, errors.New("user not found")
	}

	fmt.Println(UserIdMap)

	return UserIdMap, nil

}
