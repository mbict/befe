package expr

import (
	"context"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type HTTPFileServer interface {
	Action

	ServeFile(path string) Action

	AllowDirectoryListing(allow bool) HTTPFileServer

	WhenNotFound(actions ...Action) HTTPFileServer
	WhenFound(actions ...Action) HTTPFileServer
}

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
		notFoundActions = Actions{
			ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
				rw.WriteHeader(http.StatusNotFound)
				return true, nil
			}),
		}
	}

	notFoundHandler := notFoundActions.BuildHandler(ctx, nil)
	serveFileHandler := fs.onFoundActions.BuildHandler(ctx, func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		fileServer.ServeHTTP(rw, r)
		return true, nil
	})

	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		f, err := fileSystem.Open(r.URL.Path)
		if os.IsNotExist(err) {
			notFoundHandler(rw, r)
			return false, nil
		}
		f.Close()

		serveFileHandler(rw, r)

		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}

func (fs *fileServer) ServeFile(path string) Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
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

		return true, nil
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

func NewFileServer(path string) HTTPFileServer {
	return &fileServer{
		fileSystem: http.Dir(path),
	}
}
