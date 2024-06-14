package main

import (
	"fmt"
	bookPackage "library/Book_Information"
	usersPackage "library/Users_Personel_Informations"
	"strconv"
	"strings"
	"unicode"
)

var userIdCounter = 1

func process() {
	// router oluşturulack her bir işlem directoryde yapılacak

	var users []usersPackage.User
	var books []bookPackage.BookInformation
	var mapsOfUsersBook = make(map[string][]bookPackage.BookInformation)
	var mapsOfBooksUser = make(map[string][]usersPackage.User)

	//var mapsOfBooksUser = make(map[int][]usersPackage.UserInformations)
	for {
		var userProcess string
		var chooseUserProcess int
		fmt.Println("\n 1.Yeni kullanici ekleme \n 2.Tüm kullanicilarin listelenmesi \n 3.Yeni kitap ekleme  \n 4.Kitaplarin listelenmesi\n" +
			" 5.Kullanicilara kitap atamasi  \n 6.Kullanicilarin hangi kitaplara sahip olduğunun listelenmesi  \n 7.Kitaplarin hangi kullanicilarda olduğu  \n" +
			" 8.Kitaplarin kategoriye göre filtrelenmesi  \n 9.Kitaplarin girilen texte göre filtrelenmesi \n 0.Programın sonlandırılması")

		fmt.Println("\n --- Lütfen yapmak istediğiniz işlemi sayı ile belirtiniz ---")
		fmt.Scan(&userProcess)

		chooseUserProcess, err := strconv.Atoi(userProcess)
		if err != nil {
			fmt.Println("int değil")
			continue
		}

		// http metodlarına bakılacak tiplerine göçre işlem yapılacak bunu araştır
		switch chooseUserProcess {
		// /user, POST
		case 1:
			// DATA request body den gelecek, marshall ve unmarshall json
			// ilk validatsyon            // 400
			// başarılı ise ekleyeveksin // 200
			user := createNewUser()
			if err := user.Valid(); err != nil {
				fmt.Println("Lütfen tekrar deneyiniz.", err.Error())
			} else {
				users = append(users, user)
				userIdCounter++

				fmt.Println(users)
			}

		case 2:
			if len(users) == 0 {
				fmt.Println("Herhangi bir kullanici bulunmamaktadir önce kullanici oluşturunuz...")
				continue
			}

			PrintUsers(users)
		case 3:
			book := createNewBook()

			if err := book.Valid(); err != nil {
				fmt.Println("Lütfen tekrar deneyiniz.", err.Error())
			} else {
				books = append(books, book)
				fmt.Println(books)
			}

		case 4:
			if len(books) == 0 {
				fmt.Println("Herhangi bir kitap bulunmamaktadir önce kitap oluşturunuz...")
				continue
			}
			PrintBook(books)
		case 5:
			mapsOfUsersBook, mapsOfBooksUser = AssigneBookForUsers(users, books, mapsOfUsersBook, mapsOfBooksUser)
		case 6:
			// TODO kitabı kullanıcıya eklemeke için yeni bir numara
			if len(mapsOfUsersBook) == 0 {
				fmt.Println("Atama Bulunamadı")
				continue
			} else {
				for userId, assignedBook := range mapsOfUsersBook {
					for _, book := range assignedBook {
						fmt.Printf("\n %s id'li kullaniciya %s isbn'li %s isimli kitap atanmıştır.\n", userId, book.ISBN, book.BookName)
					}
				}
			}
		case 7:

			for bookIsbn, assignedUser := range mapsOfBooksUser {
				for _, user := range assignedUser {
					id, _ := strconv.Atoi(user.Id)
					fmt.Printf("\n %s isbn'li kitabın atandığı kullanici id'si: %d", bookIsbn, id)
				}
			}

		case 8:
			mapsOfCategory := make(map[string][]bookPackage.BookInformation)
			categoryFilter(books, mapsOfCategory)

		case 9:
			mapsofName := make(map[string][]bookPackage.BookInformation)
			nameFilter(books, mapsofName)

		case 0:
			return

		}

	}

}

//	func validtyUserInputIsDigit(userInfo string) int {
//		value, err := strconv.Atoi(userInfo)
//		if err != nil {
//			fmt.Println("int değil")
//		}
//
//		return value
//	}
func validtyUserInputIsString(userInput string) {
	for _, input := range userInput {
		if unicode.IsDigit(input) {
			fmt.Println("string değil")
		}
	}

}

