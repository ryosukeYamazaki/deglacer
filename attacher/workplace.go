package attacher

import (
	"errors"
	"log"
	"net/url"
	"strings"

	fb "github.com/huandu/facebook/v2"
	"github.com/slack-go/slack"
)

const (
	WorkplaceSuffix = ".workplace.com"
)

type WorkplaceAttacher struct {
	workplaceSession *fb.Session
}

func NewWorkplace(wpSession *fb.Session) (WorkplaceAttacher, error) {
	attacher := WorkplaceAttacher{}
	if wpSession == nil {
		e := errors.New("wpSessionがnilです。")
		return attacher, e
	}
	attacher.workplaceSession = wpSession
	return attacher, nil
}

func (a WorkplaceAttacher) SlackAttachment(url *url.URL) (slack.Attachment, error) {
	attachment := slack.Attachment{}
	permalink := "/groups/1779121192119476/permalink/3594996737198570/"
	if url.Path != permalink {
		e := errors.New("指定以外のURLです。")
		log.Printf("url.Path: %s \n.", url.Path)
		return attachment, e
	}
	postId := postIDFromPath(url.Path)
	res, _ := a.workplaceSession.Get(postId, fb.Params{"fields": "message"})
	message := res["message"]
	if message == nil {
		e := errors.New("messageがnilです。")
		return attachment, e
	}
	attachment.Text = message.(string)
	return attachment, nil
}

func postIDFromPath(workplacePath string) string {
	splitedPath := strings.Split(workplacePath, "/")
	// TODO error 処理
	return splitedPath[4]
}
