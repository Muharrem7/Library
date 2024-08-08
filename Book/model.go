package Book

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type Book struct {
	Id            int
	ISBN          string
	BookName      string
	Author        string
	Category      string
	PageCount     int
	BookCreatedAt string
}
type Assignments struct {
	BookID int
	UserID int
}
type CategoryRequest struct {
	Category string `json:"category"`
}
type AssignRequest struct {
	BookId int `json:"bookId"`
	UserId int `json:"userId"`
}
type NameRequest struct {
	Name string `json:"name"`
}
type IsbnRequest struct {
	Isbn string `json:"isbn"`
}
type MyNullString struct {
	sql.NullString
}
type ErrorResponse struct {
	ErrorType string `json:"error type"`
	Error     string `json:"error"`
}

var Categories = map[string]int{"roman": 1, "hikaye": 2, "sosyoloji": 3, "gerilim": 4, "tarih": 5, "psikoloji": 6, "aşk": 7, "çocuk": 8, "fantastik": 9, "edebiyat": 10}

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func (newBook *Book) Valid() error {

	if err := newBook.ValidateBookIsbn(); err != nil {
		return err
	}

	if err := newBook.ValidateBookName(); err != nil {
		return err
	}

	if err := newBook.ValidateAuthor(); err != nil {
		return err
	}

	if err := newBook.ValidateCategory(); err != nil {
		return err
	}

	for index, _ := range strconv.Itoa(newBook.PageCount) {
		if _, err := strconv.Atoi(strconv.Itoa(index)); err != nil {
			return errors.New("page error, cannot contain char")
		}
	}

	return nil
}

func (newBook *Book) ValidateBookIsbn() error {

	if IsBlank(newBook.ISBN) == true {
		return errors.New("ISBN is blank")
	}

	if len(newBook.ISBN) == 13 {
		for index, _ := range newBook.ISBN {
			if _, err := strconv.Atoi(string(newBook.ISBN[index])); err != nil {
				return errors.New("ISBN error, cannot contain char")
			}
		}
	} else {
		return errors.New("ISBN must 13 characters")
	}

	return nil

}

func (newBook *Book) ValidateBookName() error {
	if IsBlank(newBook.BookName) == true {
		return errors.New("Book name is blank")
	}

	for index, _ := range newBook.BookName {
		if _, err := strconv.Atoi(string(newBook.BookName[index])); err == nil {
			return errors.New("book name error, cannot contain number")
		}
	}
	return nil
}

func (newBook *Book) ValidateAuthor() error {
	if IsBlank(newBook.Author) == true {
		return errors.New("Author is blank")
	}

	for index, _ := range newBook.Author {
		if _, err := strconv.Atoi(string(newBook.Author[index])); err == nil {
			return errors.New("author error, cannot contain number")
		}
	}
	return nil
}

func (newBook *Book) ValidateCategory() error {
	if IsBlank(newBook.Category) == true {
		return errors.New("Category is blank")
	}

	for index, _ := range newBook.Category {
		if _, err := strconv.Atoi(string(newBook.Category[index])); err == nil {
			return errors.New("category error, cannot contain number")
		}
	}
	categoryFound := false

	for value, _ := range Categories {
		if value == newBook.Category {
			categoryFound = true
			break
		}
	}

	if !categoryFound {
		return errors.New("Category not found")
	}

	return nil
}

func ValidGetBookByIsbnRequest(s string) error {

	if IsBlank(s) == true {
		return errors.New(" ISBN request is blank")
	}

	if len(s) != 13 {
		return errors.New(" ISBN request length error, string length is " + strconv.Itoa(len(s)))
	}

	for index, _ := range s {
		if _, err := strconv.Atoi(string(s[index])); err != nil {
			return errors.New(" ISBN error, cannot contain char")
		}
	}

	return nil
}
func ValidGetBookByCategory(s string) (error, int) {
	if IsBlank(s) == true {
		return errors.New(" Category request is blank"), 0
	}
	isCategoryFound := false
	foundId := 0
	for category, id := range Categories {
		if category == s {
			isCategoryFound = true
			foundId = id
			break
		}
	}
	if !isCategoryFound {
		return errors.New("Category not found"), 0
	}

	return nil, foundId
}
func ValidGetBookByBookName(s string) error {

	if IsBlank(s) == true {
		return errors.New(" Book name is blank")
	}

	for index, _ := range s {
		if _, err := strconv.Atoi(string(s[index])); err == nil {
			return errors.New("book name error, cannot contain number")
		}
	}

	return nil
}
