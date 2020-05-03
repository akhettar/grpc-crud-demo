package handler

import (
	"context"
	"github.com/akhettar/grpc-crud-demo/api"
	"github.com/akhettar/grpc-crud-demo/repository"
	"github.com/golang/protobuf/ptypes/empty"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"log"
)

// BookstoreServer
type BookstoreServer struct {
	repository.Repository
}

func (b *BookstoreServer) ListShelves(context.Context, *empty.Empty) (*api.ListShelvesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListShelves not implemented")
}
func (b *BookstoreServer) CreateShelf(ctx context.Context, req *api.CreateShelfRequest) (*api.Shelf, error) {
	log.Printf("reveieved request to create a shelves with Id %d and theme %s", req.GetShelf().Id, req.GetShelf().Theme)
	payload := repository.Shelf{Theme: req.GetShelf().Theme}
	res, err := b.AddShelf(payload)
	return &api.Shelf{Id: int64(res.ID), Theme: res.Theme}, err
}
func (b *BookstoreServer) GetShelf(ctx context.Context, req *api.GetShelfRequest) (*api.Shelf, error) {
	log.Printf("received request to retrieve shelf for id %d", req.Shelf)
	shelf, e := b.FindShelf(req.Shelf)
	if e != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve the shelf", e)
	}
	return &api.Shelf{Id: int64(shelf.ID), Theme: shelf.Theme}, nil
}
func (b *BookstoreServer) DeleteShelf(ctx context.Context, req *api.DeleteShelfRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteShelf not implemented")
}
func (b *BookstoreServer) ListBooks(ctx context.Context, req *api.ListBooksRequest) (*api.ListBooksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBooks not implemented")
}

// CreateBook add book to the repository
func (b *BookstoreServer) CreateBook(ctx context.Context, req *api.CreateBookRequest) (*api.Book, error) {
	log.Printf("received request to create a book with title %s", req.GetBook().Title)
	payload := repository.Book{Author: req.GetBook().Author, Title: req.GetBook().Title, ShelfId: req.Shelf}
	res, err := b.AddBook(payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to crate a book", err)
	}
	return &api.Book{Id: int64(res.ID), Title: res.Title, Author: res.Author}, nil
}

// GetBook get a book for given Id
func (b *BookstoreServer) GetBook(ctx context.Context, req *api.GetBookRequest) (*api.Book, error) {
	log.Printf("received request to retrieve a book with id %d", req.Book)
	book, err := b.FindBook(req.Book)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve the book", err)
	}
	return &api.Book{Id: int64(book.ID), Title: book.Title, Author: book.Author}, nil
}
func (b *BookstoreServer) DeleteBook(ctx context.Context, req *api.DeleteBookRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBook not implemented")
}

// Create a new Server
func NewServer() *BookstoreServer {
	repo := repository.NewSQLiteRepository()
	return &BookstoreServer{repo}
}
