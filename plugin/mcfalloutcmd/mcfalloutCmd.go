package mcfalloutcmd

import (
	"fmt"
	"strings"

	"github.com/Tnze/go-mc/chat"
	"github.com/fsnotify/fsnotify"
	"github.com/rain931215/go-mc-api/api"
	"github.com/spf13/viper"
)

//Func is the type of command's method
type Func = func(c *api.Client, Sender string, Text string, Args []string)

//Command is contained command's name and command's method
type Command struct {
	name   string
	method Func
}

//McfalloutCmd _
type McfalloutCmd struct {
	Client    *api.Client
	whiteList []string
	cmdList   []*Command
}

// New _
func New(c *api.Client) *McfalloutCmd {
	p := new(McfalloutCmd)
	p.Client = c
	c.Event.AddEventHandler(p.main, "chat")

	file := viper.New()
	//Load whiteList
	file.SetConfigName("whiteList")
	file.SetConfigType("yaml")
	file.AddConfigPath("./plugin/mcfalloutcmd")
	file.WatchConfig()
	err := file.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	p.whiteList = file.GetStringSlice("admin")
	//熱插拔
	file.OnConfigChange(func(e fsnotify.Event) {
		p.whiteList = file.GetStringSlice("admin")
		fmt.Println("White List Change")
	})
	//Load defaultCommand
	p.loadDefaultCmd()

	return p
}

func (p *McfalloutCmd) main(msg chat.Message) bool {
	var text = msg.ClearString()
	for id := 0; id < len(p.whiteList); id++ {
		if strings.Index(text, "[收到私訊 "+p.whiteList[id]) == 0 {
			text = strings.TrimPrefix(text, "[收到私訊 "+p.whiteList[id]+"] : ")
			args := strings.Split(text, " ")
			for i := 0; i < len(p.cmdList); i++ {
				if args[0] == p.cmdList[i].name {
					text = strings.TrimPrefix(text, p.cmdList[i].name+" ")
					p.cmdList[i].method(p.Client, p.cmdList[i].name, text, args[1:])
					return false
				}
			}
			break
		}
		if text == "[廢土伺服] : "+p.whiteList[id]+" 想要傳送到 你 的位置" {
			p.Client.Chat("/tok")
			return false
		} else if text == "[廢土伺服] : "+p.whiteList[id]+" 想要你傳送到 該玩家 的位置!" {
			p.Client.Chat("/tok")
			return false
		}
	}
	return false
}

// AddCmd _
func (p *McfalloutCmd) AddCmd(name string, command Func) {
	newCommand := new(Command)
	newCommand.name = name
	newCommand.method = command
	p.cmdList = append(p.cmdList, newCommand)
}
