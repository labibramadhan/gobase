package watermillsvc

import (
	"encoding/json"
	"fmt"
	"strings"

	sql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

// DestinationTopicKey is the metadata key used to store the original destination topic of a message.
const DestinationTopicKey = "destination_topic"

// BunPostgreSQLSchemaConfig holds the configuration for a multi-table Bun PostgreSQL schema.
type BunPostgreSQLSchemaConfig struct {
	// TableNames is a list of all table names that will be used for the outbox.
	TableNames []string
	// TopicToTableMap maps a topic to a specific outbox table name from the TableNames list.
	TopicToTableMap map[string]string
}

// BunPostgreSQLSchema is a schema adapter that maps topics to a predefined set of outbox tables.
type BunPostgreSQLSchema struct {
	sql.DefaultPostgreSQLSchema
	config BunPostgreSQLSchemaConfig
}

// NewBunPostgreSQLSchema creates a new BunPostgreSQLSchema and validates the configuration.
func NewBunPostgreSQLSchema(config BunPostgreSQLSchemaConfig) (*BunPostgreSQLSchema, error) {
	if len(config.TableNames) == 0 {
		return nil, errors.New("TableNames cannot be empty")
	}
	if len(config.TopicToTableMap) == 0 {
		return nil, errors.New("TopicToTableMap cannot be empty")
	}

	tableNameSet := make(map[string]struct{}, len(config.TableNames))
	for _, name := range config.TableNames {
		tableNameSet[name] = struct{}{}
	}

	for topic, tableName := range config.TopicToTableMap {
		if _, ok := tableNameSet[tableName]; !ok {
			return nil, fmt.Errorf("table '%s' for topic '%s' is not declared in TableNames", tableName, topic)
		}
	}

	schema := &BunPostgreSQLSchema{config: config}
	// Set the custom GenerateMessagesTableName function
	schema.DefaultPostgreSQLSchema.GenerateMessagesTableName = schema.customGenerateMessagesTableName
	return schema, nil
}

// customGenerateMessagesTableName returns the underlying table name for a topic.
// For other topics with mappings, it uses the mapped table name.
// For topics without mappings, it uses the default Watermill behavior.
func (s *BunPostgreSQLSchema) customGenerateMessagesTableName(topic string) string {
	// Check if we have a custom mapping for this topic
	if tableName, ok := s.config.TopicToTableMap[topic]; ok {
		return tableName
	}

	// Fallback to default behavior
	return "watermill_" + strings.ReplaceAll(topic, ".", "_")
}

// buildInsertMarkers creates the SQL placeholders for the INSERT statement with transaction_id.
func buildInsertMarkers(count int) string {
	result := strings.Builder{}

	index := 1
	for i := 0; i < count; i++ {
		// For each message: uuid, payload, metadata, and pg_current_xact_id() for transaction_id
		result.WriteString(fmt.Sprintf("($%d,$%d,$%d,pg_current_xact_id()),", index, index+1, index+2))
		index += 3
	}

	return strings.TrimRight(result.String(), ",")
}

// InsertQuery generates the SQL to insert messages into the correct outbox table.
func (s *BunPostgreSQLSchema) InsertQuery(topic string, msgs message.Messages) (string, []interface{}, error) {
	tableName, ok := s.config.TopicToTableMap[topic]
	if !ok {
		// We only want to insert into explicitly configured outbox tables.
		return "", nil, fmt.Errorf("no outbox table configured for topic '%s'", topic)
	}

	// Include transaction_id column with pg_current_xact_id() function
	query := fmt.Sprintf(
		`INSERT INTO %s ("uuid", "payload", "metadata", "transaction_id") VALUES %s`,
		tableName,
		buildInsertMarkers(len(msgs)),
	)

	var args []interface{}
	for _, msg := range msgs {
		// Store the topic in the metadata
		if msg.Metadata == nil {
			msg.Metadata = make(message.Metadata)
		}
		msg.Metadata[DestinationTopicKey] = topic

		metadata, err := json.Marshal(msg.Metadata)
		if err != nil {
			return "", nil, errors.Wrapf(err, "could not marshal metadata for message %s", msg.UUID)
		}
		args = append(args, msg.UUID, json.RawMessage(msg.Payload), json.RawMessage(metadata))
	}

	return query, args, nil
}
