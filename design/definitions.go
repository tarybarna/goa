package design

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/dimfeld/httppath"
	"github.com/goadesign/goa/dslengine"
)

type (
	// APIDefinition defines the global properties of the API.
	APIDefinition struct {
		// Name of API
		Name string
		// Title of API
		Title string
		// Description of API
		Description string
		// Version is the version of the API described by this design.
		Version string
		// Host is the default API hostname
		Host string
		// Schemes is the supported API URL schemes
		Schemes []string
		// BasePath is the common base path to all API endpoints
		BasePath string
		// BaseParams define the common path parameters to all API endpoints
		BaseParams *AttributeDefinition
		// Consumes lists the mime types supported by the API controllers
		Consumes []*EncodingDefinition
		// Produces lists the mime types generated by the API controllers
		Produces []*EncodingDefinition
		// TermsOfService describes or links to the API terms of service
		TermsOfService string
		// Contact provides the API users with contact information
		Contact *ContactDefinition
		// License describes the API license
		License *LicenseDefinition
		// Docs points to the API external documentation
		Docs *DocsDefinition
		// Resources is the set of exposed resources indexed by name
		Resources map[string]*ResourceDefinition
		// Types indexes the user defined types by name
		Types map[string]*UserTypeDefinition
		// MediaTypes indexes the API media types by canonical identifier
		MediaTypes map[string]*MediaTypeDefinition
		// Traits available to all API resources and actions indexed by name
		Traits map[string]*dslengine.TraitDefinition
		// Responses available to all API actions indexed by name
		Responses map[string]*ResponseDefinition
		// Response template factories available to all API actions indexed by name
		ResponseTemplates map[string]*ResponseTemplateDefinition
		// Built-in responses
		DefaultResponses map[string]*ResponseDefinition
		// Built-in response templates
		DefaultResponseTemplates map[string]*ResponseTemplateDefinition
		// DSLFunc contains the DSL used to create this definition if any
		DSLFunc func()
		// Metadata is a list of key/value pairs
		Metadata dslengine.MetadataDefinition
		// SecuritySchemes lists the available security schemes available
		// to the API.
		SecuritySchemes []*SecuritySchemeDefinition
		// Security defines security requirements for all the
		// resources and actions, unless overridden by Resource or
		// Action-level Security() calls.
		Security *SecurityDefinition

		// rand is the random generator used to generate examples.
		rand *RandomGenerator
	}

	// ContactDefinition contains the API contact information.
	ContactDefinition struct {
		// Name of the contact person/organization
		Name string `json:"name,omitempty"`
		// Email address of the contact person/organization
		Email string `json:"email,omitempty"`
		// URL pointing to the contact information
		URL string `json:"url,omitempty"`
	}

	// LicenseDefinition contains the license information for the API.
	LicenseDefinition struct {
		// Name of license used for the API
		Name string `json:"name,omitempty"`
		// URL to the license used for the API
		URL string `json:"url,omitempty"`
	}

	// DocsDefinition points to external documentation.
	DocsDefinition struct {
		// Description of documentation.
		Description string `json:"description,omitempty"`
		// URL to documentation.
		URL string `json:"url,omitempty"`
	}

	// ResourceDefinition describes a REST resource.
	// It defines both a media type and a set of actions that can be executed through HTTP
	// requests.
	ResourceDefinition struct {
		// Resource name
		Name string
		// Schemes is the supported API URL schemes
		Schemes []string
		// Common URL prefix to all resource action HTTP requests
		BasePath string
		// Object describing each parameter that appears in BasePath if any
		BaseParams *AttributeDefinition
		// Name of parent resource if any
		ParentName string
		// Optional description
		Description string
		// Default media type, describes the resource attributes
		MediaType string
		// Exposed resource actions indexed by name
		Actions map[string]*ActionDefinition
		// Action with canonical resource path
		CanonicalActionName string
		// Map of response definitions that apply to all actions indexed by name.
		Responses map[string]*ResponseDefinition
		// Path and query string parameters that apply to all actions.
		Params *AttributeDefinition
		// Request headers that apply to all actions.
		Headers *AttributeDefinition
		// DSLFunc contains the DSL used to create this definition if any.
		DSLFunc func()
		// metadata is a list of key/value pairs
		Metadata dslengine.MetadataDefinition
		// Security defines security requirements for the Resource,
		// for actions that don't define one themselves.
		Security *SecurityDefinition
	}

	// EncodingDefinition defines an encoder supported by the API.
	EncodingDefinition struct {
		// MIMETypes is the set of possible MIME types for the content being encoded or decoded.
		MIMETypes []string
		// PackagePath is the path to the Go package that implements the encoder/decoder.
		// The package must expose a `EncoderFactory` or `DecoderFactory` function
		// that the generated code calls. The methods must return objects that implement
		// the goa.EncoderFactory or goa.DecoderFactory interface respectively.
		PackagePath string
		// Function is the name of the Go function used to instantiate the encoder/decoder.
		// Defaults to NewEncoder and NewDecoder respecitively.
		Function string
		// Encoder is true if the definition is for a encoder, false if it's for a decoder.
		Encoder bool
	}

	// ResponseDefinition defines a HTTP response status and optional validation rules.
	ResponseDefinition struct {
		// Response name
		Name string
		// HTTP status
		Status int
		// Response description
		Description string
		// Response body type if any
		Type DataType
		// Response body media type if any
		MediaType string
		// Response header definitions
		Headers *AttributeDefinition
		// Parent action or resource
		Parent dslengine.Definition
		// Metadata is a list of key/value pairs
		Metadata dslengine.MetadataDefinition
		// Standard is true if the response definition comes from the goa default responses
		Standard bool
	}

	// ResponseTemplateDefinition defines a response template.
	// A response template is a function that takes an arbitrary number
	// of strings and returns a response definition.
	ResponseTemplateDefinition struct {
		// Response template name
		Name string
		// Response template function
		Template func(params ...string) *ResponseDefinition
	}

	// ActionDefinition defines a resource action.
	// It defines both an HTTP endpoint and the shape of HTTP requests and responses made to
	// that endpoint.
	// The shape of requests is defined via "parameters", there are path parameters
	// parameters and a payload parameter (request body).
	// (i.e. portions of the URL that define parameter values), query string
	ActionDefinition struct {
		// Action name, e.g. "create"
		Name string
		// Action description, e.g. "Creates a task"
		Description string
		// Docs points to the API external documentation
		Docs *DocsDefinition
		// Parent resource
		Parent *ResourceDefinition
		// Specific action URL schemes
		Schemes []string
		// Action routes
		Routes []*RouteDefinition
		// Map of possible response definitions indexed by name
		Responses map[string]*ResponseDefinition
		// Path and query string parameters
		Params *AttributeDefinition
		// Query string parameters only
		QueryParams *AttributeDefinition
		// Payload blueprint (request body) if any
		Payload *UserTypeDefinition
		// Request headers that need to be made available to action
		Headers *AttributeDefinition
		// Metadata is a list of key/value pairs
		Metadata dslengine.MetadataDefinition
		// Security defines security requirements for the action
		Security *SecurityDefinition
	}

	// LinkDefinition defines a media type link, it specifies a URL to a related resource.
	LinkDefinition struct {
		// Link name
		Name string
		// View used to render link if not "link"
		View string
		// URITemplate is the RFC6570 URI template of the link Href.
		URITemplate string

		// Parent media Type
		Parent *MediaTypeDefinition
	}

	// ViewDefinition defines which members and links to render when building a response.
	// The view is a JSON object whose property names must match the names of the parent media
	// type members.
	// The members fields are inherited from the parent media type but may be overridden.
	ViewDefinition struct {
		// Set of properties included in view
		*AttributeDefinition
		// Name of view
		Name string
		// Parent media Type
		Parent *MediaTypeDefinition
	}

	// RouteDefinition represents an action route.
	RouteDefinition struct {
		// Verb is the HTTP method, e.g. "GET", "POST", etc.
		Verb string
		// Path is the URL path e.g. "/tasks/:id"
		Path string
		// Parent is the action this route applies to.
		Parent *ActionDefinition
	}

	// AttributeDefinition defines a JSON object member with optional description, default
	// value and validations.
	AttributeDefinition struct {
		// Attribute type
		Type DataType
		// Attribute reference type if any
		Reference DataType
		// Optional description
		Description string
		// Optional validations
		Validation *dslengine.ValidationDefinition
		// Metadata is a list of key/value pairs
		Metadata dslengine.MetadataDefinition
		// Optional member default value
		DefaultValue interface{}
		// Optional member example value
		Example interface{}
		// Optional view used to render Attribute (only applies to media type attributes).
		View string
		// NonZeroAttributes lists the names of the child attributes that cannot have a
		// zero value (and thus whose presence does not need to be validated).
		NonZeroAttributes map[string]bool
		// DSLFunc contains the initialization DSL. This is used for user types.
		DSLFunc func()
		// isCustomExample keeps track of whether the example is given by the user, or
		// should be automatically generated for the user.
		isCustomExample bool
	}

	// ContainerDefinition defines a generic container definition that contains attributes.
	// This makes it possible for plugins to use attributes in their own data structures.
	ContainerDefinition interface {
		// Attribute returns the container definition embedded attribute.
		Attribute() *AttributeDefinition
	}

	// ResourceIterator is the type of functions given to IterateResources.
	ResourceIterator func(r *ResourceDefinition) error

	// MediaTypeIterator is the type of functions given to IterateMediaTypes.
	MediaTypeIterator func(m *MediaTypeDefinition) error

	// UserTypeIterator is the type of functions given to IterateUserTypes.
	UserTypeIterator func(m *UserTypeDefinition) error

	// ActionIterator is the type of functions given to IterateActions.
	ActionIterator func(a *ActionDefinition) error

	// ResponseIterator is the type of functions given to IterateResponses.
	ResponseIterator func(r *ResponseDefinition) error
)

