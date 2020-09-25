package storage

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"

	"github.com/mailbadger/app/entities"
)

//CreateUser creates a new user
func (db *store) CreateUser(user *entities.User) error {
	return db.Create(user).Error
}

//UpdateUser updates the given user
func (db *store) UpdateUser(user *entities.User) error {
	return db.Save(user).Error
}

//GetUser returns an active user by id. If no user is found, an error is returned
func (db *store) GetUser(id int64) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("id = ? and active = ?", id, true).First(user).Error
	return user, err
}

//GetUserByUUID returns an user by uuid. If no user is found, an error is returned
func (db *store) GetUserByUUID(uuid string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("uuid = ?", uuid).First(user).Error
	return user, err
}

//GetUserByUsername returns a user by username. If no user is found,
//an error is returned
func (db *store) GetUserByUsername(username string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("username = ?", username).First(user).Error
	return user, err
}

//GetActiveUserByUsername returns an active user by username. If no user is found,
//an error is returned
func (db *store) GetActiveUserByUsername(username string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("username = ? and active = ?", username, true).First(user).Error
	return user, err
}

// BelongsToUser finds a resource by the given user id.
func BelongsToUser(userID int64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

func (db *store) DeleteUserByID(userID int64) error {
	tx := db.BeginTx(context.Background(), &sql.TxOptions{})

	// delete * from sends
	var send entities.Send
	err := tx.Where("user_id = ?", userID).Delete(send).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// delete * from ses_keys
	var sesKeys entities.SesKeys
	err = tx.Where("user_id = ?", userID).Delete(sesKeys).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from api_keys
	var apiKeys entities.APIKey
	err = tx.Where("user_id = ?", userID).Delete(apiKeys).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from sessions
	var sessions entities.Session
	err = tx.Where("user_id = ?", userID).Delete(sessions).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from campaigns
	var campaigns entities.Campaign
	err = tx.Where("user_id = ?", userID).Delete(campaigns).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from subscribers
	var subscribers entities.Subscriber
	err = tx.Where("user_id = ?", userID).Delete(subscribers).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from segments
	var segments entities.Segment
	err = tx.Where("user_id = ?", userID).Delete(segments).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from bounces
	var bounces entities.Bounce
	err = tx.Where("user_id = ?", userID).Delete(bounces).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from complaints
	var complaints entities.Complaint
	err = tx.Where("user_id = ?", userID).Delete(complaints).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from clicks
	var clicks entities.Click
	err = tx.Where("user_id = ?", userID).Delete(clicks).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from opens
	var opens entities.Open
	err = tx.Where("user_id = ?", userID).Delete(opens).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from deliveries
	var deliveries entities.Delivery
	err = tx.Where("user_id = ?", userID).Delete(deliveries).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete * from send_bulk_logs
	var sendBulkLogs entities.SendBulkLog
	err = tx.Where("user_id = ?", userID).Delete(sendBulkLogs).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	/*// delete * from users
	var user entities.User
	err = tx.Where("id = ?", userID).Delete(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}*/

	tx.Commit()
	return nil
}
