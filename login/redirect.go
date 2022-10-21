package login

import (
	"net/http"
	"net/url"

	"github.com/pchchv/logsrv/logging"
)

func (h *Handler) setRedirectCookie(w http.ResponseWriter, r *http.Request) {
	redirectTo := r.URL.Query().Get(h.config.RedirectQueryParameter)
	if redirectTo != "" && h.allowRedirect(r) && r.Method != "POST" {
		cookie := http.Cookie{
			Name:  h.config.RedirectQueryParameter,
			Value: redirectTo,
		}
		http.SetCookie(w, &cookie)
	}
}

func (h *Handler) allowRedirect(r *http.Request) bool {
	if !h.config.Redirect {
		return false
	}
	if !h.config.RedirectCheckReferer {
		return true
	}
	referer, err := url.Parse(r.Header.Get("Referer"))
	if err != nil {
		logging.Application(r.Header).Warnf("couldn't parse redirect url %s", err)
		return false
	}
	if referer.Host != r.Host {
		logging.Application(r.Header).Warnf("redirect from referer domain: '%s', not matching current domain '%s'", referer.Host, r.Host)
		return false
	}
	return true
}
