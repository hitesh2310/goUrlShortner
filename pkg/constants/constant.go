package constants

import (
	models "main/pkg/models/configModel"
	"sync"
)

var ApplicationConfig *models.Config

var Counter int

var GlobalMutex sync.Mutex

const Charset string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
