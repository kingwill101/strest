package validators

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"

	"gitlab.com/kingwill101/strest"
)

type ValidaorFunc func(*strest.Request, *http.Response) bool

// ValidatorInterface Validator interface
type ValidatorInterface interface {
	Validate(*strest.Request, *http.Response, *logrus.Entry) bool
}

// ValidatorMap register validators
type ValidatorMap map[string]interface{}

// Validator ff
type Validator struct {
	Validators     []ValidatorInterface
	ValidatorFuncs []ValidaorFunc
	sync.Mutex
}

func NewValidator() *Validator {
	v := &Validator{}
	return v
}

// Register add new validators
func (v *Validator) Register(val ValidatorInterface) {
	v.Lock()
	v.Validators = append(v.Validators, val)
	v.Unlock()
}

// RegisterFunc Register add new validators
func (v *Validator) RegisterFunc(val ValidaorFunc) {
	v.Lock()
	v.ValidatorFuncs = append(v.ValidatorFuncs, val)
	v.Unlock()
}

func (v *Validator) Validate(sreq *strest.Request, r *http.Response, logger *logrus.Entry) {
	// fmt.Println("Validating plugins")
	for _, val := range v.Validators {
		v.Lock()
		val.Validate(sreq, r, logger)
		v.Unlock()
	}
}

type LogPrint struct{}

func (lp *LogPrint) Validate(_ *strest.Request, r *http.Response, logger *logrus.Entry) bool {
	entry := logger.WithField("validator", "log-print")

	entry.Infof("response Status: %s", r.Status)
	entry.Infof("response Headers: %s", r.Header)
	body, _ := strest.ReadBody(r)
	entry.Infof("response Body: %s", strings.TrimRight(body, "\n"))
	return true
}

type LogIncoming struct{}

func (li *LogIncoming) Validate(_ *strest.Request, r *http.Response, logger *logrus.Entry) bool {
	entry := logger.WithField("validator", "log-incoming")

	entry.Info("request form:", r.Request.Form)
	entry.Info("response Headers:", r.Request.Header)
	body, _ := strest.ReadBody(r)
	entry.Infof("request Body22: %s", strings.TrimRight(body, "\n"))
	return true
}
