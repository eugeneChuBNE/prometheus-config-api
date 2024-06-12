package main

// ScrapeConfig represents a scrape configuration in the Prometheus configuration.
type ScrapeConfig struct {
	JobName        string              `yaml:"job_name"`
	StaticConfigs  []StaticConfig      `yaml:"static_configs"`
	ScrapeInterval string              `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout  string              `yaml:"scrape_timeout,omitempty"`
	MetricsPath    string              `yaml:"metrics_path,omitempty"`
	Params         map[string][]string `yaml:"params,omitempty"`
	RelabelConfigs []RelabelConfig     `yaml:"relabel_configs,omitempty"`
}

// StaticConfig represents the static configuration for targets in a scrape configuration.
type StaticConfig struct {
	Targets []string `yaml:"targets"`
}

// RelabelConfig represents the relabel configuration in a scrape configuration.
type RelabelConfig struct {
	SourceLabels []string `yaml:"source_labels,omitempty"`
	TargetLabel  string   `yaml:"target_label,omitempty"`
	Replacement  string   `yaml:"replacement,omitempty"`
}

// Job represents the simplified structure to be returned by the API, containing only the job name and static configurations.
type Job struct {
	JobName       string         `json:"job_name"`
	StaticConfigs []StaticConfig `json:"static_configs"`
}

// AddJobRequest represents the request body for adding a new job.
type AddJobRequest struct {
	JobName   string `json:"job_name"`
	IPAddress string `json:"ip_address"`
}
