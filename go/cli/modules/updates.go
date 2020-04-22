package modules

import (
	"net/http"

	"github.com/thepwagner/dependagot/go/common/dependagot/v1"
)

func NewUpdaterService(c *Container) dependagot_v1.UpdateService {
	return dependagot_v1.NewUpdateServiceProtobufClient(c.apiAddr, http.DefaultClient)
}
