/**
 * Team View Component
 *
 * Renders team hierarchy and capacity information.
 * Shows: team structure, capacity allocation, initiative ownership.
 */

import type {
  JourneyRoadmap,
  Team,
  TeamCapacity,
  TeamLevel,
  Initiative,
} from '../schema/journey-roadmap';

export interface TeamViewOptions {
  /** Show capacity details */
  showCapacity?: boolean;
  /** Show initiative ownership */
  showInitiatives?: boolean;
  /** Filter by team level */
  filterLevel?: TeamLevel | TeamLevel[];
  /** Expand all hierarchy levels */
  expandAll?: boolean;
  /** Custom CSS class */
  className?: string;
}

const defaultOptions: TeamViewOptions = {
  showCapacity: true,
  showInitiatives: true,
  expandAll: false,
};

interface TeamNode {
  team: Team;
  children: TeamNode[];
  initiatives: Initiative[];
  aggregateCapacity?: TeamCapacity;
}

/**
 * Render the team view as HTML string
 */
export function renderTeamView(
  roadmap: JourneyRoadmap,
  options: TeamViewOptions = {}
): string {
  const opts = { ...defaultOptions, ...options };

  if (!roadmap.teams || roadmap.teams.length === 0) {
    return '<div class="prism-team-view prism-empty">No teams defined</div>';
  }

  // Build team hierarchy
  const hierarchy = buildTeamHierarchy(roadmap.teams, roadmap.initiatives || []);

  // Apply filters
  let teams = hierarchy;
  if (opts.filterLevel) {
    const levels = Array.isArray(opts.filterLevel) ? opts.filterLevel : [opts.filterLevel];
    teams = filterByLevel(hierarchy, levels);
  }

  let html = `<div class="prism-team-view ${opts.className || ''}">`;
  html += '<h2>Team Structure</h2>';

  // Summary stats
  html += renderTeamStats(roadmap.teams, roadmap.initiatives || []);

  // Team hierarchy
  html += '<div class="prism-team-hierarchy">';
  for (const node of teams) {
    html += renderTeamNode(node, 0, opts);
  }
  html += '</div>';

  html += '</div>';
  return html;
}

function buildTeamHierarchy(teams: Team[], initiatives: Initiative[]): TeamNode[] {
  const teamMap = new Map<string, TeamNode>();

  // Create nodes for all teams
  for (const team of teams) {
    teamMap.set(team.id, {
      team,
      children: [],
      initiatives: [],
      aggregateCapacity: team.capacity ? { ...team.capacity } : undefined,
    });
  }

  // Link initiatives to teams
  for (const init of initiatives) {
    if (init.ownerTeam && teamMap.has(init.ownerTeam)) {
      teamMap.get(init.ownerTeam)!.initiatives.push(init);
    }
  }

  // Build hierarchy
  const roots: TeamNode[] = [];
  for (const node of teamMap.values()) {
    if (node.team.parentId && teamMap.has(node.team.parentId)) {
      teamMap.get(node.team.parentId)!.children.push(node);
    } else {
      roots.push(node);
    }
  }

  // Calculate aggregate capacities
  for (const root of roots) {
    calculateAggregateCapacity(root);
  }

  // Sort children by name
  const sortChildren = (nodes: TeamNode[]) => {
    nodes.sort((a, b) => a.team.name.localeCompare(b.team.name));
    for (const node of nodes) {
      sortChildren(node.children);
    }
  };
  sortChildren(roots);

  return roots;
}

function calculateAggregateCapacity(node: TeamNode): TeamCapacity {
  const capacity: TeamCapacity = {
    ftes: node.team.capacity?.ftes || 0,
    storyPointsPerSprint: node.team.capacity?.storyPointsPerSprint || 0,
    storyPointsPerQuarter: node.team.capacity?.storyPointsPerQuarter || 0,
  };

  for (const child of node.children) {
    const childCap = calculateAggregateCapacity(child);
    capacity.ftes = (capacity.ftes || 0) + (childCap.ftes || 0);
    capacity.storyPointsPerSprint = (capacity.storyPointsPerSprint || 0) + (childCap.storyPointsPerSprint || 0);
    capacity.storyPointsPerQuarter = (capacity.storyPointsPerQuarter || 0) + (childCap.storyPointsPerQuarter || 0);
  }

  node.aggregateCapacity = capacity;
  return capacity;
}

