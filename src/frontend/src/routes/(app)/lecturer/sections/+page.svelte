<script lang="ts">
  import { page } from "$app/state";
  import { meetingLabel } from "#lib/domain.ts";
  import PageHeader from "#lib/components/PageHeader.svelte";
  import TermSelect from "#lib/components/TermSelect.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
  const termId = $derived(data.termId);
  const sections = $derived(data.sections.filter((s) => s.term.id === termId));
</script>

<PageHeader title="My sections" subtitle="Your assigned classes and gradebooks.">
  <TermSelect terms={data.terms} value={termId} />
</PageHeader>
{#each sections as section}
  <article class="record">
    <div>
      <span class="course-code">{section.course.code} · Section {section.section_code}</span>
      <h2>{section.course.name}</h2>
      <p class="muted">
        {section.term.name} · {section.enrolled_count} students · {section.course.credits} credits
      </p>
      {#each section.meetings as meeting}
        <p class="small">{meetingLabel(meeting)}</p>
      {/each}
    </div>
    <a class="button" href="/lecturer/sections/{section.id}/grades">Open gradebook →</a>
  </article>
{/each}
{#if !sections.length}
  <div class="empty">No sections are assigned to you in this semester.</div>
{/if}
