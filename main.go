package main

/*
Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.
As duas requisições serão feitas simultaneamente para as seguintes APIs:
https://brasilapi.com.br/api/cep/v1/01153000 + cep
http://viacep.com.br/ws/" + cep + "/json/
Os requisitos para este desafio são:
- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)
	
type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilAPICEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	cepArg := os.Args[1]
	viaCEPChannel := make(chan *ViaCEP)
	brasilAPICEPChannel := make(chan *BrasilAPICEP)
	fmt.Fprintf(os.Stderr, "Starting search for CEP...\n")
	
	go func() {
		brasilAPICEPData, error := SearchBrasilAPICep(cepArg)
		if error != nil {
			fmt.Fprintf(os.Stderr, "Error when tried to get CEP from BrasilAPI data %v\n", error)
		}
		brasilAPICEPChannel <- brasilAPICEPData
	}()
	go func() {
		viaCEPData, error := SearchCEPViaCep(cepArg)
		if error != nil {
			fmt.Fprintf(os.Stderr, "Error when tried to get CEP data from ViaCEP %v\n", error)
		}
		viaCEPChannel <- viaCEPData
	}()

	select {
	case response := <-viaCEPChannel:
		PrintViaCEPData(*response)
	case response := <-brasilAPICEPChannel:
		PrintBrasilAPICEPData(*response)
	case <-time.After(time.Second):
		println("Timeout when tried to get CEP data...")
	}
}

func SearchBrasilAPICep(cep string) (*BrasilAPICEP, error) {
	fmt.Fprintf(os.Stderr, "Starting search on BrasilAPI api...\n")
	var data BrasilAPICEP
	client := http.Client{Timeout: time.Second}
	response, error := client.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if error != nil {
		if os.IsTimeout(error) {
			fmt.Println("Timeout error ocurred")
			return nil, error
		}
		return nil, error
	}
	defer response.Body.Close()
	body, error := io.ReadAll(response.Body)

	if error != nil {
		return nil, error
	}
	error = json.Unmarshal(body, &data)
	fmt.Println(data)
	if error != nil {
		return nil, error
	}
	return &data, nil
}

func SearchCEPViaCep(cep string) (*ViaCEP, error) {
	fmt.Fprintf(os.Stderr, "Starting search on ViaCEP api...\n")
	var data ViaCEP
	client := http.Client{Timeout: time.Second}
	response, error := client.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if error != nil {
		if os.IsTimeout(error) {
			fmt.Println("Timeout error ocurred")
			return nil, error
		}
		return nil, error
	}
	defer response.Body.Close()
	body, error := io.ReadAll(response.Body)
	if error != nil {
		return nil, error
	}
	error = json.Unmarshal(body, &data)
	if error != nil {
		return nil, error
	}
	return &data, nil
}

func PrintViaCEPData(data ViaCEP) {
	fmt.Printf("\nAddress: %v", data.Logradouro)
	fmt.Printf("\nNeighborhood: %v", data.Bairro)
	fmt.Printf("\nCity: %v", data.Unidade)
	fmt.Printf("\nState: %v", data.Estado)
	fmt.Printf("\nSource: ViaCEP")
}

func PrintBrasilAPICEPData(data BrasilAPICEP) {
	fmt.Printf("\nAddress: %v", data.Street)
	fmt.Printf("\nNeighborhood: %v", data.Neighborhood)
	fmt.Printf("\nCity: %v", data.City)
	fmt.Printf("\nState: %v", data.State)
	fmt.Printf("\nSource: BrasilAPI")
}
