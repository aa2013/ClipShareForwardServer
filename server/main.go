package main

import (
	"clipshare/db"
	"clipshare/forward"
	"clipshare/types"
	"clipshare/utils"
	"clipshare/web"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
)

var logs *utils.LogManager
var Env string

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[error] main", err)
		}
		db.CloseDB()
	}()
	db.ConnectDB("./data/app.db")
	utils.LogUtil = utils.NewLogManager(1000, "./data/logs")
	logs = utils.LogUtil
	types.AppConfig = types.ReadConfig()
	types.WatchConfig(onConfigChanged)
	go forward.StartForwardServer()
	web.StartWebServer()
}

func onConfigChanged(e fsnotify.Event) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("onConfigChanged", r)
		}
	}()
	reloadConfig := types.ReadConfig()
	reloadConfig.CheckValues()
	types.AppConfig = reloadConfig
	forward.UpdateRateLimitConfig()
	logs.Info("Config file changed:", e.Name)
}

// Restart self program
func restart() error {
	// 获取当前程序的路径
	path, err := os.Executable()
	if err != nil {
		logs.Info("failed to get executable path:", err)
		return err
	}
	// 创建一个新的命令来启动自身
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 启动新的程序实例
	if err := cmd.Start(); err != nil {
		logs.Info("failed to get executable path:", err)
		return err
	}
	// 终止当前程序
	os.Exit(0)
	return nil
}
