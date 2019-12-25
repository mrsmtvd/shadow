package repository

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/kihamo/shadow/components/ota/release"
)

type ShadowRecord struct {
	Architecture string     `json:"architecture"`
	Checksum     string     `json:"checksum"`
	Size         int64      `json:"size"`
	Version      string     `json:"version"`
	File         string     `json:"file"`
	CreatedAt    *time.Time `json:"created_at"`
}

type Shadow struct {
	*Memory

	u *url.URL
}

func NewShadow(u *url.URL) *Shadow {
	return &Shadow{
		Memory: NewMemory(),
		u:      u,
	}
}

func (r *Shadow) Update() error {
	response, err := http.Get(r.u.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var records []ShadowRecord

	if err := json.Unmarshal(body, &records); err != nil {
		return err
	}

	r.Memory.Clean()

	for _, record := range records {
		cs, err := hex.DecodeString(record.Checksum)
		if err != nil {
			return err
		}

		rl, err := release.NewHTTPFile(record.File, record.Version, cs, record.Size, record.Architecture, record.CreatedAt)
		if err != nil {
			return err
		}

		r.Memory.Add(release.NewCompress(rl))
	}

	return nil
}
