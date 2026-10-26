<script lang="ts">
  import { untrack } from "svelte";
  import { zod4Client } from "sveltekit-superforms/adapters";
  import { portalForm } from "#lib/form.ts";
  import { studentSearchSchema } from "#lib/schemas/tuition.ts";
  import PageHeader from "#lib/components/PageHeader.svelte";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  import type { PageProps } from "./$types";
  let { data, form: result }: PageProps = $props();
  const { form, message, allErrors, submitting, enhance } = portalForm(
    untrack(() => data.form),
    {
      validators: zod4Client(studentSearchSchema),
      invalidateAll: false,
    },
  );
  const students = $derived(result && "students" in result ? result.students : undefined);
</script>

<PageHeader
  title="Tuition records"
  subtitle="Find a student, confirm their identity, then review their ledger."
/>
<form class="toolbar" method="POST" use:enhance>
  <label class="grow">
    Student name, institutional ID, or email
    <input type="search" name="q" bind:value={$form.q} placeholder="Search student records" />
  </label>
  <button class="primary" disabled={$submitting}>
    {$submitting ? "Searching…" : "Find student"}
  </button>
</form>
<FormFeedback message={$message} errors={$allErrors} />
{#if students === undefined}
  <div class="empty">Search to begin. No student is selected.</div>
{:else}
  {#each students as student}
    <article class="record">
      <div>
        <h2>{student.full_name}</h2>
        <p>{student.institutional_id ?? "ID not provided"} · {student.email}</p>
      </div>
      <a class="button" href="/staff/tuition/students/{student.id}">View tuition →</a>
    </article>
  {/each}
  {#if !students.length}
    <div class="empty">No students match your search. Try an institutional ID or email.</div>
  {/if}
{/if}
