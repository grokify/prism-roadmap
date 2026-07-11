/**
 * Style exports
 *
 * CSS path for importing styles.
 */

// CSS file path for bundlers that support CSS imports
export const STYLES_PATH = './prism-roadmap.css';

// Inline styles as a string (for SSR or dynamic injection)
export const INLINE_STYLES = `
/* See prism-roadmap.css for full styles */
/* This is a minimal subset for SSR */

.prism-timeline-view,
.prism-storyboard-view,
.prism-dependency-view,
.prism-team-view {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-size: 14px;
  line-height: 1.5;
}

.prism-maturity-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
}

.prism-maturity-0 { background: #f5f5f5; color: #666; }
.prism-maturity-1 { background: #ffebee; color: #c62828; }
.prism-maturity-2 { background: #fff3e0; color: #ef6c00; }
.prism-maturity-3 { background: #fffde7; color: #f9a825; }
.prism-maturity-4 { background: #e8f5e9; color: #2e7d32; }
.prism-maturity-5 { background: #e3f2fd; color: #1565c0; }
`;
