package crm

import (
	"encoding/json"
	"testing"
)

// fieldsPayload is a realistic subset of a crm-backend
// GET /crm/v1/entities/{entity}/fields response (FieldService):
//   - the response is {"data": [...]}, fields in no particular order;
//   - values exist only on list/label fields, inactive options included;
//   - dependency is present only on a dependent (child) list field, and
//     absent (not null) otherwise;
//   - sort is an int or "".
const fieldsPayload = `{
	"data": [
		{
			"name": "Name",
			"code": "name",
			"entity_reference_id": "",
			"type": "text",
			"required": true,
			"editable": true,
			"show": true,
			"hidden": false,
			"multiple": false,
			"system_field": false,
			"base_field": true,
			"sort": "",
			"default_value": null,
			"tooltip": null
		},
		{
			"name": "Country",
			"code": "country",
			"entity_reference_id": "",
			"type": "list",
			"required": false,
			"editable": true,
			"show": true,
			"hidden": false,
			"multiple": true,
			"system_field": false,
			"base_field": false,
			"sort": 20,
			"default_value": null,
			"tooltip": null,
			"values": [
				{"title": "Ukraine", "value": "ua", "color": "", "sort": 1, "selected": false, "active": true},
				{"title": "Old Country", "value": "old", "color": "", "sort": "", "selected": false, "active": false},
				{"title": "No Active Key", "value": "na", "color": "", "sort": 2, "selected": false}
			]
		},
		{
			"name": "City",
			"code": "city",
			"entity_reference_id": "",
			"type": "list",
			"required": false,
			"editable": true,
			"show": true,
			"hidden": true,
			"multiple": false,
			"system_field": false,
			"base_field": false,
			"sort": 30,
			"default_value": null,
			"tooltip": null,
			"values": [
				{"title": "Kyiv", "value": "kyiv", "color": "", "sort": 1, "selected": true, "active": true}
			],
			"dependency": {
				"id": 4,
				"inverse_dependence": false,
				"show_all_options": false,
				"exclude_selected_options": false,
				"show_in_filters": false,
				"parent_field_code": "country"
			}
		}
	]
}`

// decodeFieldsPayload decodes fieldsPayload and returns its three fields in
// document order: a plain text field, a non-dependent list field, and a
// dependent (child) list field.
func decodeFieldsPayload(t *testing.T) []Field {
	t.Helper()
	var fields Fields
	if err := json.Unmarshal([]byte(fieldsPayload), &fields); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(fields.Data) != 3 {
		t.Fatalf("len(Data) = %d, want 3", len(fields.Data))
	}
	return fields.Data
}

// TestFieldDecode_PlainTextField pins the existing decoding behaviour for a
// field with no values and no dependency: both stay nil when the JSON omits
// the keys, and sort "" decodes to FlexInt(0).
func TestFieldDecode_PlainTextField(t *testing.T) {
	f := decodeFieldsPayload(t)[0]

	if f.Name != "Name" || f.Code != "name" || f.Type != "text" {
		t.Errorf("Name/Code/Type = %q/%q/%q, want %q/%q/%q", f.Name, f.Code, f.Type, "Name", "name", "text")
	}
	if !f.Required || !f.Editable || f.Hidden || f.Multiple {
		t.Errorf("Required/Editable/Hidden/Multiple = %v/%v/%v/%v, want true/true/false/false",
			f.Required, f.Editable, f.Hidden, f.Multiple)
	}
	if f.Sort.Int() != 0 {
		t.Errorf(`Sort = %d, want 0 (FlexInt("") == 0)`, f.Sort.Int())
	}
	if f.Values != nil {
		t.Errorf("Values = %+v, want nil", f.Values)
	}
	if f.Dependency != nil {
		t.Errorf("Dependency = %+v, want nil (key absent)", f.Dependency)
	}
}

// TestFieldDecode_ListFieldValues covers a non-dependent list field: its own
// Sort as an int, its Values' title/value, Sort as "" and as an int, and
// Active true, false and absent (nil).
func TestFieldDecode_ListFieldValues(t *testing.T) {
	f := decodeFieldsPayload(t)[1]

	if f.Code != "country" || f.Type != "list" || !f.Multiple {
		t.Errorf("Code/Type/Multiple = %q/%q/%v, want %q/%q/true", f.Code, f.Type, f.Multiple, "country", "list")
	}
	if f.Sort.Int() != 20 {
		t.Errorf("Sort = %d, want 20", f.Sort.Int())
	}
	if f.Dependency != nil {
		t.Errorf("Dependency = %+v, want nil (non-dependent field)", f.Dependency)
	}
	if len(f.Values) != 3 {
		t.Fatalf("len(Values) = %d, want 3", len(f.Values))
	}

	ua := f.Values[0]
	if ua.Title != "Ukraine" || ua.Value != "ua" {
		t.Errorf("Values[0].Title/Value = %q/%q, want %q/%q", ua.Title, ua.Value, "Ukraine", "ua")
	}
	if ua.Sort.Int() != 1 {
		t.Errorf("Values[0].Sort = %d, want 1", ua.Sort.Int())
	}
	if ua.Active == nil || !*ua.Active {
		t.Errorf("Values[0].Active = %v, want pointer to true", ua.Active)
	}

	old := f.Values[1]
	if old.Sort.Int() != 0 {
		t.Errorf(`Values[1].Sort = %d, want 0 (FlexInt("") == 0)`, old.Sort.Int())
	}
	if old.Active == nil || *old.Active {
		t.Errorf("Values[1].Active = %v, want pointer to false", old.Active)
	}

	na := f.Values[2]
	if na.Active != nil {
		t.Errorf("Values[2].Active = %v, want nil (key absent)", na.Active)
	}
}

// TestFieldDecode_Dependency covers a dependent (child) list field: Dependency
// decodes with its id and parent_field_code, and its single option is active.
func TestFieldDecode_Dependency(t *testing.T) {
	f := decodeFieldsPayload(t)[2]

	if f.Code != "city" || !f.Hidden {
		t.Errorf("Code/Hidden = %q/%v, want %q/true", f.Code, f.Hidden, "city")
	}
	if len(f.Values) != 1 || f.Values[0].Active == nil || !*f.Values[0].Active {
		t.Fatalf("Values = %+v, want one active option", f.Values)
	}

	dep := f.Dependency
	if dep == nil {
		t.Fatal("Dependency = nil, want non-nil")
	}
	if dep.ID != 4 {
		t.Errorf("Dependency.ID = %d, want 4", dep.ID)
	}
	if dep.ParentFieldCode != "country" {
		t.Errorf("Dependency.ParentFieldCode = %q, want %q", dep.ParentFieldCode, "country")
	}
	if dep.InverseDependence || dep.ShowAllOptions || dep.ExcludeSelectedOptions || dep.ShowInFilters {
		t.Errorf("Dependency bool fields = %+v, want all false", dep)
	}
}

// TestValueDecode_ActiveAbsent pins Value's zero-value decoding directly:
// with no "active" key at all, Active stays nil.
func TestValueDecode_ActiveAbsent(t *testing.T) {
	var v Value
	if err := json.Unmarshal([]byte(`{"title":"T","value":"v","color":"","sort":1,"selected":false}`), &v); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if v.Active != nil {
		t.Errorf("Active = %v, want nil (key absent)", v.Active)
	}
}
