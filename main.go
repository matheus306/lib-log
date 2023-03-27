package main

import (
	"fmt"
	"lib-log/pkg/reqclientlog"
	"net/http"
)

func main() {

	for i := 0; i < 20; i++ {
		var config reqclientlog.LogConfig
		config.IgnoredHeaders = []string{"Expires", "Age", "Content-Type", "X-Ratelimit-Limit", "Report-To"}
		config.IgnoredFields = []string{"userId"}

		// Criando um cliente HTTP com o transport personalizado
		client := &http.Client{
			Transport: &reqclientlog.LoggingTransportDTO{Transport: http.DefaultTransport, LogConfig: &config},
		}

		// Enviando a solicitação HTTP

		resp, err := client.Get("https://catfact.ninja/fact")
		if err != nil {
			fmt.Println(err)
			return
		}

		// Fechando o corpo do response
		resp.Body.Close()
	}

}
