package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	userauthpb "gitlab.com/emporia-digital/emporia-2.0/user-service/pb/auth"
	"google.golang.org/grpc"

	// In this example we use the html template engine
	"github.com/gofiber/template/html"
)

//ResponseTemplate variable
type ResponseTemplate struct {
	Status  int
	Message string
	Data    interface{}
}

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

	// AddFunc adds a function to the template's global function map.
	engine.AddFunc("greet", func(name string) string {
		return "Hello, " + name + "!"
	})

	// After you created your engine, you can pass it to Fiber's Views Engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Cors headers
	app.Use(cors.New())

	// To render a template, you can call the ctx.Render function
	// Render(tmpl string, values interface{}, layout ...string)
	app.Post("/phone-verification", PhoneVerification)

	log.Fatal(app.Listen(":5000"))
}

//PhoneVerification controller function
func PhoneVerification(c *fiber.Ctx) error {
	conn, err := grpc.Dial("192.168.1.222:6003", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	d := userauthpb.NewAuthClient(conn)
	//add grpc deadline form product wating time
	msEnv := int(60 * 1000)
	newGRPCDeadline := time.Duration(msEnv) * time.Millisecond
	_, cancel := context.WithTimeout(context.Background(), newGRPCDeadline)
	defer cancel()

	var data map[string]interface{}
	err = json.Unmarshal(c.Body(), &data)

	resAuth, err := d.MitraPhoneVerification(context.Background(), &userauthpb.PhoneVerificationReq{
		Phone:       data["phone"].(string),
		Application: data["type"].(string),
	})

	if err != nil {
		log.Println(err)
		messageError := strings.ReplaceAll(err.Error(), "rpc error: code = Unknown desc = ", "")
		errTemplate := &ResponseTemplate{
			Status:  fiber.StatusBadRequest,
			Message: messageError,
			Data:    nil,
		}

		return c.Status(fiber.StatusBadRequest).JSON(errTemplate)
	}

	var response map[string]interface{}
	err = json.Unmarshal(resAuth.Data, &response)
	successTemplate := &ResponseTemplate{
		Status:  fiber.StatusOK,
		Message: "Sukses menampilkan data verifikasi",
		Data:    response,
	}
	return c.JSON(successTemplate)
}
