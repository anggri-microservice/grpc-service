package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anggri-microservice/golang-service/internal/db"
	"github.com/joho/godotenv"
	json "github.com/json-iterator/go"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// In this example we use the html template engine
	"github.com/gofiber/template/html"
)

//ResponseTemplate variable
type ResponseTemplate struct {
	Status  int
	Message string
	Data    interface{}
}

const grpcServer string = "192.168.1.222:6003"

func main() {
	// Create a new engine by passing the template folder
	// and template extension using <engine>.New(dir, ext string)
	engine := html.New("./views", ".html")

	// We also support the http.FileSystem interface
	// See examples below to load templates from embedded files
	engine = html.NewFileSystem(http.Dir("./views"), ".html")

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false

	// Debug will print each template that is parsed, good for debugging
	engine.Debug(true) // Optional. Default: false

	// Layout defines the variable name that is used to yield templates within layouts
	engine.Layout("embed") // Optional. Default: "embed"

	// Delims sets the action delimiters to the specified strings
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters

	// After you created your engine, you can pass it to Fiber's Views Engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Cors headers
	app.Use(cors.New())

	// Load go env
	var err error
	if os.Getenv("SRV_DOT_ENV") == "true" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Connect db postgres
	go initServer(app)

	log.Fatal(app.Listen(":5000"))
}

func initServer(app *fiber.App) {
	var err error
	db.DBConn, err = db.PostgreSQL.ConnectSqlx()

	if err != nil {
		log.Println(err)
	}

	// To render a template, you can call the ctx.Render function
	// Render(tmpl string, values interface{}, layout ...string)
	app.Post("/users", Users)
}

//Users controller function
func Users(c *fiber.Ctx) error {
	//get request to interface
	var data map[string]interface{}
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		log.Println(err)
		messageError := err.Error()
		errTemplate := &ResponseTemplate{
			Status:  fiber.StatusBadRequest,
			Message: messageError,
			Data:    nil,
		}

		return c.Status(fiber.StatusBadRequest).JSON(errTemplate)
	}

	SQL := `select id, name from users`
	type Users struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	var users []Users
	err = db.DBConn.Select(&users, SQL)
	if err != nil {
		log.Println(err)
		messageError := err.Error()
		errTemplate := &ResponseTemplate{
			Status:  fiber.StatusBadRequest,
			Message: messageError,
			Data:    nil,
		}

		return c.Status(fiber.StatusBadRequest).JSON(errTemplate)
	}

	log.Println(users)

	// var response map[string]interface{}
	// err = json.Unmarshal([]byte(`{"test":"haha"}`), &response)
	// log.Println(response)
	// if err != nil {
	// 	log.Println(err)
	// 	messageError := err.Error()
	// 	errTemplate := &ResponseTemplate{
	// 		Status:  fiber.StatusBadRequest,
	// 		Message: messageError,
	// 		Data:    nil,
	// 	}

	// 	return c.Status(fiber.StatusBadRequest).JSON(errTemplate)
	// }

	successTemplate := &ResponseTemplate{
		Status:  fiber.StatusOK,
		Message: "Sukses menampilkan data verifikasi",
		// Data:    response,
	}
	return c.JSON(successTemplate)
}
