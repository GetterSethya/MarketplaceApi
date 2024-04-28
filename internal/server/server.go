package server

import (
	"log"
	"net/http"

	"github.com/GetterSethya/golangApiMarketplace/internal/datastore"
	"github.com/GetterSethya/golangApiMarketplace/internal/helper"
	"github.com/GetterSethya/golangApiMarketplace/internal/services"
	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	store      datastore.Store
}

func NewServer(addr string, store datastore.Store) *Server {

	return &Server{
		listenAddr: addr,
		store:      store,
	}
}

func (s *Server) Run() {

	// init router
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/v1").Subrouter()
	subrouter.HandleFunc("/healthCheck", func(w http.ResponseWriter, r *http.Request) {

		helper.WriteJson(w, http.StatusOK, map[string]interface{}{
			"hello": "mom 💖",
		})
	})

	// register service disini
	userService := services.NewUserService(s.store)
	userService.RegisterRoutes(subrouter)

	// register product service disini
	productService := services.NewProductService(s.store)
	productService.RegisterRoutes(subrouter)

	// register bankAccount service disini
	bankAccountService := services.NewBankAccountService(s.store)
	bankAccountService.RegisterRoutes(subrouter)

	// register transaction service disini
	transactionService := services.NewTransactionService(s.store)
	transactionService.RegisterRoutes(subrouter)

	log.Println("Server is running on:", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, subrouter))
}
