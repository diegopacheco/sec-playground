package api

import (
	"encoding/json"
	"net/http"

	"mcp-output-taint-tracker/internal/scenario"
	"mcp-output-taint-tracker/internal/session"
)

const maxBody = 1 << 20

type promptReq struct {
	Text string `json:"text"`
}

type resultReq struct {
	Tool   string         `json:"tool"`
	Args   map[string]any `json:"args"`
	Output string         `json:"output"`
}

type callReq struct {
	Tool string         `json:"tool"`
	Args map[string]any `json:"args"`
}

type createReq struct {
	Name string `json:"name"`
}

func New(store *session.Store) http.Handler {
	h := &handlers{store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.health)
	mux.HandleFunc("GET /api/policy", h.policy)
	mux.HandleFunc("GET /api/sessions", h.list)
	mux.HandleFunc("POST /api/sessions", h.create)
	mux.HandleFunc("GET /api/sessions/{id}", h.get)
	mux.HandleFunc("DELETE /api/sessions/{id}", h.remove)
	mux.HandleFunc("POST /api/sessions/{id}/prompt", h.prompt)
	mux.HandleFunc("POST /api/sessions/{id}/result", h.result)
	mux.HandleFunc("POST /api/sessions/{id}/call", h.call)
	mux.HandleFunc("POST /api/scenarios", h.scenarios)
	return mux
}

type handlers struct {
	store *session.Store
}

func (h *handlers) health(w http.ResponseWriter, r *http.Request) {
	reply(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *handlers) policy(w http.ResponseWriter, r *http.Request) {
	reply(w, http.StatusOK, h.store.Policy())
}

func (h *handlers) list(w http.ResponseWriter, r *http.Request) {
	reply(w, http.StatusOK, h.store.List())
}

func (h *handlers) create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if !decode(w, r, &req) {
		return
	}
	if req.Name == "" {
		req.Name = "manual session"
	}
	reply(w, http.StatusCreated, h.store.Create(req.Name))
}

func (h *handlers) get(w http.ResponseWriter, r *http.Request) {
	s, ok := h.store.Get(r.PathValue("id"))
	if !ok {
		fail(w, http.StatusNotFound, "session not found")
		return
	}
	reply(w, http.StatusOK, s)
}

func (h *handlers) remove(w http.ResponseWriter, r *http.Request) {
	if !h.store.Delete(r.PathValue("id")) {
		fail(w, http.StatusNotFound, "session not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) prompt(w http.ResponseWriter, r *http.Request) {
	var req promptReq
	if !decode(w, r, &req) {
		return
	}
	reply(w, http.StatusOK, h.store.Prompt(r.PathValue("id"), req.Text))
}

func (h *handlers) result(w http.ResponseWriter, r *http.Request) {
	var req resultReq
	if !decode(w, r, &req) {
		return
	}
	if req.Tool == "" {
		fail(w, http.StatusBadRequest, "tool is required")
		return
	}
	reply(w, http.StatusOK, h.store.Result(r.PathValue("id"), req.Tool, req.Args, req.Output))
}

func (h *handlers) call(w http.ResponseWriter, r *http.Request) {
	var req callReq
	if !decode(w, r, &req) {
		return
	}
	if req.Tool == "" {
		fail(w, http.StatusBadRequest, "tool is required")
		return
	}
	reply(w, http.StatusOK, h.store.Call(r.PathValue("id"), req.Tool, req.Args))
}

func (h *handlers) scenarios(w http.ResponseWriter, r *http.Request) {
	out := []session.Summary{}
	for _, sc := range scenario.All() {
		s, _ := scenario.Run(h.store, sc)
		out = append(out, session.Summary{ID: s.ID, Name: s.Name, Created: s.Created, Events: len(s.Events)})
	}
	reply(w, http.StatusCreated, out)
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBody)).Decode(v); err != nil {
		fail(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return false
	}
	return true
}

func reply(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func fail(w http.ResponseWriter, code int, msg string) {
	reply(w, code, map[string]string{"error": msg})
}
