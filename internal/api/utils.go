package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/natasha-m-oliveira/semana-tech-go-react-server/internal/store/pgstore"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (h apiHandler) readRoom(
	w http.ResponseWriter,
	r *http.Request,
) (room pgstore.Room, rawRoomID string, roomID uuid.UUID, ok bool) {
	rawRoomID = chi.URLParam(r, "room_id")
	roomID, err := uuid.Parse(rawRoomID)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return pgstore.Room{}, "", uuid.UUID{}, false
	}

	room, err = h.q.GetRoom(r.Context(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return pgstore.Room{}, "", uuid.UUID{}, false
		}

		slog.Error("failed to get room", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return pgstore.Room{}, "", uuid.UUID{}, false
	}

	return room, rawRoomID, roomID, true
}

func (h apiHandler) validateMessageID(
	w http.ResponseWriter,
	r *http.Request,
) (rawMessageID string, messageID uuid.UUID, ok bool) {
	rawMessageID = chi.URLParam(r, "message_id")
	messageID, err := uuid.Parse(rawMessageID)
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return "", uuid.UUID{}, false
	}

	return rawMessageID, messageID, true
}

func sendJSON(w http.ResponseWriter, rawData any) {
	data, _ := json.Marshal(rawData)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
