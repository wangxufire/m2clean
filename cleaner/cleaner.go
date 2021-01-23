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
		if ab, err := time.Parse("2006-01-02", args.Args.AccessedBefore); err == nil {
			accesseBefore = ab.Unix()
		}
	}

	m2Repository = args.Args.M2Path
	if m2Repository == "" {
		current, err := user.Current()
		if err != nil {
			log.Panic(err)
		}
		m2Repository = current.HomeDir + separator + ".m2" + separator + "repository"
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

	deleteCh := make(chan []string, 5)
	deleteDone := make(chan bool)
	go func() {
		for paths := range deleteCh {
			for _, path := range paths {
				removeDir(path)
			}
		}
		deleteDone <- true
	}()

	var size float32
	dryrun := args.Args.Dryrun
	sizeDone := make(chan bool)
	for _, k := range keys {
		sizeCh := make(chan string, 5)
		go func(ch <-chan string) {
			for path := range ch {
				size += dirSize(path)
			}
			sizeDone <- true
		}(sizeCh)

		var paths []string
		for _, file := range processMap[k] {
			dir := filepath.Dir(file.path)
			sizeCh <- dir
			paths = append(paths, dir)
		}
		close(sizeCh)

		sort.Sort(sort.StringSlice(paths))
		fmt.Printf("\n%s", strings.Join(paths, "\n"))

		if <-sizeDone && !dryrun {
			deleteCh <- paths
		}
	}

	fmt.Printf("\n\n******************************* Total size %fM *******************************\n", size)

	close(deleteCh)
	if dryrun {
		fmt.Println("************* No files were deleted as program was run in DRY-RUN mode **************")
	} else {
		<-deleteDone
		fmt.Println("**************************** All files are deleted done *****************************")
	}
	close(sizeDone)
	close(deleteDone)
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
		if err = os.RemoveAll(filepath.Join(dir, name)); err != nil {
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
	if fileInfo.mustIgnore() {
		return nil
	}
	return fileInfo
}

func (info *m2FileInfo) mustIgnore() bool {
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

func (info *m2FileInfo) composeM2FileInfo() {
	path := info.path
	artifactID := resolveArtifactID(path)
	info.artifactID = artifactID
	info.groupID = resolveGroupID(artifactID, path)
	info.version = filepath.Dir(path)
}

type m2FileInfo struct {
	path                         string
	file                         os.FileInfo
	artifactID, groupID, version string
}
