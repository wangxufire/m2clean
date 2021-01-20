package options

// Option from command line
type Option struct {
	M2Path          string   `short:"p" long:"path" description:"Path to m2 directory, if using a custom path"`
	AccessedBefore  string   `short:"a" long:"accessed-before" description:"Delete all libraries (even if latest version) last accessed on or before this date (2006-01-02)."`
	IgnoreArtifacts []string `short:"s" long:"ignore-artifacts" description:"artifactIds (full or part) to be ignored."`
	IgnoreGroups    []string `short:"g" long:"ignore-groups" description:"groupIds (full or part) to be ignored."`
	Dryrun          bool     `short:"r" long:"dryrun" description:"Do not delete files, just simulate and print result."`
}