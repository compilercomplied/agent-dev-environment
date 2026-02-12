package ls

import (
	"agent-dev-environment/src/api/v1/filesystem/ls"
	"agent-dev-environment/src/library/api"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func Handler(req ls.Request) (*ls.Response, error) {
	fileInfo, err := os.Stat(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "Path not found")
		}
		return nil, err
	}

	var entries []ls.FileInfo

	if !fileInfo.IsDir() {
		entries = append(entries, fileInfoToAPIFileInfo(fileInfo, req.Long))
		return &ls.Response{Entries: entries}, nil
	}

	if req.Recursive {
		err = filepath.WalkDir(req.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			info, err := d.Info()
			if err != nil {
				return err
			}

			entries = append(entries, fileInfoToAPIFileInfo(info, req.Long))
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		dirEntries, err := os.ReadDir(req.Path)
		if err != nil {
			return nil, err
		}

		for _, dirEntry := range dirEntries {
			info, err := dirEntry.Info()
			if err != nil {
				return nil, err
			}
			entries = append(entries, fileInfoToAPIFileInfo(info, req.Long))
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	return &ls.Response{Entries: entries}, nil
}

func fileInfoToAPIFileInfo(info os.FileInfo, includeLongFormat bool) ls.FileInfo {
	fileInfo := ls.FileInfo{
		Name:        info.Name(),
		IsDirectory: info.IsDir(),
		Size:        info.Size(),
	}

	if includeLongFormat {
		fileInfo.Mode = info.Mode().String()
		fileInfo.ModTime = info.ModTime()
	}

	return fileInfo
}
