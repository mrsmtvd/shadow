package repository

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kihamo/shadow/components/ota"
	"github.com/kihamo/shadow/components/ota/release"
)

type ShadowRepositoryRecord struct {
	Architecture string `json:"architecture"`
	Checksum     string `json:"checksum"`
	Size         int64  `json:"size"`
	Version      string `json:"version"`
	File         string `json:"file"`
}

type ShadowRepository struct {
	u *url.URL
}

func NewShadowRepository(u *url.URL) *ShadowRepository {
	return &ShadowRepository{
		u: u,
	}
}

func (r *ShadowRepository) Releases(arch string) ([]ota.Release, error) {
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

	var records []ShadowRepositoryRecord

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

		releases = append(releases, rl)
	}

	return releases, nil
}

func (r *ShadowRepository) ReleaseLatest(arch string) (ota.Release, error) {
	releases, err := r.Releases(arch)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, errors.New("latest release not found")
	}

	return releases[len(releases)-1], nil
}
