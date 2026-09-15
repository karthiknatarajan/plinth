package version

var (
	ProjectName = "plinth"
	ServiceName = "nil"
	AppVersion  = "0.0.1"
	AppCommit   = "nil"
)

type Version struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}
