package main

import (
    "os"
)

const PageSize int = 4096

type Pages [TableMaxPages]*Page

type Pager interface {
    getFile() *os.File
    getPages() Pages
}

func pagerOpen(filename string) *FilePager {
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
    if err != nil {
	log.Panic(err)
    }

    return newFilePager(file)
}

type FilePager struct {
    file *os.File
    fileLength int64
    pages Pages
}

func newFilePager(file *os.File) *FilePager {
    var pages Pages
    copy(pages[:], make([]*Page, TableMaxPages))
    
    n, err := file.Seek(0, 0)
    if err != nil {
	log.Panic(err)
    }

    pager := &FilePager{file, n, pages}

    return pager
}


func getPage(pager *FilePager, uint32 pageNum) *Page{
    if pageNum > TableMaxPages {
	fmt.Printf("Tried to fetch page number out of bounds. %d > %d\n", pageNum, TableMaxPages)
	os.Exit(1)
    }

    if page.pages[pageNum] == nil {
	page := newPage()
	numPages := pager.fileLength
	if pager.fileLength % PageSize != 0 {
	    numPages++
	}

	if pageNum <= numPages {
	    pager.file.Seek(pageNum * PageSize, 1)
	    page := make([]byte, PageSize)
	    _, err := pager.file.Read(page)
	    if err != nil {
		fmt.Println("Error reading file")
		os.Exit(1)
	    }
	}
    }

    return page.pages[pageNum]
}
