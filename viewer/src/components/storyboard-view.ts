/**
 * Storyboard View Component
 *
 * Renders roadmap as a series of storyboard cards, one per period.
 * Each card shows: headline, maturity changes, initiatives, impact, risks.
 */

import type {
  JourneyRoadmap,
  CapabilityJourney,
  Initiative,
} from '../schema/journey-roadmap';

export interface StoryboardViewOptions {
  /** Show maturity change details */
  showMaturityChanges?: boolean;
  /** Show initiatives list */
  showInitiatives?: boolean;
  /** Show risks section */
  showRisks?: boolean;
  /** Show confidence indicator */
  showConfidence?: boolean;
  /** Card layout: horizontal scroll or vertical stack */
  layout?: 'horizontal' | 'vertical';
  /** Custom CSS class for the container */
  className?: string;
}

const defaultOptions: StoryboardViewOptions = {
  showMaturityChanges: true,
  showInitiatives: true,
  showRisks: true,
  showConfidence: true,
  layout: 'horizontal',
};

export interface StoryboardCard {
  periodId: string;
  periodLabel: string;
  headline: string;
  story?: string;
  maturityChanges: MaturityChange[];
  initiatives: string[];
  userImpact?: string;
  risks: string[];
  overallConfidence: number;
  isCurrent: boolean;
}

interface MaturityChange {
  capabilityName: string;
  from: string;
  to: string;
}

/**
 * Render the storyboard view as HTML string
 */
export function renderStoryboardView(
  roadmap: JourneyRoadmap,
  options: StoryboardViewOptions = {}
): string {
  const opts = { ...defaultOptions, ...options };

  if (!roadmap.timeModel?.periods || roadmap.timeModel.periods.length === 0) {
    return '<div class="prism-storyboard-view prism-empty">No time periods defined</div>';
  }

  const cards = buildStoryboardCards(roadmap);
  const layoutClass = opts.layout === 'vertical' ? 'prism-storyboard-vertical' : 'prism-storyboard-horizontal';

  let html = `<div class="prism-storyboard-view ${layoutClass} ${opts.className || ''}">`;

  // Header with narrative
  if (roadmap.narrative) {
    html += renderNarrativeHeader(roadmap.narrative);
  } else if (roadmap.name) {
    html += `<h2>${escapeHtml(roadmap.name)}</h2>`;
  }

  // Cards container
  html += '<div class="prism-storyboard-cards">';
  for (const card of cards) {
    html += renderCard(card, opts);
  }
  html += '</div>';

  // Call to action
  if (roadmap.narrative?.callToAction) {
    html += `<div class="prism-call-to-action">
      <strong>Call to Action:</strong> ${escapeHtml(roadmap.narrative.callToAction)}
    </div>`;
  }

  html += '</div>';
  return html;
}

function renderNarrativeHeader(narrative: {
  title: string;
  currentState?: string;
  turningPoint?: string;
  destination?: string;
}): string {
  let html = '<div class="prism-narrative-header">';

  html += `<h2>${escapeHtml(narrative.title)}</h2>`;

  if (narrative.currentState) {
    html += `<div class="prism-narrative-section">
      <strong>Where We Are:</strong> ${escapeHtml(narrative.currentState)}
    </div>`;
  }

  if (narrative.turningPoint) {
    html += `<div class="prism-narrative-section">
      <strong>What Changed:</strong> ${escapeHtml(narrative.turningPoint)}
    </div>`;
  }

  if (narrative.destination) {
    html += `<div class="prism-narrative-section">
      <strong>Where We're Going:</strong> ${escapeHtml(narrative.destination)}
    </div>`;
  }

  html += '</div>';
  return html;
}

function renderCard(card: StoryboardCard, opts: StoryboardViewOptions): string {
  const currentClass = card.isCurrent ? 'prism-card-current' : '';
  const confClass = getConfidenceClass(card.overallConfidence);

  let html = `<div class="prism-storyboard-card ${currentClass}" data-period="${escapeHtml(card.periodId)}">`;

  // Card header
  html += `<div class="prism-card-header">
    <span class="prism-period-label">${escapeHtml(card.periodLabel)}</span>
    ${opts.showConfidence && card.overallConfidence > 0
      ? `<span class="prism-confidence-badge ${confClass}">${Math.round(card.overallConfidence * 100)}%</span>`
      : ''
    }
  </div>`;

  // Headline
  html += `<h3 class="prism-card-headline">${escapeHtml(card.headline || 'Untitled')}</h3>`;

  // Story
  if (card.story) {
    html += `<p class="prism-card-story">${escapeHtml(card.story)}</p>`;
  }

  // Maturity changes
  if (opts.showMaturityChanges && card.maturityChanges.length > 0) {
    html += '<div class="prism-maturity-changes">';
    html += '<strong>Capability Evolution:</strong>';
    html += '<ul>';
    for (const change of card.maturityChanges) {
      html += `<li>
        <span class="prism-capability-name">${escapeHtml(change.capabilityName)}</span>:
        <span class="prism-maturity-from">${escapeHtml(change.from)}</span>
        →
        <span class="prism-maturity-to">${escapeHtml(change.to)}</span>
      </li>`;
    }
    html += '</ul></div>';
  }

  // Initiatives
  if (opts.showInitiatives && card.initiatives.length > 0) {
    html += '<div class="prism-initiatives">';
    html += '<strong>Key Initiatives:</strong>';
    html += '<ul>';
    for (const init of card.initiatives) {
      html += `<li>${escapeHtml(init)}</li>`;
    }
    html += '</ul></div>';
  }

  // User impact
  if (card.userImpact) {
    html += `<div class="prism-user-impact">
      <strong>User Impact:</strong> ${escapeHtml(card.userImpact)}
    </div>`;
  }

  // Risks
  if (opts.showRisks && card.risks.length > 0) {
    html += '<div class="prism-risks">';
    html += '<strong>Risks:</strong>';
    html += '<ul class="prism-risk-list">';
    for (const risk of card.risks) {
      html += `<li class="prism-risk-item">⚠️ ${escapeHtml(risk)}</li>`;
    }
    html += '</ul></div>';
  }

  html += '</div>';
  return html;
}

