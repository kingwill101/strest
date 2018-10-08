package validators

import (
	"fmt"
	"net/http"
	"sync"

	strest "gitlab.com/kingwill101/strest"
)

type ValidaorFunc func(*strest.Request, *http.Response) bool

//Validator interface
type ValidatorInterface interface {
	Validate(*strest.Request, *http.Response) bool
}

//ValidatorMap register validators
type ValidatorMap map[string]interface{}

//Validator ff
type Validator struct {
	Validators     []ValidatorInterface
	ValidatorFuncs []ValidaorFunc
	sync.Mutex
}

func NewValidator() *Validator {
	v := &Validator{}
	return v
}

//Register add new validators
func (v *Validator) Register(val ValidatorInterface) {
	v.Lock()
	v.Validators = append(v.Validators, val)
	v.Unlock()
}

//Register add new validators
func (v *Validator) RegisterFunc(val ValidaorFunc) {
	v.Lock()
	v.ValidatorFuncs = append(v.ValidatorFuncs, val)
	v.Unlock()
}

func (v *Validator) Validate(sreq *strest.Request, r *http.Response) {
	fmt.Println("Validating plugins")
	for _, val := range v.Validators {
		v.Lock()
		val.Validate(sreq, r)
		v.Unlock()
	}
}

type LogPrint struct{}

func (lp *LogPrint) Validate(sreq *strest.Request, r *http.Response) bool {
	fmt.Println("response Status:", r.Status)
	fmt.Println("response Headers:", r.Header)
	body, _ := strest.ReadBody(r)
	fmt.Println("response Body:", string(body))
	return true
}

type LogIncoming struct{}

func (li *LogIncoming) Validate(sreq *strest.Request, r *http.Response) bool {
	fmt.Println("request form:", r.Request.Form)
	fmt.Println("response Headers:", r.Request.Header)
	body, _ := strest.ReadBody(r)
	fmt.Println("request Body:", string(body))
	return true
}
