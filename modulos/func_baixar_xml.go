package modulos

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	"conexao_mysql/db"

	_ "github.com/go-sql-driver/mysql"
)

func BaixarXmls(dataInicial, dataFinal, emissorP, emissorT, _nfe, _nfce string) error {

	var diretorioAnterior = ""
	db, err := db.GetDB()
	if err != nil {
		log.Fatal("database connection error, see .env configuration")
	}

	defer db.Close()

	//query := "select notafiscal, serie, modelo, chave_nota , xml_final, nota_cabecalho.fs_fase from nota_xml join nota_cabecalho ON nota_cabecalho.id = nota_xml.nota_cabecalho_id where xml_final is not null and data_nota between '2023-11-01' and '2023-11-30' and modelo in ('55','65') and emissor = 'P' order by modelo;"
	query := fmt.Sprintf("select xml_final, nota_cabecalho.fs_fase, nota_xml.emissor from nota_xml join nota_cabecalho ON nota_cabecalho.id = nota_xml.nota_cabecalho_id where xml_final is not null and data_nota between '%s' and '%s' and emissor in('%s', '%s') and modelo in ('%s','%s') order by data_nota;", dataInicial, dataFinal, emissorP, emissorT, _nfe, _nfce)
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

				if vigenciaXml != diretorioAnterior {
					fmt.Println("Copiando xmls referência: ", vigenciaXml)
				}
				diretorioAnterior = vigenciaXml
				//fmt.Println("Copiando arquivo:", filename)
				err = os.WriteFile(filepath, []byte(xmls.XmlFinalNF), 0644)
				if err != nil {
					log.Printf("Erro ao copiar arquivo: %v", err)
				}

				f, ferr := os.Open(filepath)
				if err != nil {
					panic(ferr)
				}
				defer f.Close()
			}
		}
	}

	return nil
}
