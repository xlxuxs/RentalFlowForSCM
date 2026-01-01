package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/inventory-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoItemRepository implements ItemRepository using MongoDB
type MongoItemRepository struct {
	coll *mongo.Collection
}

// NewMongoItemRepository creates a new MongoDB item repository
func NewMongoItemRepository(db *mongo.Database) *MongoItemRepository {
	return &MongoItemRepository{
		coll: db.Collection("rental_items"),
	}
}

func (r *MongoItemRepository) Create(ctx context.Context, item *domain.RentalItem) error {
	_, err := r.coll.InsertOne(ctx, item)
	return err
}

func (r *MongoItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RentalItem, error) {
	var item domain.RentalItem
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *MongoItemRepository) GetByOwner(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]*domain.RentalItem, int, error) {
	filter := bson.M{"owner_id": ownerID}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var items []*domain.RentalItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, 0, err
	}

	return items, int(total), nil
}

func (r *MongoItemRepository) List(ctx context.Context, offset, limit int, filters ItemFilters) ([]*domain.RentalItem, int, error) {
	filter := bson.M{}

	if filters.Category != nil {
		filter["category"] = *filters.Category
	}
	if filters.City != nil {
		filter["city"] = *filters.City
	}
	if filters.IsActive != nil {
		filter["is_active"] = *filters.IsActive
	}
	// Note: MinPrice/MaxPrice implemented if referenced in filter?
	// The postgres repo did NOT implement MinPrice/MaxPrice logic in List function shown above (step 773),
	// it only handled Category, City, IsActive in `List`.
	// However, `ItemFilters` struct HAS MinPrice/MaxPrice.
	// I will implement them here for completeness.
	if filters.MinPrice != nil || filters.MaxPrice != nil {
		priceFilter := bson.M{}
		if filters.MinPrice != nil {
			priceFilter["$gte"] = *filters.MinPrice
		}
		if filters.MaxPrice != nil {
			priceFilter["$lte"] = *filters.MaxPrice
		}
		filter["daily_rate"] = priceFilter
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sort := bson.M{"created_at": -1}
	if filters.SortBy != nil {
		switch *filters.SortBy {
		case "price_low":
			sort = bson.M{"daily_rate": 1}
		case "price_high":
			sort = bson.M{"daily_rate": -1}
		case "newest":
			sort = bson.M{"created_at": -1}
		}
	}

	opts := options.Find().
		SetSort(sort).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var items []*domain.RentalItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, 0, err
	}

	return items, int(total), nil
}

func (r *MongoItemRepository) Search(ctx context.Context, query string, filters ItemFilters, offset, limit int) ([]*domain.RentalItem, int, error) {
	// Simple regex search on title or description
	regex := bson.M{"$regex": query, "$options": "i"}
	filter := bson.M{
		"$or": []bson.M{
			{"title": regex},
			{"description": regex},
		},
	}

	// Apply other filters
	if filters.Category != nil {
		filter["category"] = *filters.Category
	}
	if filters.City != nil {
		filter["city"] = *filters.City
	}
	if filters.IsActive != nil {
		filter["is_active"] = *filters.IsActive
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sort := bson.M{"created_at": -1}
	if filters.SortBy != nil {
		switch *filters.SortBy {
		case "price_low":
			sort = bson.M{"daily_rate": 1}
		case "price_high":
			sort = bson.M{"daily_rate": -1}
		case "newest":
			sort = bson.M{"created_at": -1}
		}
	}

	opts := options.Find().
		SetSort(sort).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var items []*domain.RentalItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, 0, err
	}

	return items, int(total), nil
}

func (r *MongoItemRepository) Update(ctx context.Context, item *domain.RentalItem) error {
	update := bson.M{
		"$set": bson.M{
			"title":            item.Title,
			"description":      item.Description,
			"category":         item.Category,
			"subcategory":      item.Subcategory,
			"daily_rate":       item.DailyRate,
			"weekly_rate":      item.WeeklyRate,
			"monthly_rate":     item.MonthlyRate,
			"security_deposit": item.SecurityDeposit,
			"address":          item.Address,
			"city":             item.City,
			"latitude":         item.Latitude,
			"longitude":        item.Longitude,
			"specifications":   item.Specifications,
			"images":           item.Images,
			"is_active":        item.IsActive,
			"updated_at":       time.Now(),
		},
	}

	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": item.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrItemNotFound
	}
	return nil
}