function filterByLevel(nodes: TeamNode[], levels: TeamLevel[]): TeamNode[] {
  const result: TeamNode[] = [];

  const collect = (node: TeamNode) => {
    if (node.team.level && levels.includes(node.team.level)) {
      result.push(node);
    }
    for (const child of node.children) {
      collect(child);
    }
  };

  for (const node of nodes) {
    collect(node);
  }

  return result;
}

function renderTeamStats(teams: Team[], initiatives: Initiative[]): string {
  const levelCounts = new Map<string, number>();
  let totalFTEs = 0;

  for (const team of teams) {
    const level = team.level || 'unknown';
    levelCounts.set(level, (levelCounts.get(level) || 0) + 1);
    if (team.capacity?.ftes) {
      totalFTEs += team.capacity.ftes;
    }
  }

  const ownedInitiatives = new Set(initiatives.filter(i => i.ownerTeam).map(i => i.ownerTeam));

  return `
    <div class="prism-team-stats">
      <span class="prism-stat"><strong>${teams.length}</strong> Teams</span>
      <span class="prism-stat"><strong>${totalFTEs.toFixed(1)}</strong> FTEs</span>
      <span class="prism-stat"><strong>${ownedInitiatives.size}</strong> Teams with Initiatives</span>
      <span class="prism-stat"><strong>${initiatives.length}</strong> Total Initiatives</span>
    </div>
  `;
}

function renderTeamNode(node: TeamNode, depth: number, opts: TeamViewOptions): string {
  const indentClass = `prism-team-depth-${Math.min(depth, 5)}`;
  const hasChildren = node.children.length > 0;
  const expandedClass = opts.expandAll || depth < 2 ? 'prism-expanded' : 'prism-collapsed';

  let html = `<div class="prism-team-node ${indentClass} ${expandedClass}" data-team-id="${escapeHtml(node.team.id)}">`;

  // Team header
  html += '<div class="prism-team-header">';

  // Expand/collapse toggle
  if (hasChildren) {
    html += '<span class="prism-team-toggle">▶</span>';
  } else {
    html += '<span class="prism-team-toggle-placeholder"></span>';
  }

  // Team icon based on type
  const icon = getTeamIcon(node.team.type);
  html += `<span class="prism-team-icon">${icon}</span>`;

  // Team name and level
  html += `<span class="prism-team-name">${escapeHtml(node.team.name)}</span>`;
  if (node.team.level) {
    html += `<span class="prism-team-level">${formatLevel(node.team.level)}</span>`;
  }

  // Leader
  if (node.team.leaderName) {
    html += `<span class="prism-team-leader">👤 ${escapeHtml(node.team.leaderName)}</span>`;
  }

  html += '</div>';

  // Capacity section
  if (opts.showCapacity && (node.team.capacity || node.aggregateCapacity)) {
    html += renderCapacity(node.team.capacity, node.aggregateCapacity, hasChildren);
  }

  // Initiatives section
  if (opts.showInitiatives && node.initiatives.length > 0) {
    html += renderInitiatives(node.initiatives);
  }

  // Children
  if (hasChildren) {
    html += '<div class="prism-team-children">';
    for (const child of node.children) {
      html += renderTeamNode(child, depth + 1, opts);
    }
    html += '</div>';
  }

  html += '</div>';
  return html;
}

