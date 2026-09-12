package main

import (
	"cookapp/internal/handler"
	"cookapp/internal/repository"
	"cookapp/internal/service"
	"cookapp/internal/storage"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/cors"

	myMiddleware "cookapp/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	//Recupération .env
	godotenv.Load() // sans log.Fatal si erreur

	db_url := os.Getenv("DATABASE_URL")

	//Initialisation de la db
	db, err := sql.Open("pgx", db_url)
	if err != nil {
		fmt.Println("Erreur DB:", err)
		return
	}
	repository.InitDB(db)
	repository.SeedDB(db)

	r2Storage, err := storage.NewR2Storage()
	if err != nil {
    	log.Fatal("Erreur initialisation R2:", err)
	}

	//Initialisation routeur Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Initialisation des couches
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, r2Storage)
	friendshipRepo := repository.NewFriendshipRepository(db)
	friendshipService := service.NewFriendshipService(friendshipRepo)
	friendshipHandler := handler.NewFriendshipHandler(friendshipService)
	recipeRepo := repository.NewRecipeRepository(db)
	likesRepo := repository.NewLikesRepository(db)
	favoriteRepo := repository.NewFavoriteRepository(db)
	recipeService := service.NewRecipeService(recipeRepo, likesRepo, favoriteRepo)
	recipeHandler := handler.NewRecipeHandler(recipeService, r2Storage)
	likesService := service.NewLikesService(likesRepo)
	likesHandler := handler.NewLikesHandler(likesService)
	favoriteService := service.NewFavoriteService(favoriteRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)
	ingredientRepo := repository.NewIngredientRepository(db)
	ingredientService := service.NewIngredientService(ingredientRepo)
	ingredientHandler := handler.NewIngredientHandler(ingredientService)

	//Routes Chi
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
		r.Get("/search", userHandler.SearchUsers)
		r.Group(func(r chi.Router) {
			r.Use(myMiddleware.WithAuth)
			r.Get("/me", userHandler.GetMe)
			r.Put("/me", userHandler.UpdateMe)
			r.Delete("/me", userHandler.DeleteMe)
			r.Put("/me/avatar", userHandler.UpdateAvatar)
			r.Get("/me/recipes", recipeHandler.GetAllRecipesFromUser)
		})
	})
	r.Route("/recipes", func(r chi.Router) {
		r.Get("/", recipeHandler.GetAllRecipesShuffled)
		r.Get("/search", recipeHandler.SearchRecipes)
		r.Get("/{id}", recipeHandler.GetRecipeById)
		r.Group(func(r chi.Router) {
			r.Use(myMiddleware.WithAuth)
			r.Post("/", recipeHandler.CreateRecipe)
			r.Put("/{id}", recipeHandler.UpdateRecipe)
			r.Delete("/{id}", recipeHandler.DeleteRecipe)
			r.Post("/{id}/like", likesHandler.AddLike)
			r.Delete("/{id}/like", likesHandler.DeleteLike)
			r.Get("/{id}/likes", likesHandler.GetLikeCount)
			r.Post("/{id}/favorite", favoriteHandler.AddFavorite)
			r.Delete("/{id}/favorite", favoriteHandler.DeleteFavorite)
			r.Get("/{id}/favorites", favoriteHandler.GetFavoriteCount)
			r.Get("/{id}/isliked", likesHandler.IsLiked)
			r.Get("/{id}/isfavorited", favoriteHandler.IsFavorited)
		})
	})
	r.Get("/ingredients", ingredientHandler.GetAllIngredients)
	r.Group(func(r chi.Router) {
		r.Use(myMiddleware.WithAuth)
		r.Post("/friends/request", friendshipHandler.SendFriendRequest)
		r.Put("/friends/{id}/accept", friendshipHandler.AcceptFriendRequest)
		r.Put("/friends/{id}/decline", friendshipHandler.DeclineFriendRequest)
		r.Delete("/friends/{id}", friendshipHandler.DeleteFriend)
		r.Get("/friends", friendshipHandler.GetFriend)
		r.Get("/friends/requests", friendshipHandler.GetPendingFriendRequests)
	})

	r.Get("/admin/fix-images", func(w http.ResponseWriter, r *http.Request) {
    _, err := db.Exec("UPDATE recipes SET image_path = 'https://pub-c99b03c3b3594edf8d525f1ea103ef88.r2.dev/default-recipe.svg' WHERE image_path = '' OR image_path IS NULL")
    if err != nil {
        w.Write([]byte("Erreur: " + err.Error()))
        return
    }
    w.Write([]byte("OK"))
})

	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
