package types

// Deployment options
type DeployOptions struct {
	Appname string // Name of application
	Image   string // Name of image to deploy

	// Target options
	Podname    string            // Name of pod to deploy
	Nodename   string            // Specific nodes to deploy, if given, must belong to pod
	Entrypoint string            // Entrypoint to deploy
	CPUQuota   float64           // How many cores needed, e.g. 1.5
	Count      int               // How many containers needed, e.g. 4
	Memory     int64             // Memory for container, in bytes
	Env        []string          // Env for container
	Networks   map[string]string // Network names and specified IPs
	Raw        bool              // If use raw, launcher won't be used
}
