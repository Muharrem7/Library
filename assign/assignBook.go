package assign

import (
	"fmt"
	bookPackage "library/Book_Information"
	usersPackage "library/Users_Personel_Informations"
)

func AssigneBookForUsers(Users []usersPackage.User, Books []bookPackage.BookInformation, mapsOfUsersBook map[string][]bookPackage.BookInformation, mapsOfBooksUser map[string][]usersPackage.User) (map[string][]bookPackage.BookInformation, map[string][]usersPackage.User) {

	var userId string
	var bookIsbn string
	var isFoundBook = false
	var isFoundUser = false

	fmt.Println("Kitap atamak istediğiniz kullanicinin ID'sini harf kullanmadan giriniz: ")
	fmt.Scan(&userId)

	var foundUser usersPackage.User
	for _, user := range Users {
		if user.Id == userId {
			isFoundUser = true
			foundUser = user

		}

	}
	if isFoundUser == false {
		fmt.Println("bulunamadı")
	}

	fmt.Println("Lütfen atamak istediğiniz kitabin ISBN numarasini harf kullanmadan giriniz: ")
	fmt.Scan(&bookIsbn)

	var foundBook bookPackage.BookInformation
	for _, book := range Books {
		if book.ISBN == bookIsbn {
			isFoundBook = true
			foundBook = book
		}

	}
	if isFoundBook == false {
		fmt.Println("bulunamadı")
	}

	if isFoundUser == true && isFoundBook == true {
		mapsOfUsersBook[userId] = append(mapsOfUsersBook[userId], foundBook)
		mapsOfBooksUser[bookIsbn] = append(mapsOfBooksUser[bookIsbn], foundUser)
	}

	return mapsOfUsersBook, mapsOfBooksUser
}
