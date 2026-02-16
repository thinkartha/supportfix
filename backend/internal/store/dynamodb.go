package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/supporttickr/backend/internal/config"
	"github.com/supporttickr/backend/internal/models"
)

type DynamoStore struct {
	client            *dynamodb.Client
	usersTable        string
	orgsTable         string
	ticketsTable      string
	messagesTable     string
	timeEntriesTable  string
	conversionTable   string
	invoicesTable     string
	activitiesTable   string
}

// NewStore creates a DynamoDB store from app config (uses default AWS config).
func NewStore(ctx context.Context, cfg *config.Config) (Store, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(awsCfg)
	return &DynamoStore{
		client:            client,
		usersTable:        cfg.UsersTable,
		orgsTable:         cfg.OrgsTable,
		ticketsTable:      cfg.TicketsTable,
		messagesTable:     cfg.MessagesTable,
		timeEntriesTable:  cfg.TimeEntriesTable,
		conversionTable:   cfg.ConversionRequestsTable,
		invoicesTable:     cfg.InvoicesTable,
		activitiesTable:   cfg.ActivitiesTable,
	}, nil
}

func NewDynamoStore(cfg *DynamoConfig) (*DynamoStore, error) {
	client, err := cfg.DynamoDBClient(context.TODO())
	if err != nil {
		return nil, err
	}
	return &DynamoStore{
		client:            client,
		usersTable:        cfg.UsersTable,
		orgsTable:         cfg.OrgsTable,
		ticketsTable:      cfg.TicketsTable,
		messagesTable:     cfg.MessagesTable,
		timeEntriesTable:  cfg.TimeEntriesTable,
		conversionTable:   cfg.ConversionRequestsTable,
		invoicesTable:     cfg.InvoicesTable,
		activitiesTable:   cfg.ActivitiesTable,
	}, nil
}

type DynamoConfig struct {
	UsersTable             string
	OrgsTable              string
	TicketsTable           string
	MessagesTable          string
	TimeEntriesTable       string
	ConversionRequestsTable string
	InvoicesTable          string
	ActivitiesTable        string
	Region                 string
	DynamoDBClient         func(context.Context) (*dynamodb.Client, error)
}

func timeToStr(t time.Time) string { return t.UTC().Format(time.RFC3339) }
func strToTime(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) }

// --- Users ---
func (s *DynamoStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.usersTable),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(out.Items) == 0 {
		return nil, nil
	}
	return itemToUser(out.Items[0])
}

func (s *DynamoStore) GetUser(ctx context.Context, id string) (*models.User, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.usersTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	return itemToUser(out.Item)
}

func (s *DynamoStore) ListUsers(ctx context.Context, role, orgID string) ([]models.UserResponse, error) {
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(s.usersTable)})
	if err != nil {
		return nil, err
	}
	var users []models.UserResponse
	for _, item := range out.Items {
		u, err := itemToUser(item)
		if err != nil {
			continue
		}
		if role == "client" && orgID != "" {
			uOrg := ""
			if u.OrganizationID != nil {
				uOrg = *u.OrganizationID
			}
			if uOrg != orgID && uOrg != "" {
				continue
			}
		}
		users = append(users, u.ToResponse())
	}
	return users, nil
}

func (s *DynamoStore) CreateUser(ctx context.Context, u *models.User) error {
	item := map[string]types.AttributeValue{
		"id":            &types.AttributeValueMemberS{Value: u.ID},
		"name":          &types.AttributeValueMemberS{Value: u.Name},
		"email":         &types.AttributeValueMemberS{Value: u.Email},
		"password_hash": &types.AttributeValueMemberS{Value: u.PasswordHash},
		"role":          &types.AttributeValueMemberS{Value: u.Role},
		"avatar":        &types.AttributeValueMemberS{Value: u.Avatar},
		"created_at":   &types.AttributeValueMemberS{Value: timeToStr(u.CreatedAt)},
	}
	if u.OrganizationID != nil && *u.OrganizationID != "" {
		item["organization_id"] = &types.AttributeValueMemberS{Value: *u.OrganizationID}
	}
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.usersTable),
		Item:      item,
	})
	return err
}