// NewAPIDefinition returns a new design with built-in response templates.
func NewAPIDefinition() *APIDefinition {
	api := &APIDefinition{
		DefaultResponseTemplates: make(map[string]*ResponseTemplateDefinition),
		DefaultResponses:         make(map[string]*ResponseDefinition),
	}
	t := func(params ...string) *ResponseDefinition {
		if len(params) < 1 {
			dslengine.ReportError("expected media type as argument when invoking response template OK")
			return nil
		}
		return &ResponseDefinition{
			Name:      OK,
			Status:    200,
			MediaType: params[0],
		}
	}
	api.DefaultResponseTemplates[OK] = &ResponseTemplateDefinition{
		Name:     OK,
		Template: t,
	}
	for _, p := range []struct {
		status int
		name   string
	}{
		{100, Continue},
		{101, SwitchingProtocols},
		{200, OK},
		{201, Created},
		{202, Accepted},
		{203, NonAuthoritativeInfo},
		{204, NoContent},
		{205, ResetContent},
		{206, PartialContent},
		{300, MultipleChoices},
		{301, MovedPermanently},
		{302, Found},
		{303, SeeOther},
		{304, NotModified},
		{305, UseProxy},
		{307, TemporaryRedirect},
		{400, BadRequest},
		{401, Unauthorized},
		{402, PaymentRequired},
		{403, Forbidden},
		{404, NotFound},
		{405, MethodNotAllowed},
		{406, NotAcceptable},
		{407, ProxyAuthRequired},
		{408, RequestTimeout},
		{409, Conflict},
		{410, Gone},
		{411, LengthRequired},
		{412, PreconditionFailed},
		{413, RequestEntityTooLarge},
		{414, RequestURITooLong},
		{415, UnsupportedMediaType},
		{416, RequestedRangeNotSatisfiable},
		{417, ExpectationFailed},
		{418, Teapot},
		{500, InternalServerError},
		{501, NotImplemented},
		{502, BadGateway},
		{503, ServiceUnavailable},
		{504, GatewayTimeout},
		{505, HTTPVersionNotSupported},
	} {
		api.DefaultResponses[p.name] = &ResponseDefinition{
			Name:        p.name,
			Description: http.StatusText(p.status),
			Status:      p.status,
		}
	}
	return api
}

