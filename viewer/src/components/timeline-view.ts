/**
 * Timeline View Component
 *
 * Renders capability journeys over time periods as a timeline/heat map:
 * Period columns → Capability rows → Maturity levels as cells
 */

import type {
  JourneyRoadmap,
  CapabilityJourney,
  Period,
} from '../schema/journey-roadmap';

export interface TimelineViewOptions {
  /** Show confidence indicators */
  showConfidence?: boolean;
  /** Show commitment level badges */
  showCommitment?: boolean;
  /** Custom CSS class for the container */
  className?: string;
  /** Color scheme for maturity levels */
  colorScheme?: 'default' | 'heatmap' | 'monochrome';
}

const defaultOptions: TimelineViewOptions = {
  showConfidence: true,
  showCommitment: true,
  colorScheme: 'default',
};

/**
 * Render the timeline view as HTML string
 */
export function renderTimelineView(
  roadmap: JourneyRoadmap,
  options: TimelineViewOptions = {}
): string {
  const opts = { ...defaultOptions, ...options };

  if (!roadmap.timeModel?.periods || roadmap.timeModel.periods.length === 0) {
    return '<div class="prism-timeline-view prism-empty">No time periods defined</div>';
  }

  if (!roadmap.capabilityJourneys || roadmap.capabilityJourneys.length === 0) {
    return '<div class="prism-timeline-view prism-empty">No capability journeys defined</div>';
  }

  const periods = roadmap.timeModel.periods;
  const journeys = roadmap.capabilityJourneys;

  let html = `<div class="prism-timeline-view ${opts.className || ''}">`;

  // Header
  if (roadmap.name) {
    html += `<h2>${escapeHtml(roadmap.name)}</h2>`;
  }
  if (roadmap.vision) {
    html += `<p class="prism-vision">${escapeHtml(roadmap.vision)}</p>`;
  }

  // Timeline table
  html += '<div class="prism-timeline-container">';
  html += '<table class="prism-timeline-table">';

  // Header row with periods
  html += '<thead><tr><th class="prism-capability-header">Capability</th>';
  for (const period of periods) {
    const currentClass = period.isCurrent ? 'prism-period-current' : '';
    html += `<th class="prism-period-header ${currentClass}">${escapeHtml(period.label)}</th>`;
  }
  html += '</tr></thead>';

  // Body rows with capabilities
  html += '<tbody>';
  for (const journey of journeys) {
    html += renderCapabilityRow(journey, periods, opts);
  }
  html += '</tbody>';

  html += '</table>';
  html += '</div>';

  // Legend
  html += renderLegend(opts);

  html += '</div>';
  return html;
}

function renderCapabilityRow(
  journey: CapabilityJourney,
  periods: Period[],
  opts: TimelineViewOptions
): string {
  let html = '<tr class="prism-capability-row">';

  // Capability name cell
  html += `<td class="prism-capability-name">
    <strong>${escapeHtml(journey.name)}</strong>
    ${journey.owner ? `<br><small>${escapeHtml(journey.owner)}</small>` : ''}
  </td>`;

  // Period cells
  for (const period of periods) {
    const state = getStateForPeriod(journey, period.id);
    html += renderPeriodCell(state, period, journey, opts);
  }

  html += '</tr>';
  return html;
}

function getStateForPeriod(
  journey: CapabilityJourney,
  periodId: string
): { maturityLevel: string; confidence?: number; commitment?: string } | null {
  // Check if this is the current state
  if (journey.currentState?.periodId === periodId) {
    return {
      maturityLevel: journey.currentState.maturityLevel,
    };
  }

  // Check target states
  const target = journey.targetStates.find((t) => t.periodId === periodId);
  if (target) {
    return {
      maturityLevel: target.maturityLevel,
      confidence: target.confidence,
      commitment: target.commitment,
    };
  }

  return null;
}

function renderPeriodCell(
  state: { maturityLevel: string; confidence?: number; commitment?: string } | null,
  period: Period,
  _journey: CapabilityJourney,
  opts: TimelineViewOptions
): string {
  if (!state) {
    return '<td class="prism-period-cell prism-empty-cell">-</td>';
  }

  const levelNum = parseMaturityLevel(state.maturityLevel);
  const colorClass = getMaturityColorClass(levelNum, opts.colorScheme);
  const currentClass = period.isCurrent ? 'prism-current' : '';

  let cellContent = `<span class="prism-maturity-badge ${colorClass}">${escapeHtml(state.maturityLevel)}</span>`;

  // Add confidence indicator
  if (opts.showConfidence && state.confidence !== undefined) {
    const confPercent = Math.round(state.confidence * 100);
    const confClass = getConfidenceClass(state.confidence);
    cellContent += `<span class="prism-confidence ${confClass}" title="Confidence: ${confPercent}%">${confPercent}%</span>`;
  }

  // Add commitment badge
  if (opts.showCommitment && state.commitment) {
    const commitIcon = getCommitmentIcon(state.commitment);
    cellContent += `<span class="prism-commitment" title="${state.commitment}">${commitIcon}</span>`;
  }

  return `<td class="prism-period-cell ${currentClass}" data-level="${levelNum}">${cellContent}</td>`;
}

function renderLegend(opts: TimelineViewOptions): string {
  let html = '<div class="prism-timeline-legend">';
  html += '<strong>Legend:</strong> ';

  // Maturity levels
  const levels = ['M0', 'M1', 'M2', 'M3', 'M4', 'M5'];
  const labels = ['Ad-hoc', 'Developing', 'Defined', 'Managed', 'Optimizing', 'Leading'];

  for (let i = 0; i < levels.length; i++) {
    const colorClass = getMaturityColorClass(i, opts.colorScheme);
    html += `<span class="prism-legend-item">
      <span class="prism-maturity-badge ${colorClass}">${levels[i]}</span>
      <span class="prism-legend-label">${labels[i]}</span>
    </span>`;
  }

  // Commitment icons
  if (opts.showCommitment) {
    html += '<span class="prism-legend-separator">|</span>';
    html += `<span class="prism-legend-item">🎯 Committed</span>`;
    html += `<span class="prism-legend-item">📋 Planned</span>`;
    html += `<span class="prism-legend-item">🎯 Targeted</span>`;
    html += `<span class="prism-legend-item">⭐ Aspirant</span>`;
  }

  html += '</div>';
  return html;
}

// Helper functions

function parseMaturityLevel(level: string): number {
  const match = level.match(/M(\d+)/i);
  return match ? parseInt(match[1], 10) : 0;
}

function getMaturityColorClass(level: number, scheme?: string): string {
  if (scheme === 'monochrome') {
    return `prism-maturity-mono-${Math.min(level, 5)}`;
  }
  if (scheme === 'heatmap') {
    return `prism-maturity-heat-${Math.min(level, 5)}`;
  }
  return `prism-maturity-${Math.min(level, 5)}`;
}

function getConfidenceClass(confidence: number): string {
  if (confidence >= 0.8) return 'prism-confidence-high';
  if (confidence >= 0.5) return 'prism-confidence-medium';
  return 'prism-confidence-low';
}

function getCommitmentIcon(commitment: string): string {
  const icons: Record<string, string> = {
    committed: '🎯',
    planned: '📋',
    targeted: '🎯',
    aspirant: '⭐',
  };
  return icons[commitment] || '❓';
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}
