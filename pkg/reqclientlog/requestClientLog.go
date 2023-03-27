package reqclientlog

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"lib-log/pkg/cmd"
	"lib-log/pkg/dto"
	"log"
	"net/http"
	"time"
)

func (t *LoggingTransportDTO) RoundTrip(req *http.Request) (*http.Response, error) {

	// Imprimindo os dados do request
	var requestDTO dto.LogDTO

	requestDTO.Method = req.Method
	requestDTO.URL = req.URL.String()
	requestDTO.Headers = ocultarChavesHeader(req.Header, t)
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		requestDTO.Body = ocultarChavesBody(reqBody, t)
	}

	// Imprimindo o corpo do request
	go cmd.Init(&requestDTO)
	resultRequest, _ := json.MarshalIndent(requestDTO, "", "   ")
	log.Printf("  Request: %s\n", string(resultRequest))

	// Enviando a solicitação
	startTime := time.Now()
	resp, err := t.Transport.RoundTrip(req)
	elapsedTime := time.Since(startTime)

	//respHeadersJson := ocultarChavesHeader(resp.Header, t)

	// Imprimindo os dados do response
	var responseDTO dto.LogDTO

	responseDTO.Status = resp.Status
	responseDTO.Headers = ocultarChavesHeader(resp.Header, t)

	// Copiando o corpo do response para um buffer
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	responseDTO.Body = ocultarChavesBody(body, t)

	result, _ := json.MarshalIndent(responseDTO, "", "   ")
	log.Printf("  Response: %s\n", string(result))

	//Imprimindo o tempo de resposta
	log.Printf("  Tempo de resposta: %s\n", elapsedTime)

	go cmd.Init(&responseDTO)
	return resp, err
}

func ocultarChavesBody(jsonByte []byte, t *LoggingTransportDTO) map[string]interface{} {
	jsonMap := make(map[string]interface{})

	json.Unmarshal(jsonByte, &jsonMap)

	for _, field := range t.LogConfig.IgnoredFields {
		jsonMap[field] = "*********"
	}

	return jsonMap
}

func ocultarChavesHeader(header http.Header, t *LoggingTransportDTO) map[string][]string {
	var headers = make(map[string][]string)

	for key, values := range header {
		if t.LogConfig.IgnoredHeaders != nil && contains(t.LogConfig.IgnoredHeaders, key) {
			headers[key] = []string{"*******"}
			continue
		}
		headers[key] = values
	}

	return headers
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
