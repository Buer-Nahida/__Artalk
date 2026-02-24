package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/log"
)

type NotifyWebHookReqBody struct {
	NotifySubject string      `json:"notify_subject"`
	NotifyBody    string      `json:"notify_body"`
	Comment       interface{} `json:"comment"`
	ParentComment interface{} `json:"parent_comment"`
	Extra         map[string]string `json:"-"`
}

// WebHook 发送
func SendWebHook(url string, reqData *NotifyWebHookReqBody) {
	if (strings.Contains(url, "discord.com") || strings.Contains(url, "discordapp.com")) && reqData.Extra != nil {
		sendDiscordWebhook(url, reqData)
		return
	}

	jsonByte, _ := json.Marshal(reqData)
	result, err := http.Post(url, "application/json", bytes.NewReader(jsonByte))
	if err != nil {
		log.Error("[WebHook Push] ", "Failed to send msg:", err)
		return
	}

	if result.StatusCode != 200 {
		body, _ := io.ReadAll(result.Body)
		log.Error("[WebHook Push] Failed to send msg:", string(body))
	}

	defer result.Body.Close()
}

type discordWebhookBody struct {
	Content   string                `json:"content"`
	Username  string                `json:"username"`
	AvatarURL string                `json:"avatar_url"`
	Embeds    []discordWebhookEmbed `json:"embeds,omitempty"`
}

type discordWebhookEmbed struct {
	Description string                `json:"description"`
	Author      *discordWebhookAuthor `json:"author,omitempty"`
}

type discordWebhookAuthor struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}

func sendDiscordWebhook(url string, reqData *NotifyWebHookReqBody) {
	extra := reqData.Extra

	body := discordWebhookBody{
		Content:   reqData.Extra["content"],
		Username:  fmt.Sprintf("%s (%s)", extra["replier_nick"], extra["replier_email"]),
		AvatarURL: extra["replier_avatar"],
	}

	if extra["parent_nick"] != "" {
		body.Embeds = []discordWebhookEmbed{
			{
				Description: extra["parent_content"],
				Author: &discordWebhookAuthor{
					Name:    fmt.Sprintf("%s (%s)", extra["parent_nick"], extra["parent_email"]),
					IconURL: extra["parent_avatar"],
				},
			},
		}
	}

	jsonByte, _ := json.Marshal(body)
	result, err := http.Post(url, "application/json", bytes.NewReader(jsonByte))
	if err != nil {
		log.Error("[WebHook Push] ", "Failed to send msg:", err)
		return
	}

	if result.StatusCode < 200 || result.StatusCode >= 300 {
		respBody, _ := io.ReadAll(result.Body)
		log.Error("[WebHook Push] Failed to send msg to Discord:", string(respBody))
	}

	defer result.Body.Close()
}