// DSLName is the name of the DSL as displayed to the user during execution.
func (a *APIDefinition) DSLName() string {
	return "goa API"
}

// DependsOn returns the other roots this root depends on, nothing for APIDefinition.
func (a *APIDefinition) DependsOn() []dslengine.Root {
	return nil
}

// IterateSets calls the given iterator possing in the API definition, user types, media types and
// finally resources.
func (a *APIDefinition) IterateSets(iterator dslengine.SetIterator) {
	// First run the top level API DSL to initialize responses and
	// response templates needed by resources.
	iterator([]dslengine.Definition{a})

	// Then run the user type DSLs
	typeAttributes := make([]dslengine.Definition, len(a.Types))
	i := 0
	a.IterateUserTypes(func(u *UserTypeDefinition) error {
		u.AttributeDefinition.DSLFunc = u.DSLFunc
		typeAttributes[i] = u.AttributeDefinition
		i++
		return nil
	})
	iterator(typeAttributes)

	// Then the media type DSLs
	mediaTypes := make([]dslengine.Definition, len(a.MediaTypes))
	i = 0
	a.IterateMediaTypes(func(mt *MediaTypeDefinition) error {
		mediaTypes[i] = mt
		i++
		return nil
	})
	iterator(mediaTypes)

	// Then, the Security schemes definitions
	var securitySchemes []dslengine.Definition
	for _, scheme := range a.SecuritySchemes {
		securitySchemes = append(securitySchemes, dslengine.Definition(scheme))
	}
	iterator(securitySchemes)

	// And now that we have everything the resources.  The resource
	// lifecycle handlers dispatch to their children elements, like
	// Actions, etc..
	resources := make([]dslengine.Definition, len(a.Resources))
	i = 0
	a.IterateResources(func(res *ResourceDefinition) error {
		resources[i] = res
		i++
		return nil
	})
	iterator(resources)
}

// Reset sets all the API definition fields to their zero value except the default responses and
// default response templates.
func (a *APIDefinition) Reset() {
	n := NewAPIDefinition()
	*a = *n
}

// Context returns the generic definition name used in error messages.
func (a *APIDefinition) Context() string {
	if a.Name != "" {
		return fmt.Sprintf("API %#v", a.Name)
	}
	return "unnamed API"
}

// IterateMediaTypes calls the given iterator passing in each media type sorted in alphabetical order.
// Iteration stops if an iterator returns an error and in this case IterateMediaTypes returns that
// error.
func (a *APIDefinition) IterateMediaTypes(it MediaTypeIterator) error {
	names := make([]string, len(a.MediaTypes))
	i := 0
	for n := range a.MediaTypes {
		names[i] = n
		i++
	}
	sort.Strings(names)
	for _, n := range names {
		if err := it(a.MediaTypes[n]); err != nil {
			return err
		}
	}
	return nil
}