func (s *DynamoStore) UpdateUser(ctx context.Context, id string, name, email, role, orgID, avatar *string) error {
	var setParts []string
	attrs := map[string]types.AttributeValue{}
	names := map[string]string{"#n": "name", "#r": "role"}

	if name != nil {
		setParts = append(setParts, "#n = :name", "avatar = :avatar")
		attrs[":name"] = &types.AttributeValueMemberS{Value: *name}
		initials := ""
		for _, w := range strings.Fields(*name) {
			if len(w) > 0 {
				initials += string([]rune(w)[0])
			}
		}
		if len(initials) > 2 {
			initials = initials[:2]
		}
		attrs[":avatar"] = &types.AttributeValueMemberS{Value: strings.ToUpper(initials)}
	}
	if email != nil {
		setParts = append(setParts, "email = :email")
		attrs[":email"] = &types.AttributeValueMemberS{Value: *email}
	}
	if role != nil {
		setParts = append(setParts, "#r = :role")
		attrs[":role"] = &types.AttributeValueMemberS{Value: *role}
	}
	if orgID != nil {
		if *orgID == "" {
			// Remove organization_id when clearing (e.g. user changed to admin)
			// Handled below with REMOVE
		} else {
			setParts = append(setParts, "organization_id = :org_id")
			attrs[":org_id"] = &types.AttributeValueMemberS{Value: *orgID}
		}
	}
	if avatar != nil && name == nil {
		// Only add avatar separately when name wasn't updated (name update already sets avatar from initials)
		setParts = append(setParts, "avatar = :avatar")
		attrs[":avatar"] = &types.AttributeValueMemberS{Value: *avatar}
	}

	var updateExpr string
	if len(setParts) > 0 {
		updateExpr = "SET " + strings.Join(setParts, ", ")
	}
	if orgID != nil && *orgID == "" {
		if updateExpr != "" {
			updateExpr += " REMOVE organization_id"
		} else {
			updateExpr = "REMOVE organization_id"
		}
	}
	if updateExpr == "" {
		return nil
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.usersTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String(updateExpr),
	}
	if len(names) > 0 {
		input.ExpressionAttributeNames = names
	}
	if len(attrs) > 0 {
		input.ExpressionAttributeValues = attrs
	}
	_, err := s.client.UpdateItem(ctx, input)
	return err
}

func (s *DynamoStore) DeleteUser(ctx context.Context, id string) error {
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.usersTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}

func itemToUser(item map[string]types.AttributeValue) (*models.User, error) {
	var orgID *string
	if v, ok := item["organization_id"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			orgID = &s.Value
		}
	}
	createdAt := time.Time{}
	if v, ok := item["created_at"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			t, _ := time.Parse(time.RFC3339, s.Value)
			createdAt = t
		}
	}
	return &models.User{
		ID:             getStr(item, "id"),
		Name:           getStr(item, "name"),
		Email:          getStr(item, "email"),
		PasswordHash:   getStr(item, "password_hash"),
		Role:           getStr(item, "role"),
		OrganizationID: orgID,
		Avatar:         getStr(item, "avatar"),
		CreatedAt:      createdAt,
	}, nil
}

func getStr(item map[string]types.AttributeValue, key string) string {
	if v, ok := item[key]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			return s.Value
		}
	}
	return ""
}

func getNum(item map[string]types.AttributeValue, key string) float64 {
	if v, ok := item[key]; ok {
		if n, ok := v.(*types.AttributeValueMemberN); ok {
			var f float64
			fmt.Sscanf(n.Value, "%f", &f)
			return f
		}
	}
	return 0
}

