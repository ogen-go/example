// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"net/http"
	"strings"
)

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	s.cfg.NotFound(w, r)
}

func (s *Server) notAllowed(w http.ResponseWriter, r *http.Request, allowed string) {
	s.cfg.MethodNotAllowed(w, r, allowed)
}

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Path
	if prefix := s.cfg.Prefix; len(prefix) > 0 {
		if strings.HasPrefix(elem, prefix) {
			// Cut prefix from the path.
			elem = strings.TrimPrefix(elem, prefix)
		} else {
			// Prefix doesn't match.
			s.notFound(w, r)
			return
		}
	}
	if len(elem) == 0 {
		s.notFound(w, r)
		return
	}
	args := [1]string{}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/pet"
			if l := len("/pet"); len(elem) >= l && elem[0:l] == "/pet" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				switch r.Method {
				case "POST":
					s.handleAddPetRequest([0]string{}, w, r)
				default:
					s.notAllowed(w, r, "POST")
				}

				return
			}
			switch elem[0] {
			case '/': // Prefix: "/"
				if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "petId"
				// Leaf parameter
				args[0] = elem
				elem = ""

				if len(elem) == 0 {
					// Leaf node.
					switch r.Method {
					case "DELETE":
						s.handleDeletePetRequest([1]string{
							args[0],
						}, w, r)
					case "GET":
						s.handleGetPetByIdRequest([1]string{
							args[0],
						}, w, r)
					case "POST":
						s.handleUpdatePetRequest([1]string{
							args[0],
						}, w, r)
					default:
						s.notAllowed(w, r, "DELETE,GET,POST")
					}

					return
				}
			}
		}
	}
	s.notFound(w, r)
}

// Route is route object.
type Route struct {
	name        string
	operationID string
	count       int
	args        [1]string
}

// Name returns ogen operation name.
//
// It is guaranteed to be unique and not empty.
func (r Route) Name() string {
	return r.name
}

// OperationID returns OpenAPI operationId.
func (r Route) OperationID() string {
	return r.operationID
}

// Args returns parsed arguments.
func (r Route) Args() []string {
	return r.args[:r.count]
}

// FindRoute finds Route for given method and path.
func (s *Server) FindRoute(method, path string) (r Route, _ bool) {
	var (
		args = [1]string{}
		elem = path
	)
	r.args = args
	if elem == "" {
		return r, false
	}

	// Static code generated router with unwrapped path search.
	switch {
	default:
		if len(elem) == 0 {
			break
		}
		switch elem[0] {
		case '/': // Prefix: "/pet"
			if l := len("/pet"); len(elem) >= l && elem[0:l] == "/pet" {
				elem = elem[l:]
			} else {
				break
			}

			if len(elem) == 0 {
				switch method {
				case "POST":
					r.name = "AddPet"
					r.operationID = "addPet"
					r.args = args
					r.count = 0
					return r, true
				default:
					return
				}
			}
			switch elem[0] {
			case '/': // Prefix: "/"
				if l := len("/"); len(elem) >= l && elem[0:l] == "/" {
					elem = elem[l:]
				} else {
					break
				}

				// Param: "petId"
				// Leaf parameter
				args[0] = elem
				elem = ""

				if len(elem) == 0 {
					switch method {
					case "DELETE":
						// Leaf: DeletePet
						r.name = "DeletePet"
						r.operationID = "deletePet"
						r.args = args
						r.count = 1
						return r, true
					case "GET":
						// Leaf: GetPetById
						r.name = "GetPetById"
						r.operationID = "getPetById"
						r.args = args
						r.count = 1
						return r, true
					case "POST":
						// Leaf: UpdatePet
						r.name = "UpdatePet"
						r.operationID = "updatePet"
						r.args = args
						r.count = 1
						return r, true
					default:
						return
					}
				}
			}
		}
	}
	return r, false
}