// IterateUserTypes calls the given iterator passing in each user type sorted in alphabetical order.
// Iteration stops if an iterator returns an error and in this case IterateUserTypes returns that
// error.
func (a *APIDefinition) IterateUserTypes(it UserTypeIterator) error {
	names := make([]string, len(a.Types))
	i := 0
	for n := range a.Types {
		names[i] = n
		i++
	}
	sort.Strings(names)
	for _, n := range names {
		if err := it(a.Types[n]); err != nil {
			return err
		}
	}
	return nil
}

// IterateResponses calls the given iterator passing in each response sorted in alphabetical order.
// Iteration stops if an iterator returns an error and in this case IterateResponses returns that
// error.
func (a *APIDefinition) IterateResponses(it ResponseIterator) error {
	names := make([]string, len(a.Responses))
	i := 0
	for n := range a.Responses {
		names[i] = n
		i++
	}
	sort.Strings(names)
	for _, n := range names {
		if err := it(a.Responses[n]); err != nil {
			return err
		}
	}
	return nil
}

// GenerateExample returns a random value for the given data type.
// If the data type has validations then the example value validates them.
// GenerateExample returns the same random value for a given api name
// (the random generator is seeded after the api name).
func (a *APIDefinition) GenerateExample(dt DataType) interface{} {
	return dt.GenerateExample(a.RandomGenerator())
}

// RandomGenerator is seeded after the API name. It's used to generate examples.
func (a *APIDefinition) RandomGenerator() *RandomGenerator {
	if a.rand == nil {
		a.rand = NewRandomGenerator(a.Name)
	}
	return a.rand
}

// MediaTypeWithIdentifier returns the media type with a matching
// media type identifier. Two media type identifiers match if their
// values sans suffix match. So for example "application/vnd.foo+xml",
// "application/vnd.foo+json" and "application/vnd.foo" all match.
func (a *APIDefinition) MediaTypeWithIdentifier(id string) *MediaTypeDefinition {
	canonicalID := CanonicalIdentifier(id)
	for _, mt := range a.MediaTypes {
		if canonicalID == CanonicalIdentifier(mt.Identifier) {
			return mt
		}
	}
	return nil
}

// IterateResources calls the given iterator passing in each resource sorted in alphabetical order.
// Iteration stops if an iterator returns an error and in this case IterateResources returns that
// error.
func (a *APIDefinition) IterateResources(it ResourceIterator) error {
	names := make([]string, len(a.Resources))
	i := 0
	for n := range a.Resources {
		names[i] = n
		i++
	}
	sort.Strings(names)
	for _, n := range names {
		if err := it(a.Resources[n]); err != nil {
			return err
		}
	}
	return nil
}

// DSL returns the initialization DSL.
func (a *APIDefinition) DSL() func() {
	return a.DSLFunc
}

// Finalize sets the Consumes and Produces fields to the defaults if empty.
func (a *APIDefinition) Finalize() {
	if len(a.Consumes) == 0 {
		a.Consumes = DefaultDecoders
	}
	if len(a.Produces) == 0 {
		a.Produces = DefaultEncoders
	}
}

// NewResourceDefinition creates a resource definition but does not
// execute the DSL.
func NewResourceDefinition(name string, dsl func()) *ResourceDefinition {
	return &ResourceDefinition{
		Name:      name,
		MediaType: "plain/text",
		DSLFunc:   dsl,
	}
}

// Context returns the generic definition name used in error messages.
func (r *ResourceDefinition) Context() string {
	if r.Name != "" {
		return fmt.Sprintf("resource %#v", r.Name)
	}
	return "unnamed resource"
}

// IterateActions calls the given iterator passing in each resource action sorted in alphabetical order.
// Iteration stops if an iterator returns an error and in this case IterateActions returns that
// error.
func (r *ResourceDefinition) IterateActions(it ActionIterator) error {
	names := make([]string, len(r.Actions))
	i := 0
	for n := range r.Actions {
		names[i] = n
		i++
	}
	sort.Strings(names)
	for _, n := range names {
		if err := it(r.Actions[n]); err != nil {
			return err
		}
	}
	return nil
}

// CanonicalAction returns the canonical action of the resource if any.
// The canonical action is used to compute hrefs to resources.
func (r *ResourceDefinition) CanonicalAction() *ActionDefinition {
	name := r.CanonicalActionName
	if name == "" {
		name = "show"
	}
	ca, _ := r.Actions[name]
	return ca
}

// URITemplate returns a URI template to this resource.
// The result is the empty string if the resource does not have a "show" action
// and does not define a different canonical action.
func (r *ResourceDefinition) URITemplate() string {
	ca := r.CanonicalAction()
	if ca == nil || len(ca.Routes) == 0 {
		return ""
	}
	return ca.Routes[0].FullPath()
}

