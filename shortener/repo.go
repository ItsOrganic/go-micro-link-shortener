package shortener

type RedirectRepo interface{
    Find(code string) (*Redirect, error)
    Store(redirect *Redirect) error
}
