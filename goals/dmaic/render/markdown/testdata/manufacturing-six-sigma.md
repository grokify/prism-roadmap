# Manufacturing Quality Metrics

**Team:** Quality Engineering · **Period:** FY2025-Q1 · **Owner:** Maria Rodriguez · **Status:** Active · **Review cadence:** Weekly · **Version:** 1.0.0

## Overview

DMAIC metrics framework for manufacturing quality improvement and Six Sigma initiatives

**Overall health:** 83% · **Categories:** 3 · **Metrics:** 9

| Category | Health | Metrics | Green | Yellow | Red |
|:---------|-------:|--------:|------:|-------:|----:|
| Production Quality | 83% | 3 | 2 | 1 | 0 |
| Equipment Performance | 83% | 3 | 2 | 1 | 0 |
| Process Control | 83% | 3 | 2 | 1 | 0 |

## Production Quality

Metrics tracking production line quality and defect rates

**Owner:** Tom Chen · **Health:** 83%

| Metric | Phase | Baseline | Current | Target | Unit | Status |
|:-------|:-----:|---------:|--------:|-------:|:-----|:------:|
| First Pass Yield | Control | 85.0% | 94.2% | 98.0% | % | Green |
| Defects Per Million Opportunities | Improve | 15000 | 3400 | 2000 | DPMO | Green |
| Scrap Rate | Control | 4.5% | 1.8% | 1.0% | % | Yellow |

### Scrap Rate (Yellow)

Percentage of materials scrapped due to defects

| Measure | Value |
|:--------|------:|
| Baseline | 4.5% |
| Current | 1.8% |
| Target | 1.0% |
| Trend direction | Lower is better |
| Warning threshold | 2.0% |
| Critical threshold | 3.5% |

**Control limits**

| Limit | Value |
|:------|------:|
| Upper control limit (UCL) | 3.0% |
| Center line | 1.5% |
| Lower control limit (LCL) | 0.5% |
| Sigma | 0.50 |

## Equipment Performance

Metrics tracking equipment reliability and efficiency

**Owner:** James Wilson · **Health:** 83%

| Metric | Phase | Baseline | Current | Target | Unit | Status |
|:-------|:-----:|---------:|--------:|-------:|:-----|:------:|
| Overall Equipment Effectiveness | Improve | 65.0% | 82.5% | 85.0% | % | Green |
| Mean Time Between Failures | Measure | 120 | 340 | 500 | hours | Green |
| Unplanned Downtime | Improve | 12.0% | 4.2% | 2.0% | % | Yellow |

### Unplanned Downtime (Yellow)

Percentage of production time lost to unplanned equipment stops

| Measure | Value |
|:--------|------:|
| Baseline | 12.0% |
| Current | 4.2% |
| Target | 2.0% |
| Trend direction | Lower is better |
| Warning threshold | 5.0% |
| Critical threshold | 8.0% |

## Process Control

Statistical process control metrics for critical processes

**Owner:** Sarah Kim · **Health:** 83%

| Metric | Phase | Baseline | Current | Target | Unit | Status |
|:-------|:-----:|---------:|--------:|-------:|:-----|:------:|
| Process Capability Index (Cpk) | Control | 0.80 | 1.25 | 1.33 | Cpk | Green |
| Cycle Time Variation | Improve | 15 | 6.50 | 5 | seconds | Yellow |
| SPC Chart Compliance | Control | 60.0% | 92.0% | 100.0% | % | Green |

### Cycle Time Variation (Yellow)

Standard deviation of production cycle times

| Measure | Value |
|:--------|------:|
| Baseline | 15 |
| Current | 6.50 |
| Target | 5 |
| Trend direction | Lower is better |
| Warning threshold | 8 |
| Critical threshold | 12 |

## DMAIC Phase Distribution

| Phase | Description | Metrics |
|:------|:------------|--------:|
| **D** - Define | Define the problem and project goals | 0 |
| **M** - Measure | Measure current performance | 1 |
| **A** - Analyze | Analyze root causes | 0 |
| **I** - Improve | Implement improvements | 4 |
| **C** - Control | Sustain the gains | 4 |

## Improvement Initiatives

### Six Sigma Black Belt Project - Line 3

Reduce defect rate on production line 3 through DMAIC methodology

| Field | Value |
|:------|:------|
| **Owner** | Tom Chen |
| **Status** | In Progress |
| **Start date** | 2025-01-01 |
| **End date** | 2025-06-30 |
| **Expected impact** | Achieve 4.5 sigma level with DPMO below 3400 |
| **Actual impact** | Currently at 4.0 sigma, on track for target |

**Linked metrics:** First Pass Yield, Defects Per Million Opportunities, Scrap Rate

### Predictive Maintenance Implementation

Deploy IoT sensors and ML models for predictive equipment maintenance

| Field | Value |
|:------|:------|
| **Owner** | James Wilson |
| **Status** | In Progress |
| **Start date** | 2024-10-01 |
| **End date** | 2025-04-30 |
| **Expected impact** | Reduce unplanned downtime by 60% |

**Linked metrics:** Overall Equipment Effectiveness, Mean Time Between Failures, Unplanned Downtime

## Appendix: Process Capability Metrics

| Metric | Cp | Cpk | Sigma level | DPMO |
|:-------|---:|----:|------------:|-----:|
| First Pass Yield | 1.45 | 1.33 | 4 | 6210 |

### Control Limits Summary

| Metric | UCL | Center | LCL | In control |
|:-------|----:|-------:|----:|:----------:|
| First Pass Yield | 99.0% | 95.0% | 90.0% | Yes |
| Scrap Rate | 3.0% | 1.5% | 0.5% | Yes |
| Process Capability Index (Cpk) | 1.67 | 1.33 | 1 | Yes |