func createNewUser() usersPackage.User {
	var userName string
	var userLastName string
	var userIdentityNumber string
	var userCity string
	var userDistrict string
	var user usersPackage.User

	fmt.Println("\nLütfen eklemek istediğiniz kullancinin adini rakam olmadan giriniz: ")
	fmt.Scan(&userName)

	fmt.Println("Lütfen eklemek istediğiniz kullanicinin soyadini rakam olmadan giriniz: ")
	fmt.Scan(&userLastName)

	fmt.Println("Lütfen eklemek istediğiniz kullanicinin T.C no'sunu harf kullanmadan giriniz: ")
	fmt.Scan(&userIdentityNumber)

	fmt.Println("Lütfen eklemek istediğiniz kullanicinin şehrini rakam kullanmadan giriniz: ")
	fmt.Scan(&userCity)

	fmt.Println("Lütfen eklemek istediğiniz kullanicinin bulunduğu ilçeyi rakam kullanmadan giriniz: ")
	fmt.Scan(&userDistrict)

	idCounter := strconv.Itoa(userIdCounter)

	user = usersPackage.User{Id: idCounter, UserName: userName, UserLastName: userLastName, UserIdentityNumber: userIdentityNumber, UserAdress: usersPackage.UserAdress{UserCity: userCity, UserDistrict: userDistrict}}

	return user

}

func PrintUsers(users []usersPackage.User) {
	for _, userInfo := range users {
		fmt.Println(userInfo.GetUserInformations())
	}
}

func createNewBook() bookPackage.BookInformation {
	var bookIsbn string
	var bookName string
	var bookAuthor string
	var bookCategory string
	var bookPages string
	var book bookPackage.BookInformation

	fmt.Println("Lütfen eklemek istediğiniz kitabin 13 haneli ISBN numarasini harf kullanmadan giriniz: ")
	fmt.Scan(&bookIsbn)

	fmt.Println("Lütfen eklemek istediğiniz kitabin adini giriniz: ")
	fmt.Scan(&bookName)

	fmt.Println("Lütfen eklemek istediğiniz kitabin yazarini rakam kullanmadan giriniz: ")
	fmt.Scan(&bookAuthor)

	fmt.Println("Lütfen eklemek istediğiniz kitabin hangi kategoride olduğunu rakam kullanmadan giriniz: ")
	fmt.Scan(&bookCategory)

	fmt.Println("Lütfen eklemek istediğiniz kitabin kaç sayfa olduğunu harf kullanmadan giriniz: ")
	fmt.Scan(&bookPages)

	book = bookPackage.BookInformation{ISBN: bookIsbn, BookName: bookName, Author: bookAuthor, Category: bookCategory, Pages: bookPages}

	return book
}

func PrintBook(book []bookPackage.BookInformation) {
	for _, bookInfo := range book {
		fmt.Println(bookInfo.GetBookInformations())
	}
}
func AssigneBookForUsers(Users []usersPackage.User, Books []bookPackage.BookInformation, mapsOfUsersBook map[string][]bookPackage.BookInformation, mapsOfBooksUser map[string][]usersPackage.User) (map[string][]bookPackage.BookInformation, map[string][]usersPackage.User) {

	var userId string
	var bookIsbn string
	var isFoundBook = false
	var isFoundUser = false

	fmt.Println("Kitap atamak istediğiniz kullanicinin ID'sini harf kullanmadan giriniz: ")
	fmt.Scan(&userId)

	// TODO User --> users yerleri değişecek

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

// TODO returnlere bak
func categoryFilter(Books []bookPackage.BookInformation, mapsOfCategory map[string][]bookPackage.BookInformation) {
	var userInput string
	var isFoundBook = false

	category := bookPackage.GetCategories()

	fmt.Println(category)

	fmt.Println("Lütfen aradığınız kategoriyi rakam kullanmadan yazınız:")
	fmt.Scan(&userInput)
	validtyUserInputIsString(userInput)

	for _, book := range Books {

		if strings.ToLower(book.Category) == strings.ToLower(userInput) {
			isFoundBook = true
			mapsOfCategory[book.Category] = append(mapsOfCategory[book.Category], book)
			fmt.Println(mapsOfCategory)
		}
	}

	if !isFoundBook {
		fmt.Printf("%s Kategorisinde kitap bulunamadı lütfen tekrar deneyiniz.", userInput)

	}

}

// TODO returnlere bak
func nameFilter(Book []bookPackage.BookInformation, mapsOfName map[string][]bookPackage.BookInformation) {
	var userInput string
	var isFoundBook = false

	fmt.Println("Lütfen aradığınız kitabin adini rakam kullanmadan yazınız:")
	fmt.Scan(&userInput)
	validtyUserInputIsString(userInput)

	for _, book := range Book {
		if strings.ToLower(book.BookName) == strings.ToLower(userInput) {
			isFoundBook = true
			mapsOfName[book.BookName] = append(mapsOfName[book.BookName], book)
			fmt.Println(mapsOfName)

		}
	}
	if !isFoundBook {
		fmt.Println("Eşleşme olmadığı için tekrardan deneyiniz.")
	}

}
