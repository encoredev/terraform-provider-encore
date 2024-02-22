package provider

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"encr.dev/pkg/idents"
)

func ptr[T any](x T) *T {
	return &x
}

func getSchema(t TFType) (string, schema.SingleNestedAttribute) {
	key, docs, attDocs := t.GetDocs()
	rtn := schema.SingleNestedAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: docs,
		Attributes:          map[string]schema.Attribute{},
	}

	typ := reflect.TypeOf(t).Elem()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := field.Tag.Get("tfsdk")
		if name == "" {
			name = idents.Convert(field.Name, idents.SnakeCase)
		}
		fieldTyp := field.Type
		if fieldTyp.Kind() == reflect.Ptr {
			fieldTyp = fieldTyp.Elem()
		}
		switch fieldTyp.Kind() {
		case reflect.String:
			rtn.Attributes[name] = schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: attDocs[name],
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rtn.Attributes[name] = schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: attDocs[name],
			}
		case reflect.Bool:
			rtn.Attributes[name] = schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: attDocs[name],
			}
		case reflect.Float32, reflect.Float64:
			rtn.Attributes[name] = schema.Float64Attribute{
				Computed:            true,
				MarkdownDescription: attDocs[name],
			}
		default:
			panic(fmt.Sprintf("unsupported type %s", fieldTyp.Kind()))
		}
	}
	return key, rtn
}

func getAttributes(t TFType) map[string]attr.Value {
	rtn := make(map[string]attr.Value)
	val := reflect.ValueOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		name := field.Tag.Get("tfsdk")
		if name == "" {
			name = idents.Convert(field.Name, idents.SnakeCase)
		}
		val := val.Field(i)
		switch val.Kind() {
		case reflect.String:
			rtn[name] = types.StringValue(val.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rtn[name] = types.Int64Value(val.Int())
		case reflect.Bool:
			rtn[name] = types.BoolValue(val.Bool())
		case reflect.Float32, reflect.Float64:
			rtn[name] = types.Float64Value(val.Float())
		case reflect.Pointer:
			if val.Field(i).IsNil() {
				continue
			}
			ptrVal := val.Field(i).Elem()
			switch ptrVal.Kind() {
			case reflect.String:
				rtn[name] = types.StringPointerValue(ptr(ptrVal.String()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				rtn[name] = types.Int64PointerValue(ptr(ptrVal.Int()))
			case reflect.Bool:
				rtn[name] = types.BoolPointerValue(ptr(ptrVal.Bool()))
			case reflect.Float32, reflect.Float64:
				rtn[name] = types.Float64PointerValue(ptr(ptrVal.Float()))
			default:
				panic(fmt.Sprintf("unsupported pointer type %s", ptrVal.Kind()))
			}
		default:
			panic(fmt.Sprintf("unsupported type %s", val.Field(i).Kind()))
		}
	}
	return rtn
}

func getAttrTypes(in map[string]schema.Attribute) map[string]attr.Type {
	out := make(map[string]attr.Type, len(in))
	for k, v := range in {
		out[k] = v.GetType()
	}
	return out
}

func (s *Satisfier) GetData() (TFType, diag.Diagnostics) {
	field := reflect.ValueOf(s).Elem().FieldByName(s.Type)
	if !field.IsValid() {
		return nil, nil
	}
	switch field := field.Addr().Interface().(type) {
	case TFType:
		return field, nil
	default:
		var diags diag.Diagnostics
		diags.AddError("wrong satisfier type", fmt.Sprintf("expected %q, got %T", s.Type, field))
		return nil, diags
	}
}

type TFType interface {
	GetDocs() (subkey string, mdDesc string, attrDesc map[string]string)
}

type Need struct {
	ID         string
	TypeRef    TypeRef
	EncoreName string
	Satisfier  *Satisfier
}

func NewNeedsData(client *PlatformClient, envName string, ds []func() datasource.DataSource) *NeedsData {
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
	client     *PlatformClient
	defaultEnv string
	types      []TypeRef
}

func createSchema(tfTypes ...TFType) schema.Schema {
	ds := schema.Schema{
		MarkdownDescription: "Data source that provides information about an Encore-managed resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"env": schema.StringAttribute{
				Optional: true,
			},
		},
	}
	for _, tfType := range tfTypes {
		tfTypeName, tfTypeSchema := getSchema(tfType)
		ds.Attributes[tfTypeName] = tfTypeSchema
	}
	return ds
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
	n, diags := s.Get(ctx, typRef, envName.ValueString(), encoreName.ValueString())
	if diags.HasError() {
		return diags
	}
	if n == nil || n.Satisfier == nil {
		return nil
	}

	typ, diags := n.Satisfier.GetData()
	if diags.HasError() {
		return diags
	} else if typ == nil {
		return diags
	}
	key, sch := getSchema(typ)
	value, diags := types.ObjectValue(getAttrTypes(sch.Attributes), getAttributes(typ))
	if diags.HasError() {
		return diags
	}
	state.SetAttribute(ctx, path.Root("name"), encoreName)
	state.SetAttribute(ctx, path.Root("env"), envName)
	state.SetAttribute(ctx, path.Root(key), value)
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
	err := n.client.gql.Query(ctx, &q, map[string]interface{}{
		"appSlug": n.client.appSlug,
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
