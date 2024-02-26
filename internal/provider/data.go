package provider

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"encr.dev/pkg/idents"
)

var tfType = reflect.TypeOf((*TerraformDescription)(nil)).Elem()

func getAttribute(fieldTyp reflect.Type, desc string) (rtn schema.Attribute, diags diag.Diagnostics) {
	for fieldTyp.Kind() == reflect.Ptr {
		fieldTyp = fieldTyp.Elem()
	}
	switch fieldTyp.Kind() {
	case reflect.String:
		return schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: desc,
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: desc,
		}, nil
	case reflect.Bool:
		return schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: desc,
		}, nil
	case reflect.Float32, reflect.Float64:
		return schema.Float64Attribute{
			Computed:            true,
			MarkdownDescription: desc,
		}, nil
	case reflect.Slice:
		attribute, diags := getAttribute(fieldTyp.Elem(), desc)
		if diags.HasError() {
			return nil, diags
		}
		nestedAttrs, isNested := attribute.(schema.SingleNestedAttribute)
		if isNested {
			return schema.ListNestedAttribute{
				NestedObject:        nestedAttrs.GetNestedObject().(schema.NestedAttributeObject),
				Computed:            true,
				MarkdownDescription: desc,
			}, nil
		} else {
			return schema.ListAttribute{
				ElementType:         attribute.GetType(),
				Computed:            true,
				MarkdownDescription: desc,
			}, nil
		}
	case reflect.Struct:
		attributes, diags := getAttributes(fieldTyp)
		return schema.SingleNestedAttribute{
			Computed:            true,
			Attributes:          attributes,
			MarkdownDescription: desc,
		}, diags
	default:
		diags.AddError("Unsupported Type", fmt.Sprintf("unsupported type %s", fieldTyp))
	}
	return nil, diags
}

func containsFragment(field reflect.StructField, fragmentFilter ...string) bool {
	fragment := strings.TrimPrefix(field.Tag.Get("graphql"), "... on ")
	if len(fragmentFilter) == 0 || slices.Contains(fragmentFilter, fragment) {
		return true
	}
	return false
}

func getAttributes(typ reflect.Type, fragmentFilter ...string) (rtn map[string]schema.Attribute, diags diag.Diagnostics) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		diags.AddError("Unsupported Type", fmt.Sprintf("unsupported type %s", typ))
		return nil, diags
	}
	attDocs := map[string]string{}
	if reflect.PointerTo(typ).Implements(tfType) {
		attDocs = reflect.New(typ).Interface().(TerraformDescription).GetDocs()
	}

	rtn = make(map[string]schema.Attribute)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !containsFragment(field, fragmentFilter...) {
			continue
		}
		name := getTFName(field)
		att, diags := getAttribute(field.Type, attDocs[name])
		if diags.HasError() {
			return nil, diags
		}
		if sn, ok := att.(schema.SingleNestedAttribute); ok && name == "" {
			maps.Copy(rtn, sn.Attributes)
			continue
		}
		rtn[name] = att
	}
	return rtn, nil
}

func getValue(val reflect.Value) (rtn attr.Value, diags diag.Diagnostics) {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.String:
		return types.StringValue(val.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return types.Int64Value(val.Int()), nil
	case reflect.Bool:
		return types.BoolValue(val.Bool()), nil
	case reflect.Float32, reflect.Float64:
		return types.Float64Value(val.Float()), nil
	case reflect.Slice:
		elements := make([]attr.Value, val.Len())
		for i := 0; i < val.Len(); i++ {
			elem, diags := getValue(val.Index(i))
			if diags.HasError() {
				return nil, diags
			}
			elements[i] = elem
		}
		att, diags := getAttribute(val.Type().Elem(), "")
		if diags.HasError() {
			return nil, diags
		}
		return types.ListValue(att.GetType(), elements)
	case reflect.Struct:
		subVals, diags := getValues(val)
		if diags.HasError() {
			return nil, diags
		}
		attributes, diags := getAttributes(val.Type())
		if diags.HasError() {
			return nil, diags
		}
		return types.ObjectValue(
			getAttrTypes(attributes),
			subVals,
		)
	default:
		diags.AddError("Unsupported Type", fmt.Sprintf("unsupported type %s", val.Kind()))
		return nil, diags
	}
}

func getValues(val reflect.Value, fragmentFilter ...string) (rtn map[string]attr.Value, diags diag.Diagnostics) {
	rtn = make(map[string]attr.Value)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		diags.AddError("Unsupported Type", fmt.Sprintf("unsupported type %s", val.Kind()))
		return nil, diags
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if !containsFragment(field, fragmentFilter...) {
			continue
		}
		name := getTFName(field)
		val := val.Field(i)
		attr, diags := getValue(val)
		if diags.HasError() {
			return nil, diags
		} else if attr == nil {
			continue
		}
		if obj, ok := attr.(basetypes.ObjectValue); ok && name == "" {
			maps.Copy(rtn, obj.Attributes())
		} else {
			rtn[name] = attr
		}
	}
	return rtn, nil
}

