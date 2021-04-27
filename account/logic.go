package account

import (
	"github.com/go-kit/kit/log"
)

type service struct {
	repository Repository
	logger     log.Logger
}