// FullPath computes the base path to the resource actions concatenating the API and parent resource
// base paths as needed.
func (r *ResourceDefinition) FullPath() string {
	var basePath string
	if p := r.Parent(); p != nil {
		if ca := p.CanonicalAction(); ca != nil {
			if routes := ca.Routes; len(routes) > 0 {
				// Note: all these tests should be true at code generation time
				// as DSL validation makes sure that parent resources have a
				// canonical path.
				basePath = path.Join(routes[0].FullPath())
			}
		}
	} else {
		basePath = Design.BasePath
	}
	return httppath.Clean(path.Join(basePath, r.BasePath))
}

// Parent returns the parent resource if any, nil otherwise.
func (r *ResourceDefinition) Parent() *ResourceDefinition {
	if r.ParentName != "" {
		if parent, ok := Design.Resources[r.ParentName]; ok {
			return parent
		}
	}
	return nil
}

// DSL returns the initialization DSL.
func (r *ResourceDefinition) DSL() func() {
	return r.DSLFunc
}

// Finalize is run post DSL execution. It merges response definitions, creates implicit action
// parameters, initializes querystring parameters, sets path parameters as non zero attributes
// and sets the fallbacks for security schemes.
func (r *ResourceDefinition) Finalize() {
	r.IterateActions(func(a *ActionDefinition) error {
		a.Finalize()

		// 1. Merge response definitions
		for name, resp := range a.Responses {
			resp.Finalize()
			if pr, ok := a.Parent.Responses[name]; ok {
				resp.Merge(pr)
			}
			if ar, ok := Design.Responses[name]; ok {
				resp.Merge(ar)
			}
			if dr, ok := Design.DefaultResponses[name]; ok {
				resp.Merge(dr)
			}
		}
		// 2. Create implicit action parameters for path wildcards that dont' have one
		for _, r := range a.Routes {
			wcs := ExtractWildcards(r.FullPath())
			for _, wc := range wcs {
				found := false
				var o Object
				if all := a.Params; all != nil {
					o = all.Type.ToObject()
				} else {
					o = Object{}
					a.Params = &AttributeDefinition{Type: o}
				}
				for n := range o {
					if n == wc {
						found = true
						break
					}
				}
				if !found {
					o[wc] = &AttributeDefinition{Type: String}
				}
			}
		}
		// 3. Compute QueryParams from Params and set all path params as non zero attributes
		if params := a.Params; params != nil {
			queryParams := DupAtt(params)
			a.Params.NonZeroAttributes = make(map[string]bool)
			for _, route := range a.Routes {
				pnames := route.Params()
				for _, pname := range pnames {
					a.Params.NonZeroAttributes[pname] = true
					delete(queryParams.Type.ToObject(), pname)
				}
			}
			// (note: we may end up with required attribute names that don't correspond
			// to actual attributes cos' we just deleted them but that's probably OK.)
			a.QueryParams = queryParams
		}

		return nil
	})
}

// Context returns the generic definition name used in error messages.
func (enc *EncodingDefinition) Context() string {
	return fmt.Sprintf("encoding for %s", strings.Join(enc.MIMETypes, ", "))
}

// Context returns the generic definition name used in error messages.
func (a *AttributeDefinition) Context() string {
	return ""
}

// AllRequired returns the list of all required fields from the underlying object.
// An attribute type can be itself an attribute (e.g. a MediaTypeDefinition or a UserTypeDefinition)
// This happens when the DSL uses references for example. So traverse the hierarchy and collect
// all the required validations.
func (a *AttributeDefinition) AllRequired() (required []string) {
	if a.Validation == nil {
		return
	}
	required = a.Validation.Required
	if ds, ok := a.Type.(DataStructure); ok {
		required = append(required, ds.Definition().AllRequired()...)
	}
	return
}

// IsRequired returns true if the given string matches the name of a required
// attribute, false otherwise.
func (a *AttributeDefinition) IsRequired(attName string) bool {
	for _, name := range a.AllRequired() {
		if name == attName {
			return true
		}
	}
	return false
}

// AllNonZero returns the complete list of all non-zero attribute name.
func (a *AttributeDefinition) AllNonZero() []string {
	nzs := make([]string, len(a.NonZeroAttributes))
	i := 0
	for n := range a.NonZeroAttributes {
		nzs[i] = n
		i++
	}
	return nzs
}

// IsNonZero returns true if the given string matches the name of a non-zero
// attribute, false otherwise.
func (a *AttributeDefinition) IsNonZero(attName string) bool {
	return a.NonZeroAttributes[attName]
}

// IsPrimitivePointer returns true if the field generated for the given attribute should be a
// pointer to a primitive type. The target attribute must be an object.
func (a *AttributeDefinition) IsPrimitivePointer(attName string) bool {
	if !a.Type.IsObject() {
		panic("checking pointer field on non-object") // bug
	}
	att := a.Type.ToObject()[attName]
	if att == nil {
		return false
	}
	if att.Type.IsPrimitive() {
		return !a.IsRequired(attName) && !a.IsNonZero(attName)
	}
	return false
}

