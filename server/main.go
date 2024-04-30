package main

import (
	// "bytes"
	"net/http"
	// "io"
	"fmt"
	"io/ioutil"
	"os"
	// "log"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	// "github.com/rs/zerolog"
	"github.com/dn365/gin-zerolog"
	"github.com/rs/zerolog/log"
)

type AiResponse struct {
	Content string `json:"content"`
}

type HealthCheck struct {
	Status string `json:status`
}

func getImageAnalysis(contentType string, image []byte, question string) *genai.GenerateContentResponse {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal().Err(err)
	}
	defer client.Close()

	// For text-and-image input (multimodal), use the gemini-pro-vision model
	model := client.GenerativeModel("gemini-pro-vision")

	prompt := []genai.Part{
		genai.ImageData(contentType, image),
		genai.Text(question),
	}
	resp, err := model.GenerateContent(ctx, prompt...)

	if err != nil {
		log.Fatal().Err(err)
	}

	return resp
}

func readImage(ctx *gin.Context) (string, []byte) {
	imageType := ctx.PostForm("imageType")
	// Expected Content-Type: image/jpeg | image/png | image/svg
	log.Info().Msgf("[+] Image Type: %s", imageType)

	// 	if imageType == nil {
	// 		ctx.JSON(http.StatusBadRequest, err)
	// 		return "", []byte{}
	// 	}

	formFile, err := ctx.FormFile("image")
	openedFile, _ := formFile.Open()
	defer openedFile.Close()
	image, err := ioutil.ReadAll(openedFile)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid image err: %s", err.Error())
		return "", []byte{}
	}
	return imageType, image
}

func describeImage(ctx *gin.Context) {
	contentType, imageData := readImage(ctx)
	question := "Suppose I am blind, tell me what you can see, without saying it is a picture, while speaking naturally."
	description := readAiResponse(getImageAnalysis(contentType, imageData, question))
	log.Debug().Msgf("[+] Description: %s", description)
	ctx.IndentedJSON(http.StatusOK, AiResponse{Content: description})
}

func translateImage(ctx *gin.Context) {
	contentType, imageData := readImage(ctx)
	question := "If there is English text on the image return what you read, if there is text but it is not English return its translation, else mention if you didn't find anything."
	translation := readAiResponse(getImageAnalysis(contentType, imageData, question))
	log.Debug().Msgf("[+] Translation: %s", translation)
	ctx.IndentedJSON(http.StatusOK, AiResponse{Content: translation})
}

func healthCheck(ctx *gin.Context) {
	var message HealthCheck
	message.Status = "It works!"
	ctx.IndentedJSON(http.StatusOK, message)
}

func main() {
	router := gin.Default()
	router.Use(ginzerolog.Logger("gin"))
	router.GET("/gemini", healthCheck)
	router.POST("/gemini/describe", describeImage)
	router.POST("/gemini/translate", translateImage)
	router.Run(":8080")
}

func readAiResponse(resp *genai.GenerateContentResponse) string {
	var sentence string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				sentence += fmt.Sprint(part)
			}
		}
	}
	return sentence
}
