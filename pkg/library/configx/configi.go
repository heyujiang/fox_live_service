package configx

type ConfigI interface {
	Get(string) interface{}
	GetInt(string) int
	GetInt32(string) int32
	GetInt64(string) int64
	GetString(string) string
	GetFloat64(string) float64
	GetIntSlice(string) []int
	GetStringSlice(string) []string
	GetStringMapString(string) map[string]string
	GetStringMapStringSlice(string) map[string][]string
	GetBool(string) bool
	WatchConfig()
}
