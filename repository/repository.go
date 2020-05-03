package repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

// A book resource.
type Book struct {
	gorm.Model

	// An author of the book.
	Author string `json:"author,omitempty"`
	// A book title.
	Title string `json:"title,omitempty"`

	// The shelf Id
	ShelfId int64 `json:"shelf_id,omitempty"`
}

type Shelf struct {
	gorm.Model
	// Theme the theme of the shelve
	Theme string `json:"theme, omitempty"`
	Books []Book `gorm:"foreignkey:ShelfId"`
}

type Repository interface {
	AddShelf(s Shelf) (*Shelf, error)
	AddBook(b Book) (*Book, error)
	FindBook(id int64) (Book, error)
	FindShelf(id int64) (Shelf, error)
}

// SQLite Repository
type SQLiteRepository struct {
	*gorm.DB
}

// NewSQLiteRepository return a sequence of repository
func NewSQLiteRepository() Repository {
	db, err := gorm.Open("sqlite3", "books.db")
	if err != nil {
		log.Fatalf("failed to open database connection %v", err)
	}

	// Database migration
	db.AutoMigrate(&Shelf{}, &Book{})
	log.Println("database successfully migrated")
	return SQLiteRepository{db}
}

func (r SQLiteRepository) AddShelf(body Shelf) (*Shelf, error) {
	log.Println("inserting into database")
	db := r.DB.Create(&body)
	if db.Error != nil {
		log.Printf("error occured when saving into the database %v", db.Error)
		return nil, db.Error
	}
	return db.Value.(*Shelf), nil
}

func (r SQLiteRepository) AddBook(body Book) (*Book, error) {
	log.Println("inserting into database")
	db := r.DB.Create(&body)
	if db.Error != nil {
		log.Printf("error occured when saving into the database %v", db.Error)
		return nil, db.Error
	}
	return db.Value.(*Book), nil
}

func (r SQLiteRepository) FindBook(id int64) (Book, error) {
	book := new(Book)
	db := r.DB.Find(book, id)
	return *book, db.Error
}

func (r SQLiteRepository) FindShelf(id int64) (Shelf, error) {
	shelf := new(Shelf)
	db := r.DB.Find(shelf, id)
	return *shelf, db.Error
}
