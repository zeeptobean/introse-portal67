<script lang="ts">
  import { page } from "$app/state";
  import PageHeader from "#lib/components/PageHeader.svelte";
  import TermSelect from "#lib/components/TermSelect.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
  const available = $derived(data.results.terms.map((t) => t.term));
  const termId = $derived(data.termId);
  const selected = $derived(data.results.terms.find((t) => t.term.id === termId));
  const credits = $derived(selected?.courses.reduce((n, c) => n + c.course.credits, 0) ?? 0);
</script>

<PageHeader title="Results" subtitle="Your scores, with the calculation behind them.">
  {#if available.length}
    <TermSelect terms={available} value={termId} />
  {/if}
</PageHeader>
<div class="gpa-summary">
  <section>
    <p class="eyebrow">Semester · Estimated GPA</p>
    <!-- prettier-ignore -->
    <div class="figure">{selected?.gpa?.toFixed(2) ?? "Not available"}<span>/ 4.00</span></div>
    <p class="muted">{selected?.graded_credits ?? 0} of {credits} credits graded</p>
  </section>
  <section>
    <p class="eyebrow">Cumulative · Estimated GPA</p>
    <!-- prettier-ignore -->
    <div class="figure">{data.results.cumulative_gpa?.toFixed(2) ?? "Not available"}<span>/ 4.00</span></div>
    <p class="muted">Across available semesters · {data.results.graded_credits} graded credits</p>
  </section>
</div>
<h2>Course results</h2>
<div class="table-wrap">
  <table>
    <thead>
      <tr>
        <th>Course</th>
        <th>Section</th>
        <th class="numeric">Credits</th>
        <th class="numeric">Score / 10</th>
        <th>Letter</th>
        <th class="numeric">Points / 4</th>
      </tr>
    </thead>
    <tbody>
      {#each selected?.courses ?? [] as result}
        <tr>
          <td data-label="Course">
            <span class="course-code">{result.course.code}</span>
            <strong>{result.course.name}</strong>
          </td>
          <td data-label="Section">{result.section_code}</td>
          <td data-label="Credits" class="numeric">{result.course.credits}</td>
          <td data-label="Score / 10" class="numeric">
            {result.score?.toFixed(2) ?? "Not graded"}
          </td>
          <td data-label="Letter">{result.letter_grade ?? "—"}</td>
          <td data-label="Points / 4" class="numeric">{result.grade_points?.toFixed(2) ?? "—"}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
{#if !selected?.courses.length}
  <div class="empty">No results are available. Register for a section to see it here.</div>
{/if}
<details class="panel">
  <summary>How this estimated GPA is calculated</summary>
  <p>
    Sum of (grade points × course credits), divided by graded credits. Ungraded courses are
    excluded; a score of zero is graded and included. Results are rounded to two decimal places. All
    available graded enrollments count.
  </p>
  <p>
    Current seeded sample scale: 0–&lt;4: F / 0; 4–&lt;5: D / 1; 5–&lt;6.5: C / 2; 6.5–&lt;8: B / 3;
    8–&lt;9: A / 3.5; 9–10: A+ / 4.
  </p>
  <p class="muted small">
    This is a course prototype, not official HCMUS grading policy. Bands describe the repository
    seed “Sample 10-point scale”; results and points come from the API. If that scale is changed,
    update this explanation with it.
  </p>
</details>

<style>
  .gpa-summary {
    display: grid;
    grid-template-columns: 1fr 1fr;
    background: white;
    border: 1px solid var(--line);
    margin-bottom: 32px;
    border-radius: 8px;
  }
  .gpa-summary > section {
    padding: 28px 32px;
  }
  .gpa-summary > section + section {
    border-left: 1px solid var(--line);
  }
  .gpa-summary p:last-child {
    margin-bottom: 0;
  }
  .figure {
    font-size: 36px;
    font-weight: 600;
    letter-spacing: -1px;
    font-variant-numeric: tabular-nums;
    color: var(--ink);
  }
  .figure > span {
    font-size: 14px;
    color: var(--muted);
    font-weight: 400;
    letter-spacing: 0;
    margin-left: 12px;
  }
  @media (max-width: 620px) {
    .gpa-summary {
      grid-template-columns: 1fr;
    }
    .gpa-summary > section {
      padding: 24px;
    }
    .gpa-summary > section + section {
      border-left: 0;
      border-top: 1px solid var(--line);
    }
    .figure {
      font-size: 32px;
    }
  }
</style>
