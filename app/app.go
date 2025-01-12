package app

import (
	"net/http"
	"personal_blog/config"
	articlesC "personal_blog/controllers/articles"
	usersC "personal_blog/controllers/users"
	"personal_blog/routes"
	articlesS "personal_blog/services/articles"
	usersS "personal_blog/services/users"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type App struct {
	gormDB       *gorm.DB
	usersService usersS.UsersService
	articlesS    articlesS.ArticleService
	usersCtrl    *usersC.Controller
	aritclesCtrl *articlesC.Controller

	mux   *mux.Router
	httpH http.Handler
}

func NewApp() (*App, error) {
	app := &App{}
	db := config.InitDB()
	app.gormDB = db

	config.AutoMigrateDB(db)

	app.usersService = usersS.NewService(db)
	app.articlesS = articlesS.NewArticleGetter()

	app.aritclesCtrl = articlesC.NewController(app.articlesS)
	app.usersCtrl = usersC.NewController(app.usersService)

	app.mux = routes.SetupRouter2(app.usersCtrl, app.aritclesCtrl)
	app.httpH = routes.AcceptMiddlewares(app.mux)
	return app, nil
}

func (a *App) Run() error {
	return http.ListenAndServe(":8080", a.httpH)
}
