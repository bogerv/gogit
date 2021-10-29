package cmds

type Cmd struct {
	Path          string `json:"path"`
	CommitId      string `json:"commitId"`
	CurrentBranch string `json:"currentBranch"`
}
