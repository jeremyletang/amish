package slack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeremyletang/amish/conf"
	"github.com/jeremyletang/amish/domain"
	"github.com/nlopes/slack"
)

var (
	slackClient *slack.Client
)

type NotificationType int

const (
	StarMessage   = "Hey ! https://github.com/%s just stared %s !"
	UnStarMessage = "Wtf ! https://github.com/%s just remove his star on %s, go kick his ass !"
)

const (
	Star = iota
	UnStar
)

var (
	msg = map[NotificationType]string{
		Star:   StarMessage,
		UnStar: UnStarMessage,
	}
)

func InitSlack(conf conf.Conf) {
	if conf.Slack.Token != "" {
		slackClient = slack.New(conf.Slack.Token)
		if slackClient == nil {
			panic("cannot initialize slack client")
		}
	}
}

func Notify(ty NotificationType, users []*domain.User, repo *domain.Repository, channels []string) {
	for _, u := range users {
		att := slack.Attachment{
			AuthorIcon: "https://slack-imgs.com/?c=1&o1=wi16.he16&url=https%3A%2F%2Fgithub.com%2Fapple-touch-icon.png",
			AuthorName: "Github",
			Title:      fmt.Sprintf("%s (%s)", u.Login, u.Name),
			TitleLink:  fmt.Sprintf("https://github.com/%s", u.Login),
			Footer:     "Amish",
			Ts:         json.Number(fmt.Sprintf("%v", time.Now().Unix())),
			Pretext:    fmt.Sprintf(msg[ty], u.Login, repo.Owner+"/"+repo.Name),
			ThumbURL:   u.AvatarUrl,
		}
		for _, c := range channels {
			// fmt.Println(fmt.Sprintf(msg[ty], u.Login, repo.Owner+"/"+repo.Name))
			_, _, err := slackClient.PostMessage(
				c,
				"",
				slack.PostMessageParameters{
					EscapeText:  true,
					Username:    "amish",
					AsUser:      true,
					IconURL:     "https://ok-borg.slack.com/team/balek",
					Attachments: []slack.Attachment{att},
				},
			)
			if err != nil {
				fmt.Println("error:", err.Error())
			}
		}
	}
}
