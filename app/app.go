package app

import (
	"log"
	"net/http"

	usersS "github.com/Sekarfo/P_blog/services/users"

	"github.com/Sekarfo/P_blog/routes"

	"github.com/Sekarfo/P_blog/config"

	articlesC "github.com/Sekarfo/P_blog/controllers/articles"

	usersC "github.com/Sekarfo/P_blog/controllers/users"

	blogsS "github.com/Sekarfo/P_blog/services/blogs"

	blogsC "github.com/Sekarfo/P_blog/controllers/blogs"
	articlesS "github.com/Sekarfo/P_blog/services/articles"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	gormDB       *gorm.DB
	usersService usersS.UsersService
	articlesS    articlesS.ArticleService
	blogService  blogsS.BlogService
	usersCtrl    *usersC.Controller
	aritclesCtrl *articlesC.Controller
	blogCtrl     *blogsC.Controller

	mux   *mux.Router
	httpH http.Handler
}

func NewApp() (*App, error) {
	app := &App{}
	err := app.initEnvs()
	if err != nil {
		return nil, err
	}
	db := config.InitDB()
	app.gormDB = db

	config.AutoMigrateDB(db)

	app.usersService = usersS.NewService(db)
	app.articlesS = articlesS.NewArticleGetter()
	app.blogService = blogsS.NewService(db)

	app.aritclesCtrl = articlesC.NewController(app.articlesS)
	app.usersCtrl = usersC.NewController(app.usersService)
	app.blogCtrl = blogsC.NewController(app.blogService)

	app.mux = routes.SetupRouter(app.usersCtrl, app.aritclesCtrl, app.blogCtrl)
	app.httpH = routes.AcceptMiddlewares(app.mux)
	return app, nil
}

func (a *App) Run() error {
	return http.ListenAndServe(":8080", a.httpH)
}

func (a *App) initEnvs() error {
	err := godotenv.Load("c:/Users/trudk/GolandProjects/P_blog/.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return err
	}
	log.Println(".env file loaded successfully")
	return nil
}
