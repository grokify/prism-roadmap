/**
 * JourneyRoadmap - Zod Schema
 *
 * Hand-written Zod schema matching Go types in journey/*.go.
 * The auto-generated version doesn't handle $ref properly.
 *
 * Pipeline: Go types → JSON Schema → (manual) Zod Schema → TypeScript types
 *
 * To update: Sync with Go types in journey/*.go
 */

import { z } from 'zod';

// ============================================================================
// Time Model Types
// ============================================================================

export const TimeModelTypeSchema = z.enum([
  'quarterly',
  'monthly',
  'sprint',
  'milestone',
  'half',
  'custom',
]);
export type TimeModelType = z.infer<typeof TimeModelTypeSchema>;

export const PeriodSchema = z.object({
  id: z.string(),
  label: z.string(),
  startDate: z.string().optional(),
  endDate: z.string().optional(),
  isCurrent: z.boolean().optional(),
  description: z.string().optional(),
});
export type Period = z.infer<typeof PeriodSchema>;

export const TimeModelSchema = z.object({
  type: TimeModelTypeSchema,
  startDate: z.string().optional(),
  endDate: z.string().optional(),
  periods: z.array(PeriodSchema),
  fiscalYear: z.string().optional(),
});
export type TimeModel = z.infer<typeof TimeModelSchema>;

// ============================================================================
// Scope Types
// ============================================================================

export const ScopeSchema = z.object({
  type: z.string().optional(),
  id: z.string().optional(),
  name: z.string().optional(),
  description: z.string().optional(),
  tags: z.array(z.string()).optional(),
});
export type Scope = z.infer<typeof ScopeSchema>;

// ============================================================================
// Capability Journey Types
// ============================================================================

export const MaturityStateSchema = z.object({
  periodId: z.string(),
  maturityLevel: z.string(),
  summary: z.string().optional(),
  evidence: z.array(z.string()).optional(),
});
export type MaturityState = z.infer<typeof MaturityStateSchema>;

export const SuccessMeasureSchema = z.object({
  metric: z.string(),
  target: z.number(),
  unit: z.string().optional(),
  description: z.string().optional(),
});
export type SuccessMeasure = z.infer<typeof SuccessMeasureSchema>;

export const CommitmentLevelSchema = z.enum([
  'committed',
  'planned',
  'targeted',
  'aspirant',
]);
export type CommitmentLevel = z.infer<typeof CommitmentLevelSchema>;

export const TargetStateSchema = z.object({
  periodId: z.string(),
  maturityLevel: z.string(),
  summary: z.string().optional(),
  changes: z.array(z.string()).optional(),
  initiatives: z.array(z.string()).optional(),
  successMeasures: z.array(SuccessMeasureSchema).optional(),
  confidence: z.number().optional(),
  commitment: CommitmentLevelSchema.optional(),
  assumptions: z.array(z.string()).optional(),
  scenarioId: z.string().optional(),
});
export type TargetState = z.infer<typeof TargetStateSchema>;

export const CapabilityJourneySchema = z.object({
  id: z.string(),
  capabilityId: z.string(),
  name: z.string(),
  description: z.string().optional(),
  owner: z.string().optional(),
  currentState: MaturityStateSchema.nullable().optional(),
  targetStates: z.array(TargetStateSchema),
  desiredEndState: MaturityStateSchema.nullable().optional(),
  tags: z.array(z.string()).optional(),
  metadata: z.record(z.string(), z.any()).optional(),
});
export type CapabilityJourney = z.infer<typeof CapabilityJourneySchema>;

// ============================================================================
// Outcome Journey Types
// ============================================================================

export const OutcomeStateSchema = z.object({
  periodId: z.string().optional(),
  value: z.number(),
  unit: z.string().optional(),
  summary: z.string().optional(),
  evidence: z.array(z.string()).optional(),
  measuredAt: z.string().optional(),
});
export type OutcomeState = z.infer<typeof OutcomeStateSchema>;

export const OutcomeTargetStateSchema = z.object({
  periodId: z.string(),
  value: z.number(),
  unit: z.string().optional(),
  summary: z.string().optional(),
  confidence: z.number().optional(),
  enabledBy: z.array(z.string()).optional(),
  assumptions: z.array(z.string()).optional(),
  scenarioId: z.string().optional(),
});
export type OutcomeTargetState = z.infer<typeof OutcomeTargetStateSchema>;

export const OutcomeDirectionSchema = z.enum([
  'increase',
  'decrease',
  'maintain',
]);
export type OutcomeDirection = z.infer<typeof OutcomeDirectionSchema>;

