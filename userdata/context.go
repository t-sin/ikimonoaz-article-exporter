package userdata

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func (m Media) ToMap() interface{} {
	return map[string]string{
		"url":       m.URL,
		"thumbnail": m.Thumbnail,
	}
}

func (c Creature) ToMap() interface{} {
	return map[string]interface{}{
		"name":  c.Name,
		"place": c.Place.Name,
	}
}

func (t Tag) ToMap() interface{} {
	return map[string]string{
		"name": t.Name,
	}
}

func (c Comment) ToMap() interface{} {
	pat := regexp.MustCompile(`\n`)
	comment := pat.ReplaceAllString(c.Comment, "<br />")
	return map[string]string{
		"comment":    comment,
		"username":   c.User.Name,
		"created_at": c.CreatedAt.Format(time.RFC3339),
	}
}

func (a Article) ToMap() interface{} {
	media := make([]interface{}, len(a.MediaList))
	for i, m := range a.MediaList {
		media[i] = m.ToMap()
	}

	creatures := make([]interface{}, len(a.CreatureList))
	for i, m := range a.CreatureList {
		creatures[i] = m.ToMap()
	}

	tags := make([]string, len(a.Tags))
	for i, t := range a.Tags {
		tags[i] = t.Name
	}

	comments := make([]interface{}, len(a.CommentList))
	for i, c := range a.CommentList {
		comments[i] = c.ToMap()
	}

	title := a.Title
	if title == "" {
		title = "(無題)"
	}

	contents := a.Contents
	for _, m := range a.MediaList {
		matches := MediaUrlPat.FindStringSubmatch(m.URL)
		mediaPath := fmt.Sprintf("../media/%s_%s.%s", matches[1], matches[2], matches[3])
		contents = strings.ReplaceAll(contents, m.URL, mediaPath)
	}

	return map[string]interface{}{
		"id":          a.ID,
		"created_at":  a.CreatedAt.Format(time.RFC3339),
		"updated_at":  a.UpdatedAt.Format(time.RFC3339),
		"released_at": a.ReleasedAt.Format(time.RFC3339),
		"title":       title,
		"contents":    contents,
		"media":       media,
		"creatures":   creatures,
		"tags":        tags,
		"comments":    comments,
	}
}

func (u User) ToMap() interface{} {
	return map[string]interface{}{
		"name":              u.Name,
		"profile":           u.Profile,
		"profile_image_url": u.ProfileImageURL,
		"meister":           u.MeisterList,
		"place":             u.PlaceName,
	}
}

func (ud UserData) ToMap() interface{} {
	articles := make([]interface{}, len(ud.Articles))
	for i, a := range ud.Articles {
		articles[i] = a.ToMap()
	}

	return map[string]interface{}{
		"user":     ud.User.ToMap(),
		"articles": articles,
	}
}
