package response

import "gmserver/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
