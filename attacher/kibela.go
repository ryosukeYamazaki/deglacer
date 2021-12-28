package attacher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/Songmu/kibelasync/kibela"
	"github.com/slack-go/slack"
)

const (
	DomainSuffix = ".kibe.la"
)

var (
	noteReg     = regexp.MustCompile(`^/(?:@[^/]+|notes)/([0-9]+)`)
	fragmentReg = regexp.MustCompile(`(?i)^comment_([0-9]+)`)
	spacesReg   = regexp.MustCompile(`\s+`)
)

type KibelaAttacher struct {
	kibelaCli *kibela.Kibela
}

func New(kibelaCli *kibela.Kibela) (KibelaAttacher, error) {
	attacher := KibelaAttacher{}
	if kibelaCli == nil {
		e := errors.New("kibelaCliがnilです。")
		return attacher, e
	}
	attacher.kibelaCli = kibelaCli
	return attacher, nil
}

func (a KibelaAttacher) SlackAttachment(ctx context.Context, url *url.URL, team string) (slack.Attachment, error) {
	attachment := slack.Attachment{}
	m := noteReg.FindStringSubmatch(url.Path)
	if len(m) < 2 {
		e := errors.New("kibelaのURLにマッチしません。")
		return attachment, e
	}
	id, _ := strconv.Atoi(m[1])

	note, err := a.kibelaCli.GetNote(ctx, id)
	if err != nil {
		return attachment, err
	}
	var (
		author      = note.Author
		title       = note.Title
		text        = note.Summary
		publishedAt = note.PublishedAt
	)
	if m := fragmentReg.FindStringSubmatch(url.Fragment); len(m) > 1 {
		id, _ := strconv.Atoi(m[1])
		comment, err := a.kibelaCli.GetComment(ctx, id)
		if err != nil {
			return attachment, err
		}
		author = comment.Author
		title = fmt.Sprintf(`comment for "%s"`, title)
		text = comment.Summary
		publishedAt = comment.PublishedAt
	}

	attachment.AuthorLink = fmt.Sprintf("https://%s.kibe.la/@%s", team, author.Account)
	attachment.AuthorName = author.Account
	attachment.Title = title
	attachment.TitleLink = url.Path
	attachment.Text = spacesReg.ReplaceAllString(text, " ")
	attachment.Footer = "Kibela"
	attachment.FooterIcon = "https://cdn.kibe.la/assets/shortcut_icon-99b5d6891a0a53624ab74ef26a28079e37c4f953af6ea62396f060d3916df061.png"
	attachment.Ts = json.Number(fmt.Sprintf("%d", publishedAt.Time.Unix()))

	return attachment, nil
}
