package monitor

const (
	NodeStatusNotReachable = 0
	NodeStatusReachable    = 1
)

type LcdNodeInfoResp struct {
	NodeInfo struct {
		Network string `json:"network"`
	} `json:"default_node_info"`
}
