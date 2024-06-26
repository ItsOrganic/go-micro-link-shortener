package shortener
import (
	"errors"
	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
    "time"
)
var (
    ErrRedirectNotFound = errors.New("Redirect link not found")
    ErrRedirectInvalid = errors.New("Redirect link is invalid")
)

type redirectService struct {
    redirectRepo RedirectRepo
}
func NewRedirectService(redirectRepo RedirectRepo) RedirectService {
    return &redirectService{
        redirectRepo,
    }
}

func (r *redirectService) Find(code string) (*Redirect, error) {
    return r.redirectRepo.Find(code)

}
func (r *redirectService) Store(redirect *Redirect) error{
    if err := validate.Validate(redirect); err != nil {
        return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
    }

    redirect.Code = shortid.MustGenerate()
    redirect.CreatedAt = time.Now().UTC().Unix()
    return r.redirectRepo.Store(redirect)

}
