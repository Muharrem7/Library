package constants

import (
	"errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func DistrictEveryCity(userCity string, userDistrict string) error {

	//var returnDistrict string

	mapsOfCityDistricts := map[string][]string{

		"İstanbul": {"Esenler", "Bağcılar", "Kağıthane", "Şişli", "Levent", "Kartal", "Üsküdar", "Beşiktaş", "Kadıköy", "Pendik"},
		"Ankara":   {"Etimesgut", "Sincan", "Akyurt", "Beypazari", "Camlidere", "Evren", "Güdül", "Kalecik", "Kazan", "Nallihan"},
		"İzmir":    {"Bayraklı", "Bornova", "Buca", "Gaziemir", "Karabağlar", "Konak", "Narlıdere", "Foca", "Menderes"},
		"Samsun":   {"Çarşamba", "Bafra", "Atakum", "Canik", "Tekkeköy", "Alaçam", "Kavak", "Havza", "Ayvacık", "Asarcık"},
		"Adana":    {"Kozan", "İmamoğlu", "Karataş", "Pozantı", "Ceyhan", "Çukurova", "Feke", "Yumurtalık", "Seyhan", "Sarıçam"},
		"Edirne":   {"Enez", "Edirne", "Havsa", "İpsala", "Keşan", "Lalapaşa", "Meriç", "Süloğlu", "Uzunköprü"},
		"Elazığ":   {"Elazığ", "Ağın", "Alacakaya", "Keban", "Arıcak", "Baskil", "Kovancılar", "Maden", "Palu", "Sivrice"},
		"Van":      {"Edremit", "İpekyolu", "Tuşba", "Bahçesaray", "Başkale", "Çaldıran", "Erçiş", "Gevaş", "Saray", "Muradiye"},
		"Giresun":  {"Alucra", "Bulancak", "Çamoluk", "Çanakçı", "Dereli", "Doğankent", "Espiye", "Güce", "Görele", "Tirebolu"},
		"Ordu":     {"Altınordu", "Akkuş", "Çamaş", "Çaybaşı", "Gölköy", "Perşembe", "Kabataş", "Kabadüz", "Kumru", "Ünye"},
	}

	districtFound := false
	cityText := cases.Title(language.Turkish).String(userCity)
	districttext := cases.Title(language.Turkish).String(userDistrict)

	districts, exist := mapsOfCityDistricts[cityText]

	if exist {
		for _, district := range districts {
			if strings.ToLower(districttext) == strings.ToLower(district) {
				districtFound = true
				break
			}
		}

	}

	if !districtFound {
		return errors.New("hiçbir şey bulamadık")
	}

	// TODO: switch case kalkacak

	//switch cityText {
	//case "İstanbul":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//
	//	}
	//
	//case "Ankara":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//case "İzmir":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//case "Samsun":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//case "Adana":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//case "Edirne":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//case "Elazığ":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//case "Van":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//case "Giresun":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//case "Ordu":
	//	districts, exist := mapsOfCityDistricts[cityText]
	//
	//	if exist {
	//		for _, district := range districts {
	//			if strings.ToLower(districttext) == strings.ToLower(district) {
	//				districtFound = true
	//				break
	//			}
	//		}
	//	}
	//default:
	//	return errors.New("hiçbir şey bulamadık")
	//
	//}

	return nil
}
