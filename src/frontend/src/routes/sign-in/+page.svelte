<script lang="ts">
  import { untrack } from "svelte";
  import { zod4Client } from "sveltekit-superforms/adapters";
  import { portalForm } from "#lib/form.ts";
  import { loginSchema } from "#lib/schemas/auth.ts";
  import FormFeedback from "#lib/components/FormFeedback.svelte";
  import type { PageProps } from "./$types";
  let { data }: PageProps = $props();
  let show = $state(false);
  const { form, errors, message, allErrors, submitting, enhance } = portalForm(
    untrack(() => data.form),
    {
      validators: zod4Client(loginSchema),
    },
  );
</script>

<svelte:head><title>Sign in · Portal67</title></svelte:head>
<main class="sign-in">
  <!-- prettier-ignore -->
  <a class="brand" href="/">Portal<span>67</span></a>
  <p class="prototype">Course prototype</p>
  <div class="sign-in-heading">
    <p class="eyebrow">Your academic workspace</p>
    <h1>Sign in</h1>
    <p class="muted">Sample academic records. Not connected to the official HCMUS Portal.</p>
  </div>
  <FormFeedback message={$message} errors={$allErrors} />
  <form method="POST" use:enhance novalidate>
    <label for="email">Email</label>
    <input
      id="email"
      name="email"
      type="email"
      autocomplete="username"
      bind:value={$form.email}
      aria-invalid={!!$errors.email}
      aria-describedby="email-error"
    />
    <p id="email-error" class="field-error">{$errors.email?.join(" ") ?? ""}</p>
    <label for="password">Password</label>
    <div class="password-field">
      <input
        id="password"
        name="password"
        type={show ? "text" : "password"}
        autocomplete="current-password"
        bind:value={$form.password}
        aria-invalid={!!$errors.password}
        aria-describedby="password-error"
      />
      <button type="button" aria-pressed={show} onclick={() => (show = !show)}>
        {show ? "Hide" : "Show"}
      </button>
    </div>
    <p class="field-error" id="password-error">{$errors.password?.join(" ") ?? ""}</p>
    <button class="primary wide" disabled={$submitting}>
      {$submitting ? "Signing in…" : "Sign in"}
    </button>
  </form>
  <div class="help-note">
    <strong>Need an account?</strong>
    <p>Use the demo account provided by your project team.</p>
  </div>
</main>

<style>
  .sign-in .prototype {
    margin-bottom: 64px;
  }
  @media (max-width: 620px) {
    .sign-in .prototype {
      margin-bottom: 40px;
    }
  }
  .sign-in-heading {
    margin-bottom: 32px;
  }
  .sign-in form > label:first-child {
    margin-top: 0;
  }
  .password-field {
    display: flex;
    gap: 8px;
  }
  .password-field input {
    min-width: 0;
  }
  .password-field button {
    flex-shrink: 0;
  }
  .sign-in .wide {
    margin-top: 28px;
  }
  .help-note {
    padding-top: 24px;
    border-top: 1px solid var(--line);
    margin-top: 32px;
    font-size: 14px;
  }
  .help-note p {
    margin-top: 8px;
    color: var(--muted);
  }
</style>
