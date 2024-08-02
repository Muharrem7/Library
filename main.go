package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	bookPackage "library/Book"
	usersPackage "library/User"
	"net/http"
)

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
	http.HandleFunc("/user/users", getUsers)
	http.HandleFunc("/book", insertBook)
	http.HandleFunc("/book/books", getBooks)
	http.HandleFunc("/book/assign", insertAssign)
	http.HandleFunc("/book/assigned", getAssignedBooks) // TODO book id ye göre filtreke, user id leri döndür. ok
	http.HandleFunc("/user/assigned", getAssignedUsers) // TODO birde bir kitabı aynı kullanıcıya bir daha ataymazsın onu kontrol et. ok
	http.HandleFunc("/book/categories", getCategories)
	http.HandleFunc("/book/categories/book-category", getBookByCategories)
	http.HandleFunc("/book/book-name", getBookByBookName)
	http.HandleFunc("/user/{id}", getUserById)
	http.HandleFunc("/book/isbn", getBookByIsbn)

	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		// TODO: user address Description eklenecek
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Invalid request body:"+err.Error(), http.StatusBadRequest)
			return
		}

		if err := newUser.Valid(); err != nil {
			http.Error(w, "Invalid user data:"+err.Error(), http.StatusBadRequest)
			return
		}

		//city := usersPackage.CityRequest{City: newUser.UserAdress.UserCity}
		//district := usersPackage.DistrictRequest{District: newUser.UserAdress.UserDistrict}

		userRepo := usersPackage.NewUserRepository(DbConnect())
		userService := usersPackage.NewUserService(*userRepo)
		//cityId, err := userService.GetCities(city)
		if err != nil {
			http.Error(w, "Error inserting new user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		//districtId, err := userService.GetDistricts(district)
		if err != nil {
			http.Error(w, "Error inserting new user: "+err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		err = userService.InsertUser(newUser)
		if err != nil {
			http.Error(w, "Error inserting new user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(newUser)
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		userList := usersPackage.NewUserRepository(DbConnect())
		userListService := usersPackage.NewUserService(*userList)
		users, err := userListService.GetUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			http.Error(w, "Error encoding users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func insertBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := newBook.Valid(); err != nil {
		http.Error(w, "Invalid book data "+err.Error(), http.StatusBadRequest)
		return

	}
	category := bookPackage.CategoryRequest{Category: newBook.Category}
	bookRepo := bookPackage.NewBookRepository(DbConnect())
	bookRepoService := bookPackage.NewBookService(*bookRepo)
	categoryId, err := bookRepoService.GetCategory(category) // TODO: istersen constant çek,
	if err != nil {
		http.Error(w, "Error inserting new book: "+err.Error(), http.StatusInternalServerError)
	}
	err = bookRepoService.InsertBook(newBook, categoryId)
	if err != nil {
		http.Error(w, "Error inserting new book: "+err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(newBook)
	w.WriteHeader(http.StatusCreated)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		bookList := bookPackage.NewBookRepository(DbConnect())
		bookListService := bookPackage.NewBookService(*bookList)
		books, err := bookListService.GetBooks()
		if err != nil {
			http.Error(w, "Error selecting books", http.StatusInternalServerError)
			return // TODO error döndüğü zaman nolacak bir test edilmeli?
		}
		err = json.NewEncoder(w).Encode(books)
		w.WriteHeader(http.StatusOK)
		if err != nil {
			http.Error(w, "Error encoding books: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bookIdReq := bookPackage.AssignRequest{BookId: req.BookId}
	userIdReq := bookPackage.AssignRequest{UserId: req.UserId}

	bookIdRepo := bookPackage.NewBookRepository(DbConnect())
	bookIdService := bookPackage.NewBookService(*bookIdRepo)
	bookId, err := bookIdService.GetBookIdById(bookIdReq)
	if err != nil {
		http.Error(w, "Error while get book id"+err.Error(), http.StatusInternalServerError)
	}

	userIdRepo := usersPackage.NewUserRepository(DbConnect())
	userIdService := usersPackage.NewUserService(*userIdRepo)
	userId, err := userIdService.GetUserIdById(userIdReq)
	if err != nil {
		http.Error(w, "Error while get user id"+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode("Book assigned successfully")
}

func getAssignedBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	booksRepo := bookPackage.NewBookRepository(DbConnect())
	assignedRepo := bookPackage.NewBookService(*booksRepo)
	id, err := assignedRepo.GetAssignedBook()
	if err != nil {
		http.Error(w, "Error assigned book", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(id)
	w.WriteHeader(http.StatusOK)

}

func getAssignedUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	userRepo := usersPackage.NewUserRepository(DbConnect())
	assignedRepo := usersPackage.NewUserService(*userRepo)
	usersId, err := assignedRepo.GetAssignedUser()
	if err != nil {
		http.Error(w, "Error assigned user", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(usersId)
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
		http.Error(w, "Error getting categories: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		http.Error(w, "Error encoding categories: "+err.Error(), http.StatusInternalServerError)
	}
}

func getBookByCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req bookPackage.CategoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	bookRepo := bookPackage.NewBookRepository(DbConnect())
	getBookByCategoryService := bookPackage.NewBookService(*bookRepo)
	books, err := getBookByCategoryService.GetBookByCategories(req)
	if err != nil {
		http.Error(w, "Error getting books: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "Error encoding books: "+err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

}

func getBookByBookName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req bookPackage.NameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bookRepo := bookPackage.NewBookRepository(DbConnect())
	getBookByNameService := bookPackage.NewBookService(*bookRepo)
	books, err := getBookByNameService.GetBookByBookName(req)
	if err != nil {
		http.Error(w, "Failed to get book: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req usersPackage.IdRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userRepo := usersPackage.NewUserRepository(DbConnect())
	getBookByIdService := usersPackage.NewUserService(*userRepo)
	user, err := getBookByIdService.GetUserById(req)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func getBookByIsbn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var req bookPackage.IsbnRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	bookRepo := bookPackage.NewBookRepository(DbConnect())
	getBookByIsbnService := bookPackage.NewBookService(*bookRepo)
	book, err := getBookByIsbnService.GetBookByIsbn(req)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
