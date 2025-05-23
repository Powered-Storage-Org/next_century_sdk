package ncSdkMock

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	ncMockSample "github.com/Powered-Storage-Org/next_century_sdk/mock/sample"
)

func Run() {
	// login user
	http.HandleFunc("POST /Login", func(w http.ResponseWriter, r *http.Request) {
		// parse email password
		userInfo := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check user
		if userInfo.Email != "test" || userInfo.Password != "test" {
			http.Error(w, "invalid user", http.StatusUnauthorized)
			return
		}

		// return token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{
			Token: "test",
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// get daily reads
	// GET /api/Properties/:id/DailyReads/?date&from&to // &from=%s&to=%s is optional
	http.HandleFunc("GET /api/Properties/{id}/DailyReads/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != "test" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// parse id
		id := r.PathValue("id")
		if id != "x_1234" {
			http.Error(w, "invalid property id", http.StatusBadRequest)
			return
		}

		log.Println("Warning: date parser not implemented")

		// return daily reads
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(ncMockSample.DailyReadsSample)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// get units
	// GET /api/Properties/:id/Units"
	http.HandleFunc("GET /api/Properties/{id}/Units", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != "test" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		id := r.PathValue("id")
		if id != "x_1234" {
			http.Error(w, "invalid property id", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(ncMockSample.UnitsSample)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Addr:              ":1234",
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
