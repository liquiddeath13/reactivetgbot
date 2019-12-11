package reactivetgbot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	UpdatesConfig := tgbotapi.NewUpdate(0)
	UpdatesConfig.Timeout = 40
	IUpdateChannel := HandlePanicError(b.Unit.GetUpdatesChan(UpdatesConfig))
	UpdateChannel := IUpdateChannel.(tgbotapi.UpdatesChannel)
	for {
		select {
		case Update := <-UpdateChannel:
			if Update.Message == nil {
				continue
			}
			Message := Update.Message
			AnswerType := reflect.TypeOf(b.QABase[Message.Text]).Kind().String()
			Answer := ""
			switch AnswerType {
			case "string":
				Answer = b.QABase[Message.Text].(string)
				break
			case "func":
				Handler := b.QABase[Message.Text].(func(TGMessage) string)
				Answer = Handler(Message)
				break
			}
			b.Unit.Send(tgbotapi.NewMessage(Message.Chat.ID, Answer))
		}
	}
}

func Init(token, qafile string) *Bot {
	IInstance := HandlePanicError(tgbotapi.NewBotAPI(token))
	Instance := IInstance.(*tgbotapi.BotAPI)
	IJSONFile := HandlePanicError(os.Open(qafile))
	JSONFile := IJSONFile.(*os.File)
	defer JSONFile.Close()
	ByteContent := HandlePanicError(ioutil.ReadAll(JSONFile))
	Local := []QAPair{}
	HandlePanicError(nil, json.Unmarshal(ByteContent.([]byte), &Local))
	Result := Bot{}
	Result.Unit = Instance
	for _, Object := range Local {
		Result.QABase[Object.Question] = Object.Answer
	}
	return &Result
}

func HerokuServiceUP(Description string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(Description))
	})
	go log.Fatal(http.ListenAndServe(os.Getenv("PORT"), nil))
	HerokuUpTimer := time.NewTimer(5 * time.Minute)
	URL := "https://api.ipify.org?format=text"
	for {
		Response, Error := http.Get(URL)
		if Error != nil {
			panic(Error)
		}
		defer Response.Body.Close()
		_, Error = ioutil.ReadAll(Response.Body)
		if Error != nil {
			panic(Error)
		}
		<-HerokuUpTimer.C
	}
}
