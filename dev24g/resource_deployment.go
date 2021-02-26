package dev24g

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// Deployment structure for handling key info
type Deployment struct {
	Name  string `json:"name"`
	Stage *Stage `json:"environment_type"`
	UUID  string `json:"uuid,omitempty"`
}

// Stage structure for handing stage
type Stage struct {
	Name string `json:"name"`
}

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentCreate,
		Update: resourceDeploymentUpdate,
		Read:   resourceDeploymentRead,
		Delete: resourceDeploymentDelete,
		// TODO: implement import
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stage": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Test",
					"Staging",
					"Production",
				},
					false),
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func newDeploymentFromResource(d *schema.ResourceData) *Deployment {
	dk := &Deployment{
		Name: d.Get("name").(string),
		Stage: &Stage{
			Name: d.Get("stage").(string),
		},
	}
	return dk
}

func resourceDeploymentCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*Client)
	rvcr := newDeploymentFromResource(d)
	bytedata, err := json.Marshal(rvcr)

	if err != nil {
		return err
	}
	resp, err := client.Post(fmt.Sprintf("2.0/repositories/%s/environments/",
		d.Get("repository").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		if strings.Contains(err.Error(), "400") {
			return errors.New("Deployment already exists")
		}
		return err
	}

	var deployment Deployment

	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		return readerr
	}

	decodeerr := json.Unmarshal(body, &deployment)
	if decodeerr != nil {
		return decodeerr
	}
	d.Set("uuid", deployment.UUID)
	d.SetId(fmt.Sprintf("%s:%s", d.Get("repository"), deployment.UUID))

	return resourceDeploymentRead(d, m)
}

func resourceDeploymentRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*Client)
	resp, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/environments/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	))

	log.Printf("ID: %s", url.PathEscape(d.Id()))

	if resp.StatusCode == 200 {
		var Deployment Deployment
		body, readerr := ioutil.ReadAll(resp.Body)
		if readerr != nil {
			return readerr
		}

		decodeerr := json.Unmarshal(body, &Deployment)
		if decodeerr != nil {
			return decodeerr
		}

		d.Set("uuid", Deployment.UUID)
		d.Set("name", Deployment.Name)
		d.Set("stage", Deployment.Stage.Name)
	}

	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	rvcr := newDeploymentFromResource(d)
	bytedata, err := json.Marshal(rvcr)

	if err != nil {
		return err
	}
	resp, err := client.Put(fmt.Sprintf("2.0/repositories/%s/environments/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	), bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return nil
	}

	return resourceDeploymentRead(d, m)
}

func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, err := client.Delete(fmt.Sprintf("2.0/repositories/%s/environments/%s",
		d.Get("repository").(string),
		d.Get("uuid").(string),
	))
	return err
}
