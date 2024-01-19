package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"conexao_mysql/modulos"

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

	router.GET("/api/baixar-xmls", BaixarXmlsHandler)

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
func BaixarXmlsHandler(c *gin.Context) {
	dataInicial := c.Query("dataInicial")
	dataFinal := c.Query("dataFinal")
	emissorP := c.Query("emissorP")
	emissorT := c.Query("emissorT")
	_nfe := c.Query("_nfe")
	_nfce := c.Query("_nfce")
	// Execute a lógica para baixar XMLs
	err := modulos.BaixarXmls(dataInicial, dataFinal, emissorP, emissorT, _nfe, _nfce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responda com uma mensagem de sucesso
	c.JSON(http.StatusOK, gin.H{"message": "BaixarXmlsHandler executado com sucesso"})
}

