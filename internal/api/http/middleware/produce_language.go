package middleware

import (
	"github.com/gin-gonic/gin"
)

func ProduceLanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		supportedLanguages := []string{"en-US", "ru-RU"}
		defaultLanguage := "en-US"
		acceptLang := c.GetHeader("Accept-Language")

		for _, supportedLanguage := range supportedLanguages {
			if supportedLanguage == acceptLang {
				c.Set("lang", acceptLang)
				c.Next()
				return
			}
		}
		c.Set("lang", defaultLanguage)
		c.Next()
	}
}