export const OutcomeJourneySchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  metric: z.string().optional(),
  unit: z.string().optional(),
  direction: OutcomeDirectionSchema.optional(),
  currentState: OutcomeStateSchema.nullable().optional(),
  targetStates: z.array(OutcomeTargetStateSchema),
  tags: z.array(z.string()).optional(),
  metadata: z.record(z.string(), z.any()).optional(),
});
export type OutcomeJourney = z.infer<typeof OutcomeJourneySchema>;

// ============================================================================
// Initiative Types
// ============================================================================

export const InitiativeStatusSchema = z.enum([
  'proposed',
  'planned',
  'in_progress',
  'completed',
  'on_hold',
  'cancelled',
]);
export type InitiativeStatus = z.infer<typeof InitiativeStatusSchema>;

export const CapacitySchema = z.object({
  storyPoints: z.number().optional(),
  ftes: z.number().optional(),
  fteMonths: z.number().optional(),
  customValue: z.number().optional(),
  customUnit: z.string().optional(),
  notes: z.string().optional(),
});
export type Capacity = z.infer<typeof CapacitySchema>;

export const CapabilityAdvanceSchema = z.object({
  capabilityId: z.string().optional(),
  capabilityName: z.string().optional(),
  from: z.string(),
  to: z.string(),
  description: z.string().optional(),
});
export type CapabilityAdvance = z.infer<typeof CapabilityAdvanceSchema>;

export const LinkSchema = z.object({
  url: z.string(),
  title: z.string().optional(),
  type: z.string().optional(),
});
export type Link = z.infer<typeof LinkSchema>;

export const InitiativeSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  status: InitiativeStatusSchema.optional(),
  ownerTeam: z.string().optional(),
  contributingTeams: z.array(z.string()).optional(),
  periods: z.array(z.string()).optional(),
  advances: z.array(CapabilityAdvanceSchema).optional(),
  expectedOutcomes: z.array(z.string()).optional(),
  requiredCapacity: CapacitySchema.nullable().optional(),
  dependencies: z.array(z.string()).optional(),
  risks: z.array(z.string()).optional(),
  links: z.array(LinkSchema).optional(),
  tags: z.array(z.string()).optional(),
  metadata: z.record(z.string(), z.any()).optional(),
});
export type Initiative = z.infer<typeof InitiativeSchema>;

// ============================================================================
// Dependency Types
// ============================================================================

export const EntityTypeSchema = z.enum([
  'capability',
  'initiative',
  'team',
  'milestone',
  'outcome',
  'external',
  'decision',
]);
export type EntityType = z.infer<typeof EntityTypeSchema>;

export const EntityRefSchema = z.object({
  type: EntityTypeSchema,
  id: z.string(),
  name: z.string().optional(),
});
export type EntityRef = z.infer<typeof EntityRefSchema>;

export const DependencyTypeSchema = z.enum([
  'requires',
  'blocked_by',
  'resource',
  'external',
  'informs',
  'contributes',
]);
export type DependencyType = z.infer<typeof DependencyTypeSchema>;

export const DependencyStatusSchema = z.enum([
  'pending',
  'resolved',
  'blocked',
  'at_risk',
  'waived',
]);
export type DependencyStatus = z.infer<typeof DependencyStatusSchema>;

export const DependencyRiskSchema = z.enum([
  'low',
  'medium',
  'high',
  'critical',
]);
export type DependencyRisk = z.infer<typeof DependencyRiskSchema>;

export const DependencySchema = z.object({
  id: z.string().optional(),
  from: EntityRefSchema,
  to: EntityRefSchema,
  type: DependencyTypeSchema,
  description: z.string().optional(),
  status: DependencyStatusSchema.optional(),
  risk: DependencyRiskSchema.optional(),
  expectedResolution: z.string().optional(),
  notes: z.string().optional(),
});
export type Dependency = z.infer<typeof DependencySchema>;

// ============================================================================
// Team Types
// ============================================================================

export const TeamTypeSchema = z.enum([
  'engineering',
  'platform',
  'product',
  'design',
  'data',
  'infrastructure',
  'security',
  'qa',
  'devops',
  'sre',
  'external',
]);
export type TeamType = z.infer<typeof TeamTypeSchema>;

export const TeamLevelSchema = z.enum([
  'organization',
  'division',
  'department',
  'group',
  'team',
  'squad',
]);
export type TeamLevel = z.infer<typeof TeamLevelSchema>;

