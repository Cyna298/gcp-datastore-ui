package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// Initializes a Datastore client
func NewDatastoreClient(ctx context.Context) (*datastore.Client, error) {
	if os.Getenv("DATASTORE_DATASET") == "" {
		log.Fatal("DATASTORE_DATASET environment variable must be set to the project ID")
	}
	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		log.Fatal("DATASTORE_EMULATOR_HOST environment variable must be set to the emulator host")
	}
	if os.Getenv("DATASTORE_EMULATOR_HOST_PATH") == "" {
		log.Fatal("DATASTORE_EMULATOR_HOST_PATH environment variable must be set to the emulator host path")
	}
	if os.Getenv("DATASTORE_HOST") == "" {
		log.Fatal("DATASTORE_HOST environment variable must be set to the datastore host")
	}
	if os.Getenv("DATASTORE_PROJECT_ID") == "" {
		log.Fatal("DATASTORE_PROJECT_ID environment variable must be set to the project ID")
	}

	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetAllKinds retrieves all kinds from Datastore
func GetAllKinds(ctx context.Context, client *datastore.Client) ([]string, error) {
	query := datastore.NewQuery("__kind__").KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var kinds []string
	for _, key := range keys {
		kinds = append(kinds, key.Name)
	}
	return kinds, nil
}

type TableHeader struct {
	Name string
	Type string
}

type OutputProperty struct {
	Name    string      `json:"name"`
	Value   interface{} `json:"value"`
	TypeOf  string      `json:"type"`
	Indexed bool        `json:"indexed"`
}

type DisplayProperty struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	TypeOf  string `json:"type"`
	Indexed bool   `json:"indexed"`
}

type GeneralEntity map[string]OutputProperty
type DisplayEntity map[string]DisplayProperty

func (x *GeneralEntity) Load(ps []datastore.Property) error {
	*x = make(map[string]OutputProperty)

	for _, p := range ps {
		var value interface{} = p.Value
		if e, ok := p.Value.(*datastore.Entity); ok {
			// Convert nested entity to GeneralEntity

			fmt.Println("Start NESTED Entity------------------------------\n")
			nestedEntity := GeneralEntity{}
			nestedEntity.Load(e.Properties)
			value = nestedEntity

			fmt.Println("End NESTED Entity------------------------------\n")
		}
		(*x)[p.Name] = OutputProperty{
			Name:    p.Name,
			Value:   value,
			TypeOf:  fmt.Sprintf("%T", value),
			Indexed: !p.NoIndex,
		}

		fmt.Println("Start Entity------------------------------\n")
		fmt.Println(p.Name)
		fmt.Println(

			fmt.Sprintf("%T", value),
		)
		fmt.Println(value)

		fmt.Println("End---------------------------------\n")
	}

	return nil
}

func (x *GeneralEntity) Save() ([]datastore.Property, error) {

	return []datastore.Property{
		{
			Name:  "I",
			Value: x,
		},
	}, nil
}

func (ge GeneralEntity) GetValue(name string) (interface{}, error) {
	prop, ok := ge[name]
	if !ok {

		fmt.Println("Start Prop------------------------------\n")
		fmt.Println(ge)

		fmt.Println("End---------------------------------\n")
		return nil, fmt.Errorf("property %s not found", name)
	}

	fmt.Println("Start------------------------------\n")
	fmt.Println(prop.Name)
	fmt.Println(prop.TypeOf)
	fmt.Println(prop.Value)

	fmt.Println("End---------------------------------\n")

	switch prop.TypeOf {
	case "int64":
		return prop.Value.(int64), nil
	case "bool":
		return prop.Value.(bool), nil
	case "string":
		return prop.Value.(string), nil
	case "float64":
		return prop.Value.(float64), nil
	case "*datastore.Key":
		return prop.Value.(*datastore.Key), nil
	case "time.Time":
		return prop.Value.(time.Time), nil
	case "datastore.GeoPoint":
		return prop.Value.(datastore.GeoPoint), nil
	case "[]byte":
		return prop.Value.([]byte), nil

	case "[]interface {}":
		return prop.Value.([]interface{}), nil
	case "*datastore.Entity":
		return prop.Value.(*datastore.Entity), nil
	case "service.GeneralEntity":
		fmt.Println(prop.Value)
		return prop.Value.(GeneralEntity), nil
	case "<nil>":
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported type %s for property %s", prop.TypeOf, name)
	}
}

