package app

import (
	"fmt"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // call init of pq to register driver
	// "github.com/m1schka/go-playground/app/handler"
	"github.com/urfave/negroni"
)

type App struct {
	Router *mux.Router
	// DB     *sql.DB
}

func (a *App) initializeRoutes() {

	apiRoutes := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)
	apiRoutes.HandleFunc("/user", a.getUser)
	apiRoutes.HandleFunc("/upload", a.upload)

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	a.Router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(apiRoutes),
	))
	a.Router.HandleFunc("/healtz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Healthy!")
	})
}

func (a *App) Initialize(user, password, dbname string) {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
