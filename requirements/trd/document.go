// Package trd provides data types for structured Technical Requirements Documents.
package trd

import (
	"time"

	"github.com/grokify/prism-execution/common"
)

// Person is an alias for common.Person for backwards compatibility.
type Person = common.Person

// Approver is an alias for common.Approver for backwards compatibility.
type Approver = common.Approver

// Status is an alias for common.Status for backwards compatibility.
type Status = common.Status

// GlossaryTerm is an alias for common.GlossaryTerm for backwards compatibility.
type GlossaryTerm = common.GlossaryTerm

// CustomSection is an alias for common.CustomSection for backwards compatibility.
type CustomSection = common.CustomSection

// Status constants re-exported from common for backward compatibility.
const (
	StatusDraft      = common.StatusDraft
	StatusInReview   = common.StatusInReview
	StatusApproved   = common.StatusApproved
	StatusDeprecated = common.StatusDeprecated
)

// Document represents a complete Technical Requirements Document.
type Document struct {
	Metadata          Metadata         `json:"metadata"`
	ExecutiveSummary  ExecutiveSummary `json:"executiveSummary"`
	Architecture      Architecture     `json:"architecture"`
	TechnologyStack   TechnologyStack  `json:"technologyStack"`
	APISpecifications []APISpec        `json:"apiSpecifications,omitempty"`
	DataModel         *DataModel       `json:"dataModel,omitempty"`
	SecurityDesign    SecurityDesign   `json:"securityDesign"`
	Performance       Performance      `json:"performance"`
	Scalability       *Scalability     `json:"scalability,omitempty"`
	Deployment        Deployment       `json:"deployment"`
	Integration       []Integration    `json:"integrations,omitempty"`
	Development       *Development     `json:"development,omitempty"`
	Testing           *Testing         `json:"testing,omitempty"`

	// Optional sections
	Risks          []Risk          `json:"risks,omitempty"`
	Constraints    []Constraint    `json:"constraints,omitempty"`
	Assumptions    []Assumption    `json:"assumptions,omitempty"`
	Glossary       []GlossaryTerm  `json:"glossary,omitempty"`
	CustomSections []CustomSection `json:"customSections,omitempty"`
}

// Note: Status type and constants are defined in common/ and aliased above.

// Metadata contains document metadata.
type Metadata struct {
	ID               string       `json:"id"`
	Title            string       `json:"title"`
	Version          string       `json:"version"`
	Status           Status       `json:"status"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	Authors          []Person     `json:"authors"`
	Reviewers        []Person     `json:"reviewers,omitempty"`
	Approvers        []Approver   `json:"approvers,omitempty"`
	Tags             []string     `json:"tags,omitempty"`
	RelatedDocuments []RelatedDoc `json:"relatedDocuments,omitempty"`
}

// RelatedDoc represents a related document reference.
type RelatedDoc struct {
	Title        string `json:"title"`
	URL          string `json:"url,omitempty"`
	Relationship string `json:"relationship,omitempty"` // implements, extends, references
	Description  string `json:"description,omitempty"`
}

// ExecutiveSummary provides high-level technical overview.
type ExecutiveSummary struct {
	Purpose           string   `json:"purpose"`
	Scope             string   `json:"scope"`
	TechnicalApproach string   `json:"technicalApproach"`
	KeyDecisions      []string `json:"keyDecisions,omitempty"`
	OutOfScope        []string `json:"outOfScope,omitempty"`
}

// Architecture contains system architecture details.
type Architecture struct {
	Overview      string         `json:"overview"`
	Principles    []string       `json:"principles,omitempty"`
	Patterns      []string       `json:"patterns,omitempty"` // e.g., "Microservices", "Event-driven"
	Components    []Component    `json:"components"`
	Diagrams      []Diagram      `json:"diagrams,omitempty"`
	DataFlows     []DataFlow     `json:"dataFlows,omitempty"`
	ArchDecisions []ArchDecision `json:"architectureDecisions,omitempty"`
}

// Component represents a system component.
type Component struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Type             string   `json:"type,omitempty"` // Service, Library, Database, Queue, etc.
	Responsibilities []string `json:"responsibilities,omitempty"`
	Dependencies     []string `json:"dependencies,omitempty"` // IDs of dependent components
	Technology       string   `json:"technology,omitempty"`
	Owner            string   `json:"owner,omitempty"`
	Tags             []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// Diagram represents an architecture diagram.
type Diagram struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type,omitempty"` // C4, Sequence, ER, Flowchart
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source,omitempty"` // Mermaid, PlantUML, draw.io, etc.
}

