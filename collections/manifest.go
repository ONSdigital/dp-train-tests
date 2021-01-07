package collections

type FileCopy struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type FileDelete struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Manifest struct {
	FilesToCopy  []*FileCopy   `json:"filesToCopy"`
	UrisToDelete []*FileDelete `json:"uriToDelete"`
}
