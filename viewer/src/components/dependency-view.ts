/**
 * Dependency View Component
 *
 * Renders dependencies between entities in a roadmap.
 * Shows: blocked entity → blocking entity, type, status, risk level.
 */

import type {
  JourneyRoadmap,
  Dependency,
  EntityRef,
  DependencyType,
  DependencyStatus,
  DependencyRisk,
} from '../schema/journey-roadmap';

export interface DependencyViewOptions {
  /** Filter by dependency type */
  filterType?: DependencyType | DependencyType[];
  /** Filter by status */
  filterStatus?: DependencyStatus | DependencyStatus[];
  /** Filter by risk level */
  filterRisk?: DependencyRisk | DependencyRisk[];
  /** Show only critical path (blocked + high/critical risk) */
  criticalPathOnly?: boolean;
  /** Group by entity type */
  groupBy?: 'none' | 'from' | 'to' | 'type' | 'status';
  /** Custom CSS class */
  className?: string;
}

const defaultOptions: DependencyViewOptions = {
  criticalPathOnly: false,
  groupBy: 'none',
};

/**
 * Render the dependency view as HTML string
 */
export function renderDependencyView(
  roadmap: JourneyRoadmap,
  options: DependencyViewOptions = {}
): string {
  const opts = { ...defaultOptions, ...options };

  if (!roadmap.dependencies || roadmap.dependencies.length === 0) {
    return '<div class="prism-dependency-view prism-empty">No dependencies defined</div>';
  }

  let dependencies = [...roadmap.dependencies];

  // Apply filters
  dependencies = applyFilters(dependencies, opts);

  if (dependencies.length === 0) {
    return '<div class="prism-dependency-view prism-empty">No dependencies match the filters</div>';
  }

  let html = `<div class="prism-dependency-view ${opts.className || ''}">`;
  html += '<h2>Dependencies</h2>';

  // Summary stats
  html += renderSummaryStats(dependencies);

  // Critical path section
  const criticalDeps = dependencies.filter(
    (d) => d.status === 'blocked' && (d.risk === 'high' || d.risk === 'critical')
  );
  if (criticalDeps.length > 0) {
    html += renderCriticalPath(criticalDeps);
  }

  // Main dependency table/list
  if (opts.groupBy === 'none' || !opts.groupBy) {
    html += renderDependencyTable(dependencies);
  } else {
    html += renderGroupedDependencies(dependencies, opts.groupBy);
  }

  html += '</div>';
  return html;
}

function applyFilters(
  deps: Dependency[],
  opts: DependencyViewOptions
): Dependency[] {
  let result = deps;

  // Critical path filter
  if (opts.criticalPathOnly) {
    result = result.filter(
      (d) => d.status === 'blocked' && (d.risk === 'high' || d.risk === 'critical')
    );
  }

  // Type filter
  if (opts.filterType) {
    const types = Array.isArray(opts.filterType) ? opts.filterType : [opts.filterType];
    result = result.filter((d) => types.includes(d.type));
  }

  // Status filter
  if (opts.filterStatus) {
    const statuses = Array.isArray(opts.filterStatus) ? opts.filterStatus : [opts.filterStatus];
    result = result.filter((d) => d.status && statuses.includes(d.status));
  }

  // Risk filter
  if (opts.filterRisk) {
    const risks = Array.isArray(opts.filterRisk) ? opts.filterRisk : [opts.filterRisk];
    result = result.filter((d) => d.risk && risks.includes(d.risk));
  }

  return result;
}

function renderSummaryStats(deps: Dependency[]): string {
  const stats = {
    total: deps.length,
    blocked: deps.filter((d) => d.status === 'blocked').length,
    atRisk: deps.filter((d) => d.status === 'at_risk').length,
    resolved: deps.filter((d) => d.status === 'resolved').length,
    pending: deps.filter((d) => d.status === 'pending').length,
    critical: deps.filter((d) => d.risk === 'critical').length,
    high: deps.filter((d) => d.risk === 'high').length,
  };

  return `
    <div class="prism-dependency-stats">
      <span class="prism-stat"><strong>${stats.total}</strong> Total</span>
      <span class="prism-stat prism-stat-blocked"><strong>${stats.blocked}</strong> Blocked</span>
      <span class="prism-stat prism-stat-at-risk"><strong>${stats.atRisk}</strong> At Risk</span>
      <span class="prism-stat prism-stat-resolved"><strong>${stats.resolved}</strong> Resolved</span>
      <span class="prism-stat prism-stat-critical"><strong>${stats.critical}</strong> Critical</span>
      <span class="prism-stat prism-stat-high"><strong>${stats.high}</strong> High Risk</span>
    </div>
  `;
}

function renderCriticalPath(criticalDeps: Dependency[]): string {
  let html = '<div class="prism-critical-path">';
  html += '<h3>⚠️ Critical Path</h3>';
  html += '<p>These dependencies are blocking progress and have high/critical risk:</p>';
  html += '<ul class="prism-critical-list">';

  for (const dep of criticalDeps) {
    const riskIcon = dep.risk === 'critical' ? '🔴' : '🟠';
    html += `<li class="prism-critical-item">
      ${riskIcon}
      <strong>${formatEntityRef(dep.from)}</strong>
      blocked by
      <strong>${formatEntityRef(dep.to)}</strong>
      ${dep.description ? `<br><small>${escapeHtml(dep.description)}</small>` : ''}
    </li>`;
  }

  html += '</ul></div>';
  return html;
}

