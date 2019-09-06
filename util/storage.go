package util

import (
	"fmt"
	"git.openstack.org/openstack/golang-client/objectstorage/v1"
	"git.openstack.org/openstack/golang-client/openstack"
	"time"
)

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

	// Get a file
	_, body, err := objectstorage.GetObject(sess, Conf.OpenStackUrlContainer+"/pikachu.png")
	if err != nil {
		panicString := fmt.Sprint("GetObject Error:", err)
		panic(panicString)
	}

	print(body)
}
