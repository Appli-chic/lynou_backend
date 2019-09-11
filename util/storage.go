package util

import (
	"fmt"
	"git.openstack.org/openstack/golang-client/objectstorage/v1"
	"git.openstack.org/openstack/golang-client/openstack"
	"github.com/applichic/lynou/config"
	"net/http"
	"time"
)

var OpenstackSession *openstack.Session

func LoginToStorage() {
	creds := openstack.AuthOpts{
		AuthUrl:     config.Conf.OpenStackAuthUrl,
		ProjectName: config.Conf.OpenStackProjectName,
		Username:    config.Conf.OpenStackUsername,
		Password:    config.Conf.OpenStackPassword,
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
func DownloadObject(filePath string) (http.Header, []byte, error) {
	var header http.Header
	_, body, err := objectstorage.GetObject(OpenstackSession, config.Conf.OpenStackUrlContainer+"/"+filePath)

	// Get the content type
	if err == nil {
		header, err = objectstorage.GetObjectMeta(OpenstackSession, config.Conf.OpenStackUrlContainer+"/"+filePath)
	}

	return header, body, err
}

// Upload file in the object storage
func UploadObject(filePath string, content *[]byte, contentType string) error {
	headers := http.Header{}
	headers.Add("Content-Type", contentType)
	err := objectstorage.PutObject(OpenstackSession, content, config.Conf.OpenStackUrlContainer+"/"+filePath, headers)
	return err
}
