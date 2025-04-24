package utils

import (
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/gin-gonic/gin"
)

// ForwardRequest перенаправляет запрос на указанный сервис и возвращает ответ клиенту.
func ForwardRequest(c *gin.Context, method, url string) {
	client := resty.New()
	req := client.R()

	// Устанавливаем тело запроса
	if c.Request.Body != nil {
		req.SetBody(c.Request.Body)
	}

	// Копируем заголовки
	for k, v := range c.Request.Header {
		req.SetHeader(k, v[0])
	}

	// Отправляем запрос
	resp, err := req.Execute(method, url)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем статус код и возвращаем ответ
	c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
}
