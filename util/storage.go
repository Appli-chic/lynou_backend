package util

import (
	"fmt"
	"git.openstack.org/openstack/golang-client/objectstorage/v1"
	"git.openstack.org/openstack/golang-client/openstack"
	"net/http"
	"time"
)

var OpenstackSession *openstack.Session

func LoginToStorage() {
	creds := openstack.AuthOpts{
		AuthUrl:     Conf.OpenStackAuthUrl,
		ProjectName: Conf.OpenStackProjectName,
		Username:    Conf.OpenStackUsername,
		Password:    Conf.OpenStackPassword,
	}

	auth, err := openstack.DoAuthRequest(creds)
	if err != nil {
		fmt.Println("Error authenticating username/password:", err)
		return
	}

	if !auth.GetExpiration().After(time.Now()) {
		fmt.Println("There was an error. The auth token has an invalid expiration.")
		return
	}

	// Make a new client with these creds
	sess, err := openstack.NewSession(nil, auth, nil)
	if err != nil {
		panicString := fmt.Sprint("Error creating new Session:", err)
		panic(panicString)
	}

	OpenstackSession = sess
}

// Get a file in the object storage from the file path
func DownloadObject(filePath string) ([]byte, error) {
	_, body, err := objectstorage.GetObject(OpenstackSession, Conf.OpenStackUrlContainer+"/"+filePath)
	return body, err
}

// Upload file in the object storage
func UploadObject(filePath string, content *[]byte, contentType string) error {
	headers := http.Header{}
	headers.Add("Content-Type", contentType)
	err := objectstorage.PutObject(OpenstackSession, content, Conf.OpenStackUrlContainer+"/"+filePath, headers)
	return err
}
