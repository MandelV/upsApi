package models

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Configuration contient tous les paramètres de l'application.
// Chaque champ est rempli automatiquement depuis les variables d'environnement
// grâce aux tags `env` (lib caarlos0/env). Un fichier .env est chargé en priorité si présent.
type Configuration struct {
	VERSION     string `env:"VERSION" envDefault:"0.0.1"`
	ENV         string `env:"ENV" envDefault:"dev"`            // "dev" active le mode debug de Gin
	LOG_LEVEL   int    `env:"LOG_LEVEL" envDefault:"0"`        // Niveau zerolog : 0=debug, 1=info, 2=warn, 3=error, 4=fatal
	LISTEN_PORT int    `env:"LISTEN_PORT" envDefault:"9695"`   // Port TCP du serveur HTTP
	UPS_USER    string `env:"UPS_USER" envDefault:"eaton"`     // Nom de l'UPS dans NUT (partie <ups> de "upsc <ups>@<host>")
	UPS_HOST    string `env:"UPS_HOST" envDefault:"localhost"` // Hôte du serveur NUT (upsd)
}

// NewConfiguration charge le fichier .env (s'il existe), puis parse les variables
// d'environnement dans la struct Configuration. Panique si le parsing échoue.
func NewConfiguration() *Configuration {

	godotenv.Load(".env")

	conf := &Configuration{}

	if err := env.Parse(conf); err != nil {
		panic(err)
	}
	return conf
}
