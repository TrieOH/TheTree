package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) MarkAttendance(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	recordID, err := req.Path("record_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.MarkAttendance(r.Context(), activityID, recordID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK("Marked Attendance Successfully").Send(w)
}
