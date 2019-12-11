# dev engine for Telegram ChatBot
<img src="https://github.com/liquiddeath13/reactivetgbot/raw/master/media/blueberrycat.gif" width="150">

### Why do you need a ChatBot and why our team develop that tool?
Chatbots allow us to reduce the workload on staff by automating the processing of user requests. Therefore, it is necessary to create tools and libraries that accelerate the development.
![Useful case for providing ChatBot](https://github.com/liquiddeath13/reactivetgbot/raw/master/media/chatbot.png)
### Why should you use our solution?
Our solution allows you to forget about the code and directly engage in the creation of a processing environment
![Realization tips](https://github.com/liquiddeath13/reactivetgbot/raw/master/media/chatbot2.png)
And you only just to need fill that table with question-answer pair.
File with that pairs is using JSON format.
#### Example of question-answer pairs table
Provided in examplebase.json on this repository so now you can touch this with your own hands:)
```json

[
    {
        "Question" : "/about",
        "Answer" : "ChatBot created by liquiddeath13 in 2019 year"
    },
    {
        "Question" : "Hello, how can i contact with your company?",
        "Answer" : "Hello, you can contact us by provided mail address. Email: coolnickname@hostname.domain"
    }
]
```
As you can see, filling out the question-answer pairs table is quite simple. It is also simple like making life easier for your employees who, on duty of a job, need a lot and often communicate with customers.
### How-to use it inside Golang app
Firstly, you need to get that package with command:
```
go get "github.com/liquiddeath13/reactivetgbot"
```
And now you can use it at your application:
```golang
package main

import (
	"reactivetgbot"
)

func main() {
    BotInstance := reactivetgbot.Init("token", "/path/base.json")
    if BotInstance != nil {
        go BotInstance.Logic()
        //if we need host our application on Heroku and shouldn't think about uptime
        go reactivetgbot.HerokuServiceUP("Telegram ChatBot by liquiddeath13")
    }
}
```
### Programmer-defined handler for commands handling support
Since new version, now you are able to specify some programm logic to handle incoming message data, but now you should return from that handle only string (interface type in next versions coming).
Example:
```golang
package main

import (
	"reactivetgbot"
	"fmt"
)

func main() {
	BotInstance := reactivetgbot.Init("token", "/path/base.json")
	AskCounter := 0
	if BotInstance != nil {
		BotInstance.AppendHandler("how much peoples asked you?", func(Msg bbctg.TGMessage) string {
			AskCounter++
			return fmt.Sprintf("Hello.\n%d - so many people, who already asked me about this", AskCounter)
		})
		go BotInstance.Logic()
		go reactivetgbot.HerokuServiceUP("Telegram ChatBot by liquiddeath13")
	}
}
```
### Can you use it as standalone application?
It's impossible now, but in next versions it will be available to build application by cloning repo and just providing JSON file with question-answer pairs.
### Where should you communicate with developer?
You can register an account on [github](https://github.com/join) (if you have not already done so) and leave a question or suggestion in the [appropriate section of this repository](https://github.com/liquiddeath13/reactivetgbot/issues). If you want to get in touch in another way, then there is also [Gmail](mailto:ntlv.xca@gmail.com), [Telegram](http://t.me/s3thix), and [Vkontakte](https://vk.com/id554777800).
