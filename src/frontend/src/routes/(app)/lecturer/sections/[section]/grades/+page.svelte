<script lang="ts">
  import PageHeader from "#lib/components/PageHeader.svelte";
  import Gradebook from "#lib/components/Gradebook.svelte";
  import { timestamp } from "#lib/domain.ts";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<a class="breadcrumb" href="/lecturer/sections?term={data.section.term.id}">← My sections</a>
<PageHeader
  title={data.section.course.name}
  subtitle="{data.section.course.code} · Section {data.section.section_code} · {data.section.term
    .name}"
/>
{#key data.section.id}
  <Gradebook section={data.section} initial={data.roster} initialForm={data.form} />
{/key}
<details class="panel">
  <summary>Saved grade history</summary>
  <p class="small muted">History is refreshed after a confirmed save.</p>
  {#each data.history as revision}
    <p class="small">
      {revision.full_name} · {revision.score.toFixed(2)} · {revision.graded_by.full_name} · {timestamp(
        revision.created_at,
      )}
    </p>
  {/each}
</details>
