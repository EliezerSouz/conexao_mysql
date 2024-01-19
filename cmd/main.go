package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/EliezerSouz/conexao_mysql/modulos"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var shutdownSignal = make(chan os.Signal, 1)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("Erro ao obter caminho do executável: ", err)
	}

	exerDir := filepath.Dir(exePath)

	println(exerDir)

	//Carregar variaveis do ambiente

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Erro ao carregar variaveis de ambiente: ", err)
	}

	router := gin.Default()

	router.GET("/api/baixar-xmls", modulos.BaixarXmlsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Iniciando a aplicação na porta %s...\n", port)

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}

	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	<-shutdownSignal

	os.Exit(0)
}
