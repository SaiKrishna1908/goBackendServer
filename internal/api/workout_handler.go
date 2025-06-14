// Service Layer
package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"goBackendServer/internal/middleware"
	"goBackendServer/internal/store"
	"goBackendServer/internal/utils"
	"log"
	"net/http"
)

type WorkoutHandler struct {
	workOutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workOutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workOutStore: workOutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workOutId, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := wh.workOutStore.GetWorkoutById(workOutId)

	if err != nil {
		wh.logger.Printf("Failed to fetch workout %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workOut store.Workout
	err := json.NewDecoder(r.Body).Decode(&workOut)

	if err != nil {
		wh.logger.Printf("ERROR: decodingCreateWorkout: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Bad Request"})
		return
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		wh.logger.Printf("ERROR: Unauthorized Access")
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	workOut.UserId = currentUser.ID

	savedWorkout, err := wh.workOutStore.CreateWorkout(&workOut)
	if err != nil {
		wh.logger.Printf("ERROR: errorCreatingWorkout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": savedWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: errorUpdatingWorkout: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Bad Request"})
		return
	}

	existingWorkout, err := wh.workOutStore.GetWorkoutById(workoutId)

	if err != nil {
		wh.logger.Printf("Error, errorFetchingExistingWorkout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}

	if existingWorkout == nil {
		wh.logger.Printf("Error, workoutDoesNotExist: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Internal Server Error"})
		return
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "please login to use this feature"})
	}

	workoutUserId, err := wh.workOutStore.GetWorkoutOwner(workoutId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			wh.logger.Printf("Error, noWorkoutFound: %v", err)
			utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
			return
		}

		wh.logger.Printf("Error, errorGetWorkoutOwner %v: ", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutUserId != currentUser.ID {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized Access"})
		return
	}

	// at this point we have a existing workout
	var updatedWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duraiton_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedWorkoutRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedWorkoutRequest.Title != nil {
		existingWorkout.Title = *updatedWorkoutRequest.Title
	}

	if updatedWorkoutRequest.Description != nil {
		existingWorkout.Description = *updatedWorkoutRequest.Description
	}

	if updatedWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updatedWorkoutRequest.CaloriesBurned
	}

	if updatedWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updatedWorkoutRequest.DurationMinutes
	}

	if updatedWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updatedWorkoutRequest.Entries
	}

	err = wh.workOutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		wh.logger.Printf("Error, updatingWorkout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {

	workoutId, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("Error, deletingWorkout: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Bad Request"})
		return
	}

	currentUser := middleware.GetUser(r)

	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "please login to use this feature"})
	}

	workoutUserId, err := wh.workOutStore.GetWorkoutOwner(workoutId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			wh.logger.Printf("Error, noWorkoutFound: %v", err)
			utils.WriteJson(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
			return
		}

		wh.logger.Printf("Error, errorGetWorkoutOwner %v: ", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutUserId != currentUser.ID {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Bad Request"})
		return
	}

	err = wh.workOutStore.DeleteWorkoutById(workoutId)

	if err != nil {
		wh.logger.Printf("Error, deletingWorkout: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Bad Request"})
		return
	}

	if err == sql.ErrNoRows {
		wh.logger.Printf("Error, SqlErrorNoRows: %v", err)
		utils.WriteJson(w, http.StatusNoContent, utils.Envelope{"error": "Bad Request"})
		return
	}

	w.WriteHeader(http.StatusOK)

	utils.WriteJson(w, http.StatusOK, nil)
}
