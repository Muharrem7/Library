package Users_Personel_Informations

import (
	"errors"
	"fmt"
	"library/constants"
	"strconv"
	"strings"
	"unicode"
)

type UserInformations struct {
	Id                 string
	UserName           string
	UserLastName       string
	UserIdentityNumber string // TODO string yap, her bir karakterin numerik olduğunu kontrol et.
	UserAdress         UserAdressInformations
}

type UserAdressInformations struct {
	UserCity     string
	UserDistrict string
}

func NewUserInformations(id string, username string, userLastName string, userIdentityNumber string, userCity string, userDistrict string) UserInformations {

	user := UserInformations{
		Id:                 id,
		UserName:           username,
		UserLastName:       userLastName,
		UserIdentityNumber: userIdentityNumber,
		UserAdress:         UserAdressInformations{UserCity: userCity, UserDistrict: userDistrict},
	}
	return user
}

// TODO rename GetUserInformations --> GetUserInformation
func (user UserInformations) GetUserInformations() string {
	// TODO fmt.sprintF kullanılacak
	return fmt.Sprintf(
		"\n User ID: %s User Name: %s User Last Name: %s User Identity Number: %d User City: %s User District: %s",
		user.Id,
		user.UserName,
		user.UserLastName,
		user.UserIdentityNumber,
		user.UserAdress.UserCity,
		user.UserAdress.UserDistrict,
	)

}

func getCities() [10]string {
	return [10]string{"İstanbul", "Ankara", "İzmir", "Samsun", "Adana", "Edirne", "Elazığ", "Van", "Giresun", "Ordu"}
}

func (user *UserInformations) Valid() error {

	cities := getCities()

	for index, _ := range user.UserName {
		if _, err := strconv.Atoi(string(user.UserName[index])); err == nil {
			return errors.New("isim hatalı")
		}
	}

	for _, digit := range user.UserLastName {
		if _, err := strconv.Atoi(string(digit)); err == nil {
			return errors.New("soyisim hatalı")
		}
	}
	if len(user.UserIdentityNumber) == 11 {

		for index, _ := range user.UserIdentityNumber {
			if _, err := strconv.Atoi(string(user.UserIdentityNumber[index])); err != nil {
				return errors.New("T.C hatalı")
			}
		}
	} else {
		return errors.New("T.C 11 haneli olmalı")
	}

	// kodda kayıtlı şehir kontrolüde eklenecek
	// şehirler yeni constants olsun
	// cities array
	for _, userCity := range user.UserAdress.UserCity {
		if unicode.IsDigit(userCity) {
			return errors.New("il hatalı")
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
		return errors.New("bu ilde destek vermiyoruz")
	}

	var chooseDistrict = user.UserAdress.UserDistrict
	for _, userDistrict := range user.UserAdress.UserDistrict {
		if unicode.IsDigit(userDistrict) {
			return errors.New("ilce hatalı")
		}
	}

	err := constants.DistrictEveryCity(chooseCity, chooseDistrict)

	if err != nil {
		return errors.New("bu ilçede destek vermiyoruz")
	}

	return nil
}
