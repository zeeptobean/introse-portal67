<script lang="ts">
  import SectionPreview from "#lib/components/SectionPreview.svelte";
  import PageHeader from "#lib/components/PageHeader.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
</script>

<a class="breadcrumb" href="/student/registration?term={data.section.term.id}">← Registration</a>
<PageHeader title={data.section.course.name} subtitle={data.section.term.name} />
<div class="form-width">
  <SectionPreview initial={data.registerForm} section={data.section} timetable={data.timetable} />
</div>
