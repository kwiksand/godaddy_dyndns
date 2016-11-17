package godaddy_dyndns

import (
	"bytes"
	"encoding/json"
//	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
//    "github.com/asaskevich/govalidator"
)

//#var godaddyClient *GoDaddyDNSClient

var g *GoDaddyDNSClient

const recordURL = "https://api.godaddy.com/v1/domains/%v/records/A/%v"
const zoneURL = "https://api.godaddy.com/v1/domains/%v/records"

func init() {
    g = New()
}

type GoDaddyDNSClient struct {
    Key string
    Secret string
}

var domainURL string

func New() *GoDaddyDNSClient {
    g := new(GoDaddyDNSClient)
    g.Key = "username"
    g.Secret = "password"

    return g
}

func (g *GoDaddyDNSClient) SetKey(key string) () {
    g.Key = key
}

func (g *GoDaddyDNSClient) SetSecret(secret string) () {
    g.Secret = secret
}

func (g *GoDaddyDNSClient) doRequest(req *http.Request) (string, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http failed on %v %v: %v", req.Method, req.URL, resp.StatusCode)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
	    body, err := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected response on %v %v: %v %v %v", req.Method, req.URL, resp.StatusCode, err, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

func (g *GoDaddyDNSClient) GetPublicIP() (string, error) {
	req, err := http.NewRequest("GET", "http://myexternalip.com/raw", nil)
	if err != nil {
		return "", err
	}

	return g.doRequest(req)
}

func (g *GoDaddyDNSClient) addGodaddyHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %v:%v", g.Key, g.Secret))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
}

// Domain is the request/response struct for the domain API endpoint.
type Domain struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	Data string `json:"data"`
    Zone string `json:"zone,omitempty"`
	TTL  int    `json:"ttl"`
}

func (g *GoDaddyDNSClient) GetDNS(rootDomain string, subDomain string) (string, error) {
	var domainURL = fmt.Sprintf(recordURL, rootDomain, subDomain)
	req, err := http.NewRequest("GET", domainURL, nil)
	if err != nil {
		return "", err
	}
	g.addGodaddyHeaders(req)

    //fmt.Printf("%+v\n", req)

	body, err := g.doRequest(req)
	if err != nil {
		return "", err
	}

	var res []Domain
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return "", err
	}

    //fmt.Printf("%+v\n", res)


	if len(res) == 0 {
		return "", fmt.Errorf("got empty domains response")
	}

	return res[0].Data, nil
}

func (g *GoDaddyDNSClient) InsertDNS(addr string, rootDomain string, subDomain string) error {
	var domainURL = fmt.Sprintf(zoneURL, rootDomain)
	domains := []Domain{{
		Data: addr,
        Name: subDomain,
        Zone: "",
        Type: "A",
		TTL:  600,
	}}

	domainsBody, err := json.Marshal(domains)
	if err != nil {
		return err
	}
    log.Printf("Request - URL: " + domainURL + " Body: " +  string(domainsBody))

	req, err := http.NewRequest("PATCH", domainURL, bytes.NewReader(domainsBody))
	if err != nil {
		return err
	}
	g.addGodaddyHeaders(req)

	if _, err := g.doRequest(req); err != nil {
		return err
	}

	return nil
}

func (g *GoDaddyDNSClient) UpdateDNS(addr string, rootDomain string, subDomain string) error {
	var recordURL = fmt.Sprintf(recordURL, rootDomain, subDomain)
	domains := []Domain{{
		Data: addr,
        Name: "",
        Zone: "",
        Type: "",
		TTL:  600,
	}}

	domainsBody, err := json.Marshal(domains)
	if err != nil {
		return err
	}
    log.Printf("got empty domains response" + string(domainsBody))

	req, err := http.NewRequest("PUT", recordURL, bytes.NewReader(domainsBody))
	if err != nil {
		return err
	}
	g.addGodaddyHeaders(req)

	if _, err := g.doRequest(req); err != nil {
		return err
	}

	return nil
}
