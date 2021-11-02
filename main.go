package main

import (
	"context"
	"errors"
	"flag"
	"github.com/mbict/befe/buildin"
	"github.com/mbict/befe/dsl"
	"github.com/radovskyb/watcher"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

type Program struct {
	cancelFunc    context.CancelFunc
	handleRequest dsl.Handler
}

func (p *Program) HandleRequest(rw http.ResponseWriter, r *http.Request) {
	p.handleRequest(rw, r)
}

func (p *Program) DeInit() {
	p.cancelFunc()
}

var (
	scriptPath string
	addr       string

	program      *Program
	ready        bool
	compileError error
)

func init() {
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 1000
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000

	flag.StringVar(&scriptPath, "path", "./", "path where the script can be found (./your/program)")
	flag.StringVar(&addr, "addr", ":8080", "address and port to listen (:8080)")
	flag.Parse()
}

func main() {
	var err error

	//init program
	log.Printf("loading program from :%s", scriptPath)
	program, err = loadProgram(scriptPath)
	if err != nil {
		log.Fatalf("cannot compile program: %s", err)
	}

	//handleRequest the connections
	var e chan error
	go func() {
		log.Printf("starting engine and listen on :%s", addr)
		e <- http.ListenAndServe(addr, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			//handle internal errors
			defer func() {
				r := recover()
				if r != nil {
					var err error
					switch t := r.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = errors.New("unknown error")
					}
					log.Printf("panic in request handling: %s", err)
					http.Error(rw, "internal error", http.StatusInternalServerError)
				}
			}()

			//handleRequest connection
			program.handleRequest(rw, r)
		}))
	}()

	//watch for changes and recompile program if needed
	log.Printf("watching for file changes in dir: %s", scriptPath)
	watchScriptDir(scriptPath)

	log.Println(<-e)

	//de initialize the running program.
	program.DeInit()
}

func watchScriptDir(dir string) {
	w := watcher.New()
	w.SetMaxEvents(1)

	// Only notify rename and move events.
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Write, watcher.Create, watcher.Remove)

	go func() {
		var lock sync.Mutex
		for {
			select {
			case <-w.Event:
				lock.Lock()
				log.Printf("detected changes in script, reloading program")
				p, err := loadProgram(scriptPath)
				if err == nil {
					//we replace the program with the newly loaded one, and keep a copy of the old one

					log.Printf("new program running")
					oldProgram := program
					program = p

					//deinit the old program and clean up, shared resources, go routines etc
					log.Printf("deInit previous program")
					oldProgram.DeInit()
				} else {
					log.Printf("cannot compile new program: %s", err)
				}
				compileError = err
				lock.Unlock()
			case err := <-w.Error:
				log.Printf("error in watch script directory: %s", err)
				return
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive(dir); err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err := w.Start(time.Second); err != nil {
			log.Fatalln(err)
		}
	}()
}

func loadProgram(path string) (*Program, error) {
	interpreter, err := loadScript(path)
	if err != nil {
		return nil, err
	}

	//in case we have an init function we initialize
	vInit, err := interpreter.Eval("main.Init")
	if err == nil {
		initFunc := vInit.Interface().(func() error)
		if err := initFunc(); err != nil {
			return nil, err
		}
	}

	var handleFunc dsl.Handler

	//check if there is a dsl defined in the program endpoint
	vDsl, err := interpreter.Eval("main.Program")
	var program dsl.Action
	if err != nil {
		return nil, err
	}

	programFunc := vDsl.Interface().(func() dsl.Action)
	program = programFunc()

	ctx, cancelFunc := context.WithCancel(context.Background())
	handleFunc = program.BuildHandler(ctx, func(rw http.ResponseWriter, r *http.Request) {})

	return &Program{
		cancelFunc:    cancelFunc,
		handleRequest: handleFunc,
	}, nil
}

func loadScript(sourcePath string) (*interp.Interpreter, error) {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	i.Use(buildin.Symbols)

	dir, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		return nil, err
	}

	for _, fs := range dir {
		//if fs.IsDir() || (fs.Mode()&os.ModeSymlink) == os.ModeSymlink || strings.HasSuffix(fs.Name(), "_test.go") {
		if fs.IsDir() || (strings.HasPrefix(fs.Name(), ".") || strings.HasSuffix(fs.Name(), "_test.go") || !strings.HasSuffix(fs.Name(), ".go")) {
			log.Printf("ignoring file : %s", path.Join(sourcePath, fs.Name()))
			continue
		}

		log.Printf("loading file : %s", path.Join(sourcePath, fs.Name()))
		source, err := ioutil.ReadFile(path.Join(sourcePath, fs.Name()))
		if err != nil {
			return nil, err
		}

		_, err = i.Eval(string(source))
		if err != nil {
			return nil, err
		}
	}

	return i, nil
}
