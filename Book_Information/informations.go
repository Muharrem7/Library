package Book_Information

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type BookInformation struct {
	ISBN     string // TODO string
	BookName string
	Author   string
	Category string
	Pages    string
}

func NewBookInformation(isbn string, bookName string, author string, category string, pages string) BookInformation {
	book := BookInformation{
		ISBN:     isbn,
		BookName: bookName,
		Author:   author,
		Category: category,
		Pages:    pages,
	}

	return book
}

func (book BookInformation) GetBookInformations() string {

	return fmt.Sprintf("\n ISBN: %s  Book Name: %s  Author: %s  Category: %s  Pages: %s",
		book.ISBN,
		book.BookName,
		book.Author,
		book.Category,
		book.Pages,
	)
}

// TODO GetCategories
func GetCategories() [10]string {
	return [10]string{"Roman", "Hikaye", "Sosyoloji", " Gerilim", "Tarih", "Psikoloji", "Aşk", "Çocuk", "Fantastik", "Edebiyat"}
}

func (book *BookInformation) Valid() error {

	categorys := GetCategories()

	if len(book.ISBN) == 13 {
		for index, _ := range book.ISBN {
			if _, err := strconv.Atoi(string(book.ISBN[index])); err != nil {
				return errors.New("kitabin ISBN kodu hatalı")
			}
		}
	} else {
		return errors.New("kitabin ISBN kodu 13 haneli olmalı")
	}

	for index, _ := range book.BookName {
		if _, err := strconv.Atoi(string(book.BookName[index])); err == nil {
			return errors.New("kitabin ismi hatalı")
		}
	}

	for index, _ := range book.Author {
		if _, err := strconv.Atoi(string(book.Author[index])); err == nil {
			return errors.New("kitabin yazarinin ismi hatalı")
		}
	}

	for index, _ := range book.Category {
		if _, err := strconv.Atoi(string(book.Category[index])); err == nil {
			return errors.New("kitabin kategorisi hatalı")
		}
	}

	categoryFound := false
	for _, category := range categorys {
		if strings.ToLower(category) == strings.ToLower(book.Category) {
			categoryFound = true
			continue

		}

	}
	if !categoryFound {
		return errors.New("bu kategoride kitap ekleyemezsiniz.")
	}

	for index, _ := range book.Pages {
		if _, err := strconv.Atoi(string(book.Pages[index])); err != nil {
			return errors.New("kitabin sayfa sayısı hatalı")
		}
	}

	return nil
}
