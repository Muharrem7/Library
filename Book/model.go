package Book

import (
	"database/sql"
	"errors"
	"strconv"
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

func (book *Book) Valid() error {

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

	for index, _ := range strconv.Itoa(book.PageCount) {
		if _, err := strconv.Atoi(strconv.Itoa(index)); err != nil {
			return errors.New("page error, cannot contain char")
		}
	}

	return nil
}

//func (mns *MyNullString) Scan(value interface{}) error {
//
//	if err := mns.NullString.Scan(value); err != nil {
//		return err
//	}
//	if value == nil {
//		mns.Valid = false
//	} else {
//		mns.Valid = true
//	}
//	return nil
//
//}
