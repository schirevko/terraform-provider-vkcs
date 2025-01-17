package vrrpinterfaces

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type VRRPInterfaceCreate struct {
	VRRPInterface *CreateOpts `json:"dc_vrrp_interface"`
}

type CreateOpts struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	DCVRRPID      string `json:"dc_vrrp_id"`
	DCInterfaceID string `json:"dc_interface_id"`
	Priority      int    `json:"priority,omitempty"`
	Preempt       *bool  `json:"preempt,omitempty"`
	Master        *bool  `json:"master,omitempty"`
}

func (opts *VRRPInterfaceCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(vrrpInterfacesURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(vrrpInterfaceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type VRRPInterfaceUpdate struct {
	VRRPInterface *UpdateOpts `json:"dc_vrrp_interface"`
}

type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	Preempt     *bool  `json:"preempt,omitempty"`
	Master      *bool  `json:"master,omitempty"`
}

func (opts *VRRPInterfaceUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(vrrpInterfaceURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(vrrpInterfaceURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
