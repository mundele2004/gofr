package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"ldap-auth-demo/internal/auth/jwt"
	"ldap-auth-demo/internal/auth/ldapauth"
)

func main() {
	jwtIssuer := jwt.NewIssuer("super-secret-key", 1*time.Hour, true)

	cfg := ldapauth.Config{
		Addr:         "localhost:389",
		BaseDN:       "dc=example,dc=com",
		BindUserDN:   "cn=admin,dc=example,dc=com",
		BindPassword: "admin",
		TokenIssuer:  jwtIssuer,
		UserAttributes: []string{
			"dn", "uid", "mail", "cn", "ou", "memberOf", "telephoneNumber",
		},
		UsernameAttr:   "uid",
		EmailAttr:      "mail",
		FullNameAttr:   "cn",
		DepartmentAttr: "ou",
	}

	authenticator := ldapauth.New(&cfg)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, authenticator)
	})

	http.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		userInfoHandler(w, r, authenticator)
	})

	fmt.Println("Listening on :8080")

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func loginHandler(w http.ResponseWriter, r *http.Request, authenticator *ldapauth.Authenticator) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := authenticator.Authenticate(creds.Username, creds.Password)
	if err != nil {
		statusCode := getHTTPStatusForError(err)

		if statusCode == http.StatusInternalServerError {
			log.Printf("Server configuration error: %v", err)
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		log.Println("encode error:", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func userInfoHandler(w http.ResponseWriter, r *http.Request, authenticator *ldapauth.Authenticator) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username parameter required", http.StatusBadRequest)
		return
	}

	user, err := authenticator.GetUserInfo(username)
	if err != nil {
		statusCode := getHTTPStatusForError(err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("encode error:", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func getHTTPStatusForError(err error) int {
	switch {
	case errors.Is(err, ldapauth.ErrUserNotFound):
		return http.StatusUnauthorized
	case errors.Is(err, ldapauth.ErrMultipleUsersFound):
		return http.StatusInternalServerError
	default:
		return http.StatusUnauthorized
	}
}