func getAttrTypes(in map[string]schema.Attribute) map[string]attr.Type {
	out := make(map[string]attr.Type, len(in))
	for k, v := range in {
		out[k] = v.GetType()
	}
	return out
}

func getTFName(field reflect.StructField) string {
	name := field.Tag.Get("tf")
	if name == "" && !field.Anonymous {
		name = idents.Convert(field.Name, idents.SnakeCase)
	}
	return name
}

type TerraformDescription interface {
	GetDocs() (attrDesc map[string]string)
}

type Need struct {
	ID         string
	TypeRef    TypeRef
	EncoreName string
	Satisfier  *SatisfierQuery
}

func NewNeedsData(client PlatformClient, envName string, ds []func() datasource.DataSource) *NeedsData {
	if envName == "" {
		envName = "@primary"
	}
	n := &NeedsData{
		client:     client,
		needs:      map[string]map[TypeRef]map[string]*Need{},
		defaultEnv: envName,
	}
	for _, d := range ds {
		ds := d()
		if ds, ok := ds.(*EncoreDataSource); ok {
			n.types = append(n.types, ds.typeRef)
		}

	}
	return n
}

type NeedsData struct {
	needs      map[string]map[TypeRef]map[string]*Need
	client     PlatformClient
	defaultEnv string
	types      []TypeRef
}

func createSchema(desc string, fragments ...string) schema.Schema {
	attrs, diags := getAttributes(queryType, fragments...)
	if diags.HasError() {
		panic(diags)
	}
	attrs["name"] = schema.StringAttribute{
		MarkdownDescription: "The name of the Encore resource",
		Required:            true,
	}
	attrs["env"] = schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The environment of the Encore resource. Defaults to the provider environment",
	}
	return schema.Schema{
		MarkdownDescription: desc,
		Attributes:          attrs,
	}
}

func (s *NeedsData) SetValue(ctx context.Context, typRef TypeRef, reqCfg tfsdk.Config, state *tfsdk.State) diag.Diagnostics {
	var encoreName, envName types.String
	var diags diag.Diagnostics

	diags.Append(reqCfg.GetAttribute(ctx, path.Root("name"), &encoreName)...)
	diags.Append(reqCfg.GetAttribute(ctx, path.Root("env"), &envName)...)
	if diags.HasError() {
		return diags
	}
	if envName.ValueString() == "" {
		envName = types.StringValue(s.defaultEnv)
	}
	diags.Append(state.SetAttribute(ctx, path.Root("name"), encoreName)...)
	diags.Append(state.SetAttribute(ctx, path.Root("env"), envName)...)

	if diags.HasError() {
		return diags
	}

	n, diags := s.Get(ctx, typRef, envName.ValueString(), encoreName.ValueString())
	if diags.HasError() {
		return diags
	}
	if n == nil || n.Satisfier == nil {
		return nil
	}

	values, diags := getValues(reflect.ValueOf(n.Satisfier), n.Satisfier.Type)
	if diags.HasError() {
		return diags
	}
	for key, val := range values {
		diags.Append(state.SetAttribute(ctx, path.Root(key), val)...)
	}
	return diags
}

func (n *NeedsData) Get(ctx context.Context, typRef TypeRef, envName, encoreName string) (*Need, diag.Diagnostics) {
	envNeeds, diags := n.envNeeds(ctx, envName)
	if diags.HasError() {
		return nil, diags
	}
	if envNeeds[typRef] == nil {
		return nil, nil
	}
	return envNeeds[typRef][encoreName], nil
}

type TypeRef string

func (n *NeedsData) envNeeds(ctx context.Context, envName string) (map[TypeRef]map[string]*Need, diag.Diagnostics) {
	if envNeeds, ok := n.needs[envName]; ok {
		return envNeeds, nil
	}
	var q struct {
		App struct {
			Env struct {
				Needs []*Need `graphql:"needs(sel:{typeRefs:$types})"`
			} `graphql:"env(name: $envName)"`
		} `graphql:"app(slug: $appSlug)"`
	}
	err := n.client.GQL().Query(ctx, &q, map[string]interface{}{
		"appSlug": n.client.AppSlug(),
		"envName": envName,
		"types":   n.types,
	})
	if err != nil {
		var diags diag.Diagnostics
		if strings.Contains(err.Error(), "env not found") {
			diags.AddAttributeError(path.Root("env"), "Env not found", "The specified environment does not exist")
		} else {
			diags.AddError("Client Error", fmt.Sprintf("Unable to fetch Encore resources, got error: %s", err))
		}
		return nil, diags
	}
	envTypes := make(map[TypeRef]map[string]*Need)
	for _, need := range q.App.Env.Needs {
		if envTypes[need.TypeRef] == nil {
			envTypes[need.TypeRef] = make(map[string]*Need)
		}
		envTypes[need.TypeRef][need.EncoreName] = need
	}
	n.needs[envName] = envTypes
	return envTypes, nil
}
