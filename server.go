package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/terminaldotshop/terminal-sdk-go"
	"github.com/terminaldotshop/terminal-sdk-go/option"
	"log"
	"net/http"
	"text/template"
)

type ProductView struct {
	Name        string
	Description string
	Color       string
	VariantID   string
}

type VariantView struct {
	Name           string
	PriceFormatted string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/products", getProducts)
	http.HandleFunc("/ws", ws)

	client := terminal.NewClient(
		option.WithBearerToken("trm_test_3532f9f1592e704eadbc"), // defaults to os.LookupEnv("TERMINAL_BEARER_TOKEN")
		option.WithEnvironmentDev(),
	)

	response, err := client.Address.List(context.TODO())

	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", response.Data)
	log.Println("✅ Server listening on http://localhost:8080/products")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read Errors", err)
			break
		}
		fmt.Printf("Received", msg)

		err = conn.WriteMessage(websocket.TextMessage, []byte("Server got: "+string(msg)))
		if err != nil {
			fmt.Println("Write error", err)
			break
		}
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))

	client := terminal.NewClient(
		option.WithBearerToken("trm_test_3532f9f1592e704eadbc"), // defaults to os.LookupEnv("TERMINAL_BEARER_TOKEN")
		option.WithEnvironmentDev(),                             // defaults to option.WithEnvironmentProduction()
	)

	products, err := client.Product.List(context.TODO())
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	views := []ProductView{}
	for _, p := range products.Data {
		if !p.Tags.MarketNa || p.Subscription == "required" {
			continue
		}

		variantId := ""
		if len(p.Variants) > 0 {
			variantId = p.Variants[0].ID
		}
		views = append(views, ProductView{
			Name:        p.Name,
			Description: p.Description,
			Color:       p.Tags.Color,
			VariantID:   variantId,
		})
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, views); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
