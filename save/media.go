package save

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/xerrors"
)

type media struct {
	Name string
	Type string
	Body []byte
}

var pat *regexp.Regexp = regexp.MustCompile(`^.+/([^.]+).(.+)$`)

func downloadMedia(url string) (*media, error) {
	matches := pat.FindStringSubmatch(url)
	if len(matches) != 3 {
		return nil, xerrors.Errorf("invalid media URL: '%s'", url)
	} 

	media := media{Name: matches[1], Type: matches[2]}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, xerrors.New("Status code is not 200")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	media.Body = body

	return &media, nil
}

func saveMedia(dir string, url string) error {
	media, err := downloadMedia(url)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%smedia/%s.%s", dir, media.Name, media.Type)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Seek(0, 0)
	if _, err := f.Write(media.Body); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
