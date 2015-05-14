package aws

import(
    "fmt"
    "testing"

    "github.com/awslabs/aws-sdk-go/aws"
    "github.com/awslabs/aws-sdk-go/service/autoscaling"
    "github.com/hashicorp/terraform/helper/resource"
    "github.com/hashicorp/terraform/terraform"
)

func TestAccAWSAutoscalingScheduledUpdateGroupAction_basic(t *testing.T) {
    var scheduledAction autoscaling.ScheduledUpdateGroupAction

    resource.Test(t, resource.TestCase{
        PreCheck:       func () { testAccPreCheck(t) },
        Providers:      testAccProviders,
        CheckDestroy:   testAccCheckAWSAutoscalingScheduleUpdateGroupActionDestroy,
        Steps:          []resource.TestStep{
            resource.TestStep{
                Config: testAccAWSAutoscalingScheduleUpdateGroupActionConfig,
                Check:  resource.ComposeTestCheckFunc(
                    testAccCheckScheduleUpdateGroupActionExists("aws_autoscaling_scheduled_update_group_action.foobar", &scheduledAction),
                    resource.TestCheckResourceAttr("aws_autoscaling_scheduled_update_group_action.foobar", "recurrence", "*/5 * * * *"),
                    resource.TestCheckResourceAttr("aws_autoscaling_scheduled_update_group_action.foobar", "desired_capacity", "2"),
                ),
            },
        },
    })
}

func testAccCheckScheduleUpdateGroupActionExists(n string, scheduledAction *autoscaling.ScheduledUpdateGroupAction) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[n]
        if !ok {
            rs = rs
            return fmt.Errorf("Not found: %s", n)
        }
        conn := testAccProvider.Meta().(*AWSClient).autoscalingconn
        params := &autoscaling.DescribeScheduledActionsInput{
            AutoScalingGroupName: aws.String(rs.Primary.Attributes["autoscaling_group_name"]),
            ScheduledActionNames: []*string{aws.String(rs.Primary.ID)},
        }
        resp, err := conn.DescribeScheduledActions(params)
        if err != nil {
            return err
        }
        if len(resp.ScheduledUpdateGroupActions) == 0 {
            return fmt.Errorf("ScheduledUpdateGroupAction not found")
        }

        *scheduledAction = *resp.ScheduledUpdateGroupActions[0]

        return nil
    }
}

func testAccCheckAWSAutoscalingScheduleUpdateGroupActionDestroy(s *terraform.State) error {
    conn := testAccProvider.Meta().(*AWSClient).autoscalingconn

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "aws_autoscaling_group" {
            continue
        }

        params := autoscaling.DescribeScheduledActionsInput{
            AutoScalingGroupName: aws.String(rs.Primary.Attributes["autoscaling_group_name"]),
            ScheduledActionNames: []*string{aws.String(rs.Primary.ID)},
        }

        resp, err := conn.DescribeScheduledActions(&params)

        if err == nil {
            if len(resp.ScheduledUpdateGroupActions) != 0 &&
                *resp.ScheduledUpdateGroupActions[0].ScheduledActionName == rs.Primary.ID {
                    return fmt.Errorf("Scheduled Update Action Group Still Exists: %s", rs.Primary.ID)
            }
        }
    }

    return nil
}

var testAccAWSAutoscalingScheduleUpdateGroupActionConfig = `
resource "aws_launch_configuration" "foobar" {
    name = "terraform-test-foobar5"
    image_id = "ami-21f78e11"
    instance_type = "t1.micro"
}

resource "aws_autoscaling_group" "foobar" {
    availability_zones = ["us-west-2a"]
    name = "terraform-test-foobar5"
    max_size = 5
    min_size = 2
    health_check_grace_period = 300
    health_check_type = "ELB"
    desired_capacity = 4
    force_delete = true
    termination_policies = ["OldestInstance"]
    launch_configuration = "${aws_launch_configuration.foobar.name}"
    tag {
        key = "Foo"
        value = "foo-bar"
        propagate_at_launch = true
    }
}

resource "aws_autoscaling_scheduled_update_group_action" "foobar" {
    autoscaling_group_name = "${aws_autoscaling_group.foobar.name}"
    scheduled_action_name = "foobar3-terraform-test"
    recurrence = "*/5 * * * *"
    desired_capacity = 2
}
`
