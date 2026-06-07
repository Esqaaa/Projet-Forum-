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

	var likeExists int
	database.DB.QueryRow("SELECT COUNT(*) FROM message_likes WHERE message_id = ? AND user_id = ?", messageID, userID).Scan(&likeExists)
	
	if likeExists > 0 {
		database.DB.Exec("DELETE FROM message_likes WHERE message_id = ? AND user_id = ?", messageID, userID)
	} else {
		database.DB.Exec("DELETE FROM message_dislikes WHERE message_id = ? AND user_id = ?", messageID, userID)
		database.DB.Exec("INSERT INTO message_likes (message_id, user_id) VALUES (?, ?)", messageID, userID)
	}

	http.Redirect(w, r, "/topic/view?id="+topicIDStr, http.StatusSeeOther)
}

func DislikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetLoggedUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	messageIDStr := r.URL.Query().Get("id")
	topicIDStr := r.URL.Query().Get("topic_id")
	messageID, _ := strconv.Atoi(messageIDStr)

	var dislikeExists int
	database.DB.QueryRow("SELECT COUNT(*) FROM message_dislikes WHERE message_id = ? AND user_id = ?", messageID, userID).Scan(&dislikeExists)
	
	if dislikeExists > 0 {
		database.DB.Exec("DELETE FROM message_dislikes WHERE message_id = ? AND user_id = ?", messageID, userID)
	} else {
		database.DB.Exec("DELETE FROM message_likes WHERE message_id = ? AND user_id = ?", messageID, userID)
		database.DB.Exec("INSERT INTO message_dislikes (message_id, user_id) VALUES (?, ?)", messageID, userID)
	}

	http.Redirect(w, r, "/topic/view?id="+topicIDStr, http.StatusSeeOther)
}