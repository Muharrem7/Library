package Book

type BookService struct {
	bookRepo BookRepository
}

func NewBookService(bookRepo BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}
func (bs *BookService) GetCategory(req CategoryRequest) (int, error) {
	return bs.bookRepo.GetCategory(req)
}
func (bs *BookService) InsertBook(book Book, categoryId int) error {
	return bs.bookRepo.InsertBook(book, categoryId)
}
func (bs *BookService) GetBooks() ([]Book, error) {
	return bs.bookRepo.GetBooks()
}
func (bs *BookService) GetBookIdById(req AssignRequest) (int, error) {
	return bs.bookRepo.GetBookIdById(req)

}
func (bs *BookService) InsertAssigne(userId int, bookId int) error {
	return bs.bookRepo.InsertAssigne(userId, bookId)

}
func (bs *BookService) GetAssignedBook() ([]Assignments, error) {
	return bs.bookRepo.GetAssignedBook()
}
func (bs *BookService) GetBookByCategories(req CategoryRequest) ([]Book, error) {
	return bs.bookRepo.GetBookByCategories(req)

}
func (bs *BookService) GetBookByBookName(req NameRequest) ([]Book, error) {
	return bs.bookRepo.GetBookByBookName(req)

}
func (bs *BookService) GetCategories() ([]string, error) {
	return bs.bookRepo.GetCategories()
}
func (bs *BookService) GetBookByIsbn(req IsbnRequest) ([]Book, error) {
	return bs.bookRepo.GetBooksByIsbn(req)
}
