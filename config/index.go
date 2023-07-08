package config

type Config struct {
	MaxConnections                  int
	TimeOutConnection               int
	DelayAfterMaxConnectionsReached int
	RequestTimeout                  int
	Uris                            []string
	ShowResults                     bool
	LogLevel                        int
	Src                             string
}
