package prd

// ServiceLayer represents the architectural layer of a service.
type ServiceLayer string

const (
	LayerControlPlane   ServiceLayer = "control-plane"
	LayerExecutionPlane ServiceLayer = "execution-plane"
	LayerDataPlane      ServiceLayer = "data-plane"
	LayerGateway        ServiceLayer = "gateway"
)

// Protocol represents an API protocol.
type Protocol string

const (
	ProtocolREST      Protocol = "REST"
	ProtocolGRPC      Protocol = "gRPC"
	ProtocolGraphQL   Protocol = "GraphQL"
	ProtocolWebSocket Protocol = "WebSocket"
)

// StorageCategoryType represents storage categories.
type StorageCategoryType string

const (
	StorageMetadata      StorageCategoryType = "metadata"
	StorageArtifacts     StorageCategoryType = "artifacts"
	StorageState         StorageCategoryType = "state"
	StorageCache         StorageCategoryType = "cache"
	StorageObservability StorageCategoryType = "observability"
	StorageAudit         StorageCategoryType = "audit"
	StorageSecrets       StorageCategoryType = "secrets"
)

// SourceOfTruthLocation represents where artifacts are stored.
type SourceOfTruthLocation string

const (
	LocationGit            SourceOfTruthLocation = "git"
	LocationS3             SourceOfTruthLocation = "s3"
	LocationDatabase       SourceOfTruthLocation = "database"
	LocationSecretsManager SourceOfTruthLocation = "secrets-manager"
	LocationRegistry       SourceOfTruthLocation = "registry"
)

// RelationshipType represents document relationships.
type RelationshipType string

const (
	RelationshipChild      RelationshipType = "child"
	RelationshipParent     RelationshipType = "parent"
	RelationshipSibling    RelationshipType = "sibling"
	RelationshipImplements RelationshipType = "implements"
	RelationshipSupersedes RelationshipType = "supersedes"
	RelationshipRelated    RelationshipType = "related"
)

// DocumentType represents types of related documents.
type DocumentType string

const (
	DocTypePRD       DocumentType = "prd"
	DocTypeTRD       DocumentType = "trd"
	DocTypeMRD       DocumentType = "mrd"
	DocTypeDesignDoc DocumentType = "design-doc"
	DocTypeRFC       DocumentType = "rfc"
)

// Service represents a microservice in the architecture.
type Service struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	Responsibilities  []string     `json:"responsibilities,omitempty"`
	Language          string       `json:"language,omitempty"`
	LanguageRationale string       `json:"languageRationale,omitempty"`
	Layer             ServiceLayer `json:"layer,omitempty"`
	Protocol          Protocol     `json:"protocol,omitempty"`
	Owner             string       `json:"owner,omitempty"`
	Dependencies      []string     `json:"dependencies,omitempty"`
	Tags              []string     `json:"tags,omitempty"`
}

// APIEndpoint represents a single API endpoint.
type APIEndpoint struct {
	ID          string   `json:"id,omitempty"`
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Description string   `json:"description,omitempty"`
	RequestBody string   `json:"requestBody,omitempty"`
	Response    string   `json:"response,omitempty"`
	Auth        string   `json:"auth,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// APISpec represents an API specification for a service.
type APISpec struct {
	ID           string        `json:"id,omitempty"`
	ServiceID    string        `json:"serviceId,omitempty"`
	Name         string        `json:"name"`
	Protocol     Protocol      `json:"protocol"`
	BasePath     string        `json:"basePath,omitempty"`
	Version      string        `json:"version,omitempty"`
	Description  string        `json:"description,omitempty"`
	Endpoints    []APIEndpoint `json:"endpoints,omitempty"`
	OpenAPISpec  string        `json:"openApiSpec,omitempty"`
	ProtobufSpec string        `json:"protobufSpec,omitempty"`
}

// StorageCategory represents a storage use case.
type StorageCategory struct {
	ID         string              `json:"id,omitempty"`
	Category   StorageCategoryType `json:"category"`
	Purpose    string              `json:"purpose"`
	Technology string              `json:"technology"`
	Rationale  string              `json:"rationale,omitempty"`
	Encryption string              `json:"encryption,omitempty"`
	Retention  string              `json:"retention,omitempty"`
	PerTenant  bool                `json:"perTenant,omitempty"`
}

// SourceOfTruth defines where authoritative data lives.
type SourceOfTruth struct {
	Artifact      string                `json:"artifact"`
	Location      SourceOfTruthLocation `json:"location"`
	Rationale     string                `json:"rationale,omitempty"`
	GitOpsEnabled bool                  `json:"gitOpsEnabled,omitempty"`
}

// GitOpsConfig defines GitOps integration.
type GitOpsConfig struct {
	Enabled        bool            `json:"enabled"`
	Provider       string          `json:"provider,omitempty"`
	Workflow       string          `json:"workflow,omitempty"`
	SourcesOfTruth []SourceOfTruth `json:"sourcesOfTruth,omitempty"`
}

// OrchestrationEngine represents a workflow orchestration engine.
type OrchestrationEngine struct {
	Name      string   `json:"name"`
	UseCases  []string `json:"useCases,omitempty"`
	Language  string   `json:"language,omitempty"`
	Rationale string   `json:"rationale,omitempty"`
}

// OrchestrationConfig defines workflow orchestration configuration.
type OrchestrationConfig struct {
	ShortLived  *OrchestrationEngine `json:"shortLived,omitempty"`
	LongRunning *OrchestrationEngine `json:"longRunning,omitempty"`
	Description string               `json:"description,omitempty"`
}

// RelatedDocument references another requirements document.
type RelatedDocument struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Type         DocumentType     `json:"type"`
	Relationship RelationshipType `json:"relationship"`
	Path         string           `json:"path,omitempty"`
	URL          string           `json:"url,omitempty"`
	Description  string           `json:"description,omitempty"`
}

// AssumptionsConstraints contains assumptions and constraints.
type AssumptionsConstraints struct {
	Assumptions  []Assumption `json:"assumptions"`
	Constraints  []Constraint `json:"constraints"`
	Dependencies []Dependency `json:"dependencies,omitempty"`
}

// Dependency represents an external dependency.
type Dependency struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type,omitempty"` // API, Service, Team, Vendor
	Owner       string `json:"owner,omitempty"`
	Status      string `json:"status,omitempty"` // Available, Pending, Blocked
	DueDate     string `json:"dueDate,omitempty"`
}

