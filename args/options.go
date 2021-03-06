package args

// Option from command line
type Option struct {
	M2Path          string   `short:"p" long:"path" description:"Path to m2 directory, if using a custom path. Default is homedir/.m2/repository"`
	AccessedBefore  string   `short:"a" long:"accessed-before" description:"Delete all libraries last accessed on or before this date (2006-01-02). Default 3 month ago."`
	IgnoreArtifacts []string `short:"f" long:"ignore-artifacts" description:"ArtifactIds to be ignored."`
	IgnoreGroups    []string `short:"g" long:"ignore-groups" description:"GroupIds to be ignored."`
	Dryrun          bool     `short:"d" long:"dryrun" description:"Do not delete files, just simulate and print result."`
}