// GenerateExample returns a random instance of the attribute that validates.
func (a *AttributeDefinition) GenerateExample(r *RandomGenerator) interface{} {
	if example := newExampleGenerator(a, r).generate(); example != nil {
		return example
	}
	return a.Type.GenerateExample(r)
}

// SetExample sets the custom example. SetExample also handles the case when the user doesn't
// want any example or any auto-generated example.
func (a *AttributeDefinition) SetExample(example interface{}) bool {
	if example == nil {
		a.Example = nil
		a.isCustomExample = true
		return true
	}
	if a.Type == nil || a.Type.IsCompatible(example) {
		a.Example = example
		a.isCustomExample = true
		return true
	}
	return false
}

// finalizeExample goes through each Example and consolidates all of the information it knows i.e.
// a custom example or auto-generate for the user. It also tracks whether we've randomized
// the entire example; if so, we shall re-generate the random value for Array/Hash.
func (a *AttributeDefinition) finalizeExample(stack []*AttributeDefinition) (interface{}, bool) {
	if a.Example != nil || a.isCustomExample {
		return a.Example, a.isCustomExample
	}

	// note: must traverse each node to finalize the examples unless given
	switch true {
	case a.Type.IsArray():
		ary := a.Type.ToArray()
		example, isCustom := ary.ElemType.finalizeExample(stack)
		a.Example, a.isCustomExample = ary.MakeSlice([]interface{}{example}), isCustom
	case a.Type.IsHash():
		h := a.Type.ToHash()
		exampleK, isCustomK := h.KeyType.finalizeExample(stack)
		exampleV, isCustomV := h.ElemType.finalizeExample(stack)
		a.Example, a.isCustomExample = h.MakeMap(map[interface{}]interface{}{exampleK: exampleV}), isCustomK || isCustomV
	case a.Type.IsObject():
		// keep track of the type id, in case of a cyclical situation
		stack = append(stack, a)

		// ensure fixed ordering
		aObj := a.Type.ToObject()
		keys := make([]string, 0, len(aObj))
		for n := range aObj {
			keys = append(keys, n)
		}
		sort.Strings(keys)

		example, hasCustom, isCustom := map[string]interface{}{}, false, false
		for _, n := range keys {
			att := aObj[n]
			// avoid a cyclical dependency
			isCyclical := false
			if ssize := len(stack); ssize > 0 {
				aid := ""
				if mt, ok := att.Type.(*MediaTypeDefinition); ok {
					aid = mt.Identifier
				} else if ut, ok := att.Type.(*UserTypeDefinition); ok {
					aid = ut.TypeName
				}
				if aid != "" {
					for _, sa := range stack[:ssize-1] {
						if mt, ok := sa.Type.(*MediaTypeDefinition); ok {
							isCyclical = mt.Identifier == aid
						} else if ut, ok := sa.Type.(*UserTypeDefinition); ok {
							isCyclical = ut.TypeName == aid
						}
						if isCyclical {
							break
						}
					}
				}
			}
			if !isCyclical {
				example[n], isCustom = att.finalizeExample(stack)
			} else {
				// unable to generate any example and here we set
				// isCustom to avoid touching this example again
				// i.e. GenerateExample in the end of this func
				example[n], isCustom = nil, true
			}
			hasCustom = hasCustom || isCustom
		}
		a.Example, a.isCustomExample = example, hasCustom
	}
	// while none of the examples is custom, we generate a random value for the entire object
	if !a.isCustomExample {
		a.Example = a.GenerateExample(Design.RandomGenerator())
	}
	return a.Example, a.isCustomExample
}

// Merge merges the argument attributes into the target and returns the target overriding existing
// attributes with identical names.
// This only applies to attributes of type Object and Merge panics if the
// argument or the target is not of type Object.
func (a *AttributeDefinition) Merge(other *AttributeDefinition) *AttributeDefinition {
	if other == nil {
		return a
	}
	if a == nil {
		return other
	}
	left := a.Type.(Object)
	right := other.Type.(Object)
	if left == nil || right == nil {
		panic("cannot merge non object attributes") // bug
	}
	for n, v := range right {
		left[n] = v
	}
	return a
}

// Inherit merges the properties of existing target type attributes with the argument's.
// The algorithm is recursive so that child attributes are also merged.
func (a *AttributeDefinition) Inherit(parent *AttributeDefinition) {
	if !a.shouldInherit(parent) {
		return
	}

	a.inheritValidations(parent)
	a.inheritRecursive(parent)
}

// DSL returns the initialization DSL.
func (a *AttributeDefinition) DSL() func() {
	return a.DSLFunc
}

func (a *AttributeDefinition) inheritRecursive(parent *AttributeDefinition) {
	if !a.shouldInherit(parent) {
		return
	}

	for n, att := range a.Type.ToObject() {
		if patt, ok := parent.Type.ToObject()[n]; ok {
			if att.Description == "" {
				att.Description = patt.Description
			}
			att.inheritValidations(patt)
			if att.DefaultValue == nil {
				att.DefaultValue = patt.DefaultValue
			}
			if att.View == "" {
				att.View = patt.View
			}
			if att.Type == nil {
				att.Type = patt.Type
			} else if att.shouldInherit(patt) {
				for _, att := range att.Type.ToObject() {
					att.Inherit(patt.Type.ToObject()[n])
				}
			}
		}
	}
}

