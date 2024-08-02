package User

import (
	"errors"
	"strconv"
)

// TODO: buraya structlar gelecek

type User struct {
	UserId             int
	UserName           string
	UserLastName       string
	UserIdentityNumber string
	UserAdress         UserAdress
	UserCreatedAt      string
}

type UserAdress struct {
	UserCity              string
	UserDistrict          string
	UserAdressDescription string
}
type CityInfo struct {
	Id        int
	CityPlate int
	Districts []string
}
type IdRequest struct {
	UserId int `json:"userId"`
}
type Assignments struct {
	UserId int `json:"userId"`
	BookId int `json:"bookId"`
}

var citiesAndDistricts = map[string]CityInfo{
	"istanbul": {
		Id:        1,
		CityPlate: 34,
		Districts: []string{"Esenler", "Bağcılar", "Kağıthane", "Şişli", "Levent", "Kartal", "Üsküdar", "Beşiktaş", "Kadıköy", "Pendik"},
	},
	"ankara": {
		Id:        2,
		CityPlate: 06,
		Districts: []string{"Etimesgut", "Sincan", "Akyurt", "Beypazari", "Camlidere", "Evren", "Güdül", "Kalecik", "Kazan", "Nallihan"},
	},
	"izmir": {
		Id:        3,
		CityPlate: 35,
		Districts: []string{"Bayraklı", "Bornova", "Buca", "Gaziemir", "Karabağlar", "Konak", "Narlıdere", "Foca", "Menderes"},
	},
	"samsun": {
		Id:        4,
		CityPlate: 55,
		Districts: []string{"Çarşamba", "Bafra", "Atakum", "Canik", "Tekkeköy", "Alaçam", "Kavak", "Havza", "Ayvacık", "Asarcık"},
	},
	"adana": {
		Id:        5,
		CityPlate: 01,
		Districts: []string{"Kozan", "İmamoğlu", "Karataş", "Pozantı", "Ceyhan", "Çukurova", "Feke", "Yumurtalık", "Seyhan", "Sarıçam"},
	},
	"edirne": {
		Id:        6,
		CityPlate: 22,
		Districts: []string{"Enez", "Edirne", "Havsa", "İpsala", "Keşan", "Lalapaşa", "Meriç", "Süloğlu", "Uzunköprü"},
	},
	"elazığ": {
		Id:        7,
		CityPlate: 23,
		Districts: []string{"Elazığ", "Ağın", "Alacakaya", "Keban", "Arıcak", "Baskil", "Kovancılar", "Maden", "Palu", "Sivrice"},
	},
	"van": {
		Id:        8,
		CityPlate: 65,
		Districts: []string{"Edremit", "İpekyolu", "Tuşba", "Bahçesaray", "Başkale", "Çaldıran", "Erçiş", "Gevaş", "Saray", "Muradiye"},
	},
	"giresun": {
		Id:        9,
		CityPlate: 28,
		Districts: []string{"Alucra", "Bulancak", "Çamoluk", "Çanakçı", "Dereli", "Doğankent", "Espiye", "Güce", "Görele", "Tirebolu"},
	},
	"ordu": {
		Id:        10,
		CityPlate: 52,
		Districts: []string{"Altınordu", "Akkuş", "Çamaş", "Çaybaşı", "Gölköy", "Perşembe", "Kabataş", "Kabadüz", "Kumru", "Ünye"},
	},
}
var CityId int
var DistrictId int

func (newUser *User) Valid() error {

	for index, _ := range newUser.UserName {
		if _, err := strconv.Atoi(string(newUser.UserName[index])); err == nil {
			return errors.New("name error,cannot contain a number")
		}
	}

	for _, digit := range newUser.UserLastName {
		if _, err := strconv.Atoi(string(digit)); err == nil {
			return errors.New("lastname error,cannot contain a number")
		}
	}
	if len(newUser.UserIdentityNumber) == 11 {

		for index, _ := range newUser.UserIdentityNumber {
			if _, err := strconv.Atoi(string(newUser.UserIdentityNumber[index])); err != nil {
				return errors.New("T.C error,cannot contain a char")
			}
		}
	} else {
		return errors.New("T.C must 11 characters")
	}

	for _, digit := range newUser.UserAdress.UserCity {
		if _, err := strconv.Atoi(string(digit)); err == nil {
			return errors.New("city error,cannot contain a number")
		}
	}
	for _, digit := range newUser.UserAdress.UserDistrict {
		if _, err := strconv.Atoi(string(digit)); err == nil {
			return errors.New("district error,cannot contain a number")
		}
	}
	var districtFound = false

	cityInfo, cityExists := citiesAndDistricts[newUser.UserAdress.UserCity]

	if !cityExists {
		return errors.New("city error,city not found")
	}

	for index, districts := range cityInfo.Districts {
		if districts == newUser.UserAdress.UserDistrict {
			CityId = cityInfo.Id
			index += 1
			DistrictId = index
			districtFound = true
			break
		}
	}

	if !districtFound {
		return errors.New("city error,district not found")
	}

	return nil
}
