package main

func main() {
	config := newConfig()
	app := &App{
		config: config,
	}
	app.registerRoutes()
	app.startServer()
}
