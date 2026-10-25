package main

import "fmt"

type MaintenanceRequest struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Summary  string `json:"summary"`
}
type TenantDocument struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Kind     string `json:"kind"`
}
type InspectionReminder struct {
	ID         string `json:"id"`
	PropertyID string `json:"property_id"`
	DueDate    string `json:"due_date"`
}

func ErrorGroup(kind, id string) string { return kind + ":" + id }
func ValidateMaintenance(r MaintenanceRequest) error {
	if r.ID == "" || r.TenantID == "" || r.Summary == "" {
		return fmt.Errorf("maintenance request is incomplete")
	}
	return nil
}
func CaptureMaintenance(c *InfraiClient, r MaintenanceRequest) error {
	if err := ValidateMaintenance(r); err != nil {
		return c.Capture(map[string]any{"title": "maintenance request rejected", "message": err.Error(), "level": "error", "fingerprint": []string{ErrorGroup("maintenance", r.ID)}, "exception": err.Error(), "context": map[string]any{"request": r}})
	}
	return nil
}
