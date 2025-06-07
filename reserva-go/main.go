package main

import (
	"context"
	"log"
	"net/http"
	h "reserva-go/handlers"
	"reserva-go/services"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const httpServerAddress = ":3000"

// --- Main Function ---
func main() {
	services.InitMongoDB()
	services.InitRabbitMQ()  // This also starts the pagamentoAprovado consumer in a goroutine

	defer services.MongoClient.Disconnect(context.Background())
	defer services.RabbitMQConnection.Close()
	if services.RabbitMQChannelGlobal != nil {
		defer services.RabbitMQChannelGlobal.Close()
	}


	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/", h.HelloHandler).Methods(http.MethodGet)
	r.HandleFunc("/minhas-reservas", h.MinhasReservasHandler).Methods(http.MethodGet)
	r.HandleFunc("/destinos", h.CriarDestinoHandler).Methods(http.MethodPost)
	r.HandleFunc("/destinos/buscar", h.BuscarDestinosHandler).Methods(http.MethodPost)
	r.HandleFunc("/destinos-por-categoria", h.DestinosPorCategoriaHandler).Methods(http.MethodGet)
	r.HandleFunc("/destinos/reservar", h.ReservarDestinoHandler).Methods(http.MethodPost)
	r.HandleFunc("/destinos/cancelar", h.CancelarViagemHandler).Methods(http.MethodPost)
    r.HandleFunc("/sse", h.SSEHandler).Methods(http.MethodGet)
    r.HandleFunc("/sse/send", h.SSESendHandler).Methods(http.MethodPost)

	// CORS middleware
	// AllowedOrigins can be more specific in production
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // Or specify your frontend domain(s)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(), // Important if you use cookies/auth headers from frontend
	)

	log.Printf("🦊 Go server is running at %s", httpServerAddress)
	if err := http.ListenAndServe(httpServerAddress, corsHandler(r)); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}