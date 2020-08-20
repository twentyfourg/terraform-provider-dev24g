package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataReadRepository,

		Schema: map[string]*schema.Schema{
			"scm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "git",
			},
			"has_wiki": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"has_issues": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"website": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clone_ssh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone_https": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"pipelines_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"fork_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "allow_forks",
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataReadRepository(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	name := d.Get("name")
	if name == "" {
		return fmt.Errorf("name must not be blank")
	}
	workspace := d.Get("workspace")
	if workspace == "" {
		return fmt.Errorf("workspace must not be blank")
	}

	r, err := c.Get(fmt.Sprintf("2.0/repositories/%s/%s", workspace, name))
	if err != nil {
		return err
	}

	if r.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository not found")
	}

	if r.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching user")
	}

	var u Repository

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return err
	}

	id := strings.TrimSuffix(u.UUID, "}")
	id = strings.TrimPrefix(id, "{")
	d.SetId(id)
	d.Set("uuid", id)
	d.Set("scm", u.SCM)
	d.Set("has_wiki", u.HasWiki)
	d.Set("has_issues", u.HasIssues)
	d.Set("website", u.Website)
	d.Set("is_private", u.IsPrivate)
	d.Set("fork_policy", u.ForkPolicy)
	d.Set("language", u.Language)
	d.Set("description", u.Description)
	d.Set("name", u.Name)
	d.Set("slug", u.Slug)

	return nil
}