/**
 * Build storyboard cards from a journey roadmap
 */
export function buildStoryboardCards(roadmap: JourneyRoadmap): StoryboardCard[] {
  if (!roadmap.timeModel?.periods) {
    return [];
  }

  const cards: StoryboardCard[] = [];

  for (const period of roadmap.timeModel.periods) {
    const card: StoryboardCard = {
      periodId: period.id,
      periodLabel: period.label,
      headline: '',
      maturityChanges: [],
      initiatives: [],
      risks: [],
      overallConfidence: 0,
      isCurrent: period.isCurrent || false,
    };

    // Get narrative chapter if available
    if (roadmap.narrative?.journey) {
      const chapter = roadmap.narrative.journey.find((c) => c.periodId === period.id);
      if (chapter) {
        card.headline = chapter.headline;
        card.story = chapter.story;
        card.userImpact = chapter.userImpact;
        card.risks = chapter.risks || [];
      }
    }

    // Collect maturity changes
    if (roadmap.capabilityJourneys) {
      card.maturityChanges = collectMaturityChanges(roadmap.capabilityJourneys, period.id);
    }

    // Collect initiatives
    if (roadmap.initiatives) {
      card.initiatives = collectInitiatives(roadmap.initiatives, period.id);
    }

    // Calculate average confidence
    if (roadmap.capabilityJourneys) {
      card.overallConfidence = calculateAverageConfidence(roadmap.capabilityJourneys, period.id);
    }

    // Generate headline if not provided by narrative
    if (!card.headline) {
      card.headline = generateHeadline(card);
    }

    cards.push(card);
  }

  return cards;
}

function collectMaturityChanges(
  journeys: CapabilityJourney[],
  periodId: string
): MaturityChange[] {
  const changes: MaturityChange[] = [];

  for (const journey of journeys) {
    for (let i = 0; i < journey.targetStates.length; i++) {
      const target = journey.targetStates[i];
      if (target.periodId === periodId) {
        let fromLevel: string;
        if (i === 0 && journey.currentState) {
          fromLevel = journey.currentState.maturityLevel;
        } else if (i > 0) {
          fromLevel = journey.targetStates[i - 1].maturityLevel;
        } else {
          continue;
        }

        if (fromLevel !== target.maturityLevel) {
          changes.push({
            capabilityName: journey.name,
            from: fromLevel,
            to: target.maturityLevel,
          });
        }
      }
    }
  }

  return changes;
}

function collectInitiatives(initiatives: Initiative[], periodId: string): string[] {
  const names: string[] = [];

  for (const init of initiatives) {
    if (init.periods?.includes(periodId)) {
      names.push(init.name);
    }
  }

  return names;
}

function calculateAverageConfidence(
  journeys: CapabilityJourney[],
  periodId: string
): number {
  let totalConf = 0;
  let count = 0;

  for (const journey of journeys) {
    for (const target of journey.targetStates) {
      if (target.periodId === periodId && target.confidence !== undefined && target.confidence > 0) {
        totalConf += target.confidence;
        count++;
      }
    }
  }

  return count > 0 ? totalConf / count : 0;
}

function generateHeadline(card: StoryboardCard): string {
  if (card.maturityChanges.length > 0) {
    const capCount = card.maturityChanges.length;
    return `${capCount} Capability${capCount > 1 ? ' Advances' : ' Advance'}`;
  }
  if (card.initiatives.length > 0) {
    return card.initiatives[0];
  }
  if (card.isCurrent) {
    return 'Current State';
  }
  return card.periodLabel;
}

function getConfidenceClass(confidence: number): string {
  if (confidence >= 0.8) return 'prism-confidence-high';
  if (confidence >= 0.5) return 'prism-confidence-medium';
  return 'prism-confidence-low';
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}
