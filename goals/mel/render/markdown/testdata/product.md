# Onboarding MEL Plan

**Team:** Product · **Period:** FY2027 · **Author:** Head of Product · **Status:** Active · **Version:** 1.0.0

## Purpose

Tells the product team, leadership, and the board whether the onboarding bet is working: whether new accounts reach a reconciled month and stay, which step of the chain is holding, and what to change next.

## Results Framework

- **Impact:** Small business owners spend less time on their books (`less-time-on-books`)
  - **Long-term outcome:** New accounts are still active at 90 days (`active-90`) — Retention past the first quarter is when an account has formed the habit the product depends on.
    - **Medium-term outcome:** Connected accounts reconcile weekly (`weekly-reconcile`)
      - **Short-term outcome:** New accounts connect a bank in their first session (`first-session-connect`)
        - **Output:** Guided bank import is released to all new accounts (`guided-import-live`)

## Activities

| Activity | Produces | Lead | Timing | Status |
|----------|----------|------|--------|--------|
| Build and roll out guided bank import | Guided bank import is released to all new accounts | Onboarding squad | 2027-01-04 to 2027-03-26 | In Progress |
| Send weekly reconciliation reminders | Connected accounts reconcile weekly | Lifecycle marketing | from 2027-01-04 | In Progress |

## Indicators

| Indicator | Measures | Baseline | Target | Frequency | Source | Responsible |
|-----------|----------|----------|--------|-----------|--------|-------------|
| Hours owners spend on their books each month | Small business owners spend less time on their books | 6 hours | 4 hours by 2027-12-31 | Every six months | Owner survey | User research lead |
| Accounts still active at 90 days | New accounts are still active at 90 days | 52% | 60% by 2027-12-31 | Monthly | Product analytics | Head of Product |
| Connected accounts that reconcile weekly | Connected accounts reconcile weekly | 30% | 50% by 2027-12-31 | Weekly | Product analytics | Head of Core Product |
| New accounts that connect a bank in their first session | New accounts connect a bank in their first session | 45% | 70% by 2027-06-30 | Weekly | Product analytics | Head of Onboarding |
| New accounts offered guided import | Guided bank import is released to all new accounts | 0% | 100% by 2027-03-31 | Weekly | Feature flag service | Onboarding squad |
| Bank connection attempts that succeed | Bank aggregators keep connection success above 90%. | 93% | 90% by 2027-12-31 | Weekly | Aggregator API logs | Platform engineering lead |

## Indicator Reference Sheets

### Hours owners spend on their books each month

| | |
|---|---|
| **ID** | `ind-hours-on-books` |
| **Measures** | Small business owners spend less time on their books |
| **Definition** | Median self-reported hours per month spent on bookkeeping, from the twice-yearly owner survey. |
| **Type** | quantitative |
| **Unit** | hours |
| **Direction** | decrease |
| **Disaggregation** | Business size, Tenure |
| **Baseline** | 6 hours (2026-11-30); source: Owner survey, November 2026 |
| **Targets** | 4 hours by 2027-12-31 |
| **Data source** | Owner survey |
| **Collection method** | In-app survey to a random sample of active accounts |
| **Frequency** | Every six months |
| **Responsible** | User research lead |
| **Data quality checks** | Sample of at least 400 responses; Weighted by business size |
| **Known limitations** | Self-reported; owners tend to underestimate time |

### Accounts still active at 90 days

| | |
|---|---|
| **ID** | `ind-active-90` |
| **Measures** | New accounts are still active at 90 days |
| **Definition** | Share of accounts created in a month that have any activity in their 13th week. |
| **Type** | quantitative |
| **Unit** | % |
| **Direction** | increase |
| **Baseline** | 52% (2026-12-31) |
| **Targets** | 56% by 2027-06-30 (Midyear); 60% by 2027-12-31 |
| **Current** | 53% |
| **Data source** | Product analytics |
| **Collection method** | Monthly cohort query |
| **Frequency** | Monthly |
| **Responsible** | Head of Product |
| **Data quality checks** | Test accounts excluded; Query reviewed when the event schema changes |

### Connected accounts that reconcile weekly

| | |
|---|---|
| **ID** | `ind-weekly-reconcile` |
| **Measures** | Connected accounts reconcile weekly |
| **Definition** | Share of accounts with a connected bank that reconcile in at least three of the last four weeks. |
| **Type** | quantitative |
| **Unit** | % |
| **Direction** | increase |
| **Disaggregation** | Bank aggregator |
| **Baseline** | 30% (2026-12-31) |
| **Targets** | 50% by 2027-12-31 |
| **Current** | 31% |
| **Data source** | Product analytics |
| **Collection method** | Weekly query on reconciliation events |
| **Frequency** | Weekly |
| **Responsible** | Head of Core Product |

### New accounts that connect a bank in their first session

