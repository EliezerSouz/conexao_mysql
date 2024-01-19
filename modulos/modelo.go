package modulos

import "encoding/xml"

type NotaXml struct {
	NumNF       string `json:"notafiscal"`
	TipoEmissor string `json:"emissor"`
	ModeloNF    string `json:"modelo"`
	ChaveNF     string `json:"chave_nota"`
	XmlFinalNF  string `json:"xml_final"`
	StatusNF    string `json:"fs_fase"`
}
type NFe struct {
	XMLName xml.Name `xml:"nfeProc"`
	NFe     struct {
		XMLName xml.Name `xml:"NFe"`
		InfNFe  struct {
			Versao string `xml:"versao,attr"`
			ID     string `xml:"Id,attr"`
			Ide    struct {
				Mod string `xml:"mod"`
			} `xml:"ide"`
			Emit struct {
				CNPJ        string `xml:"CNPJ"`
				RazaoSocial string `xml:"xNome"`
				Fantasia    string `xml:"xFant"`
			} `xml:"emit"`
		} `xml:"infNFe"`
		Signature struct {
			XMLName xml.Name `xml:"Signature"`
		} `xml:"Signature"`
	} `xml:"NFe"`
	ProtNFe struct {
		Versao  string `xml:"versao,attr"`
		InfProt struct {
			TpAmb    string `xml:"tpAmb"`
			VerAplic string `xml:"verAplic"`
			ChNFe    string `xml:"chNFe"`
			DhRecbto string `xml:"dhRecbto"`
			NProt    string `xml:"nProt"`
			DigVal   string `xml:"digVal"`
			CStat    string `xml:"cStat"`
			XMotivo  string `xml:"xMotivo"`
		} `xml:"infProt"`
	} `xml:"protNFe"`
}
type NFe2 struct {
	XMLName xml.Name `xml:"NFe"`
	InfNFe  struct {
		Versao string `xml:"versao,attr"`
		ID     string `xml:"Id,attr"`
		Ide    struct {
			Mod string `xml:"mod"`
		} `xml:"ide"`
		Emit struct {
			CNPJ        string `xml:"CNPJ"`
			RazaoSocial string `xml:"xNome"`
			Fantasia    string `xml:"xFant"`
		} `xml:"emit"`
	} `xml:"infNFe"`
	Signature struct {
		XMLName xml.Name `xml:"Signature"`
	} `xml:"Signature"`
}
type NFe3 struct {
	XMLName xml.Name `xml:"nfeProc"`
	NFe     struct {
		XMLName xml.Name `xml:"NFe"`
		InfNFe  struct {
			Versao string `xml:"versao,attr"`
			ID     string `xml:"Id,attr"`
			Ide    struct {
				Mod string `xml:"mod"`
			} `xml:"ide"`
			Emit struct {
				CNPJ        string `xml:"CNPJ"`
				RazaoSocial string `xml:"xNome"`
				Fantasia    string `xml:"xFant"`
			} `xml:"emit"`
		} `xml:"infNFe"`
		Signature struct {
			XMLName xml.Name `xml:"Signature"`
		} `xml:"Signature"`
	} `xml:"NFe"`
	ProtNFe struct {
		Versao  string `xml:"versao,attr"`
		InfProt struct {
			ID       string `xml:"Id,attr"`
			TpAmb    string `xml:"tpAmb"`
			VerAplic string `xml:"verAplic"`
			ChNFe    string `xml:"chNFe"`
			DhRecbto string `xml:"dhRecbto"`
			NProt    string `xml:"nProt"`
			DigVal   string `xml:"digVal"`
			CStat    string `xml:"cStat"`
			XMotivo  string `xml:"xMotivo"`
		} `xml:"infProt"`
	} `xml:"protNFe"`
}
