package database

import (
	"database/sql"
	"log"

	"MrRSS/internal/models"
)

// UpdateArticleLabels updates the labels field for an article.
func (db *DB) UpdateArticleLabels(id int64, labels string) error {
	db.WaitForReady()
	_, err := db.Exec("UPDATE articles SET labels = ? WHERE id = ?", labels, id)
	return err
}

// GetArticlesByLabel retrieves articles that have a specific label.
func (db *DB) GetArticlesByLabel(label string, limit, offset int) ([]models.Article, error) {
	db.WaitForReady()
	query := `
		SELECT a.id, a.feed_id, a.title, a.url, a.image_url, a.audio_url, a.video_url, 
		       a.published_at, a.is_read, a.is_favorite, a.is_hidden, a.is_read_later, 
		       a.translated_title, a.labels, f.title
		FROM articles a
		JOIN feeds f ON a.feed_id = f.id
		WHERE a.labels LIKE ?
		ORDER BY a.published_at DESC 
		LIMIT ? OFFSET ?
	`
	// SQLite JSON1 extension support for proper JSON querying would be ideal,
	// but for compatibility we use LIKE pattern matching
	pattern := "%\"" + label + "\"%"
	
	rows, err := db.Query(query, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		var imageURL, audioURL, videoURL, translatedTitle, labels sql.NullString
		if err := rows.Scan(&a.ID, &a.FeedID, &a.Title, &a.URL, &imageURL, &audioURL, &videoURL, 
			&a.PublishedAt, &a.IsRead, &a.IsFavorite, &a.IsHidden, &a.IsReadLater, 
			&translatedTitle, &labels, &a.FeedTitle); err != nil {
			log.Println("Error scanning article:", err)
			continue
		}
		a.ImageURL = imageURL.String
		a.AudioURL = audioURL.String
		a.VideoURL = videoURL.String
		a.TranslatedTitle = translatedTitle.String
		a.Labels = labels.String
		articles = append(articles, a)
	}
	return articles, nil
}