// DataFlow represents a data flow between components.
type DataFlow struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Source      string `json:"source"`      // Component ID
	Destination string `json:"destination"` // Component ID
	DataType    string `json:"dataType,omitempty"`
	Protocol    string `json:"protocol,omitempty"` // HTTP, gRPC, AMQP, etc.
	Description string `json:"description,omitempty"`
}

// ArchDecision represents an Architecture Decision Record (ADR).
type ArchDecision struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Status       string   `json:"status"` // Proposed, Accepted, Deprecated, Superseded
	Context      string   `json:"context"`
	Decision     string   `json:"decision"`
	Consequences []string `json:"consequences,omitempty"`
	Alternatives []string `json:"alternatives,omitempty"`
	Date         string   `json:"date,omitempty"`
	Tags         []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// TechnologyStack defines technology choices.
type TechnologyStack struct {
	Languages      []Technology `json:"languages,omitempty"`
	Frameworks     []Technology `json:"frameworks,omitempty"`
	Databases      []Technology `json:"databases,omitempty"`
	MessageQueues  []Technology `json:"messageQueues,omitempty"`
	Caching        []Technology `json:"caching,omitempty"`
	Infrastructure []Technology `json:"infrastructure,omitempty"`
	Monitoring     []Technology `json:"monitoring,omitempty"`
	CICD           []Technology `json:"cicd,omitempty"`
	Other          []Technology `json:"other,omitempty"`
}

// Technology represents a technology choice.
type Technology struct {
	Name         string   `json:"name"`
	Version      string   `json:"version,omitempty"`
	Purpose      string   `json:"purpose"`
	Rationale    string   `json:"rationale,omitempty"`
	Alternatives []string `json:"alternatives,omitempty"`
	Constraints  []string `json:"constraints,omitempty"`
}

// APISpec represents an API specification.
type APISpec struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Type        string        `json:"type"` // REST, gRPC, GraphQL, WebSocket
	Version     string        `json:"version,omitempty"`
	Description string        `json:"description,omitempty"`
	BaseURL     string        `json:"baseUrl,omitempty"`
	Auth        string        `json:"auth,omitempty"` // OAuth2, API Key, mTLS, etc.
	Endpoints   []APIEndpoint `json:"endpoints,omitempty"`
	SpecURL     string        `json:"specUrl,omitempty"` // OpenAPI, Proto file, etc.
	RateLimit   string        `json:"rateLimit,omitempty"`
	Tags        []string      `json:"tags,omitempty"` // For filtering by topic/domain
}

// APIEndpoint represents an API endpoint.
type APIEndpoint struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Description string   `json:"description,omitempty"`
	Request     string   `json:"request,omitempty"`  // Schema reference
	Response    string   `json:"response,omitempty"` // Schema reference
	Errors      []string `json:"errors,omitempty"`
}

// DataModel contains data modeling information.
type DataModel struct {
	Overview   string      `json:"overview"`
	Entities   []Entity    `json:"entities,omitempty"`
	Diagrams   []Diagram   `json:"diagrams,omitempty"`
	DataStores []DataStore `json:"dataStores,omitempty"`
	Migrations string      `json:"migrations,omitempty"` // Migration strategy
}

// Entity represents a data entity.
type Entity struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description,omitempty"`
	Attributes    []Attribute `json:"attributes,omitempty"`
	Relationships []string    `json:"relationships,omitempty"`
	Tags          []string    `json:"tags,omitempty"` // For filtering by topic/domain
}