func getInt(item map[string]types.AttributeValue, key string) int {
	if v, ok := item[key]; ok {
		if n, ok := v.(*types.AttributeValueMemberN); ok {
			var i int
			fmt.Sscanf(n.Value, "%d", &i)
			return i
		}
	}
	return 0
}

// --- Orgs ---
func (s *DynamoStore) ListOrgs(ctx context.Context, role, orgID string) ([]models.Organization, error) {
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(s.orgsTable)})
	if err != nil {
		return nil, err
	}
	var list []models.Organization
	for _, item := range out.Items {
		o, err := itemToOrg(item)
		if err != nil {
			continue
		}
		if role == "client" && orgID != "" && o.ID != orgID {
			continue
		}
		list = append(list, *o)
	}
	return list, nil
}

func (s *DynamoStore) GetOrg(ctx context.Context, id string) (*models.Organization, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.orgsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	return itemToOrg(out.Item)
}

func (s *DynamoStore) CreateOrg(ctx context.Context, o *models.Organization) error {
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.orgsTable),
		Item: map[string]types.AttributeValue{
			"id":            &types.AttributeValueMemberS{Value: o.ID},
			"name":          &types.AttributeValueMemberS{Value: o.Name},
			"plan":          &types.AttributeValueMemberS{Value: o.Plan},
			"contact_email": &types.AttributeValueMemberS{Value: o.ContactEmail},
			"created_at":    &types.AttributeValueMemberS{Value: timeToStr(o.CreatedAt)},
		},
	})
	return err
}

func (s *DynamoStore) UpdateOrg(ctx context.Context, id string, name, plan, contactEmail *string) error {
	expr := "SET "
	attrs := map[string]types.AttributeValue{}
	if name != nil {
		expr += " #n = :name, "
		attrs[":name"] = &types.AttributeValueMemberS{Value: *name}
	}
	if plan != nil {
		expr += " plan = :plan, "
		attrs[":plan"] = &types.AttributeValueMemberS{Value: *plan}
	}
	if contactEmail != nil {
		expr += " contact_email = :ce, "
		attrs[":ce"] = &types.AttributeValueMemberS{Value: *contactEmail}
	}
	expr = strings.TrimSuffix(strings.TrimSuffix(expr, ", "), ", ")
	if expr == "SET " {
		return nil
	}
	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.orgsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          aws.String(expr),
		ExpressionAttributeNames:  map[string]string{"#n": "name"},
		ExpressionAttributeValues: attrs,
	})
	return err
}

func (s *DynamoStore) DeleteOrg(ctx context.Context, id string) error {
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.orgsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}

func itemToOrg(item map[string]types.AttributeValue) (*models.Organization, error) {
	createdAt := time.Time{}
	if v, ok := item["created_at"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			t, _ := time.Parse(time.RFC3339, s.Value)
			createdAt = t
		}
	}
	return &models.Organization{
		ID:           getStr(item, "id"),
		Name:         getStr(item, "name"),
		Plan:         getStr(item, "plan"),
		ContactEmail: getStr(item, "contact_email"),
		CreatedAt:    createdAt,
	}, nil
}

// --- Tickets ---
func (s *DynamoStore) ListTickets(ctx context.Context, status, priority, category, organizationID, assignedTo, search string) ([]models.Ticket, error) {
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(s.ticketsTable)})
	if err != nil {
		return nil, err
	}
	var list []models.Ticket
	for _, item := range out.Items {
		t, err := itemToTicket(item)
		if err != nil {
			continue
		}
		if status != "" && t.Status != status {
			continue
		}
		if priority != "" && t.Priority != priority {
			continue
		}
		if category != "" && t.Category != category {
			continue
		}
		if organizationID != "" && t.OrganizationID != organizationID {
			continue
		}
		if assignedTo != "" {
			a := ""
			if t.AssignedTo != nil {
				a = *t.AssignedTo
			}
			if a != assignedTo {
				continue
			}
		}
		if search != "" {
			if !strings.Contains(strings.ToLower(t.Title), strings.ToLower(search)) &&
				!strings.Contains(strings.ToLower(t.Description), strings.ToLower(search)) {
				continue
			}
		}
		list = append(list, *t)
	}
	return list, nil
}

