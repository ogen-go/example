package httpmiddleware

import "net/url"

// Server is a generic ogen server type.
type Server[R Route] interface {
	FindPath(method string, u *url.URL) (r R, _ bool)
}

// Route is a generic ogen route type.
type Route interface {
	Name() string
	OperationID() string
	PathPattern() string
}

// RouteFinder finds Route by given URL.
type RouteFinder func(method string, u *url.URL) (Route, bool)

// MakeRouteFinder creates RouteFinder from given server.
func MakeRouteFinder[R Route, S Server[R]](server S) RouteFinder {
	return func(method string, u *url.URL) (Route, bool) {
		return server.FindPath(method, u)
	}
}