func (a *AttributeDefinition) inheritValidations(parent *AttributeDefinition) {
	if parent.Validation == nil {
		return
	}
	if a.Validation == nil {
		a.Validation = &dslengine.ValidationDefinition{}
	}
	a.Validation.AddRequired(parent.Validation.Required)
}

func (a *AttributeDefinition) shouldInherit(parent *AttributeDefinition) bool {
	return a != nil && a.Type.ToObject() != nil &&
		parent != nil && parent.Type.ToObject() != nil
}

// Context returns the generic definition name used in error messages.
func (c *ContactDefinition) Context() string {
	if c.Name != "" {
		return fmt.Sprintf("contact %s", c.Name)
	}
	return "unnamed contact"
}

// Context returns the generic definition name used in error messages.
func (l *LicenseDefinition) Context() string {
	if l.Name != "" {
		return fmt.Sprintf("license %s", l.Name)
	}
	return "unnamed license"
}

// Context returns the generic definition name used in error messages.
func (d *DocsDefinition) Context() string {
	return fmt.Sprintf("documentation for %s", Design.Name)
}

// Context returns the generic definition name used in error messages.
func (t *UserTypeDefinition) Context() string {
	if t.TypeName != "" {
		return fmt.Sprintf("type %#v", t.TypeName)
	}
	return "unnamed type"
}

// DSL returns the initialization DSL.
func (t *UserTypeDefinition) DSL() func() {
	return t.DSLFunc
}

// Context returns the generic definition name used in error messages.
func (r *ResponseDefinition) Context() string {
	var prefix, suffix string
	if r.Name != "" {
		prefix = fmt.Sprintf("response %#v", r.Name)
	} else {
		prefix = "unnamed response"
	}
	if r.Parent != nil {
		suffix = fmt.Sprintf(" of %s", r.Parent.Context())
	}
	return prefix + suffix
}

// Finalize sets the response media type from its type if the type is a media type and no media
// type is already specified.
func (r *ResponseDefinition) Finalize() {
	if r.Type == nil {
		return
	}
	if r.MediaType != "" && r.MediaType != "plain/text" {
		return
	}
	mt, ok := r.Type.(*MediaTypeDefinition)
	if !ok {
		return
	}
	r.MediaType = mt.Identifier
}

// Dup returns a copy of the response definition.
func (r *ResponseDefinition) Dup() *ResponseDefinition {
	res := ResponseDefinition{
		Name:        r.Name,
		Status:      r.Status,
		Description: r.Description,
		MediaType:   r.MediaType,
	}
	if r.Headers != nil {
		res.Headers = DupAtt(r.Headers)
	}
	return &res
}

// Merge merges other into target. Only the fields of target that are not already set are merged.
func (r *ResponseDefinition) Merge(other *ResponseDefinition) {
	if other == nil {
		return
	}
	if r.Name == "" {
		r.Name = other.Name
	}
	if r.Status == 0 {
		r.Status = other.Status
	}
	if r.Description == "" {
		r.Description = other.Description
	}
	if r.MediaType == "" {
		r.MediaType = other.MediaType
	}
	if other.Headers != nil {
		otherHeaders := other.Headers.Type.ToObject()
		if len(otherHeaders) > 0 {
			if r.Headers == nil {
				r.Headers = &AttributeDefinition{Type: Object{}}
			}
			headers := r.Headers.Type.ToObject()
			for n, h := range otherHeaders {
				if _, ok := headers[n]; !ok {
					headers[n] = h
				}
			}
		}
	}
}

// Context returns the generic definition name used in error messages.
func (r *ResponseTemplateDefinition) Context() string {
	if r.Name != "" {
		return fmt.Sprintf("response template %#v", r.Name)
	}
	return "unnamed response template"
}

// Context returns the generic definition name used in error messages.
func (a *ActionDefinition) Context() string {
	var prefix, suffix string
	if a.Name != "" {
		suffix = fmt.Sprintf(" action %#v", a.Name)
	} else {
		suffix = " unnamed action"
	}
	if a.Parent != nil {
		prefix = a.Parent.Context()
	}
	return prefix + suffix
}

// PathParams returns the path parameters of the action across all its routes.
func (a *ActionDefinition) PathParams() *AttributeDefinition {
	obj := make(Object)
	for _, r := range a.Routes {
		for _, p := range r.Params() {
			if _, ok := obj[p]; !ok {
				obj[p] = a.Params.Type.ToObject()[p]
			}
		}
	}
	res := &AttributeDefinition{Type: obj}
	if a.HasAbsoluteRoutes() {
		return res
	}
	res = res.Merge(a.Parent.BaseParams)
	res = res.Merge(Design.BaseParams)
	if p := a.Parent.Parent(); p != nil {
		res = res.Merge(p.CanonicalAction().PathParams())
	}
	return res
}

