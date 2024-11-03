package dns

import (
	"net/http"
	"text/template"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetRecords(w http.ResponseWriter, r *http.Request) {
	records, err := h.service.repository.getDNSRecordUpdates()
	if err == nil {
		tpl := template.Must(template.New("index").Parse(recordsTmpl))

		w.WriteHeader(http.StatusOK)
		if err := tpl.Execute(w, records); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}
