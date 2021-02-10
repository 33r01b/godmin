package config

import (
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

const projectDirName = "godmin"

func BootstrapTest() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := string(re.Find([]byte(cwd)))

	_ = godotenv.Load(rootPath + "/.env.test")
}