// Attribute represents an entity attribute.
type Attribute struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
	Constraints string `json:"constraints,omitempty"`
}

// DataStore represents a data storage system.
type DataStore struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // PostgreSQL, MongoDB, Redis, S3, etc.
	Purpose     string   `json:"purpose"`
	Capacity    string   `json:"capacity,omitempty"`
	Replication string   `json:"replication,omitempty"`
	Backup      string   `json:"backup,omitempty"`
	Tags        []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// SecurityDesign contains security architecture.
type SecurityDesign struct {
	Overview         string            `json:"overview"`
	AuthN            *AuthN            `json:"authentication,omitempty"`
	AuthZ            *AuthZ            `json:"authorization,omitempty"`
	Encryption       *Encryption       `json:"encryption,omitempty"`
	NetworkSecurity  *NetworkSecurity  `json:"networkSecurity,omitempty"`
	Compliance       []string          `json:"compliance,omitempty"` // SOC2, HIPAA, PCI-DSS, etc.
	ThreatModel      []Threat          `json:"threatModel,omitempty"`
	SecurityControls []SecurityControl `json:"securityControls,omitempty"`
}

// AuthN represents authentication design.
type AuthN struct {
	Method      string `json:"method"` // OAuth2, SAML, mTLS, API Key
	Provider    string `json:"provider,omitempty"`
	MFA         bool   `json:"mfa,omitempty"`
	SessionMgmt string `json:"sessionManagement,omitempty"`
	Details     string `json:"details,omitempty"`
}