func (r *MongoItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrItemNotFound
	}
	return nil
}

// MongoAvailabilityRepository implements AvailabilityRepository using MongoDB
type MongoAvailabilityRepository struct {
	coll *mongo.Collection
}

func NewMongoAvailabilityRepository(db *mongo.Database) *MongoAvailabilityRepository {
	return &MongoAvailabilityRepository{
		coll: db.Collection("availability_slots"),
	}
}

func (r *MongoAvailabilityRepository) Create(ctx context.Context, slot *domain.AvailabilitySlot) error {
	_, err := r.coll.InsertOne(ctx, slot)
	return err
}

func (r *MongoAvailabilityRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AvailabilitySlot, error) {
	var slot domain.AvailabilitySlot
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&slot)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrSlotNotFound
		}
		return nil, err
	}
	return &slot, nil
}

func (r *MongoAvailabilityRepository) GetByItem(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time) ([]*domain.AvailabilitySlot, error) {
	filter := bson.M{
		"rental_item_id": itemID,
		"start_date":     bson.M{"$gte": startDate},
		"end_date":       bson.M{"$lte": endDate},
	}

	opts := options.Find().SetSort(bson.M{"start_date": 1})

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var slots []*domain.AvailabilitySlot
	if err := cursor.All(ctx, &slots); err != nil {
		return nil, err
	}
	return slots, nil
}

func (r *MongoAvailabilityRepository) Update(ctx context.Context, slot *domain.AvailabilitySlot) error {
	update := bson.M{
		"$set": bson.M{
			"status":     slot.Status,
			"booking_id": slot.BookingID,
		},
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": slot.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrSlotNotFound
	}
	return nil
}

func (r *MongoAvailabilityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrSlotNotFound
	}
	return nil
}

func (r *MongoAvailabilityRepository) CheckConflict(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time, excludeSlotID *uuid.UUID) (bool, error) {
	// Find any slot that overlaps and is not 'available'
	// Overlap logic: (StartA <= EndB) and (EndA >= StartB)
	filter := bson.M{
		"rental_item_id": itemID,
		"status":         bson.M{"$ne": domain.StatusAvailable},
		"start_date":     bson.M{"$lt": endDate},
		"end_date":       bson.M{"$gt": startDate},
	}

	if excludeSlotID != nil {
		filter["_id"] = bson.M{"$ne": excludeSlotID}
	}

	count, err := r.coll.CountDocuments(ctx, filter)
	return count > 0, err
}

// MongoMaintenanceRepository implements MaintenanceRepository
type MongoMaintenanceRepository struct {
	coll *mongo.Collection
}

func NewMongoMaintenanceRepository(db *mongo.Database) *MongoMaintenanceRepository {
	return &MongoMaintenanceRepository{
		coll: db.Collection("maintenance_logs"),
	}
}

func (r *MongoMaintenanceRepository) Create(ctx context.Context, log *domain.MaintenanceLog) error {
	_, err := r.coll.InsertOne(ctx, log)
	return err
}

func (r *MongoMaintenanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MaintenanceLog, error) {
	var log domain.MaintenanceLog
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&log)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrMaintenanceNotFound
		}
		return nil, err
	}
	return &log, nil
}

func (r *MongoMaintenanceRepository) GetByItem(ctx context.Context, itemID uuid.UUID, offset, limit int) ([]*domain.MaintenanceLog, int, error) {
	filter := bson.M{"rental_item_id": itemID}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*domain.MaintenanceLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil
}

func (r *MongoMaintenanceRepository) Update(ctx context.Context, log *domain.MaintenanceLog) error {
	update := bson.M{
		"$set": bson.M{
			"status":   log.Status,
			"end_date": log.EndDate,
		},
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": log.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrMaintenanceNotFound
	}
	return nil
}

func (r *MongoMaintenanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrMaintenanceNotFound
	}
	return nil
}