// AllParams returns the path and query string parameters of the action across all its routes.
func (a *ActionDefinition) AllParams() *AttributeDefinition {
	var res *AttributeDefinition
	if a.Params != nil {
		res = DupAtt(a.Params)
	} else {
		res = &AttributeDefinition{Type: Object{}}
	}
	if a.HasAbsoluteRoutes() {
		return res
	}
	res = res.Merge(a.Parent.BaseParams)
	res = res.Merge(Design.BaseParams)
	if p := a.Parent.Parent(); p != nil {
		res = res.Merge(p.CanonicalAction().AllParams())
	}
	return res
}

// HasAbsoluteRoutes returns true if all the action routes are absolute.
func (a *ActionDefinition) HasAbsoluteRoutes() bool {
	for _, r := range a.Routes {
		if !r.IsAbsolute() {
			return false
		}
	}
	return true
}

// CanonicalScheme returns the preferred scheme for making requests. Favor secure schemes.
func (a *ActionDefinition) CanonicalScheme() string {
	if a.WebSocket() {
		for _, s := range a.EffectiveSchemes() {
			if s == "wss" {
				return s
			}
		}
		return "ws"
	}
	for _, s := range a.EffectiveSchemes() {
		if s == "https" {
			return s
		}
	}
	return "http"
}

// EffectiveSchemes return the URL schemes that apply to the action. Looks recursively into action
// resource, parent resources and API.
func (a *ActionDefinition) EffectiveSchemes() []string {
	// Compute the schemes
	schemes := a.Schemes
	if len(schemes) == 0 {
		res := a.Parent
		schemes = res.Schemes
		parent := res.Parent()
		for len(schemes) == 0 && parent != nil {
			schemes = parent.Schemes
			parent = parent.Parent()
		}
		if len(schemes) == 0 {
			schemes = Design.Schemes
		}
	}
	return schemes
}

// WebSocket returns true if the action scheme is "ws" or "wss" or both (directly or inherited
// from the resource or API)
func (a *ActionDefinition) WebSocket() bool {
	schemes := a.EffectiveSchemes()
	if len(schemes) == 0 {
		return false
	}
	for _, s := range schemes {
		if s != "ws" && s != "wss" {
			return false
		}
	}
	return true
}

// Finalize creates fallback security schemes and links before rendering.
func (a *ActionDefinition) Finalize() {
	if a.Security == nil {
		a.Security = a.Parent.Security // ResourceDefinition
		if a.Security == nil {
			a.Security = Design.Security
		}
	}

	if a.Security != nil && a.Security.Scheme.Kind == NoSecurityKind {
		a.Security = nil
	}
}

// Context returns the generic definition name used in error messages.
func (l *LinkDefinition) Context() string {
	var prefix, suffix string
	if l.Name != "" {
		prefix = fmt.Sprintf("link %#v", l.Name)
	} else {
		prefix = "unnamed link"
	}
	if l.Parent != nil {
		suffix = fmt.Sprintf(" of %s", l.Parent.Context())
	}
	return prefix + suffix
}

// Attribute returns the linked attribute.
func (l *LinkDefinition) Attribute() *AttributeDefinition {
	p := l.Parent.ToObject()
	if p == nil {
		return nil
	}
	att, _ := p[l.Name]

	return att
}

// MediaType returns the media type of the linked attribute.
func (l *LinkDefinition) MediaType() *MediaTypeDefinition {
	att := l.Attribute()
	mt, _ := att.Type.(*MediaTypeDefinition)
	return mt
}

// Context returns the generic definition name used in error messages.
func (v *ViewDefinition) Context() string {
	var prefix, suffix string
	if v.Name != "" {
		prefix = fmt.Sprintf("view %#v", v.Name)
	} else {
		prefix = "unnamed view"
	}
	if v.Parent != nil {
		suffix = fmt.Sprintf(" of %s", v.Parent.Context())
	}
	return prefix + suffix
}

// Context returns the generic definition name used in error messages.
func (r *RouteDefinition) Context() string {
	return fmt.Sprintf(`route %s "%s" of %s`, r.Verb, r.Path, r.Parent.Context())
}

// Params returns the route parameters.
// For example for the route "GET /foo/:fooID" Params returns []string{"fooID"}.
func (r *RouteDefinition) Params() []string {
	return ExtractWildcards(r.FullPath())
}

// FullPath returns the action full path computed by concatenating the API and resource base paths
// with the action specific path.
func (r *RouteDefinition) FullPath() string {
	if r.IsAbsolute() {
		return httppath.Clean(r.Path[1:])
	}
	var base string
	if r.Parent != nil && r.Parent.Parent != nil {
		base = r.Parent.Parent.FullPath()
	}
	return httppath.Clean(path.Join(base, r.Path))
}

// IsAbsolute returns true if the action path should not be concatenated to the resource and API
// base paths.
func (r *RouteDefinition) IsAbsolute() bool {
	return strings.HasPrefix(r.Path, "//")
}