function renderDependencyTable(deps: Dependency[]): string {
  let html = `
    <table class="prism-dependency-table">
      <thead>
        <tr>
          <th>From (Blocked)</th>
          <th>To (Blocking)</th>
          <th>Type</th>
          <th>Status</th>
          <th>Risk</th>
          <th>Expected Resolution</th>
        </tr>
      </thead>
      <tbody>
  `;

  for (const dep of deps) {
    const statusClass = getStatusClass(dep.status);
    const riskClass = getRiskClass(dep.risk);

    html += `
      <tr class="prism-dependency-row ${statusClass} ${riskClass}">
        <td>${formatEntityRef(dep.from)}</td>
        <td>${formatEntityRef(dep.to)}</td>
        <td><span class="prism-dep-type">${formatType(dep.type)}</span></td>
        <td><span class="prism-dep-status ${statusClass}">${formatStatus(dep.status)}</span></td>
        <td><span class="prism-dep-risk ${riskClass}">${formatRisk(dep.risk)}</span></td>
        <td>${dep.expectedResolution ? escapeHtml(dep.expectedResolution) : '-'}</td>
      </tr>
    `;
  }

  html += '</tbody></table>';
  return html;
}

function renderGroupedDependencies(
  deps: Dependency[],
  groupBy: 'from' | 'to' | 'type' | 'status'
): string {
  const groups = new Map<string, Dependency[]>();

  for (const dep of deps) {
    let key: string;
    switch (groupBy) {
      case 'from':
        key = `${dep.from.type}:${dep.from.id}`;
        break;
      case 'to':
        key = `${dep.to.type}:${dep.to.id}`;
        break;
      case 'type':
        key = dep.type;
        break;
      case 'status':
        key = dep.status || 'unknown';
        break;
      default:
        key = 'all';
    }

    if (!groups.has(key)) {
      groups.set(key, []);
    }
    groups.get(key)!.push(dep);
  }

  let html = '<div class="prism-dependency-groups">';

  for (const [key, groupDeps] of groups) {
    html += `<div class="prism-dependency-group">`;
    html += `<h4>${formatGroupKey(key, groupBy)} (${groupDeps.length})</h4>`;
    html += renderDependencyTable(groupDeps);
    html += '</div>';
  }

  html += '</div>';
  return html;
}

// Formatting helpers

function formatEntityRef(ref: EntityRef): string {
  const icon = getEntityIcon(ref.type);
  const name = ref.name || ref.id;
  return `${icon} <span class="prism-entity-name">${escapeHtml(name)}</span>`;
}

function getEntityIcon(type: string): string {
  const icons: Record<string, string> = {
    capability: '🎯',
    initiative: '🚀',
    team: '👥',
    milestone: '🏁',
    outcome: '📊',
    external: '🔗',
    decision: '⚖️',
  };
  return icons[type] || '📦';
}

function formatType(type: DependencyType): string {
  const labels: Record<DependencyType, string> = {
    requires: 'Requires',
    blocked_by: 'Blocked By',
    resource: 'Resource',
    external: 'External',
    informs: 'Informs',
    contributes: 'Contributes',
  };
  return labels[type] || type;
}

function formatStatus(status?: DependencyStatus): string {
  if (!status) return '-';
  const labels: Record<DependencyStatus, string> = {
    pending: '⏳ Pending',
    resolved: '✅ Resolved',
    blocked: '🚫 Blocked',
    at_risk: '⚠️ At Risk',
    waived: '↩️ Waived',
  };
  return labels[status] || status;
}

function formatRisk(risk?: DependencyRisk): string {
  if (!risk) return '-';
  const labels: Record<DependencyRisk, string> = {
    low: '🟢 Low',
    medium: '🟡 Medium',
    high: '🟠 High',
    critical: '🔴 Critical',
  };
  return labels[risk] || risk;
}

function formatGroupKey(key: string, groupBy: string): string {
  if (groupBy === 'type') {
    return formatType(key as DependencyType);
  }
  if (groupBy === 'status') {
    return formatStatus(key as DependencyStatus);
  }
  return key;
}

function getStatusClass(status?: DependencyStatus): string {
  if (!status) return '';
  return `prism-status-${status.replace('_', '-')}`;
}

function getRiskClass(risk?: DependencyRisk): string {
  if (!risk) return '';
  return `prism-risk-${risk}`;
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

/**
 * Get dependencies blocking a specific entity
 */
export function getBlockersFor(
  roadmap: JourneyRoadmap,
  entityType: string,
  entityId: string
): Dependency[] {
  if (!roadmap.dependencies) return [];
  return roadmap.dependencies.filter(
    (d) => d.from.type === entityType && d.from.id === entityId
  );
}

/**
 * Get dependencies where an entity is blocking others
 */
export function getBlockingBy(
  roadmap: JourneyRoadmap,
  entityType: string,
  entityId: string
): Dependency[] {
  if (!roadmap.dependencies) return [];
  return roadmap.dependencies.filter(
    (d) => d.to.type === entityType && d.to.id === entityId
  );
}

/**
 * Get critical path dependencies (blocked + high/critical risk)
 */
export function getCriticalPath(roadmap: JourneyRoadmap): Dependency[] {
  if (!roadmap.dependencies) return [];
  return roadmap.dependencies.filter(
    (d) => d.status === 'blocked' && (d.risk === 'high' || d.risk === 'critical')
  );
}
