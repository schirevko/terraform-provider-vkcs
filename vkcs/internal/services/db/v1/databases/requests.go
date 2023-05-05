package databases

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type CreateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// BatchCreateOpts is used to send request to create databases
type BatchCreateOpts struct {
	Databases []CreateOpts `json:"databases"`
}

// CreateOpts represents parameters of creation of database
type CreateOpts struct {
	Name    string `json:"name" required:"true"`
	CharSet string `json:"character_set,omitempty"`
	Collate string `json:"collate,omitempty"`
}

// Map converts opts to a map (for a request body)
func (opts *BatchCreateOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

// Create performs request to create database
func Create(client *gophercloud.ServiceClient, id string, opts CreateOptsBuilder, dbmsType string) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(databasesURL(client, dbmsType, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})

	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// List performs request to list databases
func List(client *gophercloud.ServiceClient, id string, dbmsType string) pagination.Pager {
	return pagination.NewPager(client, databasesURL(client, dbmsType, id), func(r pagination.PageResult) pagination.Page {
		return Page{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// Delete performs request to delete database
func Delete(client *gophercloud.ServiceClient, id string, dbName string, dbmsType string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(databaseURL(client, dbmsType, id, dbName), &gophercloud.RequestOpts{})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
