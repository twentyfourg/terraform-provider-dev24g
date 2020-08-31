package dev24g

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

type environment struct {
	UUID string `json:"uuid"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func dataEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataReadEnvironment,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repo_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataReadEnvironment(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	environmentName := d.Get("name")
	if environmentName == "" {
		return fmt.Errorf("name must not be blank")
	}
	repoSlug := d.Get("repo_slug")
	if repoSlug == "" {
		return fmt.Errorf("repo_slug must not be blank")
	}
	workspace := d.Get("workspace")
	if workspace == "" {
		return fmt.Errorf("workspace must not be blank")
	}

	l, err := c.Get(fmt.Sprintf("2.0/repositories/%s/%s/environments/", workspace, repoSlug))
	if err != nil {
		return err
	}

	if l.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository not found")
	}

	if l.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching user")
	}

	var list []environment

	err = json.NewDecoder(l.Body).Decode(&list)
	for _, e := range list {
		fmt.Println(e.Name)
	}
	if err != nil {
		return err
	}

	return nil
}