func (s *DynamoStore) GetTicket(ctx context.Context, id string) (*models.Ticket, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.ticketsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	return itemToTicket(out.Item)
}

func (s *DynamoStore) CreateTicket(ctx context.Context, t *models.Ticket) error {
	item := map[string]types.AttributeValue{
		"id":              &types.AttributeValueMemberS{Value: t.ID},
		"title":           &types.AttributeValueMemberS{Value: t.Title},
		"description":     &types.AttributeValueMemberS{Value: t.Description},
		"status":          &types.AttributeValueMemberS{Value: t.Status},
		"priority":        &types.AttributeValueMemberS{Value: t.Priority},
		"category":        &types.AttributeValueMemberS{Value: t.Category},
		"organization_id": &types.AttributeValueMemberS{Value: t.OrganizationID},
		"created_by":      &types.AttributeValueMemberS{Value: t.CreatedBy},
		"hours_worked":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", t.HoursWorked)},
		"created_at":      &types.AttributeValueMemberS{Value: timeToStr(t.CreatedAt)},
		"updated_at":      &types.AttributeValueMemberS{Value: timeToStr(t.UpdatedAt)},
	}
	if t.AssignedTo != nil {
		item["assigned_to"] = &types.AttributeValueMemberS{Value: *t.AssignedTo}
	}
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.ticketsTable),
		Item:      item,
	})
	return err
}

func (s *DynamoStore) UpdateTicket(ctx context.Context, id string, status, priority, assignedTo *string, hoursWorked *float64) error {
	expr := "SET updated_at = :ua"
	attrs := map[string]types.AttributeValue{
		":ua": &types.AttributeValueMemberS{Value: timeToStr(time.Now().UTC())},
	}
	if status != nil {
		expr += ", #st = :status"
		attrs[":status"] = &types.AttributeValueMemberS{Value: *status}
	}
	if priority != nil {
		expr += ", priority = :priority"
		attrs[":priority"] = &types.AttributeValueMemberS{Value: *priority}
	}
	if assignedTo != nil {
		if *assignedTo == "" {
			expr = strings.TrimSuffix(expr, ", ")
			expr += " REMOVE assigned_to"
		} else {
			expr += ", assigned_to = :assigned_to"
			attrs[":assigned_to"] = &types.AttributeValueMemberS{Value: *assignedTo}
		}
	}
	if hoursWorked != nil {
		expr += ", hours_worked = :hw"
		attrs[":hw"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", *hoursWorked)}
	}
	names := map[string]string{"#st": "status"}
	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(s.ticketsTable),
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression:          aws.String(expr),
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: attrs,
	})
	return err
}

func (s *DynamoStore) UpdateTicketCategory(ctx context.Context, id, category string) error {
	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.ticketsTable),
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression: aws.String("SET category = :cat, updated_at = :ua"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cat": &types.AttributeValueMemberS{Value: category},
			":ua":  &types.AttributeValueMemberS{Value: timeToStr(time.Now().UTC())},
		},
	})
	return err
}

