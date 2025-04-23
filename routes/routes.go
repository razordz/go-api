package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"go-api/controllers"

	_ "go-api/docs"
)

func Setup() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API em Go está no ar!"))
	})

	// ✅ Redireciona /doc/api para /doc/api/index.html
	router.HandleFunc("/doc/api", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/doc/api/index.html", http.StatusMovedPermanently)
	})

	// ✅ Serve Swagger em /doc/api/
	router.PathPrefix("/doc/api/").Handler(httpSwagger.WrapHandler)

	return router
}
