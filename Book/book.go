package Book

// book
import (
	"errors"
	"fmt"
	Users "library/Users_Personel_Informations"
	"strconv"
	"strings"
)

type BookInformation struct {
	ISBN      string // TODO string
	BookName  string // JSON ANNOTATION KULLANILACAK
	Author    string
	Category  string
	PageCount int // page count inter
}

func NewBookInformation(isbn string, bookName string, author string, category string, pageCount int) BookInformation {
	book := BookInformation{
		ISBN:      isbn,
		BookName:  bookName,
		Author:    author,
		Category:  category,
		PageCount: pageCount,
	}

	return book
}

func (book BookInformation) GetBookInformations() string {

	return fmt.Sprintf("\n ISBN: %s  Book Name: %s  Author: %s  Category: %s  Pages: %d",
		book.ISBN,
		book.BookName,
		book.Author,
		book.Category,
		book.PageCount,
	)
}

func GetCategories() [10]string {
	return [10]string{"Roman", "Hikaye", "Sosyoloji", " Gerilim", "Tarih", "Psikoloji", "Aşk", "Çocuk", "Fantastik", "Edebiyat"}
}

func isEmpty(isbn string, bookName string, author string, category string, pageCount int) error {

	if len(isbn) < 3 {
		return errors.New("isbn is too short")
	}
	if len(bookName) < 3 {
		return errors.New("book name is too short")
	}
	if len(author) < 3 {
		return errors.New("author name is too short")
	}
	if len(category) < 3 {
		return errors.New("category name is too short")
	}
	if len(strconv.Itoa(pageCount)) < 3 {
		return errors.New("page count is too short")
	}

	return nil
}

func (book *BookInformation) Valid() error {

	err := isEmpty(book.ISBN, book.BookName, book.Author, book.Category, book.PageCount)
	if err != nil {
		return err
	}

	categories := GetCategories() // error dönme olasılığın var, gerek yok

	if len(book.ISBN) == 13 {
		for index, _ := range book.ISBN {
			if _, err := strconv.Atoi(string(book.ISBN[index])); err != nil {
				return errors.New("ISBN error, cannot contain char")
			}
		}
	} else {
		return errors.New("ISBN must 13 characters")
	}

	for index, _ := range book.BookName {
		if _, err := strconv.Atoi(string(book.BookName[index])); err == nil {
			return errors.New("book name error, cannot contain number")
		}
	}

	for index, _ := range book.Author {
		if _, err := strconv.Atoi(string(book.Author[index])); err == nil {
			return errors.New("author error, cannot contain number")
		}
	}

	for index, _ := range book.Category {
		if _, err := strconv.Atoi(string(book.Category[index])); err == nil {
			return errors.New("category error, cannot contain number")
		}
	}

	categoryFound := false
	for _, category := range categories {
		if strings.ToLower(category) == strings.ToLower(book.Category) {
			categoryFound = true
			continue

		}

	}
	if !categoryFound {
		return errors.New("cannot add this category")
	}

	for index, _ := range strconv.Itoa(book.PageCount) {
		if _, err := strconv.Atoi(strconv.Itoa(index)); err != nil {
			return errors.New("page error, cannot contain char")
		}
	}

	return nil
}

type AssignRequest struct {
	UserId   int    `json:"userId"`
	BookIsbn string `json:"bookIsbn"`
}

func AssignBook(users []Users.User, books []BookInformation, req AssignRequest) (Users.User, BookInformation, error) {
	var foundUser Users.User
	var foundBook BookInformation
	isUserFound, isBookFound := false, false

	for _, user := range users {
		if user.Id == req.UserId {
			foundUser = user
			isUserFound = true
			break
		}
	}

	for _, book := range books {
		if book.ISBN == req.BookIsbn {
			foundBook = book
			isBookFound = true
			break
		}
	}

	if !isUserFound || !isBookFound {
		return Users.User{}, BookInformation{}, errors.New("user or book not found")
	}

	return foundUser, foundBook, nil
}

type CategoryRequest struct {
	Category string `json:"category"`
}

func CategoryFilter(books []BookInformation, req CategoryRequest) (map[string][]BookInformation, error) {
	category := strings.ToLower(req.Category)
	isFoundCategory := false
	mapsOfCategory := make(map[string][]BookInformation)

	for _, book := range books {
		if strings.ToLower(book.Category) == strings.ToLower(category) {
			isFoundCategory = true
			mapsOfCategory[book.Category] = append(mapsOfCategory[book.Category], book)
			break
		}
	}

	if !isFoundCategory {
		return mapsOfCategory, errors.New("category not found")
	}

	return mapsOfCategory, nil
}

type NameRequest struct {
	Name string `json:"name"`
}

func NameFilter(books []BookInformation, req NameRequest) (map[string][]BookInformation, error) {

	name := strings.ToLower(req.Name)
	isFoundName := false
	mapsOfName := make(map[string][]BookInformation)

	for _, book := range books {
		if strings.ToLower(book.BookName) == strings.ToLower(name) {
			isFoundName = true
			mapsOfName[book.BookName] = append(mapsOfName[book.BookName], book)
			break
		}
	}

	if !isFoundName {
		return mapsOfName, errors.New("name not found")
	}

	return mapsOfName, nil

}

type IdRequest struct {
	Id int `json:"id"`
}

func IdFilter(users []Users.User, req IdRequest) (map[int][]Users.User, error) {

	mapsOfUserId := make(map[int][]Users.User)
	isFoundUserId := false

	for _, user := range users {
		if user.Id == req.Id {
			isFoundUserId = true
			mapsOfUserId[user.Id] = append(mapsOfUserId[user.Id], user)
			break
		}
	}

	if !isFoundUserId {
		return mapsOfUserId, errors.New("user not found")
	}

	return mapsOfUserId, nil

}

type IsbnRequest struct {
	Isbn string `json:"isbn"`
}

func IsbnFilter(books []BookInformation, req IsbnRequest) (map[string][]BookInformation, error) {
	mapsOfBookIsbn := make(map[string][]BookInformation)
	isFoundBookIsbn := false
	for _, book := range books {
		if book.ISBN == req.Isbn {
			isFoundBookIsbn = true
			mapsOfBookIsbn[book.ISBN] = append(mapsOfBookIsbn[book.ISBN], book)
			break
		}
	}

	if !isFoundBookIsbn {
		return mapsOfBookIsbn, errors.New("isbn not found")
	}

	return mapsOfBookIsbn, nil

}
