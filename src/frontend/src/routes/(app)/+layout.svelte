<script lang="ts">
  import { page, navigating } from "$app/state";
  import { enhance, type SubmitFunction } from "$app/forms";
  import { roleHome } from "#lib/domain.ts";
  import { provideDraft } from "#lib/draft.svelte.ts";
  import Notice from "#lib/components/Notice.svelte";
  import type { LayoutProps } from "./$types";
  let { data, children }: LayoutProps = $props();
  let menu: HTMLDialogElement;
  let errorMessage = $state("");
  let signingOut = $state(false);
  const links = {
    student: [
      ["Today", "/student"],
      ["Registration", "/student/registration"],
      ["Schedule", "/student/schedule/classes"],
      ["Results", "/student/results"],
      ["Tuition", "/student/tuition"],
      ["Documents", "/student/documents"],
    ],
    lecturer: [["My sections", "/lecturer/sections"]],
    staff: [
      ["Document requests", "/staff/documents"],
      ["Tuition records", "/staff/tuition"],
    ],
  };
  function active(href: string) {
    return href === "/student"
      ? page.url.pathname === href
      : page.url.pathname.startsWith(href.replace("/classes", ""));
  }
  const draft = provideDraft();
  const signOut: SubmitFunction = ({ cancel }) => {
    if (draft.dirty && !window.confirm("Sign out and discard your unsaved changes?")) {
      cancel();
      return;
    }
    draft.dirty = false;
    signingOut = true;
    return async ({ result, update }) => {
      if (result.type === "error") errorMessage = result.error.message ?? "Could not sign out.";
      else await update();
      signingOut = false;
    };
  };
</script>

{#snippet brand()}
  <!-- prettier-ignore -->
  <a class="brand" href={roleHome[data.user.role]}>Portal<span>67</span></a>
  <p class="prototype">Course prototype</p>
{/snippet}
{#snippet navigation()}
  <nav aria-label="Main navigation">
    {#each links[data.user.role] as [label, href]}
      <a
        {href}
        class:active={active(href)}
        aria-current={active(href) ? "page" : undefined}
        onclick={() => menu?.close()}
      >
        {label}
      </a>
    {/each}
  </nav>
{/snippet}
<a class="skip-link" href="#main">Skip to content</a>
<aside class="sidebar">
  <div>{@render brand()}{@render navigation()}</div>
  <div class="account-block">
    <strong>{data.user.full_name}</strong>
    <span class="role-label">{data.user.role}</span>
    <a href="/account">Account & help</a>
    <form method="POST" action="/sign-out" use:enhance={signOut}>
      <button class="text-button" disabled={signingOut}>
        {signingOut ? "Signing out…" : "Sign out"}
      </button>
    </form>
  </div>
</aside>
<div class="mobile-bar">
  <!-- prettier-ignore -->
  <div class="brand">Portal<span>67</span></div>
  <span class="role-label">{data.user.role}</span>
  <button onclick={() => menu.showModal()}>Menu</button>
</div>
<dialog class="menu-dialog" bind:this={menu} aria-label="Navigation">
  <button autofocus onclick={() => menu.close()}>Close menu</button>
  {@render brand()}{@render navigation()}
  <p>{data.user.full_name} · {data.user.role}</p>
  <a href="/account" onclick={() => menu.close()}>Account & help</a>
  <form method="POST" action="/sign-out" use:enhance={signOut}>
    <button class="text-button" disabled={signingOut} onclick={() => menu.close()}>Sign out</button>
  </form>
</dialog>
<main id="main" class="workspace" tabindex="-1" aria-busy={!!navigating.to}>
  {#if navigating.to}
    <p class="navigation-status" role="status">Loading page…</p>
  {/if}
  {#if errorMessage}
    <Notice tone="error">{errorMessage}</Notice>
  {/if}
  {@render children()}
  <footer class="workspace-footer">
    Portal67 <span>Sample academic records · Campus time UTC+7</span>
  </footer>
</main>

<style>
  .sidebar {
    width: 208px;
    background: white;
    border-right: 1px solid var(--line);
    position: fixed;
    inset: 0 auto 0 0;
    padding: 32px 20px 24px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    z-index: 10;
  }
  .sidebar nav,
  .menu-dialog nav {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .sidebar nav a,
  .menu-dialog nav a {
    padding: 11px 12px;
    text-decoration: none;
    color: var(--muted);
    font-size: 14px;
    min-height: 44px;
    border-left: 3px solid transparent;
    border-radius: 0 5px 5px 0;
  }
  .sidebar nav a.active,
  .menu-dialog nav a.active {
    background: var(--tint);
    border-left-color: var(--ink);
    color: var(--ink);
    font-weight: 600;
  }
  .account-block {
    border-top: 1px solid var(--line);
    padding-top: 20px;
    font-size: 13px;
    overflow-wrap: anywhere;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .account-block a {
    margin-top: 12px;
    min-height: 30px;
  }
  .mobile-bar {
    display: none;
  }
  @media (max-width: 900px) {
    .sidebar {
      display: none;
    }
    .mobile-bar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      padding: 16px 24px;
      border-bottom: 1px solid var(--line);
      background: white;
    }
    .mobile-bar .brand {
      font-size: 24px;
    }
    .mobile-bar .role-label {
      margin-left: auto;
    }
  }
  @media (max-width: 620px) {
    .mobile-bar {
      padding: 12px 16px;
    }
  }
  .menu-dialog {
    position: fixed;
    inset: 0 auto 0 0;
    margin: 0;
    width: 288px;
    max-width: calc(100vw - 32px);
    height: 100dvh;
    max-height: 100dvh;
    border-radius: 0;
    padding: 24px;
  }
  .menu-dialog > .brand {
    display: block;
    margin-top: 24px;
  }
  .menu-dialog .prototype {
    margin-bottom: 24px;
  }
  .menu-dialog > p {
    margin-top: 24px;
  }
  .workspace {
    margin-left: 208px;
    padding: 44px 32px 0;
    max-width: 1800px;
    min-height: 100dvh;
    outline: none;
    display: flex;
    flex-direction: column;
  }
  /* Page content comes from child routes, so it is matched globally. */
  .workspace > :global(:not(.workspace-footer)) {
    flex-shrink: 0;
  }
  @media (min-width: 1500px) {
    .workspace {
      padding-left: 48px;
      padding-right: 48px;
    }
  }
  @media (max-width: 900px) {
    .workspace {
      margin-left: 0;
      padding: 28px 24px 0;
    }
  }
  @media (max-width: 620px) {
    .workspace {
      padding: 24px 16px 0;
    }
  }
  .workspace-footer {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    color: var(--muted);
    border-top: 1px solid var(--line);
    font-size: 12px;
    padding: 24px 0;
    margin-top: 64px;
  }
  @media (max-width: 620px) {
    .workspace-footer {
      display: block;
      font-size: 11px;
      margin-top: 40px;
    }
    .workspace-footer span {
      display: block;
      margin-top: 8px;
    }
  }
  .skip-link {
    position: fixed;
    top: -80px;
    left: 24px;
    z-index: 100;
    background: var(--ink);
    color: white;
    padding: 12px 20px;
  }
  .skip-link:focus {
    top: 12px;
  }
  .navigation-status {
    position: fixed;
    top: 12px;
    right: 20px;
    background: white;
    border: 1px solid var(--line);
    padding: 8px 16px;
    z-index: 20;
    font-size: 13px;
  }
</style>