export const TeamCapacitySchema = z.object({
  ftes: z.number().optional(),
  storyPointsPerSprint: z.number().optional(),
  storyPointsPerQuarter: z.number().optional(),
  allocatedPercent: z.number().optional(),
  reservedPercent: z.number().optional(),
  effectiveDate: z.string().optional(),
  notes: z.string().optional(),
});
export type TeamCapacity = z.infer<typeof TeamCapacitySchema>;

export const TeamSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  type: TeamTypeSchema.optional(),
  level: TeamLevelSchema.optional(),
  parentId: z.string().optional(),
  childIds: z.array(z.string()).optional(),
  leaderId: z.string().optional(),
  leaderName: z.string().optional(),
  capacity: TeamCapacitySchema.nullable().optional(),
  skills: z.array(z.string()).optional(),
  costCenter: z.string().optional(),
  location: z.string().optional(),
  tags: z.array(z.string()).optional(),
  metadata: z.record(z.string(), z.any()).optional(),
});
export type Team = z.infer<typeof TeamSchema>;

// ============================================================================
// Narrative Types
// ============================================================================

export const JourneyChapterSchema = z.object({
  periodId: z.string(),
  headline: z.string(),
  story: z.string(),
  keyChanges: z.array(z.string()).optional(),
  milestones: z.array(z.string()).optional(),
  userImpact: z.string().optional(),
  risks: z.array(z.string()).optional(),
});
export type JourneyChapter = z.infer<typeof JourneyChapterSchema>;

export const RoadmapNarrativeSchema = z.object({
  title: z.string(),
  currentState: z.string().optional(),
  turningPoint: z.string().optional(),
  journey: z.array(JourneyChapterSchema).optional(),
  destination: z.string().optional(),
  callToAction: z.string().optional(),
});
export type RoadmapNarrative = z.infer<typeof RoadmapNarrativeSchema>;

// ============================================================================
// Risk Types
// ============================================================================

export const RiskSchema = z.object({
  id: z.string(),
  description: z.string(),
  probability: z.string().optional(),
  impact: z.string().optional(),
  mitigation: z.string().optional(),
  status: z.string().optional(),
  affectedPeriods: z.array(z.string()).optional(),
  tags: z.array(z.string()).optional(),
});
export type Risk = z.infer<typeof RiskSchema>;

// ============================================================================
// Scenario Types
// ============================================================================

export const ScenarioSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  isBase: z.boolean().optional(),
});
export type Scenario = z.infer<typeof ScenarioSchema>;

// ============================================================================
// Journey Roadmap (Top-Level Container)
// ============================================================================

export const JourneyRoadmapSchema = z.object({
  id: z.string(),
  type: z.string().optional(),
  name: z.string(),
  vision: z.string().optional(),
  description: z.string().optional(),
  scope: ScopeSchema.nullable().optional(),
  timeModel: TimeModelSchema.nullable().optional(),
  capabilityJourneys: z.array(CapabilityJourneySchema).optional(),
  outcomeJourneys: z.array(OutcomeJourneySchema).optional(),
  initiatives: z.array(InitiativeSchema).optional(),
  dependencies: z.array(DependencySchema).optional(),
  teams: z.array(TeamSchema).optional(),
  risks: z.array(RiskSchema).optional(),
  narrative: RoadmapNarrativeSchema.nullable().optional(),
  scenarios: z.array(ScenarioSchema).optional(),
  metadata: z.record(z.string(), z.any()).optional(),
});
export type JourneyRoadmap = z.infer<typeof JourneyRoadmapSchema>;

// ============================================================================
// Validation Helpers
// ============================================================================

/**
 * Validate a JourneyRoadmap object (throws on invalid data).
 */
export function validateJourneyRoadmap(data: unknown): JourneyRoadmap {
  return JourneyRoadmapSchema.parse(data);
}

/**
 * Safe validation helper (returns result instead of throwing).
 */
export function safeValidateJourneyRoadmap(
  data: unknown
): z.SafeParseReturnType<unknown, JourneyRoadmap> {
  return JourneyRoadmapSchema.safeParse(data);
}

// ============================================================================
// Maturity Level Constants
// ============================================================================

export const MaturityLevels = {
  M0: 'M0', // Initial/Ad-hoc
  M1: 'M1', // Developing/Repeatable
  M2: 'M2', // Defined
  M3: 'M3', // Managed/Measured
  M4: 'M4', // Optimizing
  M5: 'M5', // Innovating/Leading
} as const;

export type MaturityLevel = (typeof MaturityLevels)[keyof typeof MaturityLevels];
