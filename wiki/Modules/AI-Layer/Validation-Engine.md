# Validation Engine

## Purpose

The Validation Engine verifies AI-generated outputs before they are applied to the codebase. It checks that patches are syntactically valid, confidence scores meet thresholds, and quality gates are satisfied.

## Inputs

- AI-generated unified diff patch
- Confidence score from the LLM
- Quality gate rules for the target repository
- Solo Gate enforcement mode configuration

## Processing

### Current Implementation

The current validation is minimal:

1. **Response parsing** — The parser (`app/services/aiworker/parser.go`) extracts the diff block and confidence score from the LLM response
2. **Diff extraction** — Validates that a unified diff block exists in the response
3. **Confidence extraction** — Parses the confidence score as a floating-point value
4. **Status determination** — Sets status to `completed` if a diff was produced, `failed` otherwise

### Planned Validation Steps

The following validation stages are designed but not yet implemented:

1. **Patch syntax validation** — Verify the diff applies cleanly to the target file
2. **Confidence threshold check** — Compare against configurable minimum confidence
3. **Quality gate re-evaluation** — Run the proposed change against existing quality rules
4. **Test execution** — Apply the patch in a sandboxed environment and run tests
5. **Security scan** — Check the proposed change for new security issues

## Outputs

- Validation result (pass/fail)
- Reasons for failure (if applicable)
- Approved patch ready for application (if passed)

## Key Paths

| Purpose | Path |
|---------|------|
| Response parser | `app/services/aiworker/parser.go` |
| Solo Gate evaluation | `app/services/sologate/engine.go` |

## Status

**Concept** — Basic response parsing and diff extraction are implemented. Full validation pipeline (patch application check, test execution, security scan) is planned.

## Future Work

- Automated patch application verification
- Test suite execution against proposed patches
- Security scanning of proposed code changes
- Confidence threshold configuration per space
- Multi-stage validation pipeline
