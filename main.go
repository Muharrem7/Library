package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	bookPackage "library/Book"
	usersPackage "library/User"
	"net/http"
	"strconv"
)

// TODO boş geçemez = OK
//TODO  düzgün insert yapılmadıysa cevap dönmesin = OK
// TODO endpoint isimlerini düzelt
// TODO kategorlileri validasyon için kodda tut = OK

var newUser usersPackage.User
var newBook bookPackage.Book

const (
	username = "root"
	password = "Smorkan2"
	hostname = "127.0.0.1:3306"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func DbConnect() *sql.DB {
	db, dbErr := sql.Open("mysql", dsn("libraryDB"))
	if dbErr != nil {
		fmt.Println("db connect err:", dbErr)
	}

	return db
}

func main() {

	http.HandleFunc("/user", insertUser)
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/users/", getUserById)            // TODO /users/{id} = OK
	http.HandleFunc("/users/assign", getAssignedUsers) // TODO birde bir kitabı aynı kullanıcıya bir daha ataymazsın onu kontrol et. ok // users/{id}/books

	http.HandleFunc("/book", insertBook)
	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/books/assign", insertAssign)
	http.HandleFunc("/books/name", getBookByBookName)          //TODO buna gerek yok, /books a taşı URL Parameter
	http.HandleFunc("/books/isbn", getBookByIsbn)              //TODO  GET /book/{isbn} --> /books/1204 = OK
	http.HandleFunc("/books/categories", getCategories)        // TODO /books/categories = OK
	http.HandleFunc("/books/categories/", getBookByCategories) // TODO buna gerek yok, /books a taşı URL Parameter, TODO books olacak /books?category=roman  = OK
	http.HandleFunc("/books/assigned", getAssignedBooks)       // TODO book id ye göre filtreke, user id leri döndür. ok,        books/{id}/users

	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			errorMessage := usersPackage.ErrorResponse{ErrorType: "Invalid request Body", Error: err.Error()}
			json.NewEncoder(w).Encode(errorMessage)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := newUser.Valid(); err != nil {
			errorMessage := usersPackage.ErrorResponse{ErrorType: "Validation error", Error: err.Error()}
			json.NewEncoder(w).Encode(errorMessage)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userRepo := usersPackage.NewUserRepository(DbConnect())
		userService := usersPackage.NewUserService(*userRepo)
		err = userService.InsertUser(newUser)
		if err != nil {
			errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while inserting new user", Error: err.Error()}
			json.NewEncoder(w).Encode(errorMessage)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			err = json.NewEncoder(w).Encode(newUser)
			if err != nil {
				errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while encoding new user", Error: err.Error()}
				json.NewEncoder(w).Encode(errorMessage)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		userList := usersPackage.NewUserRepository(DbConnect())
		userListService := usersPackage.NewUserService(*userList)
		users, err := userListService.GetUsers()
		if err != nil {
			errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while getting users", Error: err.Error()}
			json.NewEncoder(w).Encode(errorMessage)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while encoding users", Error: err.Error()}
			json.NewEncoder(w).Encode(errorMessage)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func insertBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Invalid request Body", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := newBook.Valid(); err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Validation error", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	category := bookPackage.CategoryRequest{Category: newBook.Category}
	bookRepo := bookPackage.NewBookRepository(DbConnect())
	bookRepoService := bookPackage.NewBookService(*bookRepo)
	categoryId, err := bookRepoService.GetCategory(category) // TODO: istersen constant çek = OK
	if err != nil {
		errorMesage := bookPackage.ErrorResponse{ErrorType: "Error while getting category", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMesage)
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	err = bookRepoService.InsertBook(newBook, categoryId)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while inserting book", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(newBook)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while encoding new book", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		bookList := bookPackage.NewBookRepository(DbConnect())
		bookListService := bookPackage.NewBookService(*bookList)
		books, err := bookListService.GetBooks()
		if err != nil {
			errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while getting books", Error: err.Error()}
			json.NewEncoder(w).Encode(errorMessage)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(books)
		if err != nil {
			errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while encoding books", Error: err.Error()}
			json.NewEncoder(w).Encode(errorMessage)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusOK)
		}

	} else {

		json.NewEncoder(w).Encode("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func insertAssign(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req bookPackage.AssignRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Invalid request Body", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bookIdReq := bookPackage.AssignRequest{BookId: req.BookId}
	userIdReq := bookPackage.AssignRequest{UserId: req.UserId}

	bookIdRepo := bookPackage.NewBookRepository(DbConnect())
	bookIdService := bookPackage.NewBookService(*bookIdRepo)
	bookId, err := bookIdService.GetBookIdById(bookIdReq)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while getting bookId", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userIdRepo := usersPackage.NewUserRepository(DbConnect())
	userIdService := usersPackage.NewUserService(*userIdRepo)
	userId, err := userIdService.GetUserIdById(userIdReq)
	if err != nil {
		errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while getting userId", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userId == 0 || bookId == 0 {
		http.Error(w, "Error while get user/book id", http.StatusInternalServerError)
		//w.WriteHeader(http.StatusNotFound)
		return
	}
	assignRepo := bookPackage.NewBookRepository(DbConnect())
	assignRepoService := bookPackage.NewBookService(*assignRepo)
	err = assignRepoService.InsertAssigne(userId, bookId)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while assigning book", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode("Book assigned successfully")
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while encode respond event", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getAssignedBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	booksRepo := bookPackage.NewBookRepository(DbConnect())
	assignedRepo := bookPackage.NewBookService(*booksRepo)
	ids, err := assignedRepo.GetAssignedBook()
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while getting assigned books", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(ids)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while encode assigned books", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getAssignedUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	userRepo := usersPackage.NewUserRepository(DbConnect())
	assignedRepo := usersPackage.NewUserService(*userRepo)
	ids, err := assignedRepo.GetAssignedUser()
	if err != nil {
		errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while getting assigned users", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(ids)
	if err != nil {
		errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while encode assigned users", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	bookRepo := bookPackage.NewBookRepository(DbConnect())
	category := bookPackage.NewBookService(*bookRepo)
	categories, err := category.GetCategories()
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while getting categories", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while encode categories", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getBookByCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	v := r.FormValue("category")

	err, categoryId := bookPackage.ValidGetBookByCategory(v)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while validate categories", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bookRepo := bookPackage.NewBookRepository(DbConnect())
	getBookByCategoryService := bookPackage.NewBookService(*bookRepo)
	books, err := getBookByCategoryService.GetBookByCategories(categoryId)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while getting book by category", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(books); err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while encode book by category", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getBookByBookName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	v := r.FormValue("name")
	err := bookPackage.ValidGetBookByBookName(v)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while validate book name", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bookRepo := bookPackage.NewBookRepository(DbConnect())
	getBookByNameService := bookPackage.NewBookService(*bookRepo)
	books, err := getBookByNameService.GetBookByBookName(v)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while getting book by name", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while encode book by name", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	v, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while convert user id", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userRepo := usersPackage.NewUserRepository(DbConnect())
	getBookByIdService := usersPackage.NewUserService(*userRepo)
	user, err := getBookByIdService.GetUserById(v)
	if err != nil {
		errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while getting user id", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		errorMessage := usersPackage.ErrorResponse{ErrorType: "Error while encode user id", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func getBookByIsbn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	v := r.FormValue("isbn")
	err := bookPackage.ValidGetBookByIsbnRequest(v)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while validate isbn", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bookRepo := bookPackage.NewBookRepository(DbConnect())
	getBookByIsbnService := bookPackage.NewBookService(*bookRepo)
	book, err := getBookByIsbnService.GetBookByIsbn(v)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		errorMessage := bookPackage.ErrorResponse{ErrorType: "Error while encode book by isbn", Error: err.Error()}
		json.NewEncoder(w).Encode(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}
