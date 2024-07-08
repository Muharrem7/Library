package main

import (
	"encoding/json"
	"fmt"
	bookPackage "library/Book"
	usersPackage "library/User"
	"net/http"
)

var users []usersPackage.User
var books []bookPackage.Book
var usersMap = make(map[int][]bookPackage.Book)
var booksMap = make(map[string][]usersPackage.User)
var newUser usersPackage.User
var newBook bookPackage.Book

var userIdCounter = 0

// id ve isbn'e göre bilgileri getirsin.
func main() {
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/user/users", usersHandler)
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/book/books", booksHandler)
	http.HandleFunc("/book/assign", assigneBookHandler)
	http.HandleFunc("/book/assigned", assignedBooksHandler)
	http.HandleFunc("/user/assigned", assignedUsersHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/filter/categories", filterCategoriesHandler)
	http.HandleFunc("/filter/names", filterNamesHandler)
	http.HandleFunc("/filter/id", filterIdHandler)
	http.HandleFunc("/filter/isbn", filterIsbnHandler)

	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server:", err)

	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := newUser.Valid(); err != nil {
			http.Error(w, "Invalid user data "+err.Error(), http.StatusBadRequest)
			return
		}

		userIdCounter++
		newUser.Id = userIdCounter
		users = append(users, newUser)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
		fmt.Println("User created successfully", newUser.GetUserInformations())

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func usersHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(users)
		if err != nil {
			fmt.Println("Error encoding users:", err)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}

}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := newBook.Valid(); err != nil {
			http.Error(w, "Invalid book data "+err.Error(), http.StatusBadRequest)
			return

		}

		books = append(books, newBook)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newBook)
		fmt.Println("Book created successfully", newBook.GetBookInformations())

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(books)
		if err != nil {
			http.Error(w, "Error encoding books: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func assigneBookHandler(w http.ResponseWriter, r *http.Request) {

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

	foundUser, foundBook, err := bookPackage.AssignBook(users, books, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	usersMap[req.UserId] = append(usersMap[req.UserId], foundBook)
	booksMap[req.BookIsbn] = append(booksMap[req.BookIsbn], foundUser)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode("Book assigned successfully")
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	fmt.Println("Book assigned successfully to user", usersMap)
}
func assignedBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(usersMap)
		if err != nil {
			http.Error(w, "Error assigned book", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func assignedUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(booksMap)
		if err != nil {
			http.Error(w, "Error assigned user", http.StatusInternalServerError)

		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(bookPackage.GetCategories()) // json dönecek
		if err != nil {
			http.Error(w, "Error encoding categories: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func filterCategoriesHandler(w http.ResponseWriter, r *http.Request) {
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

	categoriesMap, err := bookPackage.CategoryFilter(books, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(categoriesMap)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func filterNamesHandler(w http.ResponseWriter, r *http.Request) {
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

	namesMap, err := bookPackage.NameFilter(books, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(namesMap)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func filterIdHandler(w http.ResponseWriter, r *http.Request) {
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

	UserIdMap, err := usersPackage.IdFilter(users, req)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(UserIdMap)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func filterIsbnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var req bookPackage.IsbnRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	BookIsbnMap, err := bookPackage.IsbnFilter(books, req)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(BookIsbnMap)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
