package pkg

import (
	"net/http"
	"net/netip"
	"net/url"
	"strings"

	"github.com/rotisserie/eris"
	"golang.org/x/net/idna"
)

func newHTTPClient() *http.Client {
	return &http.Client{CheckRedirect: checkRedirect}
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return eris.New("stopped after 10 redirects")
	}

	status := 0
	if req.Response != nil {
		status = req.Response.StatusCode
	}

	if status != http.StatusTemporaryRedirect && status != http.StatusPermanentRedirect {
		return eris.Errorf("refusing to follow redirect with status %d; only 307 and 308 are supported", status)
	}

	origin, err := effectiveOrigin(via[0].URL)
	if err != nil {
		return err
	}

	target, err := effectiveOrigin(req.URL)
	if err != nil {
		return err
	}

	if target != origin {
		return eris.Errorf("refusing to follow redirect from %s to %s", origin, target)
	}

	return nil
}

func effectiveOrigin(u *url.URL) (string, error) {
	scheme := strings.ToLower(u.Scheme)

	port := u.Port()
	if port == "" {
		switch scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		default:
			return "", eris.Errorf("no default port for scheme %q", scheme)
		}
	}

	host := strings.ToLower(u.Hostname())
	if addr, err := netip.ParseAddr(host); err == nil {
		host = addr.Unmap().String()
		if strings.Contains(host, ":") {
			host = "[" + host + "]"
		}
	} else {
		ascii, err := idna.Lookup.ToASCII(host)
		if err != nil {
			return "", eris.Wrapf(err, "failed to normalize host %q", u.Hostname())
		}
		host = ascii
	}

	return scheme + "://" + host + ":" + port, nil
}
