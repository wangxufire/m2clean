package cleaner

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wangxufire/m2clean/args"
)

const (
	separator = string(os.PathSeparator)
	mb        = float32(1024 * 1024)
)

var (
	m2Repository                  string
	accesseBefore                 int64
	ignoreGroups, ignoreArtifacts []string
	processMap                    = make(map[string][]*m2FileInfo)
)

func initVar() {
	if args.Args.AccessedBefore == "" {
		accesseBefore = time.Now().AddDate(0, -3, 0).Unix()
	} else {
		ab, err := time.Parse("2006-01-02", args.Args.AccessedBefore)
		// ab, err := time.Parse("2006-01-02 15:04:05", args.Args.AccessedBefore+" 00:00:00")
		if err == nil {
			accesseBefore = ab.Unix()
		}
	}

	m2Repository = args.Args.M2Path
	if m2Repository == "" {
		current, err := user.Current()
		if err != nil {
			log.Panic(err)
		}
		m2Repository = current.HomeDir + separator + ".m2/repository"
	}

	ignoreGroups = args.Args.IgnoreGroups
	ignoreArtifacts = args.Args.IgnoreArtifacts
}

// Process files
func Process() {
	initVar()

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".pom" {
			return nil
		}
		if accesseBefore == 0 {
			return nil
		}
		accesseTime := accesseBefore
		if v := atime(info); v != 0 {
			accesseTime = v
		}
		if accesseTime < accesseBefore {
			m2f := &m2FileInfo{path: path, file: info}
			process(m2f)
		}
		return nil
	}

	if err := filepath.Walk(m2Repository, walk); err != nil {
		log.Panic(err)
	}

	cap := len(processMap)
	if cap == 0 {
		fmt.Println("************** No files were deleted as find nothing **************")
		return
	}

	keys := make([]string, 0, cap)
	for k := range processMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("***************************** Directories to be deleted *****************************")

	var filePaths []string
	var size float32
	for _, k := range keys {
		files := processMap[k]
		var paths []string
		for _, file := range files {
			dir := filepath.Dir(file.path)
			size += dirSize(dir)
			paths = append(paths, dir)
		}
		filePaths = append(filePaths, paths...)
		sort.Sort(sort.StringSlice(paths))
		fmt.Printf("\n%s", strings.Join(paths, "\n"))
	}
	fmt.Printf("\n\n******************************* Total size %fM *******************************\n", size)
	if args.Args.Dryrun {
		fmt.Println("************** No files were deleted as program was run in DRY-RUN mode **************")
	} else {
		for _, path := range filePaths {
			removeDir(path)
		}
		fmt.Println("***************************** All files are deleted done *****************************")
	}
}

func dirSize(path string) float32 {
	var size int64
	walk := func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Mode().IsRegular() {
			size += info.Size()
		}
		return nil
	}
	filepath.Walk(path, walk)
	return float32(size) / mb
}

func removeDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func process(fileInfo *m2FileInfo) {
	if info := resolveM2File(fileInfo); info != nil {
		id := info.groupID + ":" + info.artifactID
		files := processMap[id]
		files = append(files, info)
		processMap[id] = files
	}
}

func resolveGroupID(artifactID, path string) string {
	i1 := strings.Index(path, "repository") + 11
	i2 := strings.LastIndex(path, artifactID) - 1
	return strings.Join(strings.Split(path[i1:i2], separator), ".")
}

func resolveArtifactID(path string) string {
	s := filepath.Dir(filepath.Dir(path))
	return s[strings.LastIndex(s, separator)+1:]
}

func resolveM2File(fileInfo *m2FileInfo) *m2FileInfo {
	fileInfo.composeM2FileInfo()
	if mustIgnore(fileInfo) {
		return nil
	}
	return fileInfo
}

func mustIgnore(info *m2FileInfo) bool {
	var contains = func(elements []string, element string) bool {
		for _, v := range elements {
			if v == element {
				return true
			}
		}
		return false
	}

	if contains(ignoreGroups, info.groupID) {
		return true
	}

	if contains(ignoreArtifacts, info.artifactID) {
		return true
	}

	return false
}

func (m2f *m2FileInfo) composeM2FileInfo() {
	path := m2f.path
	artifactID := resolveArtifactID(path)
	groupID := resolveGroupID(artifactID, path)
	m2f.artifactID = artifactID
	m2f.groupID = groupID
	m2f.version = filepath.Dir(path)
}

type m2FileInfo struct {
	path                         string
	file                         os.FileInfo
	groupID, artifactID, version string
}