func itemToTicket(item map[string]types.AttributeValue) (*models.Ticket, error) {
	var assignedTo *string
	if v, ok := item["assigned_to"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			assignedTo = &s.Value
		}
	}
	createdAt, _ := time.Parse(time.RFC3339, getStr(item, "created_at"))
	updatedAt, _ := time.Parse(time.RFC3339, getStr(item, "updated_at"))
	return &models.Ticket{
		ID:             getStr(item, "id"),
		Title:          getStr(item, "title"),
		Description:    getStr(item, "description"),
		Status:         getStr(item, "status"),
		Priority:       getStr(item, "priority"),
		Category:       getStr(item, "category"),
		OrganizationID: getStr(item, "organization_id"),
		CreatedBy:      getStr(item, "created_by"),
		AssignedTo:     assignedTo,
		HoursWorked:    getNum(item, "hours_worked"),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

// --- Messages ---
func (s *DynamoStore) GetMessagesByTicketID(ctx context.Context, ticketID string) ([]models.Message, error) {
	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.messagesTable),
		KeyConditionExpression: aws.String("ticket_id = :tid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":tid": &types.AttributeValueMemberS{Value: ticketID},
		},
	})
	if err != nil {
		return nil, err
	}
	var list []models.Message
	for _, item := range out.Items {
		m, err := itemToMessage(item)
		if err != nil {
			continue
		}
		list = append(list, *m)
	}
	return list, nil
}

func (s *DynamoStore) AddMessage(ctx context.Context, m *models.Message) error {
	internal := "false"
	if m.IsInternal {
		internal = "true"
	}
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.messagesTable),
		Item: map[string]types.AttributeValue{
			"ticket_id":   &types.AttributeValueMemberS{Value: m.TicketID},
			"id":          &types.AttributeValueMemberS{Value: m.ID},
			"user_id":     &types.AttributeValueMemberS{Value: m.UserID},
			"content":     &types.AttributeValueMemberS{Value: m.Content},
			"is_internal": &types.AttributeValueMemberS{Value: internal},
			"created_at":  &types.AttributeValueMemberS{Value: timeToStr(m.CreatedAt)},
		},
	})
	return err
}

func itemToMessage(item map[string]types.AttributeValue) (*models.Message, error) {
	createdAt, _ := time.Parse(time.RFC3339, getStr(item, "created_at"))
	internal := getStr(item, "is_internal") == "true"
	return &models.Message{
		ID:         getStr(item, "id"),
		TicketID:   getStr(item, "ticket_id"),
		UserID:     getStr(item, "user_id"),
		Content:    getStr(item, "content"),
		IsInternal: internal,
		CreatedAt:  createdAt,
	}, nil
}

// --- Time entries ---
func (s *DynamoStore) GetTimeEntriesByTicketID(ctx context.Context, ticketID string) ([]models.TimeEntry, error) {
	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.timeEntriesTable),
		KeyConditionExpression: aws.String("ticket_id = :tid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":tid": &types.AttributeValueMemberS{Value: ticketID},
		},
	})
	if err != nil {
		return nil, err
	}
	var list []models.TimeEntry
	for _, item := range out.Items {
		te, _ := itemToTimeEntry(item)
		list = append(list, *te)
	}
	return list, nil
}

func (s *DynamoStore) AddTimeEntry(ctx context.Context, te *models.TimeEntry) error {
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.timeEntriesTable),
		Item: map[string]types.AttributeValue{
			"ticket_id":   &types.AttributeValueMemberS{Value: te.TicketID},
			"id":          &types.AttributeValueMemberS{Value: te.ID},
			"user_id":     &types.AttributeValueMemberS{Value: te.UserID},
			"hours":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", te.Hours)},
			"description": &types.AttributeValueMemberS{Value: te.Description},
			"entry_date":  &types.AttributeValueMemberS{Value: te.Date},
			"created_at":  &types.AttributeValueMemberS{Value: timeToStr(te.CreatedAt)},
		},
	})
	return err
}

func itemToTimeEntry(item map[string]types.AttributeValue) (*models.TimeEntry, error) {
	createdAt, _ := time.Parse(time.RFC3339, getStr(item, "created_at"))
	return &models.TimeEntry{
		ID:          getStr(item, "id"),
		TicketID:    getStr(item, "ticket_id"),
		UserID:      getStr(item, "user_id"),
		Hours:       getNum(item, "hours"),
		Description: getStr(item, "description"),
		Date:        getStr(item, "entry_date"),
		CreatedAt:   createdAt,
	}, nil
}

