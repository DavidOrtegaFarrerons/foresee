package web

import (
	"net/http"
	"path"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func NeuteredFileSystem(fs http.FileSystem) http.FileSystem {
	return neuteredFileSystem{fs: fs}
}

func (nfs neuteredFileSystem) Open(p string) (http.File, error) {
	f, err := nfs.fs.Open(p)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := path.Join(p, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
