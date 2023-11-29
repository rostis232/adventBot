package templatedata

import (
	"github.com/rostis232/adventBot/config"
	"time"
)

type TemplateData struct {
	Config config.Config
	StringMap map[string]string
	IntMap map[string]int
	FloatMap map[string]float32
	Data map[string]any
	Flash string
	Warning string
	Error string
	Authenticated bool
	Now time.Time
}

type Keys struct {
	SkID int
	SecretKey int
	Link string
	ChatID int
}