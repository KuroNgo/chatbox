package infrastructor

import (
	"chatbox/bootstrap"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Env     *bootstrap.Database
	MongoDB *mongo.Client
}

func App() *Application {
	app := &Application{}
	app.Env = bootstrap.NewEnv()
	app.MongoDB = NewMongoDatabase(app.Env)
	return app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.MongoDB)
}