// TechnicalArchitecture contains technical design information.
type TechnicalArchitecture struct {
	Overview          string          `json:"overview"`
	SystemDiagram     string          `json:"systemDiagram,omitempty"` // URL or path to diagram
	DataModel         string          `json:"dataModel,omitempty"`     // URL or path to ERD
	IntegrationPoints []Integration   `json:"integrationPoints,omitempty"`
	TechnologyStack   TechnologyStack `json:"technologyStack,omitempty"`
	SecurityDesign    string          `json:"securityDesign,omitempty"`
	ScalabilityDesign string          `json:"scalabilityDesign,omitempty"`

	// Services lists the microservices in the architecture.
	Services []Service `json:"services,omitempty"`

	// APIs lists the API specifications.
	APIs []APISpec `json:"apis,omitempty"`

	// StorageArchitecture defines storage by category.
	StorageArchitecture []StorageCategory `json:"storageArchitecture,omitempty"`

	// GitOps defines GitOps configuration.
	GitOps *GitOpsConfig `json:"gitOps,omitempty"`

	// Orchestration defines workflow orchestration.
	Orchestration *OrchestrationConfig `json:"orchestration,omitempty"`
}

// Integration represents an external integration point.
type Integration struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"` // REST API, GraphQL, Event, Database
	Description   string `json:"description"`
	Protocol      string `json:"protocol,omitempty"`
	AuthMethod    string `json:"authMethod,omitempty"`
	DataFormat    string `json:"dataFormat,omitempty"` // JSON, XML, Protobuf
	RateLimit     string `json:"rateLimit,omitempty"`
	Documentation string `json:"documentation,omitempty"` // URL to docs
}

// TechnologyStack defines the technology choices.
type TechnologyStack struct {
	Frontend       []Technology `json:"frontend,omitempty"`
	Backend        []Technology `json:"backend,omitempty"`
	Database       []Technology `json:"database,omitempty"`
	Infrastructure []Technology `json:"infrastructure,omitempty"`
	DevOps         []Technology `json:"devops,omitempty"`
	Monitoring     []Technology `json:"monitoring,omitempty"`
}

// Technology represents a technology choice.
type Technology struct {
	Name         string   `json:"name"`
	Version      string   `json:"version,omitempty"`
	Purpose      string   `json:"purpose,omitempty"`
	Rationale    string   `json:"rationale,omitempty"`
	Alternatives []string `json:"alternatives,omitempty"` // Considered alternatives
}

// UXRequirements contains UX/UI requirements.
type UXRequirements struct {
	DesignPrinciples []string          `json:"designPrinciples,omitempty"`
	Wireframes       []Wireframe       `json:"wireframes,omitempty"`
	InteractionFlows []InteractionFlow `json:"interactionFlows,omitempty"`
	Accessibility    AccessibilitySpec `json:"accessibility,omitempty"`
	BrandGuidelines  string            `json:"brandGuidelines,omitempty"` // URL or path
	DesignSystem     string            `json:"designSystem,omitempty"`    // URL or path
}

// Wireframe represents a wireframe or mockup.
type Wireframe struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`              // Link to wireframe
	Status      string `json:"status,omitempty"` // Draft, Approved
}

// InteractionFlow represents a user interaction flow.
type InteractionFlow struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	DiagramURL  string   `json:"diagramUrl,omitempty"`
}

// AccessibilitySpec defines accessibility requirements.
type AccessibilitySpec struct {
	Standard        string   `json:"standard"` // WCAG 2.1 AA
	Requirements    []string `json:"requirements,omitempty"`
	TestingApproach string   `json:"testingApproach,omitempty"`
}
