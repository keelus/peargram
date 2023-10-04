package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func ForceError(codigoErrorHttp int, c *gin.Context, razonVisual string) {
	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "error", gin.H{
			"CodigoError": codigoErrorHttp,
			"TextoError":  http.StatusText(codigoErrorHttp),
			"Detalles":    razonVisual,
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"panelToLoad": "error",
			"CodigoError": codigoErrorHttp,
			"TextoError":  http.StatusText(codigoErrorHttp),
			"Detalles":    razonVisual,
		})
	}
}

var vConfig *viper.Viper

func LoadEnv() {

	vConfig = viper.New()
	vConfig.SetConfigType("env")
	vConfig.SetConfigFile(".env.local")

	if err := vConfig.ReadInConfig(); err != nil {
		fmt.Println(".env.local not found. Trying with .env ...")

		vConfig.SetConfigFile(".env")
		if err2 := vConfig.ReadInConfig(); err2 != nil {
			panic("Error reading the .env file")
		} else {
			fmt.Println(".env loaded")
		}
	} else {
		fmt.Println(".env.local loaded")
	}

}

func GetEnv(key string) string {
	value, ok := vConfig.Get(key).(string)
	if !ok {
		panic(fmt.Sprintf("Error getting the key '%s' FROM .env file.", key))
	}

	return value
}

func GetAbsolutePath() string {
	path, _ := os.Getwd()
	return path
}
