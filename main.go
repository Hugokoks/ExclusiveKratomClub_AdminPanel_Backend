package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	///////Common Library EKC packages
	ekc_cors "github.com/Hugokoks/kratomclub-go-common/appcors"
	ekc_db "github.com/Hugokoks/kratomclub-go-common/db"
	ekc_mid "github.com/Hugokoks/kratomclub-go-common/middlewares"
	ekc_services "github.com/Hugokoks/kratomclub-go-common/services"
	ekc_settings "github.com/Hugokoks/kratomclub-go-common/settings"
	ekc_utils "github.com/Hugokoks/kratomclub-go-common/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main(){

	isProd := os.Getenv("ENV") == "production"

	// .env nacist jen mimo produkci
	if !isProd {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: .env file not loaded")
		}
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// === DB ===
	ekc_db.InitPool("DATABASE_ADMIN_URL")
	defer ekc_db.Pool.Close()

	// === App settings ===
	settingsRows, err := ekc_db.SelectAppSettings()
	if err != nil {
		log.Fatal(err)
	}
	if err := ekc_settings.Load(context.Background(), settingsRows); err != nil {
		log.Fatal(err)
	}

	// === Gin Init ===
	r := gin.New()
	r.Use(gin.Logger(),gin.Recovery())
	r.SetTrustedProxies(nil)



	// ---- Security hlavičky ----
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		// HSTS zapni jen pokud jsi za HTTPS reverse proxy:
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Next()
	})


	// ---- Limit velikosti request body ----
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2<<20) // 2 MB
		c.Next()
	})

	// === Rate limiting per (anon) session ===
	sessLimiter := ekc_mid.NewSessionRateLimiter(1, 30, 100_000)
	go sessLimiter.CleanupStaleEntries(30 * time.Minute)

	// ===CORS Creating ===
	corsConfig := ekc_cors.CORSConfig{
		Env:        os.Getenv("ENV"), // "production" / "development"
		DevClient:  ekc_utils.LoadEnvVariable("DEVELOPMENT_CLIENT_URL"),
		AdminUI:    ekc_utils.LoadEnvVariable("ADMIN_CLIENT_URL"),
		// optional:
		// AllowHeaders: []string{...}, AllowMethods: []string{...}, MaxAge: time.Hour,
		LogDecisions: true, // dej true při ladění CORS
	}
	adminCORS := ekc_cors.NewAdmin(corsConfig)


	// ==================ADMIN API ==================

	admin := r.Group("/admin")
	admin.Use(adminCORS)
	admin.Use(ekc_mid.AnonSession())
	admin.Use(ekc_mid.SessionRateLimitMiddleware(sessLimiter))
	//Preflight
	admin.OPTIONS("/*path", func(c *gin.Context) { 
		c.Status(http.StatusNoContent)
	})


	ekc_services.ServerStart(r)

}