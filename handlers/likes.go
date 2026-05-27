package handlers

import (
	"forum/database"
	"net/http"
	"strconv"
)

func LikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	messageIDStr := r.URL.Query().Get("id")
	topicIDStr := r.URL.Query().Get("topic_id")
	messageID, _ := strconv.Atoi(messageIDStr)

	// Si déjà liké, on le retire, sinon on l'ajoute
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM message_likes WHERE message_id = ? AND user_id = ?", messageID, userID).Scan(&exists)
	
	if err == nil && exists > 0 {
		// Il a déjà liké -> On retire le like
		database.DB.Exec("DELETE FROM message_likes WHERE message_id = ? AND user_id = ?", messageID, userID)
	} else {
		// Il n'a pas encore liké -> On ajoute le like
		database.DB.Exec("INSERT INTO message_likes (message_id, user_id) VALUES (?, ?)", messageID, userID)
	}

	// On recharge la page du topic actuel
	http.Redirect(w, r, "/topic/view?id="+topicIDStr, http.StatusSeeOther)
}
