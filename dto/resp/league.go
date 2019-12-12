package resp

type CheckJoinResp struct {
	Allow bool `json:"allow"`
}

type LeagueItem struct {
	ChainId string `json:"chain_id"`
	Name string `json:"name"`
	ServerIp string `json:"server_ip"`
	ServerPort uint32 `json:"server_port"`
	Desc string `json:"desc"`
}

func NewLeagueItem(chainId, name, serverIp, desc string, port uint32) *LeagueItem {
	return &LeagueItem{
		ChainId:chainId,
		ServerIp:serverIp,
		Name:name,
		Desc:desc,
		ServerPort:port,
	}
}