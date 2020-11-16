package main

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	tb "gopkg.in/tucnak/telebot.v2"
)

func getSlug() string {
	b := make([]byte, 4)
	rand.Read(b)
	s := hex.EncodeToString(b)
	return s
}

func contains(a string, list [][]string) int {
	for i, b := range list {
		if b[0] == a {
			return i
		}
	}
	return 0
}

func setUserLanguage(user string, language string, dataFile string) {
	f, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	res := contains(user, data)
	if res != 0 {
		data[res][1] = language
	} else {
		data = append(data, []string{user, language})
	}
	f.Truncate(0)
	f.Seek(0, 0)
	w := csv.NewWriter(f)
	w.WriteAll(data)
}

func getUserLanguage(user string, dataFile string) string {
	f, err := os.OpenFile(dataFile, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	res := contains(user, data)
	if res == 0 {
		setUserLanguage(user, "en", dataFile)
		log.Println("User language not set. Setting it to default language English.")
		return "en"
	}
	return data[res][1]
}

func getTranscript(filePath string, languageCode string) string {
	ctx := context.Background()

	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_OGG_OPUS,
			SampleRateHertz: 16000,
			LanguageCode:    languageCode,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		log.Fatalf("Failed to recognize: %v", err)
	}
	if len(resp.Results) == 0 {
		return "Transcription unsuccessful 💀"
	}
	return resp.Results[0].Alternatives[0].Transcript
}

func main() {
	dataFile := "user_language.csv"
	languageMenu := &tb.ReplyMarkup{}

	languageMenu.Inline(
		languageMenu.Row(tb.Btn{
			Unique: "en",
			Text:   "English",
		}),
		languageMenu.Row(tb.Btn{
			Unique: "fr",
			Text:   "French",
		}),
		languageMenu.Row(tb.Btn{
			Unique: "de",
			Text:   "German",
		}),
		languageMenu.Row(tb.Btn{
			Unique: "ru",
			Text:   "Russian",
		}),
		languageMenu.Row(tb.Btn{
			Unique: "zh",
			Text:   "Chinese",
		}),
	)

	languageMenu.InlineKeyboard = append(languageMenu.InlineKeyboard)

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("BOT_API_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Oh, hello there. \nSet your /language and send me a voice message, I will recognize the text and send it back to you!")
		setUserLanguage(strconv.Itoa(m.Sender.ID), "en", dataFile)
	})

	b.Handle("/language", func(m *tb.Message) {
		b.Send(m.Sender, "Languages available", languageMenu)
	})

	for _, button := range languageMenu.InlineKeyboard {
		button := button
		b.Handle(&button[0], func(c *tb.Callback) {
			setUserLanguage(strconv.Itoa(c.Sender.ID), button[0].Unique, dataFile)
			message := fmt.Sprintf("Language was set to %s. Send me a voice message.", button[0].Text)
			b.Respond(c, &tb.CallbackResponse{
				Text: message,
			})
			b.Send(c.Sender, message)
		})
	}

	b.Handle(tb.OnVoice, func(m *tb.Message) {
		userID := strconv.Itoa(m.Sender.ID)
		f, err := b.FileByID(m.Voice.FileID)
		if err != nil {
			log.Fatal(err)
			return
		}
		filePath := fmt.Sprintf("%s-%s.ogg", userID, getSlug())
		b.Download(&f, filePath)
		languageCode := getUserLanguage(userID, dataFile)
		transcription := getTranscript(filePath, languageCode)
		b.Send(m.Sender, transcription)
		os.Remove(filePath)
	})

	b.Start()
}
