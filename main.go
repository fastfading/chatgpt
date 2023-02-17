package main

import (
	// "bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	gpt3 "github.com/fastfading/go-gpt3"
	"github.com/spf13/cobra"

	htgotts "github.com/hegedustibor/htgo-tts"
	handlers "github.com/hegedustibor/htgo-tts/handlers"
	voices "github.com/hegedustibor/htgo-tts/voices"
	"github.com/peterh/liner"
)

func GetResponse(client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(2000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		if len(resp.Choices) > 0 {
			txt := resp.Choices[0].Text
			fmt.Print(txt)
			BufferText(txt)
		}
	})
	if err != nil {
		fmt.Println(err)
		// os.Exit(13)
	}
	fmt.Printf("\n")
}

var bufStr string

func BufferText(txt string) {
	bufStr += txt
	if txt == "." || txt == "!" || txt == "?" {
		fmt.Println("")
		speak(bufStr)
		bufStr = ""
	}
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("Missing API KEY")
	}

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			// scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				// fmt.Print("Question (enter:'quit' to quit): \n")
				// fmt.Print(">")

				// if !scanner.Scan() {
				// 	break
				// }

				// question := scanner.Text()
				var question string
				var err error
				if question, err = line.Prompt(">"); err != nil {
					break
				}
				questionParam := validateQuestion(question)
				switch questionParam {
				case "quit":
					quit = true
				case "":
					continue

				default:
					GetResponse(client, ctx, questionParam)
				}
			}
		},
	}

	log.Fatal(rootCmd.Execute())
}

func validateQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "loop", "break", "continue", "cls", "exit", "block"}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}

func savemp3(str string) {
	speech := htgotts.Speech{Folder: "audio", Language: voices.English}
	speech.Speak(str)
}

func speak(str string) {
	speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.MPlayer{}}
	speech.Speak(str)
}

func nativespeak(str string) {
	speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
	speech.Speak(str)
}
