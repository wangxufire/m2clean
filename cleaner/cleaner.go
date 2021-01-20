package cleaner

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	options "github.com/wangxufire/m2clean/flag"
)

const (
	separator = string(os.PathSeparator)
)

const (
	accesse = iota
	common
	snapshot
	source
	javadoc
)

var (
	m2Repository                  string
	accesseBefore                 int64
	ignoreGroups, ignoreArtifacts []string
	reasonMap                     map[int]m2FileInfo
	processMap                    = make(map[string][]*m2FileInfo)
)

func init() {
	ab, err := time.Parse("2006-01-02", options.Args.AccessedBefore)
	if err != nil {
		log.Panic(err)
	}
	accesseBefore = ab.Unix()

	m2Repository = options.Args.M2Path
	if m2Repository == "" {
		current, err := user.Current()
		if err != nil {
			log.Panic(err)
		}
		m2Repository = current.HomeDir + separator + ".m2/repository"
	}

	ignoreGroups = options.Args.IgnoreGroups
	ignoreArtifacts = options.Args.IgnoreArtifacts
}

// Process files
func Process() {
	err := filepath.Walk(m2Repository,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) != ".pom" {
				return nil
			}
			// if info.ModTime().Unix() > accesseBefore {
			// 	return nil
			// }
			m2f := &m2FileInfo{path: path, file: info}
			process(m2f)
			return nil
		})
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", processMap)
}

func process(fileInfo *m2FileInfo) {
	info := resolveM2File(fileInfo)
	if info == nil {
		return
	}
	id := info.groupID + ":" + info.artifactID
	files := processMap[id]
	files = append(files, info)
	processMap[id] = files
}

func resolveGroupID(artifactID, path string) string {
	//  ~/.m2/repository/org/springframework/boot/spring-boot-starter-jdbc/2.3.4.RELEASE
	i1 := strings.Index(path, "repository") + 11
	i2 := strings.LastIndex(path, artifactID) - 1
	return strings.Join(strings.Split(path[i1:i2], separator), ".")
}

func resolveArtifactID(path string) string {
	s := filepath.Dir(filepath.Dir(path))
	return s[strings.LastIndex(s, string(os.PathSeparator))+1:]
}

func resolveM2File(fileInfo *m2FileInfo) *m2FileInfo {
	fileInfo.composeM2FileInfo()
	if mustIgnore(fileInfo) {
		return nil
	}
	// reason := checkDeleteReason(info)
	// reasonMap[reason] = info
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

func checkDeleteReason(info m2FileInfo) int {
	name := info.file.Name()

	if strings.Index(name, "-javadoc.jar") > 0 {
		return javadoc
	}
	if strings.Index(name, "-sources.jar") > 0 {
		return source
	}
	if strings.Index(name, "-SNAPSHOT.jar") > 0 {
		return snapshot
	}
	//             if (attributes.lastAccessTime().toMillis() < argData.getAccessedBefore()) {
	//                 DELETE_MAP.get(DeleteReason.ACCESS_DATE).add(file.getParentFile());
	//                 return;
	//             }

	return common
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
