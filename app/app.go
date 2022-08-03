package app

import (
	"capi/domain"
	"capi/errs"
	"capi/logger"
	"capi/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type key int

const (
	userInfo key = iota
	test
	test2
	// ...
)

func sanityCheck(){
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
	}

	for _, envKey := range envProps {
		if os.Getenv(envKey) == "" {
			logger.Fatal(fmt.Sprintf("environtment variabel %s not defined. terminating application..", envKey))
		}
	}

	logger.Info("environtment variabel loaded...")
}

func Start() {
	
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("error loading .env file...")
	}

	logger.Info("load environment variabel...")

	sanityCheck()

	dbClient := getClientDB()
	
	// * wiring
	// * setup repository
	customerRepositoryDB := domain.NewCustomerRepositoryDB(dbClient)
	accountRepositoryDB := domain.NewAccountRepositoryDB(dbClient)
	authRepositoryDB := domain.NewAuthRepositoryDB(dbClient)
	
	// * setup handle
	customerService := service.NewCustomerService(customerRepositoryDB)
	accountService := service.NewAccountService(accountRepositoryDB)
	authService := service.NewAuthService(authRepositoryDB)
	
	// * wiring
	ch := CustomerHandler{customerService}
	ah := AccountHandler{accountService}
	authH := AuthHandler{authService}

	// * create ServeRoute
	mux := mux.NewRouter()
	mux.Use(loggingMiddleware)

	authR := mux.PathPrefix("/auth").Subrouter()
	authR.HandleFunc("/login", authH.Login).Methods(http.MethodPost)

	customerR := mux.PathPrefix("/customers").Subrouter()
	customerR.HandleFunc("/{customer_id:[0-9]+}", ch.getCustomerByID).Methods("GET")
	customerR.HandleFunc("/{customer_id:[0-9]+}/accounts/{account_id:[0-9]+}", ah.MakeTransaction).Methods("POST")
	customerR.Use(authMiddleware)

	// New Account
	adminR := mux.PathPrefix("/customers").Subrouter()
	adminR.HandleFunc("", ch.getAllCustomers).Methods("GET")
	adminR.HandleFunc("/{customer_id:[0-9]+}/accounts", ah.NewAccount).Methods(http.MethodPost)
	adminR.Use(authMiddleware)
	adminR.Use(isAdminMiddleware)
	
	mux.Use(authMiddleware)
	// * starting the server
	serverAddr := os.Getenv("SERVER_ADDRESS")
	serverPort := os.Getenv("SERVER_PORT")

	logger.Info(fmt.Sprintf("Start Server on %s:%s...", serverAddr, serverPort))
	http.ListenAndServe(fmt.Sprintf("%s:%s", serverAddr, serverPort), mux)
}

func getClientDB() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPasswd, dbAddr, dbPort, dbName)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Success Connect to Database...")

	return db
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		timer := time.Now()
		next.ServeHTTP(w, r)
		logger.Info(fmt.Sprintf("%v %v %v", r.Method, r.URL, time.Since(timer)))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		authorizationHeader  := r.Header.Get("Authorization")

		// check token validation
		if !strings.Contains(authorizationHeader, "Bearer") {
			logger.Error("Token Invalid")
			errApp := errs.NewForbiddenError("Invalid Token")
			writeResponse(w, errApp.Code, errApp.AsMessage())
			return
		}

		// spilt token -> ambil tokennya buang "Bearer" nya
		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1) // me-return array
		
		token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("signing method invalid")
			}

			return []byte("rahasia"), nil
		})

		if err != nil {
			appErr := errs.NewBadRequestError(err.Error())
			writeResponse(w, appErr.Code, appErr.AsMessage())
			return
		}

		claims, ok := token.Claims.(*domain.AccessTokenClaims)
		// claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			appErr := errs.NewBadRequestError("invalid token")
			writeResponse(w, appErr.Code, appErr.AsMessage())
			return
		}

		if claims.Role == "user" {
			vars := mux.Vars(r)
			customerID := vars["customer_id"]
			accountID := vars["account_id"]

			if claims.CustomerID != customerID {
				appErr := errs.NewForbiddenError("don'thave access to this resource")
				writeResponse(w, appErr.Code, appErr.AsMessage())
				return
			}

			if accountID != "" {
				var isValidAccountID bool
				for _, a := range claims.Accounts {
					if a == accountID {
						isValidAccountID = true
					}
				}
				if !isValidAccountID {
					appErr := errs.NewForbiddenError("don'thave access to this resource")
					writeResponse(w, appErr.Code, appErr.AsMessage())
					return
				}
			}

		}

		ctx := context.WithValue(r.Context(), userInfo, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		})
}

func isAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(userInfo).(*domain.AccessTokenClaims)

		if ctx.Role != "admin" {
			appErr:= errs.NewForbiddenError("Don't have ")
			writeResponse(w, appErr.Code, appErr.AsMessage())
			return
		}
		next.ServeHTTP(w, r)
	})
	
}