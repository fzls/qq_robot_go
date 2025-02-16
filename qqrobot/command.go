package qqrobot

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Mrs4s/go-cqhttp/internal/msg"
	"github.com/pkg/errors"

	"github.com/Mrs4s/go-cqhttp/qqrobot/qinglong"

	"github.com/Mrs4s/MiraiGo/message"
	logger "github.com/sirupsen/logrus"
)

// 2021/10/02 5:25 by fzls
func (r *QQRobot) processCommand(commandStr string, m *message.GroupMessage) (rspMsg string, extraReplies []message.IMessageElement, err error) {
	var match []string
	if match = commandregexAddwhitelist.FindStringSubmatch(commandStr); len(match) == len(commandregexAddwhitelist.SubexpNames()) {
		// full_match|ruleName|qq
		ruleName := match[1]
		qq, _ := strconv.ParseInt(match[2], 10, 64)
		for _, rule := range r.Rules {
			if ruleName != rule.Config.Name {
				continue
			}
			rule.Config.ExcludeQQs = append(rule.Config.ExcludeQQs, qq)
			logger.Info("【Command】", commandStr)

			if len(rspMsg) != 0 {
				rspMsg += " | "
			}
			rspMsg += fmt.Sprintf("已将【%v】加入到规则【%v】的白名单", qq, ruleName)
		}
	} else if match = commandregexRulenamelist.FindStringSubmatch(commandStr); len(match) == len(commandregexRulenamelist.SubexpNames()) {
		for _, rule := range r.Rules {
			if _, ok := rule.Config.GroupIds[m.GroupCode]; !ok {
				continue
			}

			if len(rspMsg) == 0 {
				rspMsg += "规则集合："
			}
			rspMsg += ", " + rule.Config.Name
		}
	} else if match = commandregexBuycard.FindStringSubmatch(commandStr); len(match) == len(commandregexBuycard.SubexpNames()) {
		now := time.Now()
		endTime, _ := time.Parse("2006-01-02", r.Config.Robot.SellCardEndTime)
		if !r.Config.Robot.EnableSellCard || now.After(endTime) {
			return "目前尚未启用卖卡功能哦", nil, nil
		}

		qq := match[1]
		cardIndex := match[2]

		logger.Infof("开始调用卖卡脚本~")
		cmd := exec.Command("python", "sell_cards.py",
			"--run_remote",
			"--target_qq", qq,
			"--card_index", cardIndex,
		)
		cmd.Dir = "D:\\_codes\\Python\\djc_helper_public"
		out, err := cmd.Output()

		if err != nil {
			return "", nil, err
		}

		err = json.Unmarshal(out, &rspMsg)
		if err != nil {
			return "", nil, err
		}

		if strings.Contains(rspMsg, "成功发送以下卡片") {
			image, err := r._makeLocalImage("https://z3.ax1x.com/2020/12/16/r1yWZT.png")
			if err == nil {
				extraReplies = append(extraReplies, image)
			}
		}
	} else if match = commandregexQuerycard.FindStringSubmatch(commandStr); len(match) == len(commandregexQuerycard.SubexpNames()) {
		logger.Infof("开始查询卡片信息~")
		cmd := exec.Command("python", "sell_cards.py",
			"--run_remote",
			"--query",
		)
		cmd.Dir = "D:\\_codes\\Python\\djc_helper_public"
		out, err := cmd.Output()

		if err != nil {
			return "", nil, err
		}

		err = json.Unmarshal(out, &rspMsg)
		if err != nil {
			return "", nil, err
		}
	} else if match = commandRegexMusic.FindStringSubmatch(commandStr); len(match) == len(commandRegexMusic.SubexpNames()) {
		// full_match|听歌关键词|musicName
		musicName := match[2]

		musicElem, err := r.makeMusicShareElement(musicName, message.QQMusic)
		if err != nil {
			return fmt.Sprintf("没有找到歌曲：%v", musicName), nil, nil
		}

		rspMsg = fmt.Sprintf("请欣赏歌曲：%v", musicName)
		extraReplies = append(extraReplies, musicElem)
	} else if match = commandRegexQinglongChart.FindStringSubmatch(commandStr); len(match) == len(commandRegexQinglongChart.SubexpNames()) {
		// full_match|参数
		queryParam := match[1]

		cookieInfo := qinglong.QueryCookieInfo(queryParam)
		if cookieInfo == nil {
			return fmt.Sprintf("未找到相关cookie：%v，请使用 pt_pin 或者 备注 来进行查询", queryParam), nil, nil
		}

		chartPath := qinglong.QueryChartPath(cookieInfo)
		extraReplies = append(extraReplies, &msg.LocalImage{File: chartPath})
		extraReplies = append(extraReplies, message.NewText(cookieInfo.ToChatMessage()))
	} else if match = commandRegexQinglongSummary.FindStringSubmatch(commandStr); len(match) == len(commandRegexQinglongSummary.SubexpNames()) {
		// full_match|参数
		queryParam := match[1]

		cookieInfo := qinglong.QueryCookieInfo(queryParam)
		if cookieInfo == nil {
			return fmt.Sprintf("未找到相关cookie：%v，请使用 pt_pin 或者 备注 来进行查询", queryParam), nil, nil
		}

		summary := qinglong.QuerySummary(cookieInfo)
		extraReplies = append(extraReplies, message.NewText(summary))
		extraReplies = append(extraReplies, message.NewText(cookieInfo.ToChatMessage()))
	} else if match = commandRegexQinglongCookieExpired.FindStringSubmatch(commandStr); len(match) == len(commandRegexQinglongCookieExpired.SubexpNames()) {
		// full_match|参数
		queryParam := match[1]

		cookieInfo := qinglong.QueryCookieInfo(queryParam)
		if cookieInfo == nil {
			return fmt.Sprintf("未找到相关cookie：%v，请使用 pt_pin 或者 备注 来进行查询", queryParam), nil, nil
		}

		result := qinglong.QueryCookieExpired(cookieInfo)
		extraReplies = append(extraReplies, message.NewText(result))
		extraReplies = append(extraReplies, message.NewText(cookieInfo.ToChatMessage()))
	} else {
		return "", nil, errors.Errorf("没有找到该指令哦")
	}

	return rspMsg, extraReplies, nil
}
