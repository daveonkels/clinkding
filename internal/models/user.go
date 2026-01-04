package models

type UserProfile struct {
	Theme                    string `json:"theme"`
	BookmarkDateDisplay      string `json:"bookmark_date_display"`
	BookmarkLinkTarget       string `json:"bookmark_link_target"`
	WebArchiveIntegration    string `json:"web_archive_integration"`
	EnableSharing            bool   `json:"enable_sharing"`
	EnablePublicSharing      bool   `json:"enable_public_sharing"`
	EnableFavicons           bool   `json:"enable_favicons"`
	EnablePreviewImages      bool   `json:"enable_preview_images"`
	DisplayURL               bool   `json:"display_url"`
	DisplayViewedDate        bool   `json:"display_viewed_date"`
	PermanentNotes           bool   `json:"permanent_notes"`
	SearchPreferences        SearchPreferences `json:"search_preferences"`
}

type SearchPreferences struct {
	Sort           string `json:"sort"`
	SharedBookmarks string `json:"shared"`
	UnreadOnly     bool   `json:"unread_only"`
}
