package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary(envConfig EnvConfig) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromURL(envConfig.CloudinaryUrl)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	cld.Config.URL.Secure = true

	return cld
}
