package postmark

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	postmarkSDK "github.com/keighl/postmark"
)

func doStreamRequests(method string, path string, stream Stream, serverToken string) (diag.Diagnostics, Stream) {
	var diags diag.Diagnostics

	client := &http.Client{}

	req, err := http.NewRequest(method, "https://api.postmarkapp.com/message-streams/"+path, nil)
	if err != nil {
		return diag.FromErr(err), Stream{}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Server-Token", serverToken)
	if strings.Contains(path, "archive") {
		req.Header.Add("Content-Length", "0")
	}

	body, err := json.Marshal(stream)
	if err != nil {
		return diag.FromErr(err), Stream{}
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err), Stream{}
	}
	defer res.Body.Close()

	if err != nil {
		return diag.FromErr(err), Stream{}
	}
	err = json.NewDecoder(res.Body).Decode(&stream)
	if err != nil {
		return diag.FromErr(err), Stream{}
	}
	return diags, stream
}

func doDomainRequests(method string, path string, domain Domain, m interface{}) (diag.Diagnostics, Domain) {
	var diags diag.Diagnostics

	client := &http.Client{}
	c := m.(*postmarkSDK.Client)

	req, err := http.NewRequest(method, "https://api.postmarkapp.com/domains/"+path, nil)
	if err != nil {
		return diag.FromErr(err), Domain{}
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Account-Token", c.AccountToken)

	body, err := json.Marshal(domain)
	if err != nil {
		return diag.FromErr(err), Domain{}
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err), Domain{}
	}
	defer res.Body.Close()

	if err != nil {
		return diag.FromErr(err), Domain{}
	}
	err = json.NewDecoder(res.Body).Decode(&domain)
	if err != nil {
		return diag.FromErr(err), Domain{}
	}
	return diags, domain
}
