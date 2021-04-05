package main

import (
	"context"
	"fmt"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/urfave/cli/v2"
)

func pageRuleEnable(c *cli.Context) error {
	if err := checkFlags(c, "zone"); err != nil {
		return err
	}
	zoneName := c.String("zone")

	if err := checkFlags(c, "rule"); err != nil {
		return err
	}
	ruleID := c.String("rule")

	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rule, err := api.PageRule(context.Background(), zoneID, ruleID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rule.Status = "enabled"
	if err := api.ChangePageRule(context.Background(), zoneID, ruleID, rule);
	err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func pageRuleDisable(c *cli.Context) error {
	if err := checkFlags(c, "zone"); err != nil {
		return err
	}
	zoneName := c.String("zone")

	if err := checkFlags(c, "rule"); err != nil {
		return err
	}
	ruleID := c.String("rule")

	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rule, err := api.PageRule(context.Background(), zoneID, ruleID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rule.Status = "disabled"
	if err := api.ChangePageRule(context.Background(), zoneID, ruleID, rule);
	err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func pageRuleList(c *cli.Context) error {
	if err := checkFlags(c, "zone"); err != nil {
		return err
	}
	zone := c.String("zone")

	zoneID, err := api.ZoneIDByName(zone)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rules, err := api.ListPageRules(context.Background(), zoneID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("%3s %-32s %-8s %s\n", "Pri", "ID", "Status", "URL")
	for _, r := range rules {
		var settings []string
		fmt.Printf("%3d %s %-8s %s\n", r.Priority, r.ID, r.Status, r.Targets[0].Constraint.Value)
		for _, a := range r.Actions {
			var s string
			switch v := a.Value.(type) {
			case int:
				s = fmt.Sprintf("%s: %d", cloudflare.PageRuleActions[a.ID], v)
			case float64:
				s = fmt.Sprintf("%s: %.f", cloudflare.PageRuleActions[a.ID], v)
			case map[string]interface{}:
				s = fmt.Sprintf("%s: %.f - %s", cloudflare.PageRuleActions[a.ID], v["status_code"], v["url"])
			case nil:
				s = fmt.Sprintf("%s", cloudflare.PageRuleActions[a.ID])
			default:
				vs := fmt.Sprintf("%s", v)
				s = fmt.Sprintf("%s: %s", cloudflare.PageRuleActions[a.ID], strings.Title(strings.Replace(vs, "_", " ", -1)))
			}
			settings = append(settings, s)
		}
		fmt.Println("   ", strings.Join(settings, ", "))
	}

	return nil
}