// --- Conversion requests ---
func (s *DynamoStore) GetConversionByTicketID(ctx context.Context, ticketID string) (*models.ConversionRequest, error) {
	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.conversionTable),
		IndexName:              aws.String("ticket-id-index"),
		KeyConditionExpression: aws.String("ticket_id = :tid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":tid": &types.AttributeValueMemberS{Value: ticketID},
		},
	})
	if err != nil || len(out.Items) == 0 {
		return nil, nil
	}
	return itemToConversion(out.Items[0])
}

func (s *DynamoStore) GetConversionByID(ctx context.Context, id string) (*models.ConversionRequest, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.conversionTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil || out.Item == nil {
		return nil, nil
	}
	return itemToConversion(out.Item)
}

func (s *DynamoStore) CreateConversionRequest(ctx context.Context, cr *models.ConversionRequest) error {
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.conversionTable),
		Item: map[string]types.AttributeValue{
			"id":                &types.AttributeValueMemberS{Value: cr.ID},
			"ticket_id":         &types.AttributeValueMemberS{Value: cr.TicketID},
			"proposed_type":     &types.AttributeValueMemberS{Value: cr.ProposedType},
			"reason":            &types.AttributeValueMemberS{Value: cr.Reason},
			"internal_approval": &types.AttributeValueMemberS{Value: cr.InternalApproval},
			"client_approval":   &types.AttributeValueMemberS{Value: cr.ClientApproval},
			"proposed_by":       &types.AttributeValueMemberS{Value: cr.ProposedBy},
			"created_at":       &types.AttributeValueMemberS{Value: timeToStr(cr.CreatedAt)},
		},
	})
	return err
}

func (s *DynamoStore) UpdateConversionRequest(ctx context.Context, id string, internalApproval, clientApproval *string) error {
	expr := "SET "
	attrs := map[string]types.AttributeValue{}
	if internalApproval != nil {
		expr += " internal_approval = :ia, "
		attrs[":ia"] = &types.AttributeValueMemberS{Value: *internalApproval}
	}
	if clientApproval != nil {
		expr += " client_approval = :ca, "
		attrs[":ca"] = &types.AttributeValueMemberS{Value: *clientApproval}
	}
	expr = strings.TrimSuffix(strings.TrimSuffix(expr, ", "), ", ")
	if expr == "SET " {
		return nil
	}
	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(s.conversionTable),
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression:          aws.String(expr),
		ExpressionAttributeValues: attrs,
	})
	return err
}

func (s *DynamoStore) ListConversionRequestsPending(ctx context.Context) ([]models.ConversionRequest, error) {
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(s.conversionTable),
		FilterExpression: aws.String("(internal_approval = :p OR client_approval = :p)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":p": &types.AttributeValueMemberS{Value: "pending"},
		},
	})
	if err != nil {
		return nil, err
	}
	var list []models.ConversionRequest
	for _, item := range out.Items {
		cr, _ := itemToConversion(item)
		if cr != nil {
			list = append(list, *cr)
		}
	}
	return list, nil
}

func itemToConversion(item map[string]types.AttributeValue) (*models.ConversionRequest, error) {
	createdAt, _ := time.Parse(time.RFC3339, getStr(item, "created_at"))
	return &models.ConversionRequest{
		ID:               getStr(item, "id"),
		TicketID:         getStr(item, "ticket_id"),
		ProposedType:     getStr(item, "proposed_type"),
		Reason:           getStr(item, "reason"),
		InternalApproval: getStr(item, "internal_approval"),
		ClientApproval:   getStr(item, "client_approval"),
		ProposedBy:       getStr(item, "proposed_by"),
		CreatedAt:        createdAt,
	}, nil
}

