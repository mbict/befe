package dsl

import (
	"context"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type noListingFilesystem struct {
	fs http.FileSystem
}

func (fs noListingFilesystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	//if is dir, try to check if there is an index file
	if fi.IsDir() {
		if _, err = fs.Open(filepath.Join(path, "index.html")); err != nil {
			return nil, os.ErrNotExist
		}
		return f, nil
	}

	return f, nil
}

type HTTPFileServer interface {
	Action
	ServeFile(path string) Action

	AllowDirectoryListing(allow bool) HTTPFileServer

	WhenNotFound(actions ...Action) HTTPFileServer
	WhenFound(actions ...Action) HTTPFileServer
}

type fileServer struct {
	fileSystem            http.FileSystem
	allowDirectoryListing bool
	notFoundActions       Actions
	onFoundActions        Actions
}

func (fs *fileServer) BuildHandler(ctx context.Context, next Handler) Handler {

	fileSystem := fs.fileSystem
	if fs.allowDirectoryListing == false {
		fileSystem = noListingFilesystem{fileSystem}
	}
	fileServer := http.FileServer(fileSystem)

	notFoundActions := fs.notFoundActions
	if len(notFoundActions) == 0 {
		notFoundActions = Actions{ResponseCode(404)}
	}
	notFoundHandler := notFoundActions.BuildHandler(ctx, emptyHandler)
	serveFileHandler := fs.onFoundActions.BuildHandler(ctx, fileServer.ServeHTTP)

	return func(w http.ResponseWriter, r *http.Request) {
		f, err := fileSystem.Open(r.URL.Path)
		if os.IsNotExist(err) {
			notFoundHandler(w, r)
			return
		}
		f.Close()

		serveFileHandler(w, r)

		next(w, r)
	}
}

func (fs *fileServer) ServeFile(path string) Action {
	return ActionBuilder(func(ctx context.Context, next Handler) Handler {
		return func(rw http.ResponseWriter, req *http.Request) {
			f, err := fs.fileSystem.Open(path)
			if err != nil {
				panic(err)
			}

			fi, err := f.Stat()
			if err != nil {
				panic(err)
			}

			//read the first few bytes, if we need to peek for the content type
			buf := make([]byte, 1024)
			numRead, _ := f.Read(buf)

			ctype := mime.TypeByExtension(filepath.Ext(path))
			if ctype == "" {
				ctype = http.DetectContentType(buf[0:numRead])
			}

			rw.Header().Add("Content-Type", ctype)
			rw.Header().Add("Content-Length", strconv.Itoa(int(fi.Size())))
			rw.WriteHeader(200)

			rw.Write(buf[0:numRead])
			if _, err := io.Copy(rw, f); err != nil {
				panic(err)
			}
			defer f.Close()

			next(rw, req)
		}
	})
}

func (fs *fileServer) AllowDirectoryListing(allow bool) HTTPFileServer {
	fs.allowDirectoryListing = allow
	return fs
}

func (fs *fileServer) WhenNotFound(actions ...Action) HTTPFileServer {
	fs.notFoundActions = append(fs.notFoundActions, actions...)
	return fs
}

func (fs *fileServer) WhenFound(actions ...Action) HTTPFileServer {
	fs.onFoundActions = append(fs.onFoundActions, actions...)
	return fs
}

func FileServer(path string) HTTPFileServer {
	return &fileServer{
		fileSystem: http.Dir(path),
	}
}
