package cloudinary

type UploadImage struct {
	ImageURL string `bson:"image_url" json:"image_url"`
	AssetID  string `bson:"asset_id" json:"asset_id"`
}
