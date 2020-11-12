package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

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

func main() {
	dataFile := "user_language.csv"
	var languageMenu = &tb.ReplyMarkup{}
	var fr = languageMenu.Data("French", "fr")
	var en = languageMenu.Data("English", "en")
	var es = languageMenu.Text("Spanish")
	languageMenu.Inline(
		languageMenu.Row(fr),
		languageMenu.Row(en),
		languageMenu.Row(es),
	)

	log.Println(languageMenu.InlineKeyboard)
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
		b.Send(m.Sender, "Oh, hello there. \nSet your /language")
	})

	b.Handle("/language", func(m *tb.Message) {
		setUserLanguage(strconv.Itoa(m.Sender.ID), "en", dataFile)
		b.Send(m.Sender, "test inline keyboard", languageMenu)

	})

	b.Handle(&fr, func(c *tb.Callback) {
		b.Respond(c, &tb.CallbackResponse{
			Text: "Language was set",
		})
		b.Send(c.Sender, "test message")
	})

	b.Start()
}
