package modules

import (
	"net/http"

	"github.com/github/dependabot/go/common/dependabot/v1"
)

func NewUpdaterService(c *Container) dependabot_v1.UpdateService {
	return dependabot_v1.NewUpdateServiceJSONClient(c.apiAddr, http.DefaultClient)
}
