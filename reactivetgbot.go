package reactivetgbot

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type (
	QAPair struct {
		Question string `json:"Question"`
		Answer   string `json:"Answer"`
	}
	Bot struct {
		Unit      *tgbotapi.BotAPI
		Token     string
		LastError error
		QABase    map[string]interface{}
	}
	TGMessage *tgbotapi.Message
)

var Host = ""

func HandlePanicError(arg interface{}, Error error) interface{} {
	if Error != nil {
		log.Panic(Error.Error())
	}
	return arg
}

func HandleInfoError(arg interface{}, Error error) interface{} {
	if Error != nil {
		log.Print(Error.Error())
	}
	return arg
}

func (b *Bot) AppendHandler(Question string, Handler func(Msg TGMessage) string) {
	b.QABase[Question] = Handler
}

func (b *Bot) Logic() {
	for {
		select {
		case Update := <-b.Unit.ListenForWebhook("/" + b.Token):
			if Update.Message == nil {
				continue
			}
			Message := Update.Message
			if b.QABase[Message.Text] == nil {
				continue
			}
			AnswerType := reflect.TypeOf(b.QABase[Message.Text]).Kind()
			Answer := ""
			switch AnswerType {
			case reflect.String:
				Answer = b.QABase[Message.Text].(string)
				break
			case reflect.Func:
				Handler := b.QABase[Message.Text].(func(TGMessage) string)
				Answer = Handler(Message)
				break
			}
			HandleInfoError(b.Unit.Send(tgbotapi.NewMessage(Message.Chat.ID, Answer)))
		}
	}
}

func Init(token, dictionary string) *Bot {
	newTelegramBot := Bot{}
	IInstance := HandlePanicError(tgbotapi.NewBotAPI(token))
	newTelegramBot.Unit = IInstance.(*tgbotapi.BotAPI)
	newTelegramBot.QABase = make(map[string]interface{})
	if dictionary != "" {
		IJSONFile := HandlePanicError(os.Open(dictionary))
		JSONFile := IJSONFile.(*os.File)
		defer JSONFile.Close()
		ByteContent := HandlePanicError(ioutil.ReadAll(JSONFile))
		var Local []QAPair
		HandlePanicError(nil, json.Unmarshal(ByteContent.([]byte), &Local))
		for _, Object := range Local {
			newTelegramBot.QABase[Object.Question] = Object.Answer
		}
	}
	return &newTelegramBot
}

func (b *Bot) HerokuUsage(Description string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Host = strings.Split(strings.Split(strings.Split(r.RequestURI, "//")[1], "/")[0], ".")[0]
		println(Host)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		HandleInfoError(w.Write([]byte(Description)))
	})
	go log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	HandleInfoError(http.Get(":" + os.Getenv("PORT")))
	HandlePanicError(http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=https://%s.herokuapp.com/%s", b.Token, Host, b.Token)))
}
