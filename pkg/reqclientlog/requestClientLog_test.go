package reqclientlog

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	slice := []string{"foo", "bar", "baz"}
	element := "bar"
	if !contains(slice, element) {
		t.Errorf("contains(%v, %s) = false, expected true", slice, element)
	}

	slice2 := []string{"apple", "banana", "cherry"}
	element2 := "pear"
	if contains(slice2, element2) {
		t.Errorf("contains(%v, %s) = true, expected false", slice2, element2)
	}
}

func TestOcultarChavesHeader(t *testing.T) {

	var config LogConfig
	config.IgnoredHeaders = []string{"Expires", "Age", "Content-Type", "X-Ratelimit-Limit", "Report-To"}
	config.IgnoredFields = []string{"userId"}

	// Configuração do teste
	header := make(http.Header)
	header.Add("Authorization", "Bearer my-secret-token")
	header.Add("Content-Type", "application/json")
	header.Add("X-Api-Key", "my-api-key")
	ignoredHeaders := []string{"Authorization", "X-Api-Key"}

	transport := &LoggingTransportDTO{Transport: http.DefaultTransport, LogConfig: &config}
	transport.LogConfig.IgnoredHeaders = ignoredHeaders

	// Executa a função
	headers := ocultarChavesHeader(header, transport)

	// Verifica se as chaves estão ocultas
	expected := map[string][]string{
		"Authorization": {"*******"},
		"Content-Type":  {"application/json"},
		"X-Api-Key":     {"*******"},
	}
	if !reflect.DeepEqual(headers, expected) {
		t.Errorf("ocultarChavesHeader(%v, %v) = %v, expected %v", header, transport, headers, expected)
	}
}

func TestOcultarChavesBody(t *testing.T) {
	// Configuração do teste
	jsonStr := `{
        "name": "John Doe",
        "email": "john.doe@example.com",
        "password": "my-secret-password",
        "address": {
            "street": "123 Main St",
            "city": "Anytown",
            "state": "CA",
            "zip": "12345"
        }
    }`
	jsonByte := []byte(jsonStr)
	ignoredFields := []string{"password"}

	var config LogConfig
	config.IgnoredHeaders = []string{"Expires", "Age", "Content-Type", "X-Ratelimit-Limit", "Report-To"}
	config.IgnoredFields = ignoredFields

	transport := &LoggingTransportDTO{Transport: http.DefaultTransport, LogConfig: &config}

	// Executa a função
	jsonMap := ocultarChavesBody(jsonByte, transport)

	// Verifica se os campos estão ocultos
	expected := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john.doe@example.com",
		"password": "*********",
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "Anytown",
			"state":  "CA",
			"zip":    "12345",
		},
	}
	if !reflect.DeepEqual(jsonMap, expected) {
		t.Errorf("ocultarChavesBody(%s, %v) = %v, expected %v", jsonStr, transport, jsonMap, expected)
	}
}

func TestRoundTrip(t *testing.T) {

	var config LogConfig
	config.IgnoredHeaders = []string{"Expires", "Age", "Content-Type", "X-Ratelimit-Limit", "Report-To"}
	config.IgnoredFields = []string{"userId"}

	// Configuração do teste
	client := &http.Client{
		Transport: &LoggingTransportDTO{Transport: http.DefaultTransport, LogConfig: &config},
	}

	// Cria um request HTTP para testar
	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer secret-token")
	req.Header.Set("Content-Type", "application/json")
	reqBody := `{"username": "johndoe", "password": "secret"}`
	req.Body = ioutil.NopCloser(strings.NewReader(reqBody))

	var response *http.Response

	// Executa a função para o request
	response, err = client.Do(req)

	if err != nil {
		t.Errorf("RoundTrip(%v) returned error: %v", req, err)
	}

	if response.StatusCode != 200 {
		t.Errorf("response(%s), expected %v", response.Status, 200)
	}
}
