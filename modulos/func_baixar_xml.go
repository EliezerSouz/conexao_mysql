package modulos

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/EliezerSouz/conexao_mysql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func BaixarXmlsHandler(c * gin.Context){
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

func BaixarXmls(dataInicial, dataFinal, emissorP, emissorT, _nfe, _nfce string) error {

	fmt.Println("creating zip archive...")
	//archive, err := os.Create("archive.zip")
	archive, err := os.Create("Xmls.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)

	db, err := db.GetDB()
	if err != nil {
		log.Fatal("database connection error, see .env configuration")
	}

	defer db.Close()

	//query := "select notafiscal, serie, modelo, chave_nota , xml_final, nota_cabecalho.fs_fase from nota_xml join nota_cabecalho ON nota_cabecalho.id = nota_xml.nota_cabecalho_id where xml_final is not null and data_nota between '2023-11-01' and '2023-11-30' and modelo in ('55','65') and emissor = 'P' order by modelo;"
	query := fmt.Sprintf("select xml_final, nota_cabecalho.fs_fase, nota_xml.emissor from nota_xml join nota_cabecalho ON nota_cabecalho.id = nota_xml.nota_cabecalho_id where xml_final is not null and data_nota between '%s' and '%s' and emissor in('%s', '%s') and modelo in ('%s','%s') order by nota_xml.notafiscal;", dataInicial, dataFinal, emissorP, emissorT, _nfe, _nfce)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Erro ao executar a consulta: %v", err)
	}
	defer rows.Close()

	var filename, modNfe, vigenciaXml, chaveAcesso, cnpjTerceiro string

	LimparXmls()

	for rows.Next() {

		diretorio := "fs-xmls"
		var xmls NotaXml

		err := rows.Scan(&xmls.XmlFinalNF, &xmls.StatusNF, &xmls.TipoEmissor)
		if err != nil {
			log.Printf("Erro ao escanear linha: %v", err)
			continue
		}
		if xmls.XmlFinalNF != "" {
			if strings.Contains(xmls.XmlFinalNF, "<nfeProc") {

				if xmls.StatusNF == "19" {
					xmlData := []byte(xmls.XmlFinalNF)
					if strings.Contains(xmls.XmlFinalNF, "<chNFe>") {
						var readXml NFe3
						errors := xml.Unmarshal(xmlData, &readXml)
						if err != nil {
							log.Fatal(errors)
						}

						modNfe = readXml.NFe.InfNFe.Ide.Mod
						chaveAcesso = readXml.ProtNFe.InfProt.ChNFe
						vigenciaXml = "20" + chaveAcesso[2:6]
						filename = chaveAcesso + ".xml"
					} else {

						var readXml NFe2

						errors := xml.Unmarshal(xmlData, &readXml)
						if err != nil {
							log.Fatal(errors)
						}

						modNfe = readXml.InfNFe.Ide.Mod
						chaveAcesso = readXml.InfNFe.ID
						vigenciaXml = "20" + chaveAcesso[5:9]
						filename = chaveAcesso + ".xml"
					}

				} else if xmls.StatusNF != "19" {

					xmlData := []byte(xmls.XmlFinalNF)
					var readXml NFe

					errors := xml.Unmarshal(xmlData, &readXml)
					if err != nil {
						log.Fatal(errors)
					}

					modNfe = readXml.NFe.InfNFe.Ide.Mod
					chaveAcesso = readXml.ProtNFe.InfProt.ChNFe
					cnpjTerceiro = readXml.NFe.InfNFe.Emit.CNPJ
					vigenciaXml = "20" + chaveAcesso[2:6]
					filename = chaveAcesso + ".xml"

				}
				var tipoEmissor = ""
				if xmls.TipoEmissor == "P" {
					tipoEmissor = "Própria"
					diretorio += "/" + vigenciaXml + "/" + tipoEmissor

					if modNfe == "55" {
						diretorio += "/NFe"
					} else if modNfe == "65" {
						diretorio += "/NFCe"
					}
					if xmls.StatusNF == "07" {
						diretorio += "/Canceladas"
					} else if xmls.StatusNF == "19" {
						diretorio += "/Denegadas"
					} else {
						diretorio += "/Autorizadas"
					}
				} else if xmls.TipoEmissor == "T" {
					tipoEmissor = "Terceiros"
					diretorio += "/" + vigenciaXml + "/" + tipoEmissor + "/CNPJ-" + cnpjTerceiro

				}

				_, erro1 := os.Stat(diretorio)
				if os.IsNotExist(erro1) {
					errDir := os.MkdirAll(diretorio, 0755)
					if errDir != nil {
						log.Fatalf("Erro ao criar diretório: %v", errDir)
					}
				}

				filepath := fmt.Sprintf("%s/%s", diretorio, filename)
				fmt.Println("Copiando arquivo:", filename)
				err = os.WriteFile(filepath, []byte(xmls.XmlFinalNF), 0644)
				if err != nil {
					log.Printf("Erro ao copiar arquivo: %v", err)
				}

				f, ferr := os.Open(filepath)
				if err != nil {
					panic(ferr)
				}
				defer f.Close()
				w1, err := zipWriter.Create(filepath)
				if err != nil {
					println("erro1")
					panic(err)
				}
				if _, err := io.Copy(w1, f); err != nil {
					println("erro2")
					panic(err)
				}
			}
		}
	}

	fmt.Println("Fechando arquivo zip...")
	zipWriter.Flush()
	zipWriter.Close()
	return nil
}
