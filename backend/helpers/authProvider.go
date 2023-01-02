package helpers

var ProviderFieldMapping = map[string]string{
	"github": "github_id",
	"google": "google_id",
}

type NextAuthProvider string

const (
	Github NextAuthProvider = "github"
	Google NextAuthProvider = "google"
)
