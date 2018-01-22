package api

import (
	"net/http"
)

// TODO: write your own scanner to Scan nested psql data types, example arrays

// implement interface??
// can use `db: <prefix>_<field_name>` ??
type Package struct {
	ID         uint32     `db:"pkg_id" json:"id"`
	Name       string     `db:"pkg_name" json:"name"`
	Version    string     `db:"pkg_version" json:"version"`
	Maintainer string     `db:"pkg_maintainer" json:"maintainer"`
	Summary    NullString `db:"pkg_summary" json:"summary,omitempty"` // make omitempty works
	Comments   NullString `db:"pkg_comments" json:"comments"`
	Labels     NullString `db:"labels" json:"labels"`
}

func (s *Service) listPackages(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var pkgs []Package
	err := s.DB.Select(&pkgs,
		`
		SELECT
			pkg_id, pkg_name, pkg_version, pkg_maintainer,
			pkg_summary, pkg_comments, json_agg(lb_name) as labels
		FROM packages
		LEFT JOIN packages_labels ON pl_pkg_id = pkg_id
		LEFT JOIN labels ON lb_id = pl_lb_id
		GROUP BY pkg_id
		`,
	)
	return pkgs, 200, err
}

func (s *Service) addPackage(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {

	return "", 200, nil
}

func (s *Service) editPackage(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {

	return "", 200, nil
}

func (s *Service) removePackage(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {

	return "", 200, nil
}
