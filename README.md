# Mir
[![Build Status](https://api.travis-ci.com/alimy/mir.svg?branch=master)](https://travis-ci.com/alimy/mir)
[![codecov](https://codecov.io/gh/alimy/mir/branch/master/graph/badge.svg)](https://codecov.io/gh/alimy/mir)
[![GoDoc](https://godoc.org/github.com/alimy/mir?status.svg)](https://pkg.go.dev/github.com/alimy/mir/v2)
[![sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/alimy/mir)

Mir is used for register handler to http router(eg: [Gin](https://github.com/gin-gonic/gin), [Chi](https://github.com/go-chi/chi), [Echo](https://github.com/labstack/echo), [Iris](https://github.com/kataras/iris), [Macaron](https://github.com/go-macaron/macaron), [Mux](https://github.com/gorilla/mux), [httprouter](https://github.com/julienschmidt/httprouter))
 depends on struct tag string info that defined in logic object's struct type field.
 
 ### Usage
 
 * Generate a simple template project
 
 ```
% go get github.com/alimy/mir/mirc/v2@latest
% mirc new -d mir-examples
% tree mir-examples
mir-examples
├── Makefile
├── README.md
├── go.mod
├── main.go
└── mirc
    ├── main.go
    └── routes
        ├── site.go
        ├── v1
        │   └── site.go
        └── v2
            └── site.go

% cd mir-examples
% make generate
 ```
 
 * Custom route info just use struct tag. eg:
 
```go
// file: mirc/routes/site.go

package routes

import (
	"github.com/alimy/mir/v2"
	"github.com/alimy/mir/v2/engine"
)

func init() {
	engine.AddEntry(new(Site))
}

// Site mir's struct tag define
type Site struct {
	Chain    mir.Chain `mir:"-"`
	Index    mir.Get   `mir:"/index/"`
	Articles mir.Get   `mir:"/articles/:category/"`
}
```

* Invoke mir's generator to generate interface. eg:

```
% cat mirc/main.go
package main

import (
	"log"

	"github.com/alimy/mir/v2/core"
	"github.com/alimy/mir/v2/engine"

	_ "github.com/alimy/mir/v2/examples/mirc/routes"
	_ "github.com/alimy/mir/v2/examples/mirc/routes/v1"
	_ "github.com/alimy/mir/v2/examples/mirc/routes/v2"
)

//go:generate go run main.go
func main() {
	log.Println("generate code start")
	opts := core.Options{
		core.RunMode(core.InSerialMode),
		core.GeneratorName(core.GeneratorGin),
		core.SinkPath("./gen"),
	}
	if err := engine.Generate(opts); err != nil {
		log.Fatal(err)
	}
	log.Println("generate code finish")
}
```

* Then generate interface from routes info defined above

```go
% make generate
% cat mirc/gen/api/site.go
// Code generated by go-mir. DO NOT EDIT.

package api

import (
	"github.com/gin-gonic/gin"
)

// Site mir's struct tag define
type Site interface {
	Chain() gin.HandlersChain
	Index(c *gin.Context)
	Articles(c *gin.Context)
}

// RegisterSiteServant register site to gin
func RegisterSiteServant(e *gin.Engine, s Site) {
	router := e

	// use chain for router
	middlewares := s.Chain()
	router.Use(middlewares...)

	// register route info to router
	router.Handle("GET", "/index/", s.Index)
	router.Handle("GET", "/articles/:category/", s.Articles)
}
```

* Register interface to router

```go
package main

import (
	"log"

	"github.com/gin-gonic/gin"

	api "github.com/alimy/mir/v2/examples/mirc/gen/api"
	v1 "github.com/alimy/mir/v2/examples/mirc/gen/api/v1"
	v2 "github.com/alimy/mir/v2/examples/mirc/gen/api/v2"
	"github.com/alimy/mir/v2/examples/servants"
)

func main() {
	e := gin.New()

	// register servants to engine
	registerServants(e)

	// start servant service
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}
}

func registerServants(e *gin.Engine) {
	// register default group routes
	api.RegisterSiteServant(e, servants.EmptySiteWithNoGroup{})

	// register routes for group v1
	v1.RegisterSiteServant(e, servants.EmptySiteV1{})

	// register routes for group v2
	v2.RegisterSiteServant(e, servants.EmptySiteV2{})
}
```

* Build application and run

```shell
% make run
```

### Reference Project
 * [examples](examples)  
 Just a simple exmples project for explain how to use [Mir](https://github.com/alimy/mir).
 
 * [mir-covid19](https://github.com/alimy/mir-covid19)  
 COVID-19 Live Updates of Tencent Health is developed to track the live updates of COVID-19, including the global pandemic trends, domestic live updates, and overseas live updates. This project is just a go version of [TH_COVID19_International](https://github.com/Tencent/TH_COVID19_International) for a guide of how to use [Mir](https://github.com/alimy/mir) to develop web application.

