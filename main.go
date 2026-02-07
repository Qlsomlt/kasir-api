package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"os"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Kategori struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	Description string `json:"description"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie", Harga: 1500, Stok: 10},
	{ID: 2, Nama: "KitKat", Harga: 8000, Stok: 28},
	{ID: 3, Nama: "LifeBoy", Harga: 5000, Stok: 41},
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

var kategori = []Kategori{
	{ID: 1, Nama: "Makanan", Description: "Kategori makanan ringan dan berat"},
	{ID: 2, Nama: "Minuman", Description: "Kategori minuman dingin dan hangat"},
	{ID: 3, Nama: "Sabun Mandi", Description: "Kategori Peralatan Sabun Mandi dan Cuci Muka"},
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		getProdukByID(w, r)
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Product error", http.StatusNotFound)
}

func updateProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updatedProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updatedProduk)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			produk[i] = updatedProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk[i])
			return
		}
	}
	http.Error(w, "Produk tidak ada", http.StatusNotFound)
}

func deleteProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {

			produk = append(produk[:1], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			return
		}

		http.Error(w, "Produk tidak ada", http.StatusNotFound)
	}
}

func getKategoriID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		getKategoriID(w, r)
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	for _, k := range kategori {
		if k.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(k)
			return
		}
	}
	http.Error(w, "Product error", http.StatusNotFound)
}

func updateKategoriByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var updatedKategori Kategori
	err = json.NewDecoder(r.Body).Decode(&updatedKategori)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range kategori {
		if kategori[i].ID == id {
			kategori[i] = updatedKategori

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategori[i])
			return
		}
	}
	http.Error(w, "Kategori tidak ada", http.StatusNotFound)
}

func deleteKategoriByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	for i, k := range kategori {
		if k.ID == id {

			kategori = append(kategori[:1], kategori[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			return
		}

		http.Error(w, "Kategori tidak ada", http.StatusNotFound)
	}
}

func main() {

	// 1. Tell Viper
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Note: .env file not found, using system environment variables")
	}

	// Check "DATABASE_URL" first, then fallback to "DB_CONN"
	dbConn := viper.GetString("DATABASE_URL")
	if dbConn == "" {
		dbConn = viper.GetString("DB_CONN")
	}
	if dbConn == "" {
		log.Fatal("Error: Database connection string is empty. Check your .env file or Environment Variables.")
	}

	// Connect
	ctx := context.Background()
	fmt.Printf("Attempting to connect to: %s\n", dbConn) // Debug line

	conn, err := pgx.Connect(ctx, dbConn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)
	fmt.Println("Successfully connected to the database!")

	// GET localhost:8080/api/produk/
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProdukByID(w, r)
		} else if r.Method == "DELETE" {
			deleteProdukByID(w, r)
		}
	})
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)

		} else if r.Method == "POST" {
			// baca data dari request
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// masukkin data ke dalam variable produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	// GET localhost:8080/api/kategori/
	http.HandleFunc("/api/kategori/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getKategoriID(w, r)
		} else if r.Method == "PUT" {
			updateKategoriByID(w, r)
		} else if r.Method == "DELETE" {
			deleteKategoriByID(w, r)
		}
	})

	// Category Endpoints
	http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategori)
		} else if r.Method == "POST" {
			var kategoriBaru Kategori
			err := json.NewDecoder(r.Body).Decode(&kategoriBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			kategoriBaru.ID = len(kategori) + 1
			kategori = append(kategori, kategoriBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(kategoriBaru)
		}
	})

	// 3. Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := conn.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "ERROR", "message": "Database unreachable"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK", "message": "API and DB running"})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Kasir API!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
