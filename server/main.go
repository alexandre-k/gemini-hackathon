package main

import "bytes"
import "io"
import "os"
import "log"
// import "strings"
import "context"
import "fmt"
import "github.com/google/generative-ai-go/genai"
import "google.golang.org/api/option"
// import "github.com/hegedustibor/htgo-tts"
// import "github.com/hegedustibor/htgo-tts/voices"
// import htgotts "github.com/hegedustibor/htgo-tts"
// import handlers "github.com/hegedustibor/htgo-tts/handlers"
// import voices "github.com/hegedustibor/htgo-tts/voices"
import "github.com/kataras/iris/v12"



func getImageAnalysis(contentType string, image []byte, question string) *genai.GenerateContentResponse {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-and-image input (multimodal), use the gemini-pro-vision model
	model := client.GenerativeModel("gemini-pro-vision")

	// imgData, err := os.ReadFile("signal-2021-06-23-140030.jpg")
	// if err != nil {
	//   log.Fatal(err)
	// }
	prompt := []genai.Part{
		genai.ImageData(contentType, image),
		genai.Text(question),
	}
	resp, err := model.GenerateContent(ctx, prompt...)

	if err != nil {
	  log.Fatal(err)
	}

	return resp
}


// func readText(content string) {
// 	fmt.Println("Read text: ", content)
// 	// speech := htgotts.Speech{Folder: "audio", Language: voices.English}
// 	// speech.Speak(content)

// 	speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
// 	fmt.Println("Speak.")
// 	sentences := strings.Split(content, ".")
// 	for _, sentence := range sentences {
// 		speech.Speak(sentence)
// 	}
// }


type AiResponse struct {
	Content string `json:"content"`
}

type Image struct {
	content []byte
	contentType string

}


func readImage(ctx iris.Context) (string, []byte) {
	contentType := ctx.FormValue("imageType")
	fmt.Println("CONTENT TYPE <<<<<<<<<<<<<< ")
	fmt.Println(contentType)
	file, _, err := ctx.FormFile("image")
	defer file.Close()
	// Expected Content-Type: image/jpeg | image/png | image/svg
	if err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return "", []byte{}
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return "", []byte{}
	}
	fmt.Println("IMAGE <<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println(buf.Bytes())
	return contentType, buf.Bytes()
}

func describeImage(ctx iris.Context) {
	contentType, imageData := readImage(ctx)
	question := "Suppose I am blind, tell me what you can see, without saying it is a picture, while speaking naturally."
	description := readAiResponse(getImageAnalysis(contentType, imageData, question))
	fmt.Println("description")
	fmt.Println(description)
	ctx.JSON(AiResponse{Content: description})
}


func translateImage(ctx iris.Context) {
	contentType, imageData := readImage(ctx)
	question := "If there is English text on the image return what you read, if there is text but it is not English return its translation, else mention if you didn't find anything."
	translation := readAiResponse(getImageAnalysis(contentType, imageData, question))
	fmt.Println("translation")
	fmt.Println(translation)
	ctx.JSON(AiResponse{Content: translation})
}

type HealthCheck struct {
	Status string `json:status`
}

func healthCheck(ctx iris.Context) {
	var h HealthCheck
	h.Status = "It works!"
	ctx.JSON(h)
}

func main () {
	app := iris.New()
	geminiApi := app.Party("/gemini")
	{
		geminiApi.Get("/", healthCheck)
		geminiApi.Post("/describe", describeImage)
		geminiApi.Post("/translate", translateImage)
	}
	app.Listen(":8080")

	// image := os.Args[1]
	// readText(description)
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

