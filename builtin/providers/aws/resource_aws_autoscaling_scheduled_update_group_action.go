package aws

import (
	"fmt"
	"log"
    "time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/autoscaling"
)

func resourceAwsAutoscalingScheduledUpdateGroupAction() *schema.Resource {
    return &schema.Resource{
        Create: resourceAwsAutoscalingScheduledUpdateGroupActionCreate,
        Read:   resourceAwsAutoscalingScheduledUpdateGroupActionRead,
        Update: resourceAwsAutoscalingScheduledUpdateGroupActionUpdate,
        Delete: resourceAwsAutoscalingScheduledUpdateGroupActionDelete,

        Schema: map[string]*schema.Schema{
            "autoscaling_group_name": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
            },
            "desired_capacity": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
            },
            "end_time": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
            },
            "max_size": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
            },
            "min_size": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
            },
            "recurrence": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
            },
            "scheduled_action_name": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
            },
            "start_time": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
            },
        },
    }
}

func resourceAwsAutoscalingScheduledUpdateGroupActionCreate(d *schema.ResourceData, meta interface{}) error {
    autoscalingconn := meta.(*AWSClient).autoscalingconn

    params := getAwsAutoscalingScheduledUpdateGroupActionInput(d)

    log.Printf("[DEBUG] AutoScaling PutScheduledUpdateGroupAction: %#v", params)
    _, err := autoscalingconn.PutScheduledUpdateGroupAction(&params)
    if err != nil {
        return fmt.Errorf("Error putting scheduled update group action: %s", err)
    }

    d.SetId(d.Get("scheduled_action_name").(string))

    return resourceAwsAutoscalingScheduledUpdateGroupActionRead(d, meta)
}

func resourceAwsAutoscalingScheduledUpdateGroupActionRead(d *schema.ResourceData, meta interface{}) error {
    p, err := getAwsAutoscalingScheduledUpdateGroupAction(d, meta)
    if err != nil {
        return err
    }
    if p == nil {
        return nil
    }

    log.Printf("[DEBUG] Read Scheduled Update Group Action")
    d.Set("autoscaling_group_name", p.AutoScalingGroupName)
    d.Set("desired_capacity", p.DesiredCapacity)
    d.Set("end_time", p.EndTime)
    d.Set("max_size", p.MaxSize)
    d.Set("min_size", p.MinSize)
    d.Set("recurrence", p.Recurrence)
    d.Set("scheduled_action_name", p.ScheduledActionName)
    d.Set("start_time", p.StartTime)

    return nil
}

func resourceAwsAutoscalingScheduledUpdateGroupActionUpdate(d *schema.ResourceData, meta interface{}) error {
    autoscalingconn := meta.(*AWSClient).autoscalingconn

    params := getAwsAutoscalingScheduledUpdateGroupActionInput(d)

    log.Printf("[DEBUG] AutoScaling PutScheduledUpdateGroupAction: %#v", params)
    _, err := autoscalingconn.PutScheduledUpdateGroupAction(&params)
    if err != nil {
        return fmt.Errorf("Error putting scheduled update group action: %s", err)
    }

    return resourceAwsAutoscalingScheduledUpdateGroupActionRead(d, meta)
}

func resourceAwsAutoscalingScheduledUpdateGroupActionDelete(d *schema.ResourceData, meta interface{}) error {
    autoscalingconn := meta.(*AWSClient).autoscalingconn
    p, err := getAwsAutoscalingScheduledUpdateGroupAction(d, meta)
    if err != nil {
        return err
    }
    if p == nil {
        return nil
    }

    params := autoscaling.DeleteScheduledActionInput{
        AutoScalingGroupName:   aws.String(d.Get("autoscaling_group_name").(string)),
        ScheduledActionName:    aws.String(d.Get("scheduled_action_name").(string)),
    }
    if _, err := autoscalingconn.DeleteScheduledAction(&params); err != nil {
        return fmt.Errorf("Autoscaling Scheduled Action Delete: %s", err)
    }

    d.SetId("")
    return nil
}

func getAwsAutoscalingScheduledUpdateGroupActionInput(d *schema.ResourceData) autoscaling.PutScheduledUpdateGroupActionInput {
    var params = autoscaling.PutScheduledUpdateGroupActionInput{
        AutoScalingGroupName:   aws.String(d.Get("autoscaling_group_name").(string)),
        ScheduledActionName:    aws.String(d.Get("scheduled_action_name").(string)),
    }

    if v, ok := d.GetOk("desired_capacity"); ok {
        params.DesiredCapacity = aws.Long(int64(v.(int)))
    }
    if v, ok := d.GetOk("end_time"); ok {
        t, err := time.Parse(time.RFC3339, v.(string))
        if err != nil {
            // handle error
        }
        params.EndTime = aws.Time(t)
    }
    if v, ok := d.GetOk("max_size"); ok {
        params.MaxSize = aws.Long(int64(v.(int)))
    }
    if v, ok := d.GetOk("min_size"); ok {
        params.MinSize = aws.Long(int64(v.(int)))
    }
    if v, ok := d.GetOk("recurrence"); ok {
        params.Recurrence = aws.String(v.(string))
    }
    if v, ok := d.GetOk("start_time"); ok {
        t, err := time.Parse(time.RFC3339, v.(string))
        if err != nil {
            // handle error
        }
        params.StartTime = aws.Time(t)
    }

    return params
}

func getAwsAutoscalingScheduledUpdateGroupAction(d *schema.ResourceData, meta interface{}) (*autoscaling.ScheduledUpdateGroupAction, error) {
    autoscalingconn := meta.(*AWSClient).autoscalingconn

    params := autoscaling.DescribeScheduledActionsInput{
        AutoScalingGroupName: aws.String(d.Get("autoscaling_group_name").(string)),
        ScheduledActionNames: []*string{aws.String(d.Get("scheduled_action_name").(string))},
    }

    log.Printf("[DEBUG] AutoScaling Scheduled Action Describe Params: %#v", params)
    resp, err := autoscalingconn.DescribeScheduledActions(&params)
    if err != nil {
        return nil, fmt.Errorf("Error retrieving scheduled actions: %s", err)
    }

    // find scheduled action
    scheduled_action_name := d.Get("scheduled_action_name")
    for idx, suga := range resp.ScheduledUpdateGroupActions{
        if *suga.ScheduledActionName == scheduled_action_name {
            return resp.ScheduledUpdateGroupActions[idx], nil
        }
    }

    d.SetId("")
    return nil, nil
}
