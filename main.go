package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	// Parse all templates from the "./views" folder
	// where extension is ".html" and parse them
	// using the standard `html/template` package.
	tmpl := iris.HTML("./views", ".html")
	// Set custom delimeters.
	tmpl.Delims("{{", "}}")
	// Enable re-build on local template files changes.
	tmpl.Reload(true)

	// Default template funcs are:
	//
	// - {{ urlpath "myNamedRoute" "pathParameter_ifNeeded" }}
	// - {{ render "header.html" }}
	// and partial relative path to current page:
	// - {{ render_r "header.html" }}
	// - {{ yield }}
	// - {{ current }}
	// Register a custom template func:
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})

	// Register the view engine to the views,
	// this will load the templates.
	app.RegisterView(tmpl)

	// Method:    GET
	// Resource:  http://localhost:8080
	app.Get("/", func(ctx iris.Context) {
		// Bind: {{.message}} with "Hello world!"
		ctx.ViewData("message", "Hello world!")
		// Render template file: ./views/hi.html
		ctx.View("hi.html")
	})

	app.Listen(":8080")
}
