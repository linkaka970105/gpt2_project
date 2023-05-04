package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	openai "github.com/sashabaranov/go-openai"
	"time"
)

func GetGptResp(uid int, msg string) (content string, err error) {
	// 获取历史聊天
	token, err := getTokenByUID(uid)
	if err != nil {
		return
	}
	var history []openai.ChatCompletionMessage
	rs, err := RedisCli().Get(token).Result()
	if err != nil && err != redis.Nil {
		return
	}
	var expireTime time.Duration
	if rs == "" {
		history = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是ChatGPT, 一个由OpenAI训练的大型语言模型, 你旨在回答并解决人们的任何问题，并且可以使用多种语言与人交流。",
			},
		}
	} else {
		err = json.Unmarshal([]byte(rs), &history)
		expireTime, err = RedisCli().TTL(token).Result()
		if err != nil {
			return
		}
	}
	// 返回结果
	client := GptCli()
	query := append([]openai.ChatCompletionMessage{}, history...)
	query = append(query, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 2000,
			Messages:  query,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	content = resp.Choices[0].Message.Content
	// 记录历史
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})
	b, err := json.Marshal(history)
	if err != nil {
		return
	}
	if expireTime.Seconds() <= 0 {
		expireTime = 24 * 2600 * time.Second
	}
	err = RedisCli().Set(token, string(b), expireTime).Err()
	return
}
