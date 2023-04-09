package mapper

import (
	"fmt"
	"strconv"
)

func ModelIdToRepoId(ID string) (int64, error) {
	repoID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse id: %s", ID)
	}
	return repoID, nil
}

func RepoIdToModelId(ID int64) string {
	return strconv.FormatInt(ID, 10)
}
