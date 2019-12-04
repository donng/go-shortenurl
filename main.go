package main

func main() {
	app := App{}
	app.Init(getEnv())
	app.Run(":8080")
}
