// Copyright (C) 2015 NTT Innovation Institute, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

import "github.com/cloudwan/gohan/util"

//Property is a definition of each Property
type Property struct {
	ID, Title, Description string
	Type, Format           string
	Properties             []Property
	Items                  *Property
	Relation               string
	RelationColumn         string
	RelationProperty       string
	Unique                 bool
	Nullable               bool
	SQLType                string
	OnDeleteCascade        bool
	Default                interface{}
	DefaultMask            interface{}
	Indexed                bool
}

//PropertyMap is a map of Property
type PropertyMap map[string]Property

//NewProperty is a constructor for Property type
func NewProperty(id, title, description, typeID, format, relation, relationColumn, relationProperty, sqlType string, unique, nullable, onDeleteCascade bool, properties []Property, items *Property, defaultValue interface{}, indexed bool) Property {
	Property := Property{
		ID:               id,
		Title:            title,
		Format:           format,
		Description:      description,
		Type:             typeID,
		Relation:         relation,
		RelationColumn:   relationColumn,
		RelationProperty: relationProperty,
		Unique:           unique,
		Nullable:         nullable,
		Default:          defaultValue,
		Properties:       properties,
		Items:            items,
		SQLType:          sqlType,
		OnDeleteCascade:  onDeleteCascade,
		Indexed:          indexed,
	}
	return Property
}

//NewPropertyFromObj make Property  from obj
func NewPropertyFromObj(id string, rawTypeData interface{}, required bool) *Property {
	typeData, _ := rawTypeData.(map[string]interface{})
	title, _ := typeData["title"].(string)
	description, _ := typeData["description"].(string)
	var typeID string
	nullable := false
	switch typeData["type"].(type) {
	case string:
		typeID = typeData["type"].(string)
	case []interface{}:
		for _, typeInt := range typeData["type"].([]interface{}) {
			// type can be either string or list of string. we allow for any type and optional null
			// in order to retrieve right type, we need to skip null
			if typeInt.(string) != "null" {
				typeID = typeInt.(string)
			} else {
				nullable = true
			}
		}
	}
	format, _ := typeData["format"].(string)
	relation, _ := typeData["relation"].(string)
	relationColumn, _ := typeData["relation_column"].(string)
	relationProperty, _ := typeData["relation_property"].(string)
	unique, _ := typeData["unique"].(bool)
	cascade, _ := typeData["on_delete_cascade"].(bool)
	defaultValue, _ := typeData["default"]
	if !required && defaultValue == nil {
		nullable = true
	}
	sqlType, _ := typeData["sql"].(string)
	indexed, _ := typeData["indexed"].(bool)

	var items *Property
	if itemsRaw, hasItems := typeData["items"]; hasItems {
		items = NewPropertyFromObj("[]", itemsRaw, true)
	}

	properties := []Property{}
	if typeID == "object" {
		properties = parseSubproperties(typeData)
	}

	property := NewProperty(id, title, description, typeID, format, relation, relationColumn, relationProperty,
		sqlType, unique, nullable, cascade, properties, items, defaultValue, indexed)
	property.generateDefaultMask()
	return &property
}

func parseSubproperties(typeData map[string]interface{}) []Property {
	required, _ := typeData["required"].([]string)
	properties, _ := typeData["properties"].(map[string]interface{})

	parsedProperties := []Property{}

	for innerPropertyID, innerPropertyDesc := range properties {
		isRequired := util.ContainsString(required, innerPropertyID)
		parsedProperty := NewPropertyFromObj(innerPropertyID, innerPropertyDesc, isRequired)
		parsedProperties = append(parsedProperties, *parsedProperty)
	}

	return parsedProperties
}

func (p *Property) generateDefaultMask() {
	if p.Default != nil {
		p.DefaultMask = p.Default
		return
	}
	if p.Type != "object" {
		p.DefaultMask = nil
		return
	}

	var defaultMask map[string]interface{}
	for _, prop := range p.Properties {
		prop.generateDefaultMask()
		if prop.DefaultMask != nil {
			if defaultMask == nil {
				defaultMask = map[string]interface{}{}
			}
			defaultMask[prop.ID] = prop.DefaultMask
		}
	}

	p.DefaultMask = defaultMask
}
