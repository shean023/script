package main

import (
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// 配置你的 Telegram Bot Token
var botToken = "7336EVP3Z8ieMgqzINJgz"

// 设置监听的源群组 ID 和目标群组 ID 列表
var sourceGroupID int64 = -4668108 // 监听的源群组的 chat_id
var targetGroups = []int64{
	-4724966,
}

func main() {
	// 创建一个 Telegram Bot 实例
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// 设置 Webhook 或者轮询更新
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 获取消息更新
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// 监听消息
	for update := range updates {
		// 检查消息是否来自指定的源群组
		//if update.Message != nil && update.Message.Chat.ID == sourceGroupID {
		if update.Message.Chat.ID == sourceGroupID {
			// 转发消息到多个目标群组
			for _, groupID := range targetGroups {
				// 发送相同的消息到目标群组
				msg := tgbotapi.NewForward(groupID, update.Message.Chat.ID, update.Message.MessageID)
				sentMessage, err := bot.Send(msg)

				if err != nil {
					log.Printf("转发失败到群组 %d: %v", groupID, err)
				} else {
					log.Printf("成功转发消息到群组 %d", groupID)
					go func(msgID int, chatID int64) {
						// 等待 1 分钟
						time.Sleep(1 * time.Minute)

						// 删除转发的消息
						deleteMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
						_, err := bot.Send(deleteMsg)
						if err != nil {
							log.Printf("删除消息失败: %v", err)
						} else {
							log.Printf("成功删除群组 %d 中的消息: %d", chatID, msgID)
						}
					}(sentMessage.MessageID, groupID)
				}
			}
		}
	}
}