// --- Invoices ---
func (s *DynamoStore) ListInvoices(ctx context.Context, role, orgID string) ([]models.Invoice, error) {
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(s.invoicesTable)})
	if err != nil {
		return nil, err
	}
	var list []models.Invoice
	for _, item := range out.Items {
		inv, _ := itemToInvoice(item)
		if role == "client" && orgID != "" && inv.OrganizationID != orgID {
			continue
		}
		list = append(list, *inv)
	}
	return list, nil
}

func (s *DynamoStore) CreateInvoice(ctx context.Context, inv *models.Invoice) error {
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.invoicesTable),
		Item: map[string]types.AttributeValue{
			"id":               &types.AttributeValueMemberS{Value: inv.ID},
			"organization_id":  &types.AttributeValueMemberS{Value: inv.OrganizationID},
			"month":            &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", inv.Month)},
			"year":             &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", inv.Year)},
			"tickets_closed":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", inv.TicketsClosed)},
			"total_hours":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", inv.TotalHours)},
			"rate_per_hour":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", inv.RatePerHour)},
			"total_amount":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", inv.TotalAmount)},
			"status":          &types.AttributeValueMemberS{Value: inv.Status},
			"created_at":      &types.AttributeValueMemberS{Value: timeToStr(inv.CreatedAt)},
		},
	})
	return err
}

func (s *DynamoStore) UpdateInvoiceStatus(ctx context.Context, id, status string) error {
	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.invoicesTable),
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression: aws.String("SET #st = :status"),
		ExpressionAttributeNames: map[string]string{"#st": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: status},
		},
	})
	return err
}

func itemToInvoice(item map[string]types.AttributeValue) (*models.Invoice, error) {
	createdAt, _ := time.Parse(time.RFC3339, getStr(item, "created_at"))
	return &models.Invoice{
		ID:             getStr(item, "id"),
		OrganizationID: getStr(item, "organization_id"),
		Month:          getInt(item, "month"),
		Year:           getInt(item, "year"),
		TicketsClosed:  getInt(item, "tickets_closed"),
		TotalHours:     getNum(item, "total_hours"),
		RatePerHour:    getNum(item, "rate_per_hour"),
		TotalAmount:    getNum(item, "total_amount"),
		Status:         getStr(item, "status"),
		CreatedAt:      createdAt,
	}, nil
}

// --- Activities ---
func (s *DynamoStore) ListActivities(ctx context.Context, limit int) ([]models.ActivityItem, error) {
	if limit <= 0 {
		limit = 50
	}
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.activitiesTable),
		Limit:     aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, err
	}
	var list []models.ActivityItem
	for _, item := range out.Items {
		a, _ := itemToActivity(item)
		list = append(list, *a)
	}
	return list, nil
}

func (s *DynamoStore) CreateActivity(ctx context.Context, a *models.ActivityItem) error {
	item := map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: a.ID},
		"type":        &types.AttributeValueMemberS{Value: a.Type},
		"description": &types.AttributeValueMemberS{Value: a.Description},
		"user_id":     &types.AttributeValueMemberS{Value: a.UserID},
		"created_at":  &types.AttributeValueMemberS{Value: timeToStr(a.CreatedAt)},
	}
	if a.TicketID != nil {
		item["ticket_id"] = &types.AttributeValueMemberS{Value: *a.TicketID}
	}
	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.activitiesTable),
		Item:      item,
	})
	return err
}

func itemToActivity(item map[string]types.AttributeValue) (*models.ActivityItem, error) {
	var ticketID *string
	if v, ok := item["ticket_id"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			ticketID = &s.Value
		}
	}
	createdAt, _ := time.Parse(time.RFC3339, getStr(item, "created_at"))
	return &models.ActivityItem{
		ID:          getStr(item, "id"),
		Type:        getStr(item, "type"),
		Description: getStr(item, "description"),
		UserID:      getStr(item, "user_id"),
		TicketID:    ticketID,
		CreatedAt:   createdAt,
	}, nil
}
