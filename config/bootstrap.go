package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

const projectDirName = "godmin"

func Bootstrap() {
	env := os.Getenv("GODMIN_ENV")
	if "" == env {
		env = "dev"
	}

	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := string(re.Find([]byte(cwd)))

	_ = godotenv.Load(rootPath + "/.env." + env + ".local")

	if env != "test" {
		_ = godotenv.Load(rootPath + "/.env.local")
	}

	_ = godotenv.Load(rootPath + "/.env." + env)

	// The Original .env
	if err := godotenv.Load(rootPath + "/.env"); err != nil {
		log.WithFields(log.Fields{
			"cause": err,
			"cwd":   cwd,
		}).Fatal("Problem loading .env file")

		os.Exit(-1)
	}
}