| | |
|---|---|
| **ID** | `ind-first-session-connect` |
| **Measures** | New accounts connect a bank in their first session |
| **Type** | quantitative |
| **Unit** | % |
| **Direction** | increase |
| **Disaggregation** | Bank aggregator, Signup channel |
| **Baseline** | 45% (2026-12-31) |
| **Targets** | 70% by 2027-06-30 |
| **Current** | 51% |
| **Data source** | Product analytics |
| **Collection method** | Weekly funnel query |
| **Frequency** | Weekly |
| **Responsible** | Head of Onboarding |

### New accounts offered guided import

| | |
|---|---|
| **ID** | `ind-guided-import-rollout` |
| **Measures** | Guided bank import is released to all new accounts |
| **Type** | quantitative |
| **Unit** | % |
| **Direction** | increase |
| **Baseline** | 0% (2026-12-31) |
| **Targets** | 100% by 2027-03-31 |
| **Current** | 100% |
| **Data source** | Feature flag service |
| **Collection method** | Rollout percentage read from the flag |
| **Frequency** | Weekly |
| **Responsible** | Onboarding squad |

### Bank connection attempts that succeed

| | |
|---|---|
| **ID** | `ind-connection-success` |
| **Measures** | Bank aggregators keep connection success above 90%. |
| **Definition** | Share of bank connection attempts that the aggregator completes without an error. |
| **Type** | quantitative |
| **Unit** | % |
| **Direction** | maintain |
| **Disaggregation** | Bank aggregator |
| **Baseline** | 93% (2026-12-31) |
| **Targets** | 90% by 2027-12-31 (Floor) |
| **Data source** | Aggregator API logs |
| **Collection method** | Daily error-rate query |
| **Frequency** | Weekly |
| **Responsible** | Platform engineering lead |

## Assumptions to Monitor

| Assumption | Results that depend on it | How it is monitored |
|------------|---------------------------|---------------------|
| Connecting a bank, not reconciling itself, is what stops new accounts from forming the weekly habit. | Connected accounts reconcile weekly | Compare weekly reconciliation of accounts that connected in the first session with those that connected later. Indicators: Connected accounts that reconcile weekly, New accounts that connect a bank in their first session |
| Bank aggregators keep connection success above 90%. | New accounts connect a bank in their first session | Aggregator status pages and the weekly connection-failure report. Indicators: Bank connection attempts that succeed |

## Evaluation Plan

| Evaluation | Type | Timing | Evaluator | Lead | Status |
|------------|------|--------|-----------|------|--------|
| Onboarding baseline | baseline | 2026-11-01 to 2026-12-31 | internal | User research lead | Completed |
| Guided import usability study | formative | 2027-01-11 to 2027-03-26 | internal | Onboarding designer | In Progress |
| First-year results evaluation | impact | 2027-10-01 to 2028-01-31 | external | Head of Product | Planned |

### Onboarding baseline

**Method:** Funnel analysis of the last four cohorts and 20 owner interviews

1. Where do new accounts stall before reaching a reconciled month?

### Guided import usability study

**Method:** Moderated sessions with 12 new owners, repeated after each major change

1. Can owners connect a bank without help on the first try?
2. What makes owners abandon the connection step?

### First-year results evaluation

**Method:** Comparison with a 10% holdout cohort, plus the owner survey

**Budget:** Research agency engagement

1. Did guided import raise 90-day activity, compared with the holdout cohort?
2. Did owners who formed the weekly habit spend less time on their books?

## Learning Agenda

| Question | Why it matters | How it will be answered | Lead | Due | Status |
|----------|----------------|-------------------------|------|-----|--------|
| Does connecting a bank early lead accounts to reconcile weekly, or do committed owners simply do both? | If early connection is not the cause, the team should invest in the review screen instead of onboarding. | Holdout cohort comparison; Owner interviews | Head of Product | 2027-06-30 | Open |
| At which step of connecting a bank do owners give up, and why? | Decides which part of guided import to change first. | Funnel analysis; Usability sessions | Head of Onboarding | 2027-03-31 | Answered |

**At which step of connecting a bank do owners give up, and why?** Most owners who give up do so at bank login, when their bank asks for a second factor the flow did not explain.

## Reviews and Reporting

| Review | Frequency | Participants | Lead | Next | Questions taken up |
|--------|-----------|--------------|------|------|--------------------|
| Monthly metrics review | Monthly | Product leads, Data analyst | Head of Product | 2027-02-05 |  |
| Quarterly learning review | Quarterly | Product leads, Design, User research, Customer success | Head of Product | 2027-04-02 | `lq-connect-cause`, `lq-abandon-step` |

| Audience | Receives | Frequency | Channel | Prepared by | Indicators |
|----------|----------|-----------|---------|-------------|------------|
| Leadership team | Onboarding dashboard | Monthly | Leadership meeting | Data analyst | Accounts still active at 90 days; New accounts that connect a bank in their first session; Connected accounts that reconcile weekly |
| Board | Quarterly product update | Quarterly | Board deck | Head of Product | Accounts still active at 90 days; Hours owners spend on their books each month |

## Budget

About 5% of the product team's time, plus the research agency engagement for the first-year evaluation
