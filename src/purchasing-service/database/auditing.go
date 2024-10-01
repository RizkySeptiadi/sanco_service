package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"purchasing_service/structs"

	"gorm.io/gorm"
)

const AuditTableName = "sanco_audits"

// BeforeCreate callback to register the audit callback before create
func beforeCreate(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("audit_before_create", func(tx *gorm.DB) {
		// Skip auditing if operating on the audit table itself
		if tx.Statement.Table == AuditTableName {
			return
		}
		fmt.Println("Before create callback triggered")
	})
}

// AfterCreate callback to log after record creation
func afterCreate(db *gorm.DB) {
	db.Callback().Create().After("gorm:create").Register("audit_after_create", func(tx *gorm.DB) {
		if tx.Error != nil || tx.Statement.Table == AuditTableName {
			return
		}

		// Capture the created entity's details
		if tx.Statement != nil && tx.Statement.Model != nil {
			userID := int64(123) // Replace with actual user ID extraction logic

			createdData, _ := json.Marshal(tx.Statement.Model)

			auditLog := structs.Sanco_audit{
				UserID:        userID,
				AuditableType: tx.Statement.Table, // Table name
				AuditableID:   0,                  // Replace with ID if available
				Event:         "create",
				NewValues:     string(createdData),
				CreatedAt:     time.Now(),
			}

			// Store audit log in DB
			err := tx.Session(&gorm.Session{NewDB: true}).Create(&auditLog).Error
			if err != nil {
				log.Printf("Failed to create audit log: %v\n", err)
			}
		}
	})
}

// AfterUpdate callback to log after record update

// BeforeUpdate callback to capture old values
func beforeUpdate(db *gorm.DB) {
	db.Callback().Update().Before("gorm:update").Register("audit_before_update", func(tx *gorm.DB) {
		if tx.Error != nil || tx.Statement.Table == AuditTableName {
			return
		}

		// Extract the ID from the URL or context
		// Assuming that `tx.Statement.Context` contains the request context (Gin Gonic as an example)
		id := tx.Statement.Context.Value("id").(string) // Assuming you have middleware that stores ID in context

		// Query the current record from the database using the extracted ID
		var oldRecord map[string]interface{}
		err := tx.Session(&gorm.Session{NewDB: true}).Table(tx.Statement.Table).Where("id = ?", id).Find(&oldRecord).Error
		if err != nil {
			log.Printf("Failed to fetch old record for audit log: %v\n", err)
			return
		}

		// Serialize the old record data to JSON
		oldValues, _ := json.Marshal(oldRecord)

		// Store the old values in the context for the afterUpdate callback
		tx.InstanceSet("audit_old_values", string(oldValues))

		fmt.Println("Before update callback: Fetched old values")
	})
}

// AfterUpdate callback to log after record update
func afterUpdate(db *gorm.DB) {
	db.Callback().Update().After("gorm:update").Register("audit_after_update", func(tx *gorm.DB) {
		if tx.Error != nil || tx.Statement.Table == AuditTableName {
			return
		}

		// Fetch old values from the beforeUpdate context
		var oldValues string
		if v, ok := tx.InstanceGet("audit_old_values"); ok {
			oldValues = v.(string)
		}

		// Serialize the updated data
		updatedData, _ := json.Marshal(tx.Statement.Model)
		userID := int64(123) // Replace with actual user ID extraction logic

		// Create the audit log with both old and new values
		auditLog := structs.Sanco_audit{
			UserID:        userID,
			AuditableType: tx.Statement.Table,
			AuditableID:   0, // Replace with actual ID if available
			Event:         "update",
			OldValues:     oldValues,
			NewValues:     string(updatedData),
			CreatedAt:     time.Now(),
		}

		// Store audit log in the DB
		err := tx.Session(&gorm.Session{NewDB: true}).Create(&auditLog).Error
		if err != nil {
			log.Printf("Failed to create audit log: %v\n", err)
		}
	})
}

// AfterDelete callback to log after record deletion
func afterDelete(db *gorm.DB) {
	db.Callback().Delete().After("gorm:delete").Register("audit_after_delete", func(tx *gorm.DB) {
		if tx.Error != nil || tx.Statement.Table == AuditTableName {
			return
		}

		if tx.Statement != nil && tx.Statement.Model != nil {
			deletedData, _ := json.Marshal(tx.Statement.Model)
			userID := int64(123) // Example: Replace with actual logic

			auditLog := structs.Sanco_audit{
				UserID:        userID,
				AuditableType: tx.Statement.Table,
				AuditableID:   0, // Replace with actual ID
				Event:         "delete",
				OldValues:     string(deletedData),
				CreatedAt:     time.Now(),
			}

			// Store audit log in DB
			err := tx.Session(&gorm.Session{NewDB: true}).Create(&auditLog).Error
			if err != nil {
				log.Printf("Failed to create audit log: %v\n", err)
			}
		}
	})
}

// Register all callbacks
func RegisterCallbacks(db *gorm.DB) {
	beforeCreate(db)
	afterCreate(db)
	beforeUpdate(db)
	afterUpdate(db)
	afterDelete(db)
}
