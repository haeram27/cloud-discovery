package apis

import (
	"context"
	"fmt"
	"log"
	"testing"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// implicit uses Application Default Credentials to authenticate.
func TestImplicit(t *testing.T) {
	projectID := "your project id"
	ctx := context.Background()

	// For API packages whose import path is starting with "cloud.google.com/go",
	// such as cloud.google.com/go/storage in this case, if there are no credentials
	// provided, the client library will look for credentials in the environment.

	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/path/to/service-account-key.json")
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer storageClient.Close()

	it := storageClient.Buckets(ctx, projectID)
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(bucketAttrs.Name)
	}

	// For packages whose import path is starting with "google.golang.org/api",
	// such as google.golang.org/api/cloudkms/v1, use NewService to create the client.
	kmsService, err := cloudkms.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_ = kmsService
}

// explicit reads credentials from the specified path.
func TestExplicit(t *testing.T) {
	jsonKeyPath := "/path/to/service_account_key.json"
	projectID := "your project id"
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonKeyPath))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("Buckets:")
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(battrs.Name)
	}
}

// explicitDefault finds the default credentials.
//
// It is very uncommon to need to explicitly get the default credentials in Go.
// Most of the time, client libraries can use Application Default Credentials
// without having to pass the credentials in directly. See implicit above.

// OAuth2 scopes used by this API.
// Catalog of all all valid scope names at https://developers.google.com/identity/protocols/oauth2/scopes
const (
	// View and manage your data across Google Cloud Platform services
	CloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

	// View your data across Google Cloud Platform services
	CloudPlatformReadOnlyScope = "https://www.googleapis.com/auth/cloud-platform.read-only"
)

func TestExplicitDefault(t *testing.T) {
	projectID := "your project id"
	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, CloudPlatformReadOnlyScope)
	if err != nil {
		log.Fatal(err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("Buckets:")
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(battrs.Name)
	}
}
