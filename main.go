package main

// curl localhost:8080/admin -u admin:secret
// run in terminal: ADMIN_PASSWORD=secret go run main.go
import (
	"net/http"
	"rest-api-go/admin"
	"rest-api-go/controllers"
)

func main() {
	admin := admin.NewAdminPortal()
	coastersHandlers := controllers.NewCoasterHandlers()
	http.HandleFunc("/coasters", coastersHandlers.Coasters)
	http.HandleFunc("/coaster/", coastersHandlers.GetCoaster)
	http.HandleFunc("/admin", admin.Handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err.Error())
	}
}
