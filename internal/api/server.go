package api

import (
	"log"
	"net/http"
	"time"

	"github.com/AndreLiberato/bitbank/internal/service"
)

// NewRouter monta o roteador HTTP com todos os endpoints do banco, delegando
// para o AccountHandler. Usa o roteamento por método/padrão do net/http (Go 1.22+).
func NewRouter(svc *service.AccountService) http.Handler {
	handler := NewAccountHandler(svc)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /banco/conta/{$}", handler.Create)
	mux.HandleFunc("PUT /banco/conta/transferencia", handler.Transfer)
	mux.HandleFunc("PUT /banco/conta/rendimento", handler.RenderJuros)
	mux.HandleFunc("GET /banco/conta/{id}", handler.Get)
	mux.HandleFunc("GET /banco/conta/{id}/saldo", handler.Balance)
	mux.HandleFunc("PUT /banco/conta/{id}/credito", handler.Credit)
	mux.HandleFunc("PUT /banco/conta/{id}/debito", handler.Debit)

	return logging(mux)
}

// Serve inicia o servidor HTTP da API REST no endereço informado.
func Serve(addr string, svc *service.AccountService) error {
	server := &http.Server{
		Addr:              addr,
		Handler:           NewRouter(svc),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("BitBank REST API ouvindo em http://localhost%s", addr)
	return server.ListenAndServe()
}

// logging é um middleware simples de log de acesso.
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
