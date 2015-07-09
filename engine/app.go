package webengine



func init() {
	// bootstrap


}

type App struct {
	Name string
}

func (app *App) SetName(name string) {
	app.Name = name
	return
}

func (app *App) GetName() string {
	return app.Name
}
