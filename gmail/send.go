package gmail

import (
	"encoding/base64"

	log "github.com/cihub/seelog"
	"google.golang.org/api/gmail/v1"
)

func SendMail() {
	from := "amish@gmail.com"
	to := "letang.jeremy@gmail.com"
	body := "hello world"

	msgRaw := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: Am I Still Hype daily report here !!!" + "\r\n\r\n" +
		body + "\r\n"

	msg := &gmail.Message{
		Raw: base64.StdEncoding.EncodeToString([]byte(msgRaw)),
	}
	call := service.Users.Messages.Send("me", msg)
	_, err := call.Do()
	if err != nil {
		log.Errorf("unable to send mail: %s", err.Error())
	}

}
