package main

import (
	"encoding/json"
	"fmt"
	bookPackage "library/Book_Information"
	usersPackage "library/Users_Personel_Informations"
	"net/http"
)

var users []usersPackage.User
var books []bookPackage.BookInformation
var mapsOfUsersBook = make(map[int][]bookPackage.BookInformation)
var mapsOfBooksUser = make(map[string][]usersPackage.User)
var newUser usersPackage.User
var newBook bookPackage.BookInformation

var userIdCounter = 0

// id ve isbn'e göre bilgileri getirsin.
func main() {
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/assign-book", assignBookHandler)
	http.HandleFunc("/assigned-books", assignedBooksHandler)
	http.HandleFunc("/assigned-users", assignedUsersHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/filter-categories", filterCategoriesHandler)
	http.HandleFunc("/filter-names", filterNamesHandler)
	http.HandleFunc("/filter-id", filterUserIdHandler)
	http.HandleFunc("/filter-isbn", filterBookIsbnHandler)

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

func assignBookHandler(w http.ResponseWriter, r *http.Request) {

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

	mapsOfUsersBook[req.UserId] = append(mapsOfUsersBook[req.UserId], foundBook)
	mapsOfBooksUser[req.BookIsbn] = append(mapsOfBooksUser[req.BookIsbn], foundUser)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode("Book assigned successfully")
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	fmt.Println("Book assigned successfully to user", mapsOfUsersBook)
}
func assignedBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(mapsOfUsersBook)
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
		err := json.NewEncoder(w).Encode(mapsOfBooksUser)
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

	mapsOfCategory, err := bookPackage.CategoryFilter(books, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(mapsOfCategory)
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

	mapsOfName, err := bookPackage.NameFilter(books, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(mapsOfName)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func filterUserIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var req bookPackage.IdRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	mapsOfUserId, err := bookPackage.IdFilter(users, req)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(mapsOfUserId)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func filterBookIsbnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var req bookPackage.IsbnRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	mapsOfBookIsbn, err := bookPackage.IsbnFilter(books, req)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(mapsOfBookIsbn)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
