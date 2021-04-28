package userdata

import (
	"regexp"
	"time"
)

// 記事中の写真・動画。動画の場合はサムネイル画像がある。
type Media struct {
	URL       string `json:"url"`
	Thumbnail string `json:"movie_thumbnail"`
}

var MediaUrlPat *regexp.Regexp = regexp.MustCompile(`^.+/([^/]+?)/([^.]+?).([^.]+)$`)

type Place struct {
	Name string `json:"name"`
}

// 記事のいきもの情報。名前といる場所（園館）。
type Creature struct {
	Name  string `json:"name"`
	Place Place  `json:"place"`
}

// 記事のタグ情報。
type Tag struct {
	Name string `json:"name"`
}

type CommentUser struct {
	Name string `json:"user_name"`
}

// 記事のコメントデータ。
type Comment struct {
	Comment   string      `json:"comment"`
	User      CommentUser `json:"user"`
	CreatedAt time.Time   `json:"create_date"`
}

// 記事データ。
type Article struct {
	ID           int        `json:"article_id"`
	CreatedAt    time.Time  `json:"create_date"`
	UpdatedAt    time.Time  `json:"edit_date"`
	ReleasedAt   time.Time  `json:"release_date"`
	Title        string     `json:"title"`
	Contents     string     `json:"contents"`
	MediaList    []Media    `json:"media"`
	CreatureList []Creature `json:"creatures"`
	Tags         []Tag      `json:"tags"`
	CommentList  []Comment  `json:"comments"`
}

// いきものマイスター
type Meister struct {
	Name string `json:"name"`
}

// ユーザ情報。
type User struct {
	Name            string    `json:"user_name"`
	Profile         string    `json:"profile"`
	ProfileImageURL string    `json:"profile_image_url"`
	MeisterList     []Meister `json:"meister"`    // マイスター一覧
	PlaceName       string    `json:"place_name"` // よく行く園館
}

// このプログラムで扱うすべてのユーザ情報。
type UserData struct {
	User     User      `json:"user"`
	Articles []Article `json:"articles"`
}
