package postmark

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