function renderCapacity(
  direct?: TeamCapacity | null,
  aggregate?: TeamCapacity,
  hasChildren?: boolean
): string {
  let html = '<div class="prism-team-capacity">';

  if (direct) {
    html += '<span class="prism-capacity-item">';
    if (direct.ftes) {
      html += `<strong>${direct.ftes}</strong> FTEs`;
    }
    if (direct.storyPointsPerSprint) {
      html += ` | <strong>${direct.storyPointsPerSprint}</strong> pts/sprint`;
    }
    if (direct.allocatedPercent !== undefined) {
      const available = 100 - (direct.allocatedPercent || 0) - (direct.reservedPercent || 0);
      html += ` | <strong>${available.toFixed(0)}%</strong> available`;
    }
    html += '</span>';
  }

  if (hasChildren && aggregate && aggregate.ftes !== direct?.ftes) {
    html += `<span class="prism-capacity-aggregate">(Total: ${aggregate.ftes?.toFixed(1)} FTEs)</span>`;
  }

  html += '</div>';
  return html;
}

function renderInitiatives(initiatives: Initiative[]): string {
  let html = '<div class="prism-team-initiatives">';
  html += '<strong>Owns:</strong> ';
  html += initiatives.map(i => {
    const statusIcon = getInitiativeStatusIcon(i.status);
    return `${statusIcon} ${escapeHtml(i.name)}`;
  }).join(', ');
  html += '</div>';
  return html;
}

// Helper functions

function getTeamIcon(type?: string): string {
  const icons: Record<string, string> = {
    engineering: '💻',
    platform: '🏗️',
    product: '📦',
    design: '🎨',
    data: '📊',
    infrastructure: '🔧',
    security: '🔒',
    qa: '✅',
    devops: '🔄',
    sre: '📟',
    external: '🔗',
  };
  return icons[type || ''] || '👥';
}

function formatLevel(level: TeamLevel): string {
  const labels: Record<TeamLevel, string> = {
    organization: 'Org',
    division: 'Division',
    department: 'Dept',
    group: 'Group',
    team: 'Team',
    squad: 'Squad',
  };
  return labels[level] || level;
}

function getInitiativeStatusIcon(status?: string): string {
  const icons: Record<string, string> = {
    proposed: '💡',
    planned: '📋',
    in_progress: '🔄',
    completed: '✅',
    on_hold: '⏸️',
    cancelled: '❌',
  };
  return icons[status || ''] || '📌';
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
 * Get team by ID from roadmap
 */
export function getTeamById(roadmap: JourneyRoadmap, teamId: string): Team | undefined {
  return roadmap.teams?.find(t => t.id === teamId);
}

/**
 * Get all teams at a specific level
 */
export function getTeamsByLevel(roadmap: JourneyRoadmap, level: TeamLevel): Team[] {
  return roadmap.teams?.filter(t => t.level === level) || [];
}

/**
 * Get initiatives owned by a team
 */
export function getTeamInitiatives(roadmap: JourneyRoadmap, teamId: string): Initiative[] {
  return roadmap.initiatives?.filter(i => i.ownerTeam === teamId) || [];
}

/**
 * Calculate aggregate capacity for a team and its descendants
 */
export function aggregateTeamCapacity(roadmap: JourneyRoadmap, teamId: string): TeamCapacity {
  const result: TeamCapacity = { ftes: 0, storyPointsPerSprint: 0, storyPointsPerQuarter: 0 };

  const addCapacity = (id: string) => {
    const team = roadmap.teams?.find(t => t.id === id);
    if (!team) return;

    if (team.capacity) {
      result.ftes = (result.ftes || 0) + (team.capacity.ftes || 0);
      result.storyPointsPerSprint = (result.storyPointsPerSprint || 0) + (team.capacity.storyPointsPerSprint || 0);
      result.storyPointsPerQuarter = (result.storyPointsPerQuarter || 0) + (team.capacity.storyPointsPerQuarter || 0);
    }

    // Add children
    for (const child of roadmap.teams || []) {
      if (child.parentId === id) {
        addCapacity(child.id);
      }
    }
  };

  addCapacity(teamId);
  return result;
}
