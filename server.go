package main

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/jroyal/drafthouse-seat-finder/drafthouse"
	"github.com/labstack/echo"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
)

// Command Line Options
var (
	port    = flag.String("port", "8080", "Port to run the server on")
	local   = flag.Bool("local", false, "Run the server only on localhost")
	urlBase = flag.String("urlBase", "", "For reverse proxy support")
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	flag.Parse()
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	base := strings.Trim(*urlBase, "/")
	index := "/" + base
	if base != "" {
		base = "/" + base
	}

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("baseUrl", base)
			c.Set("indexUrl", index)
			return h(c)
		}
	})

	e.Static(fmt.Sprintf("%s", index), "public")

	// These are the two main routes used by the UI
	e.GET(fmt.Sprintf("%s", index), drafthouse.HandleIndex)
	e.POST(fmt.Sprintf("%s/seats", base), drafthouse.HandleSeats)

	// These are fun convienience routes that I used for testing. Eventually I might clean these out
	e.GET(fmt.Sprintf("%s/films", base), drafthouse.HandleGetFilms)
	e.GET(fmt.Sprintf("%s/movies/:film-slug", base), drafthouse.HandleGetSingleMovie)

	if *local {
		e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%s", *port)))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
	}

}
