package validations

import (
	"fmt"

	"github.com/alopez-2018459/go-the-field/internal/models"
)

func ValidatePostContent(content *models.PostContent) error {

	text := content.Text

	media := content.Media

	err := IsStringEmpty(text)

	if err != nil && len(media) == 0 || media == nil {
		return fmt.Errorf("Post content is required")

	}
	return nil
}
