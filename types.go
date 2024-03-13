package main

type ImportedPost struct {
	ID string `json:"id"`
}

type ImportedPosts struct {
	ImportedPosts []ImportedPost `json:"posts"`
}

type FacebookPosts struct {
	Data []struct {
		ID          string `json:"id,omitempty"`
		Attachments struct {
			Data []struct {
				Description string `json:"description,omitempty"`
				Media       struct {
					Image struct {
						Height int    `json:"height,omitempty"`
						Src    string `json:"src,omitempty"`
						Width  int    `json:"width,omitempty"`
					} `json:"image,omitempty"`
				} `json:"media,omitempty"`
				Target struct {
					ID  string `json:"id,omitempty"`
					URL string `json:"url,omitempty"`
				} `json:"target,omitempty"`
				Type string `json:"type,omitempty"`
				URL  string `json:"url,omitempty"`
			} `json:"data,omitempty"`
		} `json:"attachments,omitempty"`
		CreatedTime  string `json:"created_time,omitempty"`
		FullPicture  string `json:"full_picture,omitempty"`
		Message      string `json:"message,omitempty"`
		PermalinkURL string `json:"permalink_url,omitempty"`
	} `json:"data,omitempty"`
	Paging struct {
		Cursors struct {
			Before string `json:"before,omitempty"`
			After  string `json:"after,omitempty"`
		} `json:"cursors,omitempty"`
		Next string `json:"next,omitempty"`
	} `json:"paging,omitempty"`
}

