package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/magalhaesgustavo/cloud-run/pkg/cep"
	"github.com/magalhaesgustavo/cloud-run/pkg/weather"
)

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	
	router.Route("/{cep}", func(r chi.Router) {
		
		r.Use(checkCepMiddleware)
		r.Get("/", handleGetTemperatureByCEP)
	})
	
	log.Println("Iniciando o servidor na porta 8080")
	http.ListenAndServe(":8080", router)
}

func checkCepMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cep := chi.URLParam(r, "cep")

		if cep == "" || len(cep) == 0 {
			http.Error(w, "CEP is required", http.StatusBadRequest)
			return
		}

		if len(cep) != 8 {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		for _, d := range cep {
			if d < '0' || d > '9' {
				http.Error(w, "invalid CEP", http.StatusUnprocessableEntity)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func handleGetTemperatureByCEP(w http.ResponseWriter, r *http.Request) {
	cepReq := chi.URLParam(r, "cep")

	log.Println(cepReq)
	address, err := cep.GetAddressFromViaCEP(cepReq)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	log.Println(address.Localidade)

	weatherResponse, err := weather.GetWeather(address.Localidade)
	log.Println(weatherResponse.Current.TempC)
	if err != nil {
		http.Error(w, "can not find weather", http.StatusNotFound)
		return
	}

	temperature := TemperatureResponse{
		TempC: weatherResponse.Current.TempC,
		TempF: weather.CelsiusToFahrenheit(weatherResponse.Current.TempC),
		TempK: weather.CelsiusToKelvin(weatherResponse.Current.TempC),
	}
	json.NewEncoder(w).Encode(temperature)
}