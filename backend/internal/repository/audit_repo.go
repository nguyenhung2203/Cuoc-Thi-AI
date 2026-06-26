package repository

import (
	"log"

	"backend/internal/models"
)

type AuditRepository struct {
	// Data access fields
}

func NewAuditRepository() *AuditRepository {
	return &AuditRepository{}
}

// Insert simulates inserting an audit log into the database
func (r *AuditRepository) Insert(al *models.AuditLog) error {
	log.Printf("[db] INSERT INTO audit_logs (action, actor_id, resource_id, metadata, ip_address) VALUES ('%s', '%s', '%s', '%v', '%s')",
		al.Action, al.ActorID, al.ResourceID, al.Metadata, al.IPAddress)
	return nil
}
