package repository

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kihamo/shadow/components/ota"
	"github.com/kihamo/shadow/components/ota/release"
)

type ShadowRecord struct {
	Architecture string `json:"architecture"`
	Checksum     string `json:"checksum"`
	Size         int64  `json:"size"`
	Version      string `json:"version"`
	File         string `json:"file"`
}

type Shadow struct {
	u *url.URL
}

func NewShadow(u *url.URL) *Shadow {
	return &Shadow{
		u: u,
	}
}

func (r *Shadow) Releases(arch string) ([]ota.Release, error) {
	r.u.Query().Set("architecture", arch)

	response, err := http.Get(r.u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var records []ShadowRecord

	if err := json.Unmarshal(body, &records); err != nil {
		return nil, err
	}

	releases := make([]ota.Release, 0, len(records))
	for _, record := range records {
		cs, err := hex.DecodeString(record.Checksum)
		if err != nil {
			return nil, err
		}

		rl, err := release.NewHTTPFile(record.File, record.Version, cs, record.Size, record.Architecture)
		if err != nil {
			return nil, err
		}

		releases = append(releases, release.NewCompress(rl))
	}

	return releases, nil
}
