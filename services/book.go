package services

import "github.com/pyxvlad/proiect-ipdp/database/types"

type BookService struct {
}

func NewBookService() BookService {
	return BookService{}
}

func (b *BookService) CreateBook(account types.AccountID, title string, ) {
}

