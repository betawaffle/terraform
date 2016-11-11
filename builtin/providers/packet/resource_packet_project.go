package packet

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketProject() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketProjectCreate,
		Read:   resourcePacketProjectRead,
		Update: resourcePacketProjectUpdate,
		Delete: resourcePacketProjectDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"payment_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePacketProjectCreate(d *schema.ResourceData, meta interface{}) error {
	createRequest := &packngo.ProjectCreateRequest{
		Name:          d.Get("name").(string),
		PaymentMethod: d.Get("payment_method").(string),
	}

	project, err := findOrCreateProject(meta.(*packngo.Client), createRequest)
	if err != nil {
		return err
	}

	d.SetId(project.ID)

	return resourcePacketProjectRead(d, meta)
}

func resourcePacketProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	if d.Id() == "" {
		projects, _, err := client.Projects.List()
		if err != nil {
			return friendlyError(err)
		}
		name := d.Get("name").(string)
		for _, project := range projects {
			if project.Name == name {
				d.Set("id", project.ID)
				d.Set("created", project.Created)
				d.Set("updated", project.Updated)
				return nil
			}
		}
		return errors.New("no project id")
	}

	project, _, err := client.Projects.Get(d.Id())
	if err != nil {
		err = friendlyError(err)

		// If the project somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")

			return nil
		}

		return err
	}

	d.Set("id", project.ID)
	d.Set("name", project.Name)
	d.Set("created", project.Created)
	d.Set("updated", project.Updated)

	return nil
}

func resourcePacketProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	updateRequest := &packngo.ProjectUpdateRequest{
		ID:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("payment_method"); ok {
		updateRequest.PaymentMethod = attr.(string)
	}

	_, _, err := client.Projects.Update(updateRequest)
	if err != nil {
		return friendlyError(err)
	}

	return resourcePacketProjectRead(d, meta)
}

func resourcePacketProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, err := client.Projects.Delete(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}

func findOrCreateProject(client *packngo.Client, createRequest *packngo.ProjectCreateRequest) (*packngo.Project, error) {
	projects, _, err := client.Projects.List()
	if err != nil {
		return nil, friendlyError(err)
	}
	for _, project := range projects {
		if project.Name == createRequest.Name {
			return &project, nil
		}
	}

	project, _, err := client.Projects.Create(createRequest)
	if err != nil {
		return nil, friendlyError(err)
	}
	return project, nil
}