func (ge GeneralEntity) GetString(name string) (string, error) {
	value, err := ge.GetValue(name)

	fmt.Println("\nGET STRING Start------------------------------")
	fmt.Println(value)
	fmt.Println(err)
	fmt.Println(fmt.Sprintf("%T", value))

	fmt.Println("GET STRING End---------------------------------\n")
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case int64:
		return strconv.FormatInt(v, 10), nil
	case bool:
		return strconv.FormatBool(v), nil
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case *datastore.Key:
		return v.String(), nil // Assuming Key has a String method to give a meaningful representation
	case time.Time:
		return v.Format(time.RFC3339), nil
	case datastore.GeoPoint:
		return fmt.Sprintf("Lat: %f, Lng: %f", v.Lat, v.Lng), nil
	case []byte:
		return string(v), nil // May need to encode to base64 if binary data
	case *datastore.Entity:
		fmt.Println("Entity")

		fmt.Println(value.(*datastore.Entity))
		return "NOT IMPLEMENTED FOR ENTITY", nil
		// return fmt.Sprintf("Entity with Kind: %s", v.Key.Kind), nil // Simplistic representation
	case GeneralEntity:
		// nestedVal, err := v.GetValue()
		// fmt.Println("nested Val")
		//
		// fmt.Println("nested Val err")
		// fmt.Println(err)
		// if err != nil {
		// 	return "", err
		// }
		return fmt.Sprintf("%s", value), nil

	case []interface{}:
		var parts []string
		for _, item := range v {
			part, err := stringifyInterface(item)
			if err != nil {
				return "", fmt.Errorf("error converting array item to string: %s", err)
			}
			parts = append(parts, part)
		}
		return "[" + strings.Join(parts, ", ") + "]", nil
	case nil:
		return "NULL", nil
	default:
		return "", fmt.Errorf("unsupported type for conversion to string")
	}
}
func stringifyInterface(v interface{}) (string, error) {
	switch v := v.(type) {
	case int64:
		fmt.Println("INT")
		return strconv.FormatInt(v, 10), nil
	case bool:

		fmt.Println("BOOL")
		return strconv.FormatBool(v), nil
	case string:

		fmt.Println("STR")
		return v, nil
	case float64:

		fmt.Println("FL")
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case *datastore.Key:

		fmt.Println("KEy")
		return v.String(), nil
	case time.Time:

		fmt.Println("Time")
		return v.Format(time.RFC3339), nil
	case datastore.GeoPoint:
		fmt.Println("point")
		return fmt.Sprintf("Lat: %f, Lng: %f", v.Lat, v.Lng), nil
	case []byte:

		fmt.Println("Byte")
		return base64.StdEncoding.EncodeToString(v), nil
	case *datastore.Entity:

		fmt.Println("Entity")
		if v == nil {
			return "", nil
		}

		if v.Key == nil {
			return "", nil
		}
		return fmt.Sprintf("Entity with Kind: %s", v.Key.Kind), nil
	case nil:
		return "NULL", nil
	case []interface{}:

		fmt.Println("INTE")
		return "Yee", nil
	default:

		fmt.Println("DEF")
		return "", fmt.Errorf("unsupported type %s in array", reflect.TypeOf(v))
	}
}
func (ge GeneralEntity) ToDisplayEntity() (DisplayEntity, error) {
	de := make(DisplayEntity)
	for key, outputProp := range ge {
		stringValue, err := ge.GetString(key)
		if err != nil {
			return nil, fmt.Errorf("Display error converting key %s to string: %s", key, err)
		}
		de[key] = DisplayProperty{
			Name:    outputProp.Name,
			Value:   stringValue,
			TypeOf:  outputProp.TypeOf,
			Indexed: outputProp.Indexed,
		}
	}
	return de, nil
}

// GetAllEntities retrieves entities of a specific kind from Datastore
func GetAllEntities(ctx context.Context, client *datastore.Client, kind string, sortKey string, sortDirection string, limit int, cursorStr string) ([]DisplayEntity, string, error) {
	query := datastore.NewQuery(kind).Limit(limit)
	if sortKey != "" {
		if sortDirection == "desc" {
			query = query.Order("-" + sortKey)

		} else {
			query = query.Order(sortKey)
		}
	}
	if cursorStr != "" {
		cursor, err := datastore.DecodeCursor(cursorStr)
		if err != nil {
			return nil, "", err
		}
		query = query.Start(cursor)
	}

	var entities []DisplayEntity
	it := client.Run(ctx, query)

	for {
		var entity GeneralEntity = make(map[string]OutputProperty)
		key, err := it.Next(&entity)

		entity["key"] = OutputProperty{
			Name:    "key",
			Value:   key.String(),
			TypeOf:  "string",
			Indexed: true,
		}
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("err", err)
			return nil, "", err
		}

		ds, err := entity.ToDisplayEntity()

		fmt.Println("\nStart EEEE")
		fmt.Println(ds)
		fmt.Println("END EEEE\n")

		if err != nil {
			fmt.Println("err", err)
			return nil, "", err
		}
		entities = append(entities, ds)
	}

	nextCursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}
	if limit > len(entities) {
		return entities, "", nil
	}

	return entities, nextCursor.String(), nil
}
