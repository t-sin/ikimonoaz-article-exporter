package save

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/xerrors"

	"ikimonoaz-exporter/userdata"
)

type media struct {
	Name string
	Type string
	Body []byte
}

func downloadMedia(url string) (*media, error) {
	matches := userdata.MediaUrlPat.FindStringSubmatch(url)
	if len(matches) != 4 {
		return nil, xerrors.Errorf("invalid media URL: '%s'", url)
	}

	name := fmt.Sprintf("%s_%s", matches[1], matches[2])
	media := media{Name: name, Type: matches[3]}

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

	path := fmt.Sprintf("%s%cmedia%c%s.%s", dir, os.PathSeparator, os.PathSeparator, media.Name, media.Type)
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
