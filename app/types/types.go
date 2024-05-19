package types

type DBItem struct {
	Value  string
	Expiry int64
}

type ServerState struct {
	DB               map[string]DBItem
	Role             string
	MasterReplID     string
	MasterReplOffset int
	MasterHost       string
	MasterPort       string
}
