package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Label struct {
	ID      uint32     `db:"lb_id" json:"id"`
	Name    string     `db:"lb_name" json:"name"`
	Summary NullString `db:"lb_summary" json:"summary"`
	Color   NullString `db:"lb_color" json:"color"`
	Active  bool       `db:"lb_active" json:"active"`
	PKGSCnt int        `db:"pkgs_cnt" json:"pkgs_cnt"`
}

func (s *Service) listLabels(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var labels []Label
	err := s.DB.Select(&labels,
		`
		SELECT
			lb_id, lb_name, lb_summary,
			lb_color, lb_active, COUNT(pl_id) as pkgs_cnt 
		FROM labels
		LEFT JOIN packages_labels ON pl_lb_id = lb_id
		GROUP BY (lb_id)
		ORDER BY pkgs_cnt DESC
		`,
	)
	return labels, 200, err
}

func (s *Service) listLabelPackages(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var output []struct {
		PKG          Package    `db:",prefix=pkg_." json:"pkg"`
		FieldsValues NullString `db:"fields_values" json:"fields_values"`
	}
	vars := mux.Vars(req)
	err := s.DB.Select(&output,
		`
		SELECT
			pkg_id, pkg_name, pkg_version, pkg_maintainer,
			json_agg(lfv_value) as fields_values
		FROM labels
		LEFT JOIN labels_fields on lb_id = lbf_lb_id
		LEFT JOIN packages_labels on lb_id = pl_lb_id
		LEFT JOIN packages on pkg_id = pl_pkg_id
		LEFT JOIN labels_fields_values ON lfv_lbf_id = lbf_id and lfv_pl_id = pl_id
		WHERE lb_id = $1
		GROUP BY pkg_id
		`,
		vars["id"],
	)

	return output, 200, err
}

func (s *Service) getLabel(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	// TODO: fix problem with NullString for returning one object
	var output []struct {
		Label              Label      `db:",prefix=lb_." json:"label"`
		FieldsNames        NullString `db:"fields_names" json:"fields_names"`
		FieldsDescriptions NullString `db:"fields_descriptions" json:"fields_descriptions"`
	}
	vars := mux.Vars(req)
	err := s.DB.Select(&output,
		`
		SELECT
			lb_id, lb_name, lb_summary, lb_color, lb_active,
			json_agg(lbf_name) as fields_names,
			json_agg(lbf_description) as fields_descriptions,
			(
				SELECT 
					COUNT(pl_id)
				FROM packages_labels
				WHERE pl_lb_id = $1
			) as pkgs_cnt
		FROM labels
		LEFT JOIN labels_fields ON lbf_lb_id = lb_id
		GROUP BY lb_id
		HAVING lb_id = $1
		`,
		vars["id"],
	)
	return output, 200, err
}

func (s *Service) addLabel(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var (
		p struct {
			Name    string
			Summary string
			Color   string
			Fields  []struct {
				Name        string
				Description string
			}
		}
		id int
	)
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		return "", 200, err
	}
	tx := s.DB.MustBegin()
	tx.Get(&id, "INSERT INTO labels (lb_name, lb_summary, lb_color) VALUES ($1, $2, $3) RETURNING lb_id",
		p.Name, p.Summary, p.Color)
	for _, f := range p.Fields {
		tx.MustExec("INSERT INTO labels_fields (lbf_lb_id, lbf_name, lbf_description) VALUES ($1, $2, $3)",
			id, f.Name, f.Description)
	}
	tx.Commit()
	return id, 200, nil
}

func (s *Service) attachLabel(rw http.ResponseWriter, req *http.Request) (interface{}, int, error) {
	var (
		p struct {
			PKGID   int `json:"pkg_id"`
			LabelID int `json:"label_id"`
			Fields  []struct {
				ID    int
				Value string
			}
		}
		id int
	)

	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		return "", 200, err
	}
	tx := s.DB.MustBegin()
	tx.Get(&id, "INSERT INTO packages_labels (pl_pkg_id, pl_lb_id) VALUES ($1, $2) RETURNING pl_id",
		p.PKGID, p.LabelID)
	for _, f := range p.Fields {
		tx.MustExec("INSERT INTO labels_fields_values (lfv_pl_id, lfv_lbf_id, lfv_value) VALUES ($1, $2, $3)",
			id, f.ID, f.Value)
	}
	tx.Commit()
	return id, 200, nil
}