// AuthZ represents authorization design.
type AuthZ struct {
	Model    string   `json:"model"` // RBAC, ABAC, ReBAC
	Policies string   `json:"policies,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	Details  string   `json:"details,omitempty"`
}

// Encryption represents encryption design.
type Encryption struct {
	AtRest    string `json:"atRest,omitempty"`    // AES-256, etc.
	InTransit string `json:"inTransit,omitempty"` // TLS 1.3, mTLS
	KeyMgmt   string `json:"keyManagement,omitempty"`
	Details   string `json:"details,omitempty"`
}

// NetworkSecurity represents network security design.
type NetworkSecurity struct {
	Firewall      string `json:"firewall,omitempty"`
	WAF           string `json:"waf,omitempty"`
	DDoS          string `json:"ddosProtection,omitempty"`
	NetworkPolicy string `json:"networkPolicy,omitempty"`
	Segmentation  string `json:"segmentation,omitempty"`
	Details       string `json:"details,omitempty"`
}

// Threat represents a security threat.
type Threat struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category,omitempty"` // STRIDE categories
	Description string   `json:"description"`
	Likelihood  string   `json:"likelihood,omitempty"`
	Impact      string   `json:"impact,omitempty"`
	Mitigation  string   `json:"mitigation"`
	Tags        []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// SecurityControl represents a security control.
type SecurityControl struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Category       string   `json:"category,omitempty"`
	Description    string   `json:"description"`
	Implementation string   `json:"implementation,omitempty"`
	Tags           []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// Performance contains performance requirements and design.
type Performance struct {
	Overview      string            `json:"overview,omitempty"`
	Requirements  []PerfRequirement `json:"requirements"`
	Benchmarks    []Benchmark       `json:"benchmarks,omitempty"`
	Optimizations []string          `json:"optimizations,omitempty"`
}

// PerfRequirement represents a performance requirement.
type PerfRequirement struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Metric      string   `json:"metric"` // Latency, Throughput, etc.
	Target      string   `json:"target"` // e.g., "< 100ms p99"
	Priority    string   `json:"priority,omitempty"`
	Measurement string   `json:"measurement,omitempty"`
	Tags        []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// Benchmark represents a performance benchmark.
type Benchmark struct {
	Name     string `json:"name"`
	Scenario string `json:"scenario"`
	Result   string `json:"result,omitempty"`
	Date     string `json:"date,omitempty"`
}

// Scalability contains scalability design.
type Scalability struct {
	Overview        string  `json:"overview"`
	HorizontalScale string  `json:"horizontalScaling,omitempty"`
	VerticalScale   string  `json:"verticalScaling,omitempty"`
	LoadBalancing   string  `json:"loadBalancing,omitempty"`
	AutoScaling     string  `json:"autoScaling,omitempty"`
	Limits          []Limit `json:"limits,omitempty"`
}

// Limit represents a system limit.
type Limit struct {
	Name         string `json:"name"`
	Value        string `json:"value"`
	Rationale    string `json:"rationale,omitempty"`
	Configurable bool   `json:"configurable,omitempty"`
}

// Deployment contains deployment architecture.
type Deployment struct {
	Overview       string        `json:"overview"`
	Environments   []Environment `json:"environments"`
	Strategy       string        `json:"strategy,omitempty"`       // Blue-green, Canary, Rolling
	Infrastructure string        `json:"infrastructure,omitempty"` // Kubernetes, VMs, Serverless
	Regions        []string      `json:"regions,omitempty"`
	HA             string        `json:"highAvailability,omitempty"`
	DR             string        `json:"disasterRecovery,omitempty"`
}

// Environment represents a deployment environment.
type Environment struct {
	Name        string `json:"name"` // Development, Staging, Production
	Purpose     string `json:"purpose,omitempty"`
	URL         string `json:"url,omitempty"`
	Resources   string `json:"resources,omitempty"`
	AccessLevel string `json:"accessLevel,omitempty"`
}

// Integration represents an external integration.
type Integration struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`      // API, SDK, Webhook, File
	Direction     string   `json:"direction"` // Inbound, Outbound, Bidirectional
	Protocol      string   `json:"protocol,omitempty"`
	AuthMethod    string   `json:"authMethod,omitempty"`
	DataFormat    string   `json:"dataFormat,omitempty"`
	Frequency     string   `json:"frequency,omitempty"` // Real-time, Batch, On-demand
	Description   string   `json:"description,omitempty"`
	Documentation string   `json:"documentation,omitempty"`
	Tags          []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// Development contains development standards.
type Development struct {
	CodingStandards string   `json:"codingStandards,omitempty"`
	BranchStrategy  string   `json:"branchStrategy,omitempty"`
	CodeReview      string   `json:"codeReview,omitempty"`
	Documentation   string   `json:"documentation,omitempty"`
	Tools           []string `json:"tools,omitempty"`
}

// Testing contains testing strategy.
type Testing struct {
	Strategy     string   `json:"strategy"`
	UnitTests    string   `json:"unitTests,omitempty"`
	Integration  string   `json:"integrationTests,omitempty"`
	E2E          string   `json:"e2e_tests,omitempty"`
	Performance  string   `json:"performanceTests,omitempty"`
	Security     string   `json:"securityTests,omitempty"`
	Coverage     string   `json:"coverageRequirements,omitempty"`
	Environments []string `json:"testEnvironments,omitempty"`
}

// Risk represents a technical risk.
type Risk struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Probability string   `json:"probability"`
	Impact      string   `json:"impact"`
	Mitigation  string   `json:"mitigation"`
	Owner       string   `json:"owner,omitempty"`
	Status      string   `json:"status,omitempty"`
	Tags        []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// Constraint represents a technical constraint.
type Constraint struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"` // Technical, Resource, Time, Budget
	Description string   `json:"description"`
	Impact      string   `json:"impact,omitempty"`
	Rationale   string   `json:"rationale,omitempty"`
	Tags        []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// Assumption represents a technical assumption.
type Assumption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Rationale   string `json:"rationale,omitempty"`
	Validated   bool   `json:"validated,omitempty"`
	Risk        string `json:"risk,omitempty"`
}

// Note: GlossaryTerm and CustomSection types are defined in common/ and aliased above.
