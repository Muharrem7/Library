package Book

import (
	"database/sql"
	"fmt"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (bo *BookRepository) GetCategory(req CategoryRequest) (int, error) {
	categories := req.Category
	rows, err := bo.db.Query("Select id From categories where name = ? ", categories)
	if err != nil {
		return 0, fmt.Errorf("Error while querying categories: %v", err)
	}
	defer rows.Close()
	var id int
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("Error while querying categories: %v", err)
		}
	}
	return id, nil
}

func (bo *BookRepository) InsertBook(book Book, categoryId int) error {

	// TODO: kategoriler db ye eklenecek, ondan sonra integer ID si alınacak
	stmt, err := bo.db.Prepare("INSERT INTO books (isbn,name,author,category_id,page) VALUES (?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("Error while inserting book: %v", err)

	}
	defer stmt.Close()
	_, err = stmt.Exec(book.ISBN, book.BookName, book.Author, categoryId, book.PageCount)
	if err != nil {
		return fmt.Errorf("Error while execute query: %v", err)
	}

	return nil
}

func (bo *BookRepository) GetBooks() ([]Book, error) {

	rows, err := bo.db.Query("SELECT * FROM books ORDER BY created_at desc ")
	if err != nil {
		return nil, fmt.Errorf("Error while querying book: %v", err)
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.ISBN, &book.BookName, &book.Author, &book.Category, &book.PageCount, &book.BookCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Error while querying book: %v", err)
		}
		books = append(books, book)
	}
	return books, nil

}

func (bo *BookRepository) GetBookIdById(req AssignRequest) (int, error) {
	bookId := req.BookId
	rows, err := bo.db.Query("select id from books where Id = ?", bookId)
	if err != nil {
		return 0, fmt.Errorf("Error while querying book: %v", err)
	}
	var id int
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("Error while querying book: %v", err)
		}
	}

	return id, nil

}

func (bo *BookRepository) InsertAssigne(userId int, bookId int) error {

	stmt, err := bo.db.Prepare("INSERT INTO assignments (user_id, book_id) VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("Error while inserting assigne: %v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId, bookId)
	if err != nil {
		return fmt.Errorf("Error while inserting assigne: %v", err)
	}
	return nil

}
func (bo *BookRepository) GetAssignedBook() ([]Assignments, error) {
	rows, err := bo.db.Query("SELECT  book_id, user_id FROM assignments ")
	if err != nil {
		return nil, fmt.Errorf("Error while querying book: %v", err)
	}
	defer rows.Close()
	var books []Assignments
	for rows.Next() {
		var id Assignments
		err := rows.Scan(&id.BookID, &id.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning assigned book: %v", err)
		}
		books = append(books, id)
	}
	return books, nil
}

func (bo *BookRepository) GetBookByCategories(req CategoryRequest) ([]Book, error) {
	getCategoryId, err := bo.GetCategory(req)
	if err != nil {
		return nil, fmt.Errorf("Error getting request: %v", err)
	}
	rows, err := bo.db.Query("Select * From books where category_id = ? ", getCategoryId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying book: %v", err)
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.ISBN, &book.BookName, &book.Author, &book.Category, &book.PageCount, &book.BookCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Error while scaning book: %v", err)
		}
		books = append(books, book)

	}
	return books, nil
}

func (bo *BookRepository) GetBookByBookName(req NameRequest) ([]Book, error) {
	name := req.Name
	rows, err := bo.db.Query("select * from books where name = ? ", name)
	if err != nil {
		return nil, fmt.Errorf("Error while querying book: %v", err)
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var book Book

		err := rows.Scan(&book.Id, &book.ISBN, &book.BookName, &book.Author, &book.Category, &book.PageCount, &book.BookCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Error while scaning book: %v", err)
		}

		books = append(books, book)
	}

	return books, nil

}
func (bo *BookRepository) GetCategories() ([]string, error) {
	rows, err := bo.db.Query("select name from categories")
	if err != nil {
		return nil, fmt.Errorf("Error while querying categories: %v", err)
	}

	var categories []string
	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("Error while querying categories: %v", err)
		}
		categories = append(categories, name)
	}
	return categories, nil
}
func (bo *BookRepository) GetBooksByIsbn(req IsbnRequest) ([]Book, error) {
	isbn := req.Isbn
	rows, err := bo.db.Query("Select * from Books where isbn = ?", isbn)
	if err != nil {
		return nil, fmt.Errorf("Error while querying book: %v", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.ISBN, &book.BookName, &book.Author, &book.Category, &book.PageCount, &book.BookCreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning book: %v", err)
		}
		books = append(books, book)
	}
	return books, nil
}
