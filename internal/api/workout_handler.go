package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
}

func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
}

func (wh *WorkoutHandler) HandleWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workOutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "this is the workout id: %d\n", workOutId)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Created a workout\n")
}
