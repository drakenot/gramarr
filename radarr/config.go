package radarr

type Config struct {
	Hostname           string `json:"hostname"`
	APIKey             string `json:"apiKey"`
	Port               int    `json:"port"`
	URLBase            string `json:"urlBase"`
	SSL                bool   `json:"ssl"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	MaxResults         int    `json:"maxResults"`
	RestrictedFolders  []int  `json:"restrictedFolders"`
	RestrictedProfiles []int  `json:"restrictedProfiles"`
}
