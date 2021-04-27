package gui

import (
	"fmt"
	"math/rand"
	"time"
)

var messages = []string{
	"バクが頭を振ってダッシュしています",
	"バクが木をかじっています",
	"ニホンリスがクルミをかじっています",
	"ニホンリスが地面に寝そべっています",
}

func chooseMessage() string {
	idx := rand.Intn(len(messages))
	return messages[idx]
}

func generateDots(n int) string {
	dots := ""
	num := n % 4
	for i := 0; i < num; i++ {
		dots = dots + "."
	}
	return dots
}

func calculateStatus(s *state, c <-chan bool) {
	t := time.NewTicker(1 * time.Second)

	s.status.Set("エクスポート中です")

	for {
		select {
		case <-t.C:
			if s.dotCount == 3 {
				s.dotCount = 0
			} else {
				s.dotCount += 1
			}

			msg := fmt.Sprintf("エクスポート中です%s", generateDots(s.dotCount))
			s.status.Set(msg)

		case <-c:
			s.status.Set("エクスポート完了です！！！")
			t.Stop()
			return
		}
	}
}
